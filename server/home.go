/**
 * @Author: cyj19
 * @Date: 2021/12/2 16:43
 */

package server

import (
	"encoding/json"
	"fmt"
	"github.com/cyj19/chatroom/global"
	"github.com/cyj19/chatroom/logic"
	"html/template"
	"log"
	"net/http"
)

// 主页路由处理
func homeHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		_, _ = fmt.Fprintf(w, "模板解析错误")
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		_, _ = fmt.Fprintf(w, "模板执行错误")
		return
	}
}

func userListHandleFunc(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	users := logic.Broadcaster.GetUserList()
	b, err := json.Marshal(users)

	if err != nil {
		log.Println("json.Marshal error:", err)
		_, _ = fmt.Fprint(w, `[]`)
	} else {
		_, _ = fmt.Fprint(w, string(b))
	}
}
