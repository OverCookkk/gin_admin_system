//go:build wireinject

package app

import (
	"gin_admin_system/internal/app/dao"
	"github.com/google/wire"

	"gin_admin_system/internal/app/api"
	"gin_admin_system/internal/app/router"
	"gin_admin_system/internal/app/service"
)

func BuildWireInject() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		dao.RepoSet,
		service.ServiceSet,
		api.ApiSet,
		router.RouterSet,
		InitGinEngine,
		InitCasbin,
		CasbinAdapterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
