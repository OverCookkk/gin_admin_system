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
	"github.com/dchest/captcha"
	"github.com/google/wire"
	"net/http"
)

var LoginSrvSet = wire.NewSet(wire.Struct(new(LoginSrv), "*"))

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

func (l *LoginSrv) GenerateToken(ctx context.Context, userID string) (*types.LoginTokenInfo, error) {
	token, err := l.Auth.GenerateToken(ctx, userID)
	if err != nil {
		return nil, errors.New("GenerateToken failed")
	}

	return &types.LoginTokenInfo{AccessToken: token}, nil
}

func (l *LoginSrv) DestroyToken(ctx context.Context, tokenString string) error {
	err := l.Auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return err
	}
	return nil
}

func (l *LoginSrv) GetCaptcha(ctx context.Context, length int) (*types.LoginCaptcha, error) {
	captchaID := captcha.NewLen(length)
	return &types.LoginCaptcha{
		CaptchaID: captchaID,
	}, nil
}

func (l *LoginSrv) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	err := captcha.WriteImage(w, captchaID, width, height)
	if err != nil {
		return err
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")

	return nil
}

func (l *LoginSrv) GetLoginInfo(ctx context.Context, userID uint64) (*types.UserLoginInfo, error) {
	if types.CheckIsRootUser(ctx, userID) { // ROOT用户
		root := types.GetRootUser()
		return &types.UserLoginInfo{
			UserName: root.UserName,
			RealName: root.RealName,
		}, nil
	}

	user, err := l.UserRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.New("not found")
	} else if user.Status != 1 {
		return nil, errors.New("user forbidden")
	}

	userInfo := &types.UserLoginInfo{
		UserID:   string(userID),
		UserName: user.UserName,
		RealName: user.RealName,
		// Roles:    nil,
	}

	// 先查询该用户拥有哪些roleID
	userRoleResult, err := l.UserRoleRepo.Query(ctx, types.UserRoleQueryReq{
		PaginationParam: types.PaginationParam{},
		UserID:          userID,
	})
	if err != nil {
		return nil, err
	}

	// 取出该user所拥有的所有roleID，再去获取每个角色的信息
	if roleIDs := userRoleResult.Data.ToRoleIDs(); len(roleIDs) > 0 {
		roleResult, err := l.RoleRepo.Query(ctx, types.RoleQueryReq{
			PaginationParam: types.PaginationParam{},
			IDs:             roleIDs,
			Status:          1,
		})
		if err != nil {
			return nil, err
		}
		userInfo.Roles = roleResult.Data
	}

	return userInfo, nil
}
