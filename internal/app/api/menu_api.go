package api

import (
	"gin_admin_system/internal/app/response"
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

// Query 查询全部菜单大概信息
func (m *MenuApi) Query(c *gin.Context) {
	var req types.MenuQueryReq
	if err := c.ShouldBindQuery(req); err != nil {
		// 参数错误
		response.JsonError(c, err)
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
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, result)
}

// QueryMenuTree 返回菜单树，包括孩子树
func (m *MenuApi) QueryMenuTree(c *gin.Context) {

}

// Get 获取单个菜单的具体信息
func (m *MenuApi) Get(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}
	menuItem, err := m.MenuSrv.Get(c.Request.Context(), id)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, menuItem)
}

func (m *MenuApi) Create(c *gin.Context) {
	var item types.Menu
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	}
	idResult, err := m.MenuSrv.Create(c.Request.Context(), item)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, idResult)
}

func (m *MenuApi) Update(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}

	var item types.Menu
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	}
	err = m.MenuSrv.Update(c.Request.Context(), id, item)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (m *MenuApi) Delete(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}
	err = m.MenuSrv.Delete(c.Request.Context(), id)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (m *MenuApi) Enable(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}
	err = m.MenuSrv.UpdateStatus(c.Request.Context(), id, 1)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (m *MenuApi) Disable(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}
	err = m.MenuSrv.UpdateStatus(c.Request.Context(), id, 2)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}
