package cast

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/cast"
)

// ToString 转换为字符串
func ToString(i interface{}) string {
	return cast.ToString(i)
}

// ToInt 转换为int
func ToInt(i interface{}) int {
	return cast.ToInt(i)
}

// ToInt64 转换为int64
func ToInt64(i interface{}) int64 {
	return cast.ToInt64(i)
}

// ToFloat64 转换为float64
func ToFloat64(i interface{}) float64 {
	return cast.ToFloat64(i)
}

// ToBool 转换为bool
func ToBool(i interface{}) bool {
	return cast.ToBool(i)
}

// ToTime 转换为time.Time
func ToTime(i interface{}) time.Time {
	return cast.ToTime(i)
}

// ToStringSlice 转换为[]string
func ToStringSlice(i interface{}) []string {
	return cast.ToStringSlice(i)
}

// ToIntSlice 转换为[]int
func ToIntSlice(i interface{}) []int {
	return cast.ToIntSlice(i)
}

// ToStringMap 转换为map[string]interface{}
func ToStringMap(i interface{}) map[string]interface{} {
	return cast.ToStringMap(i)
}

// ToStringMapString 转换为map[string]string
func ToStringMapString(i interface{}) map[string]string {
	return cast.ToStringMapString(i)
}

// ToJSON 转换为JSON字符串
func ToJSON(i interface{}) (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串转换
func FromJSON(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), v)
}

// ToUint 转换为uint
func ToUint(i interface{}) uint {
	return cast.ToUint(i)
}

// ToUint64 转换为uint64
func ToUint64(i interface{}) uint64 {
	return cast.ToUint64(i)
}

// ToFloat32 转换为float32
func ToFloat32(i interface{}) float32 {
	return cast.ToFloat32(i)
}

// SafeToString 安全转换为字符串，带默认值
func SafeToString(i interface{}, defaultValue string) string {
	if i == nil {
		return defaultValue
	}
	result := cast.ToString(i)
	if result == "" {
		return defaultValue
	}
	return result
}

// SafeToInt 安全转换为int，带默认值
func SafeToInt(i interface{}, defaultValue int) int {
	if i == nil {
		return defaultValue
	}
	return cast.ToInt(i)
}

// SafeToInt64 安全转换为int64，带默认值
func SafeToInt64(i interface{}, defaultValue int64) int64 {
	if i == nil {
		return defaultValue
	}
	return cast.ToInt64(i)
}

// SafeToFloat64 安全转换为float64，带默认值
func SafeToFloat64(i interface{}, defaultValue float64) float64 {
	if i == nil {
		return defaultValue
	}
	return cast.ToFloat64(i)
}

// SafeToBool 安全转换为bool，带默认值
func SafeToBool(i interface{}, defaultValue bool) bool {
	if i == nil {
		return defaultValue
	}
	return cast.ToBool(i)
}

// IsEmpty 检查值是否为空
func IsEmpty(i interface{}) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	default:
		return false
	}
}

// ToPointer 转换为指针
func ToPointer(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		return i
	}
	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	return ptr.Interface()
}

// FromPointer 从指针获取值
func FromPointer(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return i
	}
	if v.IsNil() {
		return nil
	}
	return v.Elem().Interface()
}

// StringToInt 字符串转整数，带错误处理
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StringToInt64 字符串转int64，带错误处理
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// StringToFloat64 字符串转float64，带错误处理
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// StringToBool 字符串转bool，带错误处理
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
