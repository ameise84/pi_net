package helper

type SocketError struct {
	msg string
}

func (s *SocketError) Error() string {
	return s.msg
}

var (
	NetErrUnknown         = &SocketError{"net err: unknown"}
	NetErrAddress         = &SocketError{"net err: invalid address"}
	NetErrAddressFamily   = &SocketError{"net err: invalid address family"}
	NetErrAddressNotIPV4  = &SocketError{"net err: non-IPv4 address"}
	NetErrAddressNotIPV6  = &SocketError{"net err: non-IPv6 address"}
	NetErrAddressIPV6Zone = &SocketError{"net err: IPv6 zone"}
	NetErrNoListener      = &SocketError{"net err: no listener"}
	NetErrNoHandler       = &SocketError{"net err: no callback handler"}
	NetErrOptionsType     = &SocketError{"net err: options type illegal"}
	NetErrOptionsParams   = &SocketError{"net err: options params illegal"}
	NetErrConnDenied      = &SocketError{"net err: app connection denied"}
)
