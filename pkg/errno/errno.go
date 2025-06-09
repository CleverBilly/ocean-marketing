package errno

import "fmt"

// Errno 定义错误码结构
type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error 实现error接口
func (err Errno) Error() string {
	return err.Message
}

// String 返回错误信息
func (err Errno) String() string {
	return fmt.Sprintf("Err - code: %d, message: %s", err.Code, err.Message)
}

// 预定义的错误码
var (
	// 系统错误
	OK                  = Errno{Code: 0, Message: "OK"}
	InternalServerError = Errno{Code: 10001, Message: "内部服务器错误"}
	ErrBind             = Errno{Code: 10002, Message: "请求参数错误"}
	ErrValidation       = Errno{Code: 10003, Message: "参数校验失败"}
	ErrDatabase         = Errno{Code: 10004, Message: "数据库错误"}
	ErrRedis            = Errno{Code: 10005, Message: "Redis错误"}
	ErrEncrypt          = Errno{Code: 10006, Message: "加密错误"}
	ErrLimitExceed      = Errno{Code: 10007, Message: "请求频率超限"}

	// 认证授权错误
	ErrTokenInvalid     = Errno{Code: 20001, Message: "Token无效"}
	ErrTokenExpired     = Errno{Code: 20002, Message: "Token已过期"}
	ErrTokenNotFound    = Errno{Code: 20003, Message: "Token不存在"}
	ErrPermissionDenied = Errno{Code: 20004, Message: "权限不足"}
	ErrUnauthorized     = Errno{Code: 20005, Message: "未授权"}

	// 用户相关错误
	ErrUserNotFound      = Errno{Code: 30001, Message: "用户不存在"}
	ErrUserAlreadyExist  = Errno{Code: 30002, Message: "用户已存在"}
	ErrPasswordIncorrect = Errno{Code: 30003, Message: "密码错误"}
	ErrUserDisabled      = Errno{Code: 30004, Message: "用户已禁用"}

	// 资源相关错误
	ErrResourceNotFound     = Errno{Code: 40001, Message: "资源不存在"}
	ErrResourceAlreadyExist = Errno{Code: 40002, Message: "资源已存在"}
	ErrResourceConflict     = Errno{Code: 40003, Message: "资源冲突"}

	// 业务逻辑错误
	ErrBusiness = Errno{Code: 50001, Message: "业务逻辑错误"}
)

// New 创建新的错误码
func New(code int, message string) Errno {
	return Errno{
		Code:    code,
		Message: message,
	}
}

// IsErrno 判断是否为Errno类型错误
func IsErrno(err error) bool {
	_, ok := err.(Errno)
	return ok
}

// DecodeErr 解析错误，返回错误码和错误信息
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case Errno:
		return typed.Code, typed.Message
	default:
		return InternalServerError.Code, err.Error()
	}
}
