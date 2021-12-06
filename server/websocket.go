/**
 * @Author: cyj19
 * @Date: 2021/12/2 16:49
 */

package server

import (
	"github.com/cyj19/chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// WebSocketHandleFunc websocket路由处理
func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	// 如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了 InsecureSkipVerify 选项（通过第三个参数 AcceptOptions 设置）。
	// 换句话说，默认情况下，它不允许跨源请求。如果发生错误，Accept 将始终写入适当的响应
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println(err)
	}

	// 1、构建新用户
	nickname := req.FormValue("nickname")
	token := req.FormValue("token")
	// 校验昵称长度 2 ~ 20
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal:", nickname)
		_ = wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：2-20"))
		_ = conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}
	// 判断能否加入
	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("该昵称已存在")
		_ = wsjson.Write(req.Context(), conn, logic.NewErrorMessage("该昵称已存在"))
		_ = conn.Close(websocket.StatusUnsupportedData, "nickname exists")
		return
	}
	user := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	// 2、开启发送消息的goroutine
	go user.SendMessage(req.Context())

	// 3、给当前用户发送欢迎消息
	user.MessageChannel <- logic.NewWelcomeMessage(user)

	// 4、通知其他用户有新用户进入
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)

	// 5、将当前用户加入到广播器的用户列表中
	logic.Broadcaster.UserEntering(user)
	log.Println("user: ", nickname, "join chat")

	// 6、开启接收来自客户端的消息
	err = user.ReceiveMessage(req.Context())

	// 7、用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user: ", nickname, "leave chat")

	// 根据读取时的错误，选择不同的close
	if err == nil {
		_ = conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client err: ", err)
		_ = conn.Close(websocket.StatusInternalError, "Read from client error")
	}

}
