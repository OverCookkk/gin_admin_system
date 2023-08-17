package service

import (
	"gin_admin_system/internal/app/dao/menu"
	"github.com/google/wire"
)

var MenuSrvSet = wire.NewSet(wire.Struct(new(MenuSrv), "*"))

type MenuSrv struct {
	MenuRepo *menu.MenuRepo
}
