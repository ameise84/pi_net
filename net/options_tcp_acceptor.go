package net

type TcpAceptorOptions interface {
	ReusePort() bool
	ReadTimeout() int64
	WriteTimeout() int64
	SetReusePort(reusePort bool) TcpAceptorOptions
	SetReadTimeout(val int64) TcpAceptorOptions
	SetWriteTimeout(val int64) TcpAceptorOptions
}

func DefaultTcpAcceptorOptions() TcpAceptorOptions {
	return &tcpAcceptorOptions{
		reusePort:    false,
		readTimeout:  ReadTimeout,
		writeTimeout: WriteTimeout,
	}
}

type tcpAcceptorOptions struct {
	reusePort    bool  //端口地址重用
	readTimeout  int64 // 连接静默时间
	writeTimeout int64 // 写入等待时间
}

func (o *tcpAcceptorOptions) ReusePort() bool {
	return o.reusePort
}

func (o *tcpAcceptorOptions) ReadTimeout() int64 {
	return o.readTimeout
}

func (o *tcpAcceptorOptions) WriteTimeout() int64 {
	return o.writeTimeout
}

func (o *tcpAcceptorOptions) SetReusePort(reusePort bool) TcpAceptorOptions {
	o.reusePort = reusePort
	return o
}

func (o *tcpAcceptorOptions) SetReadTimeout(val int64) TcpAceptorOptions {
	if val < 0 {
		panic("invalid read timeout")
	}
	o.readTimeout = val
	return o
}

func (o *tcpAcceptorOptions) SetWriteTimeout(val int64) TcpAceptorOptions {
	if val < 0 {
		panic("invalid write timeout")
	}
	o.writeTimeout = val
	return o
}
