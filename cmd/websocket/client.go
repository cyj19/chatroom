/**
 * @Author: cyj19
 * @Date: 2021/12/2 11:00
 */

package main

import (
	"context"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// 拨号
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8888/ws", nil)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close(websocket.StatusInternalError, "内部错误")
	}()
	// 给服务器发送数据
	err = wsjson.Write(ctx, conn, "Hello Websocket Server")
	if err != nil {
		// 客户端直接停止程序
		panic(err)
	}
	// 读取服务器发送的数据
	var v interface{}
	err = wsjson.Read(ctx, conn, &v)
	if err != nil {
		panic(err)
	}
	log.Printf("接收服务器数据：%v \n", v)
}
