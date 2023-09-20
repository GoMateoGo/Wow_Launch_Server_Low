package wownet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是测试负责测试datapack 拆包和封包的单元测试
func TestDataPack(t *testing.T) {
	//模拟的服务器
	// 1.创建socketTcp
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建一个go 承载负责客从客户端处理业务
	go func() {
		// 2.从客户端读取数据, 拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error:", err)
			}

			go func(conn net.Conn) {
				//处理客户端请求
				//------> 拆包过程 <------
				//定义一个拆包对象
				dp := NewDataPack()
				for {

					// 1.第一次从conn读,把head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err = io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println(" read head error:", err)
						break
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error:", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						//说明msg里边有数据长度的, 需要第二次读取(读具体data)
						// 2.第二次读,根据head中的datalen在读取data中内容

						//断言成具体类型
						msg := msgHead.(*Message)
						//根据msg的长度开辟空间
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度,再次从io流中读取具体data
						_, err = io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack error:", err)
							return
						}

						//完整的一个消息已经读取完毕
						fmt.Println("->-> 读取完的数据 Id:", msg.Id, "数据长度:", msg.DataLen, "数据:", string(msg.Data))
					}
				}

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("测试: 客户端链接失败...")
		return
	}

	//创建一个封包的过程

	dp := NewDataPack()

	//模拟黏包过程, 封装2个msg一同发送
	//封装第一个msg 1包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'w', 'o', 'w', 's', 'f'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("客户端打包1失败:", err)
	}
	//封装第二个msg 2包
	msg2 := &Message{
		Id:      2,
		DataLen: 8,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', 'w', 'o', 'w'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("客户端打包2失败:", err)
	}
	//将2个包放粘在一起
	sendData1 = append(sendData1, sendData2...)

	//一次性发送给服务端
	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("客户端发送数据包失败...", err)
		return
	}

	//客户端阻塞
	select {}
}
