package net

import "github.com/ameise84/pi_net/net/packet"

type TcpAgentOptions interface {
	PacketSize() uint32
	ConnHoldSize() uint32
	ReadTimeout() int64
	WriteTimeout() int64
	PackHandler() packet.PackHandler
	SetPacketSize(val uint32) TcpAgentOptions
	SetConnHoldSize(val uint32) TcpAgentOptions
	SetReadTimeout(val int64) TcpAgentOptions
	SetWriteTimeout(val int64) TcpAgentOptions
	SetPackHandler(val packet.PackHandler) TcpAgentOptions
}

func DefaultTcpAgentOptions() TcpAgentOptions {
	return &tcpAgentOptions{
		packetSize:   PacketSize,
		connHoldSize: 4,
		readTimeout:  0,
		writeTimeout: WriteTimeout,
		packHandle:   packet.DefaultPackHandler,
	}
}

type tcpAgentOptions struct {
	packetSize   uint32 //消息包最大长度
	connHoldSize uint32 //conn对象缓冲池回收阈值
	readTimeout  int64  // 连接静默时间
	writeTimeout int64  // 写入等待时间
	packHandle   packet.PackHandler
}

func (o *tcpAgentOptions) PacketSize() uint32 {
	return o.packetSize
}

func (o *tcpAgentOptions) ConnHoldSize() uint32 {
	return o.connHoldSize
}

func (o *tcpAgentOptions) ReadTimeout() int64 {
	return o.readTimeout
}

func (o *tcpAgentOptions) WriteTimeout() int64 {
	return o.writeTimeout
}

func (o *tcpAgentOptions) PackHandler() packet.PackHandler {
	return o.packHandle
}

func (o *tcpAgentOptions) SetPacketSize(val uint32) TcpAgentOptions {
	if val == 0 {
		panic("invalid packet size")
	}
	o.packetSize = val
	return o
}

func (o *tcpAgentOptions) SetConnHoldSize(val uint32) TcpAgentOptions {
	if val == 0 {
		panic("invalid conn hold size")
	}
	o.connHoldSize = val
	return o
}

func (o *tcpAgentOptions) SetReadTimeout(val int64) TcpAgentOptions {
	o.readTimeout = val
	return o
}

func (o *tcpAgentOptions) SetWriteTimeout(val int64) TcpAgentOptions {
	o.writeTimeout = val
	return o
}

func (o *tcpAgentOptions) SetPackHandler(val packet.PackHandler) TcpAgentOptions {
	if val == nil {
		panic("invalid packHandler")
	}
	o.packHandle = val
	return o
}
