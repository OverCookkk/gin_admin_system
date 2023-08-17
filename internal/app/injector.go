package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

// 初始化后最终生成的对象
type Injector struct {
	GinEngine *gin.Engine
}
