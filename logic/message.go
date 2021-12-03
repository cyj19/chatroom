/**
 * @Author: cyj19
 * @Date: 2021/12/3 11:00
 */

package logic

type Message struct {
	User    *User  `json:"user"`    // 哪个用户发送的消息
	Content string `json:"content"` // 消息内容
}

func NewMessage(user *User, content string) *Message {
	return &Message{
		User:    user,
		Content: content,
	}
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Content: "欢迎进入聊天室",
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Content: user.Nickname + "进入聊天室",
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Content: user.Nickname + "离开聊天室",
	}
}

func NewErrorMessage(errStr string) *Message {
	return &Message{
		User:    sysUser,
		Content: errStr,
	}
}
