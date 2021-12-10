/**
 * @Author: cyj19
 * @Date: 2021/12/8 16:37
 */

// 数据脱敏处理

package global

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type UserResponse struct {
	UID      int       `json:"uid"`
	Nickname string    `json:"nickname"`
	EnterAt  time.Time `json:"enter_at"`
	Addr     string    `json:"addr"`
}

// ProcessSensitiveData 数据脱敏 结构体转结构体
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

func Struct2Json(src interface{}) string {
	data, err := json.Marshal(src)
	if err != nil {
		fmt.Printf("Struct2Json error:%v \n", err)
	}
	return string(data)
}

// Json2Struct 参数dst必须是指针
func Json2Struct(src string, dst interface{}) error {
	err := json.Unmarshal([]byte(src), dst)
	if err != nil {
		fmt.Printf("Json2Struct error:%v \n", err)
	}
	return err
}

// Struct2Struct struct转struct 一般通过json tag进行比较转换，因此可以通过struct，json互转实现
func Struct2Struct(src, dst interface{}) {
	str := Struct2Json(src)
	_ = Json2Struct(str, dst)
}
