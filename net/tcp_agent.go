package net

type TcpAgentHandler interface {
	OnAgentConnect(agent TcpAgent, ctx ConnContext, conn Conn)
	OnAgentConnectFailed(agent TcpAgent, tag Tag, ctx ConnContext, err error)
	OnAgentClosedConn(agent TcpAgent, tag Tag, ctx ConnContext)
}

type TcpAgent interface {
	Start() error
	Stop()
	ConnectSync(tag Tag, ctx ConnContext, addr string) (Conn, error)
	ConnectAsync(tag Tag, ctx ConnContext, addr string) error
	DisConnectAll()
}
