package dhttp

import (
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// 调用
func HandCallDll() {
	dll, err := syscall.LoadDLL("server.dll")
	if err != nil {
		SendMessageBox("错误", "未找到server.dll")
		utils.Logger.Error(fmt.Sprintf("获取dll失败:%s", err))
		panic(fmt.Sprintf("获取dll失败:%s", err))
		return
	}

	proc, err := dll.FindProc("ExpireTime")
	if err != nil {
		SendMessageBox("错误", "server.dll加载错误")
		utils.Logger.Error(fmt.Sprintf("没有找到对应函数:%s", err))
		panic(fmt.Sprintf("没有找到对应函数:%s", err))
		return
	}

	go func() {
		for {
			call, _, err := proc.Call(GetExpireTime())
			if err != nil {
				//fmt.Println("call 的err:", err)
			}
			apiTime := GetApiUnixTime()
			res := CalTime(apiTime, int64(call))
			if res == false {
				utils.Logger.Error(fmt.Sprintf("登录网关已过期"))
				os.Exit(0)
			}
			remainTime := int64(call) - apiTime
			time.Sleep(1 * time.Second)
			utils.RemainTimeSecond = remainTime
			//fmt.Println("距离过期还有:", utils.RemainTimeSecond, "秒")
		}
	}()
}

// 发送win自带窗体
func SendMessageBox(title, msg string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	msgBox := user32.NewProc("MessageBoxW")
	msgBox.Call(IntPtr(0), strToPtr(msg), strToPtr(title), IntPtr(0))
}

// 计算时间
func CalTime(apiTime int64, ExpireTime int64) bool {
	if apiTime >= ExpireTime {
		return false
	}
	return true
}

// 计算过期时间
func GetExpireTime() uintptr {
	return 0
}

// 字符串转换为uint16指针
func strToPtr(s string) uintptr {
	b, _ := syscall.UTF16PtrFromString(s)
	return uintptr(unsafe.Pointer(b))
}

// int类型转uintptr
func IntPtr(i int) uintptr {
	return uintptr(i)
}
