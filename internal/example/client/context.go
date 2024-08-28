package main

import (
	"github.com/ameise84/logger"
	"github.com/ameise84/pi_common/str_conv"
	"github.com/ameise84/pi_net/net"
)

type CliContext struct{}

func (c *CliContext) OnRecv(conn net.Conn, bytes []byte) ([]byte, error) {
	logger.Info(str_conv.ToString(bytes))
	return nil, nil
}
