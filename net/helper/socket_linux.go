//go:build linux

package helper

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"strconv"
)

func SockAddrToString(sa unix.Sockaddr) string {
	switch saImpl := (sa).(type) {
	case *unix.SockaddrInet4:
		return net.JoinHostPort(net.IP(saImpl.Addr[:]).String(), strconv.Itoa(saImpl.Port))
	case *unix.SockaddrInet6:
		return net.JoinHostPort(net.IP(saImpl.Addr[:]).String(), strconv.Itoa(saImpl.Port))
	default:
		return fmt.Sprintf("(unknown - %T)", saImpl)
	}
}

func IpToSockAddr(family int, ip net.IP, port int, zone string) (sockAddr unix.Sockaddr, err error) {
	switch family {
	case unix.AF_INET:
		if len(ip) == 0 {
			ip = net.IPv4zero
		}
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, NetErrAddressNotIPV4
		}
		sa := &unix.SockaddrInet4{Port: port}
		copy(sa.Addr[:], ip4)
		return sa, nil
	case unix.AF_INET6:
		if len(ip) == 0 || ip.Equal(net.IPv6zero) {
			ip = net.IPv6zero
		}
		ip6 := ip.To16()
		if ip6 == nil {
			return nil, NetErrAddressNotIPV6
		}

		var idx int
		if len(zone) == 0 {
			idx = 0
		} else {
			idx, err = strconv.Atoi(zone)
			if err != nil {
				return nil, NetErrAddressIPV6Zone
			}
		}

		sa := &unix.SockaddrInet6{Port: port, ZoneId: uint32(idx)}
		copy(sa.Addr[:], ip6)
		return sa, nil
	}
	return nil, NetErrAddressFamily
}

func CheckAddrFamily(network string) (int, bool) {
	switch network {
	case "tcp4", "udp4", "ip4":
		return unix.AF_INET, false
	case "tcp6", "udp6", "ip6":
		return unix.AF_INET6, true
	case "tcp", "udp", "ip":
		return unix.AF_INET, false
	default:
		return unix.AF_INET, false
	}
}

func ssoLinger(fd int, onOff int32, linger int32) error {
	return unix.SetsockoptLinger(fd, unix.SOL_SOCKET, unix.SO_LINGER, &unix.Linger{Onoff: onOff, Linger: linger})
}

func ssoUnLinger(fd int) error {
	//sso_linger(fd, 1, 0); // 关闭时立刻关闭，丢弃所有缓冲区数据，无TIME_WAIT状态
	//sso_linger(fd, 1, 1); // 关闭时，超时前尽量发送数据
	//sso_linger(fd, 0, x); // 关闭时，有数据残留，将数据发送完毕等待对方确认后关闭(默认行为)
	return ssoLinger(fd, 1, 0)
}

func SockNoDelay(fd, noDelay int) error {
	return unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_NODELAY, noDelay)
}

func CloseSocket(fd int, flag ShutDownFlag, linger bool) {
	if fd != InvalidSocket {
		if !linger {
			_ = ssoUnLinger(fd)
		}

		if flag != ShutNil {
			_ = unix.Shutdown(fd, int(flag))
		}

		_ = unix.Close(fd)
	}
}

func SockKeepAliveVal(sock int, isOnOff bool, idle, interval, count int) (err error) {
	if isOnOff {
		if err = unix.SetsockoptInt(sock, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 1); err != nil {
			return err
		}
		if err = unix.SetsockoptInt(sock, unix.SOL_SOCKET, unix.TCP_KEEPIDLE, idle); err != nil {
			return err
		}
		if err = unix.SetsockoptInt(sock, unix.SOL_SOCKET, unix.TCP_KEEPINTVL, interval); err != nil {
			return err
		}
		if err = unix.SetsockoptInt(sock, unix.SOL_SOCKET, unix.TCP_KEEPCNT, count); err != nil {
			return err
		}
	} else {
		if err = unix.SetsockoptInt(sock, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 0); err != nil {
			return err
		}
	}
	return nil
}
