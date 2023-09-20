package wownet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"gitee.com/mrmateoliu/wow_launch.git/wowiface"
)

// 封包,拆包的具体模块
type DataPack struct{}

// 拆包封包示实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头长度的方法
func (d *DataPack) GetHeadLen() uint32 {

	//DataLen uint32 (4字节) + ID uint32 (4字节)
	return 8
}

// 封包方法
// |datalen|msgId|data|
func (d *DataPack) Pack(msg wowiface.IMessage) ([]byte, error) {
	//创建一个byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将DataLen 写进dataBuff中
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	//将MsgId 写进dataBuff中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	//将Data数据 写进dataBuff中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法(将包的head信息读出来,再根据head信息里的data长度在进行读)
func (d *DataPack) UnPack(binaryData []byte) (wowiface.IMessage, error) {
	//创建一个从输入二进制的IoReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息, 得到datalen和msgid

	//创建一个Message结构
	msg := &Message{}

	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出了允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("接收到的数据包尺寸过大")
	}

	return msg, nil
}
