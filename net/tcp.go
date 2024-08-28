package net

func NewTcpServer(ops ...TcpServerOptions) (TcpServer, error) {
	if ops == nil || ops[0] == nil {
		ops = []TcpServerOptions{DefaultTcpServerOptions()}
	}
	return newStdTcpServer(ops[0])
}

func NewTcpAgent(callback TcpAgentHandler, opts ...TcpAgentOptions) (TcpAgent, error) {
	if opts == nil {
		opts = []TcpAgentOptions{DefaultTcpAgentOptions()}
	}
	return newStdTcpAgent(callback, opts[0])
}
