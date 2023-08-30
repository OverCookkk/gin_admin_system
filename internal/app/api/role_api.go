package api

import (
	"gin_admin_system/internal/app"
	"gin_admin_system/internal/app/service"
	"gin_admin_system/internal/app/types"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var RoleApiSet = wire.NewSet(wire.Struct(new(RoleApi)), "*")

type RoleApi struct {
	RoleSrv *service.RoleSrv
}

func (m *RoleApi) Query(c *gin.Context) {
	var req types.RoleQueryReq
	if err := c.ShouldBindQuery(req); err != nil {
		// 参数错误
		// app.ReturnWithDetailed()
		return
	}

	// 此处封装了NewOrderFields和NewOrderField两个函数，巧妙在NewOrderField函数参数使用...传切片的特性，使得可以直接生成一个切片结构体
	result, err := m.RoleSrv.Query(c.Request.Context(), req, types.RoleQueryOptions{
		OrderFields: types.NewOrderFields(
			types.NewOrderField("sequence", types.OrderByDesc),
		),
	})
	if err != nil {
		return
	}
	app.OkWithData(result, c)
}
