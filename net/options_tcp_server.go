package net

import "github.com/ameise84/pi_net/net/packet"

type TcpServerOptions interface {
	PacketSize() uint32
	ConnHoldSize() uint32
	PackHandler() packet.PackHandler
	SetPacketSize(val uint32) TcpServerOptions
	SetConnHoldSize(val uint32) TcpServerOptions
	SetPackHandler(val packet.PackHandler) TcpServerOptions
}

func DefaultTcpServerOptions() TcpServerOptions {
	return &tcpServerOptions{
		packetSize:   PacketSize,
		connHoldSize: 256,
		packHandle:   packet.DefaultPackHandler,
	}
}

type tcpServerOptions struct {
	packetSize   uint32 // 消息包最大长度(打包以后得长度)
	connHoldSize uint32 // conn对象缓冲池回收阈值
	packHandle   packet.PackHandler
}

func (o *tcpServerOptions) PacketSize() uint32 {
	return o.packetSize
}

func (o *tcpServerOptions) ConnHoldSize() uint32 {
	return o.connHoldSize
}

func (o *tcpServerOptions) PackHandler() packet.PackHandler {
	return o.packHandle
}

func (o *tcpServerOptions) SetPacketSize(val uint32) TcpServerOptions {
	if val == 0 {
		panic("invalid packet size")
	}
	o.packetSize = val
	return o
}

func (o *tcpServerOptions) SetConnHoldSize(val uint32) TcpServerOptions {
	if val == 0 {
		panic("invalid conn hold size")
	}
	o.connHoldSize = val
	return o
}

func (o *tcpServerOptions) SetPackHandler(handle packet.PackHandler) TcpServerOptions {
	if handle == nil {
		panic("invalid packHandler")
	}
	o.packHandle = handle
	return o
}
