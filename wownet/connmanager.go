package wownet

import (
	"errors"
	"fmt"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
	"sync"
)

/*
	连接管理模块
*/

type ConnManager struct {
	connections map[uint32]wowiface.IConnection //管理的链接信息集合
	connLock    sync.RWMutex                    //保护连接集合的读写锁
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]wowiface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn wowiface.IConnection) {
	//保护共享资源,  加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入到ConnManager中
	connMgr.connections[conn.GetConnId()] = conn
	if utils.GlobalObject.Develop {
		fmt.Println("添加链接完成链接Id", conn.GetConnId(), "当前链接数量:", connMgr.Len())
	}
}

// 删除链接
func (connMgr *ConnManager) Remove(conn wowiface.IConnection) {
	//保护共享资源,  加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接器
	delete(connMgr.connections, conn.GetConnId())
	if utils.GlobalObject.Develop {
		fmt.Println("已删除一个连接,连接Id为:", conn.GetConnId(), "当前连接数量:", connMgr.Len())
	}
}

// 根据connId获取链接
func (connMgr *ConnManager) Get(connId uint32) (wowiface.IConnection, error) {
	//保护共享资源,  加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connId]; ok {
		//找到了
		return conn, nil
	} else {
		//没找到
		return nil, errors.New(fmt.Sprintf("没有找到链接Id:%d", connId))
	}
}

// 获取服务端管理UI自身链接
func (connMgr *ConnManager) GetServerOwner() (wowiface.IConnection, error) {
	//保护共享资源,  加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	mac, err := utils.GetMACAddress()
	if err != nil {
		return nil, errors.New("没有找到在线管理员")
	}
	for _, v := range connMgr.connections {
		if v.GetConnMac() == mac {
			return v, nil
		}
	}
	return nil, errors.New("没有找到在线管理员")
}

// 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并停止所有链接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源,  加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除Conn,并停止Conn工作
	for connId, conn := range connMgr.connections {
		//停止删除
		conn.Stop()
		//删除
		delete(connMgr.connections, connId)
	}
	fmt.Println("所有连接全部删除完成,当前连接数量:", connMgr.Len())
}
