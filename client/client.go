package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端

func main() {
	fmt.Println("客户端启动...")
	time.Sleep(1 * time.Second)
	// 1. 链接远程服务器,得到conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("客户端链接失败...", err)
		return
	}

	for {
		// 2.链接调用write写数据

		_, err := conn.Write([]byte("你好.."))
		if err != nil {
			fmt.Println("客户端写入数据错误...", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("客户端读取服务端返回数据错误:", err)
			return
		}
		fmt.Printf("服务器返回数据:%s,长度:%d\n", buf, cnt)

		//cpu阻塞一下.
		time.Sleep(1 * time.Second)
	}
}
