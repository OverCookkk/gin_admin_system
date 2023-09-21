package service

import (
	"context"
	"errors"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/dao/role"
	"gin_admin_system/internal/app/dao/user"
	"gin_admin_system/internal/app/types"
	"gin_admin_system/pkg/auth"
	"gin_admin_system/pkg/util/hash"
	"github.com/google/wire"
)

var LoginSet = wire.NewSet(wire.Struct(new(LoginSrv), "*"))

type LoginSrv struct {
	Auth           auth.JWTAuth
	UserRepo       *user.UserRepo
	UserRoleRepo   *user.UserRoleRepo
	RoleRepo       *role.RoleRepo
	RoleMenuRepo   *role.RoleMenuRepo
	MenuRepo       *menu.MenuRepo
	MenuActionRepo *menu.MenuActionRepo
}

func (l *LoginSrv) Verify(ctx context.Context, userName, password string) (*types.User, error) {
	result, err := l.UserRepo.Query(ctx, types.UserQueryReq{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.New("not found user_name")
	}
	userItem := result.Data[0]
	if userItem.Password != hash.SHA1String(password) {
		return nil, errors.New("password incorrect")
	} else if userItem.Status != 1 {
		return nil, errors.New("user forbidden")
	}

	return &userItem, nil
}

func (l *LoginSrv) GenerateToken(ctx context.Context, userID string) (string, error) {
	token, err := l.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return "", errors.New("GenerateToken failed")
	}

	return token, nil
}

func (l *LoginSrv) DestroyToken(ctx context.Context, tokenString string) error {
	err := l.Auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return err
	}
	return nil
}
