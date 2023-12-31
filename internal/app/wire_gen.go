// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"gin_admin_system/internal/app/api"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/dao/role"
	"gin_admin_system/internal/app/dao/user"
	"gin_admin_system/internal/app/router"
	"gin_admin_system/internal/app/service"
)

// Injectors from wire.go:

func BuildWireInject() (*Injector, func(), error) {
	jwtAuth, cleanup, err := InitAuth()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := InitGormDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	roleRepo := &role.RoleRepo{
		DB: db,
	}
	roleMenuRepo := &role.RoleMenuRepo{
		DB: db,
	}
	menuActionResourceRepo := &menu.MenuActionResourceRepo{
		DB: db,
	}
	userRepo := &user.UserRepo{
		DB: db,
	}
	userRoleRepo := &user.UserRoleRepo{
		DB: db,
	}
	casbinAdapter := &CasbinAdapter{
		RoleRepo:         roleRepo,
		RoleMenuRepo:     roleMenuRepo,
		MenuResourceRepo: menuActionResourceRepo,
		UserRepo:         userRepo,
		UserRoleRepo:     userRoleRepo,
	}
	syncedEnforcer, cleanup3, err := InitCasbin(casbinAdapter)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	menuRepo := &menu.MenuRepo{
		DB: db,
	}
	menuActionRepo := &menu.MenuActionRepo{
		DB: db,
	}
	loginSrv := &service.LoginSrv{
		Auth:           jwtAuth,
		UserRepo:       userRepo,
		UserRoleRepo:   userRoleRepo,
		RoleRepo:       roleRepo,
		RoleMenuRepo:   roleMenuRepo,
		MenuRepo:       menuRepo,
		MenuActionRepo: menuActionRepo,
	}
	loginAPI := &api.LoginAPI{
		LoginSrv: loginSrv,
	}
	menuSrv := &service.MenuSrv{
		MenuRepo:               menuRepo,
		MenuActionRepo:         menuActionRepo,
		MenuActionResourceRepo: menuActionResourceRepo,
	}
	menuApi := &api.MenuApi{
		MenuSrv: menuSrv,
	}
	roleSrv := &service.RoleSrv{
		Enforcer:               syncedEnforcer,
		RoleRepo:               roleRepo,
		RoleMenuRepo:           roleMenuRepo,
		UserRepo:               userRepo,
		MenuActionResourceRepo: menuActionResourceRepo,
	}
	roleApi := &api.RoleApi{
		RoleSrv: roleSrv,
	}
	userSrv := &service.UserSrv{
		Enforcer:     syncedEnforcer,
		UserRepo:     userRepo,
		UserRoleRepo: userRoleRepo,
	}
	userApi := &api.UserApi{
		UserSrv: userSrv,
	}
	routerRouter := &router.Router{
		Auth:           jwtAuth,
		CasbinEnforcer: syncedEnforcer,
		LoginApi:       loginAPI,
		MenuApi:        menuApi,
		RoleApi:        roleApi,
		UserApi:        userApi,
	}
	engine := InitGinEngine(routerRouter)
	injector := &Injector{
		GinEngine: engine,
	}
	return injector, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
