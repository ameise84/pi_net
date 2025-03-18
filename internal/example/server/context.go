package main

import (
	"github.com/ameise84/pi_net/net"
)

type SrvContext struct{}

func (c *SrvContext) OnRecv(conn net.Conn, bytes []byte) ([]byte, error) {
	rsp := make([]byte, len(bytes))
	copy(rsp, bytes)
	conn.SendAsync(rsp)
	return nil, nil
}
