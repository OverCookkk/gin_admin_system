package api

import (
	"gin_admin_system/internal/app/response"
	"gin_admin_system/internal/app/service"
	"gin_admin_system/internal/app/types"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"strconv"
)

var RoleApiSet = wire.NewSet(wire.Struct(new(RoleApi), "*"))

type RoleApi struct {
	RoleSrv *service.RoleSrv
}

func (r *RoleApi) Query(c *gin.Context) {
	var req types.RoleQueryReq
	if err := c.ShouldBindQuery(req); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	}

	// 此处封装了NewOrderFields和NewOrderField两个函数，巧妙在NewOrderField函数参数使用...传切片的特性，使得可以直接生成一个切片结构体
	result, err := r.RoleSrv.Query(c.Request.Context(), req, types.RoleQueryOptions{
		OrderFields: types.NewOrderFields(
			types.NewOrderField("sequence", types.OrderByDesc),
		),
	})
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, result)
}

func (r *RoleApi) Get(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}
	menuItem, err := r.RoleSrv.Get(c.Request.Context(), id)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, menuItem)
}

func (r *RoleApi) Create(c *gin.Context) {
	var item types.Role
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	}
	idResult, err := r.RoleSrv.Create(c, item)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, idResult)
}

func (r *RoleApi) Update(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}

	var item types.Role
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	}
	err = r.RoleSrv.Update(c, id, item)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (r *RoleApi) Delete(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}

	err = r.RoleSrv.Delete(c, id)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (r *RoleApi) Enable(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}

	err = r.RoleSrv.UpdateStatus(c, id, 1)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (r *RoleApi) Disable(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseUint(idVal, 10, 64)
	if err != nil {
		id = 0
	}

	err = r.RoleSrv.UpdateStatus(c, id, 2)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}
