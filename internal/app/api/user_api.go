package api

import (
	"gin_admin_system/internal/app/response"
	"gin_admin_system/internal/app/service"
	"gin_admin_system/internal/app/types"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"strconv"
)

var UserApiSet = wire.NewSet(wire.Struct(new(UserApi), "*"))

type UserApi struct {
	UserSrv *service.UserSrv
}

func (u *UserApi) Query(c *gin.Context) {
	var req types.UserQueryReq
	if err := c.ShouldBindQuery(req); err != nil {
		// 参数错误
		response.JsonError(c, err)
	}

	// req.Pagination = true

	// todo:UserSrv.QueryShow
	// result, err := u.UserSrv.QueryShow(c.Request.Context(), req)
	// if err != nil {
	// 	return
	// }
	// response.OkWithData(result, c)
}

func (u *UserApi) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		id = 0
	}
	item, err := u.UserSrv.Get(c.Request.Context(), id)
	if err != nil {
		response.JsonError(c, err)
	}
	response.JsonData(c, item.CleanSecure())
}

func (u *UserApi) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item types.User
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	} else if item.Password == "" {
		response.JsonError(c, errors.New("password not empty"))
		return
	}

	result, err := u.UserSrv.Create(ctx, item)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonData(c, result)
}

func (u *UserApi) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item types.User
	if err := c.ShouldBindJSON(&item); err != nil {
		// 参数错误
		response.JsonError(c, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		id = 0
	}
	err = u.UserSrv.Update(ctx, id, item)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (u *UserApi) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		id = 0
	}
	err = u.UserSrv.Delete(ctx, id)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (u *UserApi) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		id = 0
	}
	err = u.UserSrv.UpdateStatus(ctx, id, 1)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}

func (u *UserApi) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 16)
	if err != nil {
		id = 0
	}
	err = u.UserSrv.UpdateStatus(ctx, id, 2)
	if err != nil {
		response.JsonError(c, err)
		return
	}
	response.JsonSuccess(c)
}
