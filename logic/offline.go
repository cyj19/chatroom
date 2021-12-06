/**
 * @Author: cyj19
 * @Date: 2021/12/6 14:48
 */

// 离线消息处理
// 问题：所有用户最近消息和某用户消息有可能重复

package logic

import (
	"container/ring"
	"github.com/spf13/viper"
)

type offlineProcessor struct {
	n          int                   // 消息数量
	recentRing *ring.Ring            // 保存所有用户最近的n条消息
	userRing   map[string]*ring.Ring // 保存某用户离线消息
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num")
	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),
		userRing:   make(map[string]*ring.Ring),
	}
}

// Save 保存离线消息
func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}
	o.recentRing.Value = msg
	// 移到下一个
	o.recentRing = o.recentRing.Next()

	for _, nickname := range msg.Ats {
		nickname = nickname[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		// 如果没有该用户的离线消息列表，新建一个
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}
		r.Value = msg
		// 移到下一个节点
		o.userRing[nickname] = r.Next()
	}
}

// Send 取出存储的离线消息
func (o *offlineProcessor) Send(user *User) {
	// 取所有用户的最近消息
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChannel <- value.(*Message)
		}
	})

	// 新用户直接返回
	if user.isNew {
		return
	}

	// 取该用户特有的离线消息
	if r, ok := o.userRing[user.Nickname]; ok {
		r.Do(func(value interface{}) {
			if value != nil {
				user.MessageChannel <- value.(*Message)
			}
		})
		delete(o.userRing, user.Nickname)
	}

}
