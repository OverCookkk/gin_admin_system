package api

import (
	"gin_admin_system/internal/app"
	"gin_admin_system/internal/app/service"
	"gin_admin_system/internal/app/types"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"strconv"
)

var MenuApiSet = wire.NewSet(wire.Struct(new(MenuApi), "*"))

type MenuApi struct {
	MenuSrv *service.MenuSrv
}

func (m *MenuApi) Query(c *gin.Context) {
	var req types.MenuQueryReq
	if err := c.ShouldBindQuery(req); err != nil {
		// 参数错误
		// app.ReturnWithDetailed()
		return
	}

	// 此处封装了NewOrderFields和NewOrderField两个函数，巧妙在NewOrderField函数参数使用...传切片的特性，使得可以直接生成一个切片结构体
	result, err := m.MenuSrv.Query(c.Request.Context(), req, types.MenuQueryOptions{
		OrderFields: types.NewOrderFields(
			types.NewOrderField("sequence", types.OrderByDesc),
			types.NewOrderField("id", types.OrderByDesc),
		),
	})
	if err != nil {
		return
	}
	app.OkWithData(result, c)
}

// QueryMenuTree 返回菜单树，包括孩子树
func (m *MenuApi) QueryMenuTree(c *gin.Context) {

}

func (m *MenuApi) Get(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}
	menuItem, err := m.MenuSrv.Get(c.Request.Context(), id)
	if err != nil {
		return
	}
	app.OkWithData(menuItem, c)
}

func (m *MenuApi) Create(c *gin.Context) {
	var item types.Menu
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		// app.ReturnWithDetailed()
		return
	}
	idResult, err := m.MenuSrv.Create(c.Request.Context(), item)
	if err != nil {
		return
	}
	app.OkWithData(idResult, c)
}
