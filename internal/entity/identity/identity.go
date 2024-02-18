package identity

import (
	"context"
	"fmt"
	"github.com/EvgeniyBudaev/love-server/internal/config"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/EvgeniyBudaev/love-server/internal/useCase/user"
	"github.com/Nerzal/gocloak/v13"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

type Identity struct {
	BaseUrl      string
	Realm        string
	ClientId     string
	ClientSecret string
	logger       logger.Logger
}

func NewIdentity(config *config.Config, l logger.Logger) *Identity {
	return &Identity{
		BaseUrl:      config.BaseUrl,
		Realm:        config.Realm,
		ClientId:     config.ClientId,
		ClientSecret: config.ClientSecret,
		logger:       l,
	}
}

func (i *Identity) loginRestApiClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(i.BaseUrl)
	token, err := client.LoginClient(ctx, i.ClientId, i.ClientSecret, i.Realm)
	if err != nil {
		i.logger.Debug("error unable to login the rest client by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to login the rest client")
	}
	return token, nil
}

func (i *Identity) CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error) {
	token, err := i.loginRestApiClient(ctx)
	if err != nil {
		return nil, err
	}
	client := gocloak.NewClient(i.BaseUrl)
	isUniqueMobileNumber, err := i.validateMobileNumbers(ctx, (*user.Attributes)["mobileNumber"], token, client)
	if err != nil {
		i.logger.Debug("error get users for validation mobile number is invalid by path"+
			" entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "get users for validation mobile number is invalid")
	}
	if !isUniqueMobileNumber {
		i.logger.Debug("error mobile number must be unique by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.New("mobile number must be unique")
	}
	userId, err := client.CreateUser(ctx, token.AccessToken, i.Realm, user)
	if err != nil {
		i.logger.Debug("error unable to create the user by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to create the user")
	}
	err = client.SetPassword(ctx, token.AccessToken, userId, i.Realm, password, false)
	if err != nil {
		i.logger.Debug("error unable to set the password for the user by path"+
			" entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to set the password for the user")
	}
	var roleNameLowerCase = strings.ToLower(role)
	roleKeycloak, err := client.GetRealmRole(ctx, token.AccessToken, i.Realm, roleNameLowerCase)
	if err != nil {
		i.logger.Debug("error unable to get role by name by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, fmt.Sprintf("unable to get role by name: '%v'", roleNameLowerCase))
	}
	err = client.AddRealmRoleToUser(ctx, token.AccessToken, i.Realm, userId, []gocloak.Role{
		*roleKeycloak,
	})
	if err != nil {
		i.logger.Debug("error unable to add a realm role to user by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to add a realm role to user")
	}
	userKeycloak, err := client.GetUserByID(ctx, token.AccessToken, i.Realm, userId)
	if err != nil {
		i.logger.Debug("error unable to get recently created user by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to get recently created user")
	}
	return userKeycloak, nil
}

func (i *Identity) UpdateUser(ctx context.Context, user gocloak.User) (*gocloak.User, error) {
	token, err := i.loginRestApiClient(ctx)
	if err != nil {
		return nil, err
	}
	client := gocloak.NewClient(i.BaseUrl)
	isUniqueMobileNumber, err := i.validateMobileNumbers(ctx, (*user.Attributes)["mobileNumber"], token, client)
	if err != nil {
		i.logger.Debug("error get users for validation mobile number is invalid by path"+
			" entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "get users for validation mobile number is invalid")
	}
	if !isUniqueMobileNumber {
		i.logger.Debug("error mobile number must be unique by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.New("mobile number must be unique")
	}
	err = client.UpdateUser(ctx, token.AccessToken, i.Realm, user)
	if err != nil {
		i.logger.Debug("error unable to update the user by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to update the user")
	}
	if user.ID == nil {
		return nil, errors.New("user ID is nil")
	}
	userKeycloak, err := client.GetUserByID(ctx, token.AccessToken, i.Realm, *user.ID)
	if err != nil {
		i.logger.Debug("error unable to get recently created user by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to get recently created user")
	}
	return userKeycloak, nil
}

func (i *Identity) DeleteUser(ctx context.Context, user gocloak.User) error {
	token, err := i.loginRestApiClient(ctx)
	if err != nil {
		return nil
	}
	client := gocloak.NewClient(i.BaseUrl)
	err = client.DeleteUser(ctx, token.AccessToken, i.Realm, *user.ID)
	if err != nil {
		i.logger.Debug("error unable to delete the user by path entity/identity/identity.go", zap.Error(err))
		return errors.Wrap(err, "unable to delete the user")
	}
	return nil
}

func (i *Identity) RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	client := gocloak.NewClient(i.BaseUrl)
	rptResult, err := client.RetrospectToken(ctx, accessToken, i.ClientId, i.ClientSecret, i.Realm)
	if err != nil {
		i.logger.Debug("error unable to retrospect token by path entity/identity/identity.go", zap.Error(err))
		return nil, errors.Wrap(err, "unable to retrospect token")
	}
	return rptResult, nil
}

func (i *Identity) GetUserList(ctx context.Context, query user.QueryParamsUserList) ([]*gocloak.User, error) {
	token, err := i.loginRestApiClient(ctx)
	if err != nil {
		return nil, err
	}
	client := gocloak.NewClient(i.BaseUrl)
	users, err := client.GetUsers(ctx, token.AccessToken, i.Realm, gocloak.GetUsersParams{Search: &query.Search})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (i *Identity) validateMobileNumbers(ctx context.Context, mobileNumberList []string, token *gocloak.JWT, client *gocloak.GoCloak) (bool, error) {
	uniqueMap := make(map[string]bool) // Для хранения уникальности каждого номера
	users, err := client.GetUsers(ctx, token.AccessToken, i.Realm, gocloak.GetUsersParams{})
	if err != nil {
		return false, err
	}
	// Инициализируем карту уникальности
	for _, mobileNumber := range mobileNumberList {
		uniqueMap[mobileNumber] = true
	}
	// Проверяем существующих пользователей
	for _, user := range users {
		mobileAttr, exists := (*user.Attributes)["mobileNumber"]
		if exists && len(mobileAttr) > 0 {
			// Если номер встречается у существующего пользователя, устанавливаем значение в карте уникальности как false
			if _, isInList := uniqueMap[mobileAttr[0]]; isInList {
				uniqueMap[mobileAttr[0]] = false
			}
		}
	}
	// Проверяем, все ли значения в карте уникальности равны true
	for _, isUnique := range uniqueMap {
		if !isUnique {
			return false, nil
		}
	}
	return true, nil
}
