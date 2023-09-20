package wownet

import "gitee.com/mrmateoliu/wow_launch.git/wowiface"

type Request struct {
	//已经和客户端建立好的链接 conn
	conn wowiface.IConnection
	//客户端请求的数据
	msg wowiface.IMessage
}

// 得到当前链接
func (r *Request) GetConnection() wowiface.IConnection {

	return r.conn
}

// 得到请求的消息数据
func (r *Request) GetData() []byte {

	return r.msg.GetData()
}

// 得到消息Id
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
