package dhttp

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"net/http"
	"strings"
	"sync"
)

var (
	SChan chan bool
)

// 去除端口号部分
func removePortFromAddress(addr string) string {
	// 使用冒号分割地址，取第一个部分作为IP地址
	parts := strings.Split(addr, ":")
	return parts[0]
}

func RunHttp() {
	SChan = make(chan bool)

	// 定义处理HTTP请求的处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var wg sync.WaitGroup
			wg.Add(1)
			go HandGetRequest(w, r, &wg)
			wg.Wait() // 等待Goroutine完成
		}
	})
	// 启动HTTP服务器并监听端口
	port := utils.GlobalObject.HttpPort
	fmt.Printf("[补丁下载服务器信息]:\nIp:%s\n端口: %d\n", utils.GlobalObject.Host, port)
	fmt.Println("----------------------")
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", utils.GlobalObject.Host, port), nil)
		if err != nil {
			fmt.Println("http服务器启动失败..:", err)
		}
	}()
}
