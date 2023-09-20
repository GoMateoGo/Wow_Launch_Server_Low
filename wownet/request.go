package wownet

import "gitee.com/mrmateoliu/wow_launch.git/wowiface"

type request struct {
	//已经和客户端建立好的链接 conn
	conn wowiface.IConnection
	//客户端请求的数据
	data []byte
}

// 得到当前链接
func (r *request) GetConnection() wowiface.IConnection {

	return r.conn
}

// 得到请求的消息数据
func (r *request) GetData() []byte {

	return r.data
}
