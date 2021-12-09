/**
 * @Author: cyj19
 * @Date: 2021/12/8 16:37
 */

// 数据脱敏处理

package global

import (
	"reflect"
	"time"
)

type UserResponse struct {
	UID      int       `json:"uid"`
	Nickname string    `json:"nickname"`
	EnterAt  time.Time `json:"enter_at"`
	Addr     string    `json:"addr"`
}

// ProcessSensitiveData 数据脱敏
func ProcessSensitiveData(dst, src interface{}) interface{} {
	if dst == nil || src == nil {
		return nil
	}
	var dstValue reflect.Value
	var srcValue reflect.Value

	dstType := reflect.TypeOf(dst)
	srcType := reflect.TypeOf(src)

	if dstType.Kind() == reflect.Ptr && srcType.Kind() == reflect.Ptr {
		dstValue = reflect.ValueOf(dst).Elem()
		srcValue = reflect.ValueOf(src).Elem()
		dstType = dstType.Elem()
		srcType = srcType.Elem()
	} else if dstType.Kind() == reflect.Struct && srcType.Kind() == reflect.Struct {
		dstValue = reflect.ValueOf(dst)
		srcValue = reflect.ValueOf(src)
	} else {
		return nil
	}

	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)
		dstFieldValue := dstValue.Field(i)
		dstTag := dstField.Tag.Get("json")

		for j := 0; j < srcType.NumField(); j++ {
			srcField := srcType.Field(j)
			srcTag := srcField.Tag.Get("json")

			if dstTag == srcTag {
				srcFieldValue := srcValue.Field(j)
				dstFieldValue.Set(srcFieldValue)
				break
			}
		}
	}

	return dst
}
