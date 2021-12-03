/**
 * @Author: cyj19
 * @Date: 2021/12/2 16:43
 */

package server

import (
	"fmt"
	"github.com/cyj19/chatroom/global"
	"html/template"
	"net/http"
)

// 主页路由处理
func homeHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Fprintf(w, "模板解析错误")
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		_, _ = fmt.Fprintf(w, "模板执行错误")
		return
	}
}
