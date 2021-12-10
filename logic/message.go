/**
 * @Author: cyj19
 * @Date: 2021/12/3 11:00
 */

package logic

import (
	"github.com/spf13/cast"
	"time"
)

const (
	MsgTypeNormal      = iota // 普通 用户消息
	MsgTypeUserWelcome        // 当前用户欢迎消息
	MsgTypeUserEnter          // 用户进入
	MsgTypeUserLeave          // 用户离开
	MsgTypeError              // 错误消息
)

type Message struct {
	User           *User     `json:"user"`             // 哪个用户发送的消息
	Type           int       `json:"type"`             // 消息类型
	Content        string    `json:"content"`          // 消息内容
	MsgTime        time.Time `json:"msg_time"`         // 消息产生时间
	ClientSendTime time.Time `json:"client_send_time"` // 客户端发送时间
	Ats            []string  `json:"ats"`              // 消息 @ 了谁
}

func NewMessage(user *User, content string, clientSendTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if clientSendTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientSendTime))
	}

	return message
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserWelcome,
		Content: user.Nickname + "你好，欢迎进入聊天室！",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	// token不传递给其他用户
	u := *user
	u.Token = ""
	return &Message{
		User:    &u,
		Type:    MsgTypeUserEnter,
		Content: user.Nickname + "加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.Nickname + "离开了聊天室",
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    sysUser,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
