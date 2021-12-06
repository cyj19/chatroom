/**
 * @Author: cyj19
 * @Date: 2021/12/2 16:05
 */

// 正式版websocket聊天室启动程序
package main

import (
	"fmt"
	"github.com/cyj19/chatroom/global"
	"github.com/cyj19/chatroom/server"
	"log"
	"net/http"
)

var (
	addr   = ":8888"
	banner = `
    ____              _____
   |     |    |   /\     |
   |     |____|  /  \    | 
   |     |    | /----\   |
   |____ |    |/      \  |

    ChatRoom，start on：%s
`
)

func init() {
	global.Init()
}

func main() {
	fmt.Printf(banner+"\n", addr)
	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
