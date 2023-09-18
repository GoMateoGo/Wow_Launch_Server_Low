package main

import "gitee.com/mrmateoliu/wow_launch.git/wownet"

func main() {
	// 1. 创建server句柄
	s := wownet.NewServer("wow-launch")
	// 2. 启动服务器
	s.Server()
}
