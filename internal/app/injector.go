package app

import (
	"gin_admin_system/internal/app/dao/menu"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

// 初始化后最终生成的对象
type Injector struct {
	// GinEngine *gin.Engine
	MenuRepo menu.MenuRepo
}
