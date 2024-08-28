package main

import (
	"github.com/ameise84/pi_common/sys"
	"github.com/ameise84/pi_net/net"
)

type srvHandler struct {
}

func (s *srvHandler) OnNewConnect(tag net.Tag, conn net.Conn) net.ConnContext {
	return &SrvContext{}
}

func (s *srvHandler) OnClosedConnect(context net.ConnContext) {
}

func main() {
	s, _ := net.NewTcpServer()
	_, _ = s.AddAcceptor(&srvHandler{}, "1", ":12345", false)
	_ = s.Start()
	sys.WaitKillSigint()
}
