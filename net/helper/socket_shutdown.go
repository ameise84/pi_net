//go:build !windows
// +build !windows

package helper

import "golang.org/x/sys/unix"

type SockCloseReason uint32

const (
	CloseReasonNil       SockCloseReason = iota // 不关闭
	CloseReasonError                            // 异常关闭
	CloseReasonAccept                           // 服务器接受连接时触发
	CloseReasonConnect                          // 连接到服务器时触发
	CloseReasonSend                             // 发送消息时
	CloseReasonPending                          // 发送消息拥堵
	CloseReasonRecv                             // 接收消息时
	CloseReasonPrivilege                        // 紧急消息时
	CloseReasonClose                            // 正常关闭或断开连接
	CloseReasonShutdown                         // 服务关闭
)

type ShutDownFlag uint32

const (
	ShutR   ShutDownFlag = unix.SHUT_RD
	ShutW                = unix.SHUT_WR
	ShutRW               = unix.SHUT_RDWR
	ShutNil              = 0xFF
)
