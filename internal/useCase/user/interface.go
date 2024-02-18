package user

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
)

type Identity interface {
	CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error)
	UpdateUser(ctx context.Context, user gocloak.User) (*gocloak.User, error)
	DeleteUser(ctx context.Context, user gocloak.User) error
	GetUserList(ctx context.Context, query QueryParamsUserList) ([]*gocloak.User, error)
}
