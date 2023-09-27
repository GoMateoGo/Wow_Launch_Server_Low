package wownet

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex // 声明一个互斥锁

var BanInstance *Ban

func init() {
	BanInstance = NewBan()
	go LoopCheckExpireBan()
}

// 封禁结构
type Ban struct {
	BanList map[string]*BannedList
}

type BannedList struct {
	MacAddr    string
	ExpireTime int64
}

// 创建指针
func NewBan() *Ban {
	return &Ban{
		BanList: make(map[string]*BannedList),
	}
}

// 增加封禁
func (b *Ban) AddBan(ip, mac string, expireTime int64) {
	if AddBanDataToTxt(ip, mac, strconv.FormatInt(expireTime, 10)) {
		b.BanList[ip] = &BannedList{
			MacAddr:    mac,
			ExpireTime: expireTime,
		}
	}
}

// 加载封禁
func (b *Ban) LoadBanData(ip, mac string, expireTime int64) {
	b.BanList[ip] = &BannedList{
		MacAddr:    mac,
		ExpireTime: expireTime,
	}
}

// 获取封禁用户信息,使用的时候需要断言成: make(map[string]*BannedList)
func (b *Ban) GetBanByInfo(ip, Mac string) interface{} {
	result := make(map[string]*BannedList)
	for k, v := range b.BanList {
		if k == ip || v.MacAddr == Mac {
			result[k] = v
			return result
		}
	}
	return nil
}

// 移除封禁
func (b *Ban) RemoveBan(ipOrMac string) bool {
	if ModifyBanDataInTxt(ipOrMac) {
		for k := range b.BanList {
			if k == ipOrMac || b.BanList[k].MacAddr == ipOrMac {
				delete(b.BanList, k)
				return true
			}
		}
	}
	return false
}

// 轮询删除过期
func LoopCheckExpireBan() {
	for {
		if len(BanInstance.BanList) > 0 {
			for k, v := range BanInstance.BanList {
				if v.ExpireTime < time.Now().Unix() {
					BanInstance.RemoveBan(k)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// 从ban.txt中删除被ban信息
func ModifyBanDataInTxt(mac string) bool {
	// 打开文件以供读写
	fileName := "ban.txt"
	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("无法打开文件11:", err)
		return false
	}
	defer file.Close()

	// 创建一个 Scanner 来逐行读取原始文件内容
	scanner := bufio.NewScanner(file)

	// 创建一个缓冲器，用于保存修改后的内容
	var modifiedContent []string

	// 遍历每一行
	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否包含目标 MAC 地址
		if strings.Contains(line, mac) {
			continue // 跳过包含目标 MAC 地址的行
		}

		// 将非匹配行添加到缓冲器
		modifiedContent = append(modifiedContent, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时发生错误:", err)
		return false
	}

	// 关闭原始文件
	file.Close()

	// 打开文件以供写入（清空原有内容）
	file, err = os.Create(fileName)
	if err != nil {
		fmt.Println("无法创建文件:", err)
		return false
	}
	defer file.Close()

	// 将修改后的内容写入文件
	for _, line := range modifiedContent {
		_, err := file.WriteString(line + "\n")
		fmt.Println("写入一行:", line)
		if err != nil {
			fmt.Println("无法写入文件:", err)
			return false
		}
	}

	return true
}

// 新增ban信息到txt文件中
func AddBanDataToTxt(ip, mac, expireTime string) bool {
	fileName := "ban.txt"
	lineToAdd := ip + "#" + mac + "#" + expireTime // 你要添加的行

	// 加锁，确保只有一个 goroutine 可以访问文件
	mu.Lock()
	defer mu.Unlock()

	// 打开文件以供追加
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("无法打开文件22:", err)
		return false
	}
	defer file.Close()

	// 写入行到文件末尾
	_, err = file.WriteString(lineToAdd + "\n")
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return false
	}
	return true
}
