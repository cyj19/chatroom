/**
 * @Author: cyj19
 * @Date: 2021/12/3 10:58
 */

package logic

import "github.com/cyj19/chatroom/global"

// 广播器
type broadcaster struct {
	users                 map[string]*User // 在线用户
	enteringChannel       chan *User
	leavingChannel        chan *User
	messageChannel        chan *Message // 广播消息通道
	checkUserChannel      chan string   // 判断用户昵称是否存在
	checkUserCanInChannel chan bool     // 判断用户是否可以加入聊天室
	requestUserChannel    chan struct{}
	usersChannel          chan []*User
}

// Broadcaster 广播器全局只能有一个，应该使用单例，饿汉式
var Broadcaster = &broadcaster{
	users:                 make(map[string]*User),
	enteringChannel:       make(chan *User), // 要不要缓冲
	leavingChannel:        make(chan *User),
	messageChannel:        make(chan *Message, global.MessageQueueLen),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
	requestUserChannel:    make(chan struct{}),
	usersChannel:          make(chan []*User),
}

// Start 开启广播
func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringChannel:
			// 添加用户
			b.users[user.Nickname] = user
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// 删除用户
			delete(b.users, user.Nickname)
			// 避免goroutine泄露
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:

			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			// 判断用户是否已存在
			_, ok := b.users[nickname]
			b.checkUserCanInChannel <- !ok
		case <-b.requestUserChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}
			b.usersChannel <- userList
		}
	}
}

// UserEntering 新用户进入
func (b *broadcaster) UserEntering(user *User) {
	b.enteringChannel <- user
}

// UserLeaving 用户离开
func (b *broadcaster) UserLeaving(user *User) {
	b.leavingChannel <- user
}

// Broadcast 广播消息
func (b *broadcaster) Broadcast(msg *Message) {
	b.messageChannel <- msg
}

// CanEnterRoom 判断用户是否可加入聊天室
func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

// GetUserList 获取在线用户列表
func (b *broadcaster) GetUserList() []*User {
	b.requestUserChannel <- struct{}{}
	return <-b.usersChannel
}
