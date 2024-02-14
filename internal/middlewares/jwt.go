package middlewares

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/EvgeniyBudaev/love-server/internal/config"
	r "github.com/EvgeniyBudaev/love-server/internal/handler/http/api/v1/response"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/EvgeniyBudaev/love-server/internal/shared/enums"
	"github.com/Nerzal/gocloak/v13"
	contribJwt "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	golangJwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

type TokenRetrospector interface {
	RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error)
}

func NewJwtMiddleware(config *config.Config, tokenRetrospector TokenRetrospector, logger logger.Logger) fiber.Handler {
	base64Str := config.RealmRS256PublicKey
	publicKey, err := parseKeycloakRSAPublicKey(base64Str, logger)
	if err != nil {
		logger.Debug("error while NewJwtMiddleware. Error in parseKeycloakRSAPublicKey", zap.Error(err))
		panic(err)
	}
	return contribJwt.New(contribJwt.Config{
		SigningKey: contribJwt.SigningKey{
			JWTAlg: contribJwt.RS256,
			Key:    publicKey,
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return successHandler(c, tokenRetrospector, logger)
		},
	})
}

func successHandler(c *fiber.Ctx, tokenRetrospector TokenRetrospector, logger logger.Logger) error {
	jwtToken := c.Locals("user").(*golangJwt.Token)
	claims := jwtToken.Claims.(golangJwt.MapClaims)
	var ctx = c.UserContext()
	var contextWithClaims = context.WithValue(ctx, enums.ContextKeyClaims, claims)
	c.SetUserContext(contextWithClaims)
	rptResult, err := tokenRetrospector.RetrospectToken(ctx, jwtToken.Raw)
	if err != nil {
		logger.Debug("error while successHandler. Error in RetrospectToken", zap.Error(err))
		return err
	}
	if !*rptResult.Active {
		err := fmt.Errorf("token is not active")
		logger.Debug("error while successHandler. Error in parseKeycloakRSAPublicKey", zap.Error(err))
		return r.WrapError(c, err, http.StatusUnauthorized)
	}
	return c.Next()
}

func parseKeycloakRSAPublicKey(base64Str string, logger logger.Logger) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		logger.Debug("error while parseKeycloakRSAPublicKey. Error in DecodeString", zap.Error(err))
		return nil, err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		logger.Debug("error while parseKeycloakRSAPublicKey. Error in ParsePKIXPublicKey", zap.Error(err))
		return nil, err
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	err = fmt.Errorf("unexpected key type %T", publicKey)
	logger.Debug("error while parseKeycloakRSAPublicKey", zap.Error(err))
	return nil, err
}
