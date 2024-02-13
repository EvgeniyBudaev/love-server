package user

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
)

type Identity interface {
	CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error)
	GetUserList(ctx context.Context, query QueryParamsUserList) ([]*gocloak.User, error)
}
