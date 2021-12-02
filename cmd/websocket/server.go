/**
 * @Author: cyj19
 * @Date: 2021/12/2 10:39
 */

// http和websocket 服务端监听同一端口
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Http hello")
	})

	// 还是一个http请求
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			_ = conn.Close(websocket.StatusInternalError, "内部出错")
		}()

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		// 读取请求数据
		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("接收到客户端数据：%v \n", v)

		// 回复客户端
		err = wsjson.Write(ctx, conn, "Hello Websocket Client")
		if err != nil {
			log.Println(err)
			return
		}
		// 正常关闭连接
		_ = conn.Close(websocket.StatusNormalClosure, "")

	})

	log.Fatal(http.ListenAndServe(":8888", nil))
}
