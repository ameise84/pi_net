package main

import (
	"github.com/ameise84/logger"
	"github.com/ameise84/pi_common/sys"
	"github.com/ameise84/pi_net/net"
)

type cliHandler struct {
}

func (c *cliHandler) OnAgentConnect(agent net.TcpAgent, ctx net.ConnContext, conn net.Conn) {

}

func (c *cliHandler) OnAgentConnectFailed(agent net.TcpAgent, tag net.Tag, ctx net.ConnContext, err error) {

}

func (c *cliHandler) OnAgentClosedConn(agent net.TcpAgent, tag net.Tag, ctx net.ConnContext) {

}

func main() {
	c, _ := net.NewTcpAgent(&cliHandler{})
	_ = c.Start()
	cc, err := c.ConnectSync("2", &CliContext{}, ":12345")
	if err != nil {
		logger.Error(err)
		return
	}
	_ = cc.SendSync([]byte("123"))
	sys.WaitKillSigint()
}
