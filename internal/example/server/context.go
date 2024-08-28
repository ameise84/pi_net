package main

import (
	"github.com/ameise84/logger"
	"github.com/ameise84/pi_common/str_conv"
	"github.com/ameise84/pi_net/net"
)

type SrvContext struct{}

func (c *SrvContext) OnRecv(conn net.Conn, bytes []byte) ([]byte, error) {
	logger.Info(str_conv.ToString(bytes))
	rsp := make([]byte, len(bytes))
	copy(rsp, bytes)
	conn.SendAsync(rsp)
	return nil, nil
}
