package api

import (
    "gin_admin_system/internal/app/service"
    "github.com/gin-gonic/gin"
    "github.com/google/wire"
)

var MenuApiSet = wire.NewSet(wire.Struct(new(MenuApi), "*"))

type MenuApi struct {
    MenuSrv *service.MenuSrv
}

func (m *MenuApi) Create(c *gin.Context) {

}
