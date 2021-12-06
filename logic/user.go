/**
 * @Author: cyj19
 * @Date: 2021/12/3 10:59
 */

package logic

import (
	"context"
	"errors"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sync/atomic"
	"time"
)

// 用户需要做什么？1、接收消息 2、发送消息

var globalID uint32

type User struct {
	UID            int             `json:"uid"`      // 用户id
	Nickname       string          `json:"nickname"` // 用户昵称
	EnterAt        time.Time       `json:"enterAt"`  // 加入时间
	Addr           string          `json:"addr"`     // IP地址
	MessageChannel chan *Message   `json:"-"`        // 消息通道  ps:chan无法转为json输出
	conn           *websocket.Conn // websocket连接
}

// 系统用户，用于发送系统消息
var sysUser = &User{}

// NewUser 工厂方法，创建用户实例
func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	user := &User{
		UID:            int(atomic.AddUint32(&globalID, 1)),
		Nickname:       nickname,
		EnterAt:        time.Now(),
		Addr:           addr,
		MessageChannel: make(chan *Message, 32),
		conn:           conn,
	}
	return user
}

// CloseMessageChannel 避免goroutine泄露
func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

// SendMessage 给客户端发送消息
func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		err := wsjson.Write(ctx, u.conn, msg)
		if err != nil {
			log.Println("server write error: ", err)
		}
	}
}

// ReceiveMessage 接收来自客户端的消息
func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			// 判断是否正常关闭连接
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}

			return err
		}

		// 发送消息到聊天室
		sendMsg := NewMessage(u, receiveMsg["content"])
		Broadcaster.Broadcast(sendMsg)
	}
}
