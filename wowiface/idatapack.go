package wowiface

/*

----------TLV(type,length,value) 格式的包处理模块----------
1. 拆包
	1)先读取head固定长度, 在读消息内容长度和消息类型
	2)在根据消息内容的长度,再次从conn中的进行读取内容
2. 封装
	1)写msg的长度
	2)写msg的id
	3)写msg的内容
*/

type IDataPack interface {
	//获取包的头长度的方法
	GetHeadLen() uint32

	//封包方法
	Pack(msg IMessage) ([]byte, error)

	//拆包方法
	UnPack([]byte) (IMessage, error)
}
