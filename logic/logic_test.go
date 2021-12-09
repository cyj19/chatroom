/**
 * @Author: cyj19
 * @Date: 2021/12/8 16:50
 */

package logic

import (
	"github.com/cyj19/chatroom/global"
	"log"
	"testing"
	"time"
)

func TestProcessSensitiveData(t *testing.T) {
	dst := &global.UserResponse{}
	src := &User{
		UID:      1,
		Nickname: "user1",
		EnterAt:  time.Now(),
		Addr:     "127.0.0.1:7777",
		Token:    "zsdfvffdgrefgdfxghfdgh",
	}
	log.Println(global.ProcessSensitiveData(dst, src))
}
