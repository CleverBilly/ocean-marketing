package middleware

import (
	"fmt"
	"reflect"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"ocean-marketing/pkg/errno"
	"ocean-marketing/pkg/response"
)

// Validation 创建验证中间件
func Validation(reqStruct interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的新实例
		reqType := reflect.TypeOf(reqStruct)
		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
		reqValue := reflect.New(reqType)
		reqInterface := reqValue.Interface()

		// 绑定请求数据
		if err := c.ShouldBindJSON(reqInterface); err != nil {
			zap.L().Error("参数绑定失败", zap.Error(err))
			response.Error(c, errno.ErrBind)
			c.Abort()
			return
		}

		// 执行govalidator验证
		valid, err := govalidator.ValidateStruct(reqInterface)
		if err != nil {
			zap.L().Error("参数验证失败", zap.Error(err))
			response.Error(c, errno.New(errno.ErrValidation.Code, err.Error()))
			c.Abort()
			return
		}

		if !valid {
			zap.L().Error("参数验证失败: 结构体验证不通过")
			response.Error(c, errno.ErrValidation)
			c.Abort()
			return
		}

		// 执行自定义验证方法（如果存在）
		if validator, ok := reqInterface.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				zap.L().Error("自定义验证失败", zap.Error(err))
				response.Error(c, errno.New(errno.ErrValidation.Code, err.Error()))
				c.Abort()
				return
			}
		}

		// 将验证后的数据存储到上下文中
		c.Set("validatedData", reqInterface)
		c.Next()
	}
}

// GetValidatedData 从上下文中获取已验证的数据
func GetValidatedData(c *gin.Context, dst interface{}) error {
	data, exists := c.Get("validatedData")
	if !exists {
		return fmt.Errorf("no validated data found in context")
	}

	// 使用反射复制数据
	srcValue := reflect.ValueOf(data)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr || dstValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst must be a pointer to struct")
	}
	dstValue = dstValue.Elem()

	if srcValue.Type() != dstValue.Type() {
		return fmt.Errorf("type mismatch: src=%s, dst=%s", srcValue.Type(), dstValue.Type())
	}

	dstValue.Set(srcValue)
	return nil
}

// ValidationQuery 查询参数验证中间件
func ValidationQuery(reqStruct interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的新实例
		reqType := reflect.TypeOf(reqStruct)
		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
		reqValue := reflect.New(reqType)
		reqInterface := reqValue.Interface()

		// 绑定查询参数
		if err := c.ShouldBindQuery(reqInterface); err != nil {
			zap.L().Error("查询参数绑定失败", zap.Error(err))
			response.Error(c, errno.ErrBind)
			c.Abort()
			return
		}

		// 执行govalidator验证
		valid, err := govalidator.ValidateStruct(reqInterface)
		if err != nil {
			zap.L().Error("查询参数验证失败", zap.Error(err))
			response.Error(c, errno.New(errno.ErrValidation.Code, err.Error()))
			c.Abort()
			return
		}

		if !valid {
			zap.L().Error("查询参数验证失败: 结构体验证不通过")
			response.Error(c, errno.ErrValidation)
			c.Abort()
			return
		}

		// 执行自定义验证方法（如果存在）
		if validator, ok := reqInterface.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				zap.L().Error("查询参数自定义验证失败", zap.Error(err))
				response.Error(c, errno.New(errno.ErrValidation.Code, err.Error()))
				c.Abort()
				return
			}
		}

		// 将验证后的数据存储到上下文中
		c.Set("validatedData", reqInterface)
		c.Next()
	}
}

// ValidationForm 表单验证中间件
func ValidationForm(reqStruct interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建请求结构体的新实例
		reqType := reflect.TypeOf(reqStruct)
		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
		reqValue := reflect.New(reqType)
		reqInterface := reqValue.Interface()

		// 绑定表单数据
		if err := c.ShouldBind(reqInterface); err != nil {
			zap.L().Error("表单参数绑定失败", zap.Error(err))
			response.Error(c, errno.ErrBind)
			c.Abort()
			return
		}

		// 执行govalidator验证
		valid, err := govalidator.ValidateStruct(reqInterface)
		if err != nil {
			zap.L().Error("表单参数验证失败", zap.Error(err))
			response.Error(c, errno.New(errno.ErrValidation.Code, err.Error()))
			c.Abort()
			return
		}

		if !valid {
			zap.L().Error("表单参数验证失败: 结构体验证不通过")
			response.Error(c, errno.ErrValidation)
			c.Abort()
			return
		}

		// 执行自定义验证方法（如果存在）
		if validator, ok := reqInterface.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				zap.L().Error("表单参数自定义验证失败", zap.Error(err))
				response.Error(c, errno.New(errno.ErrValidation.Code, err.Error()))
				c.Abort()
				return
			}
		}

		// 将验证后的数据存储到上下文中
		c.Set("validatedData", reqInterface)
		c.Next()
	}
}
