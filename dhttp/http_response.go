package dhttp

import (
	"archive/zip"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// 是否已开启socket信号
var sw = false

// 处理Post请求
func HandPostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "请求体为空", http.StatusBadRequest)
		return
	}

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

	//停止服务
	if res == macAddr && cmd == "StopSocketServer" {
		os.Exit(0)
	}

	http.Error(w, "无效的请求", http.StatusBadRequest)
}

func HandMd5CheckRequest(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "body为空", http.StatusInternalServerError)
		return
	}
	fmt.Println("post请求uri", strings.Contains(r.URL.Path, "/check"))
	if !strings.Contains(r.URL.Path, "/check") {
		http.Error(w, "请求路径错误", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "读取body失败", http.StatusInternalServerError)
	}
	defer r.Body.Close()

	fmt.Println("md5:", string(body))
	res := utils.CheckMD5(string(body))
	if res {
		return
	} else {
		http.Error(w, "md5不符合,进行更新", http.StatusBadRequest)
	}
}

var mutex sync.Mutex // 创建一个互斥锁
func HandDownLoadRequest(w http.ResponseWriter, r *http.Request) {

	rootDir := "down_load" // 设置下载根目录
	zipFileName := "downloaded.zip"

	// 获取临时ZIP文件的大小
	if fileInfo, err := os.Stat(zipFileName); err == nil {
		// 设置HTTP头部，告诉浏览器该文件是一个ZIP文件
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFileName))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

		// 发送ZIP文件给客户端
		http.ServeFile(w, r, zipFileName)

		return
	}

	mutex.Lock()         // 加锁
	defer mutex.Unlock() // 解锁
	// 创建一个临时ZIP文件
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zipWriter := zip.NewWriter(zipFile)

	// 递归遍历目录并将文件和文件夹打包
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对于根目录的路径
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// 使用斜杠分隔符确保路径在不同操作系统上都有效
		relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")

		if info.IsDir() {
			// 不需要在ZIP文件中创建文件夹，只需在路径后面添加斜杠
			relPath += "/"
		} else {
			// 创建ZIP文件中的文件
			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			// 打开文件并将其内容写入ZIP文件
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 刷新临时ZIP文件，确保所有数据都写入磁盘
	err = zipWriter.Flush()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 关闭临时ZIP文件，确保文件写入磁盘
	err = zipWriter.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取临时ZIP文件的大小
	fileInfo, err := os.Stat(zipFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 关闭临时ZIP文件句柄
	err = zipFile.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置HTTP头部，告诉浏览器该文件是一个ZIP文件
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", zipFileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 发送ZIP文件给客户端
	http.ServeFile(w, r, zipFileName)
}
