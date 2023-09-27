package conf

import (
	"bufio"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wownet"
	"os"
	"strconv"
	"strings"
)

// 读取ban.txt文件内容
func ReadBanList() {
	fileName := "ban.txt"
	// 尝试打开文件
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，创建一个新文件
			file, createErr := os.Create(fileName)
			if createErr != nil {
				fmt.Println("无法创建文件:", createErr)
				utils.Logger.Error("无法创建ban.txt文件,如已打开该文件,请关闭后再试")
				panic(createErr)
				return
			}
			defer file.Close()
			return
		} else {
			utils.Logger.Error("无法打开ban.txt文件,如目录下没有该文件,请手动创建")
			panic(err)
			return
		}
	} else {
		defer file.Close()

		// 创建一个 Scanner 来逐行读取文件内容
		scanner := bufio.NewScanner(file)

		// 遍历每一行
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text()) // 去除两端的空格和不可见字符
			fields := strings.Split(line, "#")
			if len(fields) != 3 {
				utils.Logger.Error("读取ban.txt文件,无效行:", line)
				continue
			}

			ip := fields[0]
			mac := fields[1]
			unixTimeStr := fields[2]

			// 解析 Unix 时间戳
			unixTime, err := strconv.ParseInt(unixTimeStr, 10, 64)
			if err != nil {
				fmt.Println("Ban文件:无效的 Unix 时间戳:", unixTimeStr)
				continue
			}
			// 读取到内容中
			wownet.BanInstance.LoadBanData(ip, mac, unixTime)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("读取文件时发生错误:", err)
			panic(err)
		}
	}
}
