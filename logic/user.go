/**
 * @Author: cyj19
 * @Date: 2021/12/3 10:59
 */

package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cyj19/chatroom/global"
	"github.com/spf13/cast"
	"io"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

// 用户需要做什么？1、接收消息 2、发送消息

var globalID uint32

type User struct {
	UID            int             `json:"uid"`      // 用户id
	Nickname       string          `json:"nickname"` // 用户昵称
	EnterAt        time.Time       `json:"enter_at"` // 加入时间
	Addr           string          `json:"addr"`     // IP地址
	MessageChannel chan *Message   `json:"-"`        // 消息通道  ps:chan无法转为json输出
	Token          string          `json:"token"`
	conn           *websocket.Conn // websocket连接
	isNew          bool            // 是否是新用户
}

// 系统用户，用于发送系统消息
var sysUser = &User{}

// NewUser 工厂方法，创建用户实例
func NewUser(conn *websocket.Conn, token, nickname, addr string) *User {
	user := &User{
		Nickname:       nickname,
		EnterAt:        time.Now(),
		Addr:           addr,
		MessageChannel: make(chan *Message, 32),
		Token:          token,
		conn:           conn,
	}

	if user.Token != "" {
		// 解析token，获取uid
		uid, err := parseTokenAndValidate(user.Token, user.Nickname)
		if err == nil {
			user.UID = uid
		}
	} else {
		// 新用户
		user.UID = int(atomic.AddUint32(&globalID, 1))
		user.Token = genToken(user.UID, user.Nickname)
		user.isNew = true
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
			var closeErr websocket.CloseError
			// 判断是否正常关闭连接
			if errors.As(err, &closeErr) || errors.Is(err, io.EOF) {
				return nil
			}

			return err
		}

		// 发送消息到聊天室
		sendMsg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])
		sendMsg.Content = FilterSensitiveWords(sendMsg.Content)

		// 处理@
		// [^\s@]匹配不包含@的非空字符
		reg := regexp.MustCompile(`@[^\s@]{2,20}`)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
		Broadcaster.Broadcast(sendMsg)
	}
}

// 生成token
func genToken(uid int, nickname string) string {
	message := fmt.Sprintf("%s%s%d", nickname, global.TokenSecret, uid)

	messageMAC := macSha256([]byte(message), []byte(global.TokenSecret))
	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

// 解析token和校验
func parseTokenAndValidate(token, nickname string) (int, error) {
	// 偏移量
	pos := strings.LastIndex(token, "uid")
	// base64解码
	messageMAC, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	uid := cast.ToInt(token[pos+3:])
	message := fmt.Sprintf("%s%s%d", nickname, global.TokenSecret, uid)
	ok := validateMAC([]byte(message), messageMAC, []byte(global.TokenSecret))
	if ok {
		return uid, nil
	}

	return 0, errors.New("token is illegal")
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func validateMAC(message, messageMAC, secret []byte) bool {
	expectedMAC := macSha256(message, secret)
	return hmac.Equal(messageMAC, expectedMAC)
}
