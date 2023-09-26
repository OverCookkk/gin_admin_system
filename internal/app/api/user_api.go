package api

import (
	"gin_admin_system/internal/app/service"
	"gin_admin_system/internal/app/types"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var UserApiSet = wire.NewSet(wire.Struct(new(UserApi), "*"))

type UserApi struct {
	UserSrv *service.UserSrv
}

func (u *UserApi) Query(c *gin.Context) {
	var req types.UserQueryReq
	if err := c.ShouldBindQuery(req); err != nil {
		// 参数错误
		// response.ReturnWithDetailed()
		return
	}

	// req.Pagination = true

	// todo:UserSrv.QueryShow
	// result, err := u.UserSrv.QueryShow(c.Request.Context(), req)
	// if err != nil {
	// 	return
	// }
	// response.OkWithData(result, c)
}
