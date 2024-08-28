package net

type Conn interface {
	ID() uint64
	Tag() Tag
	SetContext(ctx ConnContext)
	Context() ConnContext
	LocalAddr() string
	RemoteAddr() string
	SendSync([]byte) error //在收到OnSendOver消息前,send对象不可释放,不可复用
	SendAsync([]byte) error
	IsConnected() bool
	Close() error
	CloseAndWaitRecv() error
	TrafficData() string
}

type ConnContext interface {
	OnRecv(Conn, []byte) ([]byte, error) //OnRecv的数据不可交给其他对象复用.如需复用请自行拷贝
}
