package dhttp

import (
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

// 是否已开启socket信号
var sw = false

// 处理Post请求
func HandGetRequest(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "无法读取请求数据", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//切割字符串
	parts := strings.Split(string(body), "#")
	var mac string
	var cmd string
	if len(parts) == 2 {
		mac = parts[0] // 第一个部分 mac地址
		cmd = parts[1] // 第二个部分 指令 开启/停止
	}

	res := mac
	macAddr, err := utils.GetMACAddress()

	//启动服务
	if res == macAddr && cmd == "RunSocketServer" {
		if sw == false {
			utils.SelfMac = res
			//说明是服务端后台管理UI发送的信号,这里不需要下载
			SChan <- true
			sw = true
			w.WriteHeader(http.StatusOK)
		}
		return
	}

	//停止服务
	if res == macAddr && cmd == "StopSocketServer" {
		if sw == true {
			os.Exit(0)
		}
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	return
	//filePath := "down_load/1.txt" // 相对路径
	//fmt.Println(filePath)
	//file, err := os.Open(filePath)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("未找到指定文件: %s", err), http.StatusInternalServerError)
	//	return
	//}
	//defer file.Close()
	//
	//// 设置响应头，指定文件名
	//w.Header().Set("Content-Disposition", "attachment; filename=1.txt")
	//
	//// 将文件内容复制到响应主体
	//_, err = io.Copy(w, file)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("无法复制文件内容: %s", err), http.StatusInternalServerError)
	//	return
	//}
}
