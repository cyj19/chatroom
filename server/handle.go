/**
 * @Author: cyj19
 * @Date: 2021/12/2 16:08
 */

// Package server 服务端处理
package server

import (
	"github.com/cyj19/chatroom/logic"
	"net/http"
)

func RegisterHandle() {

	// 广播消息处理
	go logic.Broadcaster.Start()

	// 主页路由
	http.HandleFunc("/", homeHandleFunc)
	// 用户列表路由
	http.HandleFunc("/user_list", userListHandleFunc)
	// websocket路由
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
