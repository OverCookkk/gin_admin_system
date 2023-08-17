// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"gin_admin_system/internal/app/dao/menu"
)

// Injectors from wire.go:

func BuildWireInject() (*Injector, func(), error) {
	db, cleanup, err := InitGormDB()
	if err != nil {
		return nil, nil, err
	}
	menuRepo := menu.MenuRepo{
		DB: db,
	}
	injector := &Injector{
		MenuRepo: menuRepo,
	}
	return injector, func() {
		cleanup()
	}, nil
}