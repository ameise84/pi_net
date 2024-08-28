package packet

import (
	"io"
)

var DefaultPackHandler PackHandler = &defaultPackHandler{}

type Writer interface {
	io.Writer
	GetPacketSize() int //每个消息包大小
}

// PackHandler 打包拆包接口
type PackHandler interface {
	UnPack(readBuff, writeBuff []byte) (readSize int, out []byte, err error) //解包
	Pack(readBuff, writeBuff []byte) (writeSize int, err error)              //封包
}

// defaultPackHandler 默认 PackHandler
type defaultPackHandler struct{}

// UnPack 拆包
func (d *defaultPackHandler) UnPack(readBuff, writeBuff []byte) (readSize int, out []byte, err error) {
	packSize := len(readBuff)
	if packSize < HeadSize {
		return 0, nil, nil
	}
	h := NewHeadWarp(readBuff[:HeadSize])
	wn := int(h.BodySize())
	readSize = HeadSize + wn
	if packSize < readSize {
		return 0, nil, nil
	}
	if wn > len(writeBuff) {
		return 0, nil, io.ErrShortWrite
	}
	copy(writeBuff[0:], readBuff[HeadSize:readSize])
	return readSize, writeBuff[:wn], nil
}

// Pack 封包
func (d *defaultPackHandler) Pack(readBuff, writeBuff []byte) (writeSize int, err error) {
	dataSize := len(readBuff)
	writeSize = dataSize + HeadSize
	if writeSize > len(writeBuff) {
		return 0, io.ErrShortWrite
	}
	h := NewHeadWarp(writeBuff[:HeadSize])
	h.Init()
	h.SetBodySize(uint16(dataSize))
	copy(writeBuff[HeadSize:], readBuff[0:])
	return writeSize, nil
}
