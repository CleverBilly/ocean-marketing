package response

import (
	"net/http"

	"ocean-marketing/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	code, message := errno.DecodeErr(err)
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithCode 带状态码的错误响应
func ErrorWithCode(c *gin.Context, httpCode int, err error) {
	code, message := errno.DecodeErr(err)
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
	})
}

// Custom 自定义响应
func Custom(c *gin.Context, httpCode int, code int, message string, data interface{}) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, err error) {
	ErrorWithCode(c, http.StatusBadRequest, err)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, err error) {
	ErrorWithCode(c, http.StatusUnauthorized, err)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, err error) {
	ErrorWithCode(c, http.StatusForbidden, err)
}

// NotFound 404错误
func NotFound(c *gin.Context, err error) {
	ErrorWithCode(c, http.StatusNotFound, err)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, err error) {
	ErrorWithCode(c, http.StatusInternalServerError, err)
}

// PageResponse 分页响应结构
type PageResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
		Data: PageResponse{
			List:  list,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}
