package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JsonResult struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"Message"`
	Success bool        `json:"success"`
}

func JsonData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, JsonResult{
		Code:    0,
		Data:    data,
		Success: true,
	})
}

func JsonSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, JsonResult{
		Code:    0,
		Data:    nil,
		Success: true,
	})
}

func JsonError(c *gin.Context, err error) {
	if e, ok := err.(*CodeError); ok {
		c.JSON(http.StatusOK, JsonResult{
			Code:    e.Code,
			Message: e.Message,
			Data:    e.Data,
			Success: false,
		})
	}
	c.JSON(http.StatusOK, JsonResult{
		Code:    0,
		Message: err.Error(),
		Data:    nil,
		Success: false,
	})
}

func JsonErrorMsg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, JsonResult{
		Code:    0,
		Message: message,
		Data:    nil,
		Success: false,
	})
}

func JsonErrorCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, JsonResult{
		Code:    code,
		Message: message,
		Data:    nil,
		Success: false,
	})
}

func JsonErrorData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, JsonResult{
		Code:    code,
		Message: message,
		Data:    data,
		Success: false,
	})
}

// func Ok(c *gin.Context) {
// 	Result(errmsg.SUCCSE, map[string]interface{}{}, errmsg.GetErrMsg(errmsg.SUCCSE), c)
// }
//
// func OkWithMessage(message string, c *gin.Context) {
// 	Result(errmsg.SUCCSE, map[string]interface{}{}, message, c)
// }
//
// func OkWithData(data interface{}, c *gin.Context) {
// 	Result(errmsg.SUCCSE, data, "查询成功", c)
// }
//
// func OkWithDetailed(data interface{}, message string, c *gin.Context) {
// 	Result(errmsg.SUCCSE, data, message, c)
// }
//
// func Fail(c *gin.Context) {
// 	Result(errmsg.ERROR, map[string]interface{}{}, "操作失败", c)
// }
//
// func FailWithMessage(message string, c *gin.Context) {
// 	Result(errmsg.ERROR, map[string]interface{}{}, message, c)
// }
//
// // 返回data结构体
// func FailWithDetailed(data interface{}, message string, c *gin.Context) {
// 	Result(errmsg.ERROR, data, message, c)
// }
//
// func ReturnWithDetailed(code int, data interface{}, message string, c *gin.Context) {
// 	Result(code, data, message, c)
// }
//
// func ReturnWithMessage(code int, message string, c *gin.Context) {
// 	Result(code, map[string]interface{}{}, message, c)
// }
