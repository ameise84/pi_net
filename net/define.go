package net

import (
	"github.com/ameise84/pi_net/net/packet"
	"sync/atomic"
)

var _gConnInstId atomic.Uint64

type Port = uint16
type Tag = string //标签

type Mode string

const (
	TCP  Mode = "tcp"
	UDP  Mode = "udp"
	UNIX Mode = "unix"
)

const (
	IpV4 = "0.0.0.0"
	IpV6 = "[::]"
)

const (
	ReadTimeout  = 30
	WriteTimeout = 10
)

const (
	PacketSize  = 8192
	MsgSize     = PacketSize - packet.HeadSize
	CompressLen = 512
)
