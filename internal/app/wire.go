//go:build wireinject

package app

import "github.com/google/wire"
import "gin_admin_system/internal/app/dao/menu"

func BuildWireInject() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		menu.MenuSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
