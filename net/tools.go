package net

import (
	"errors"
	"github.com/ameise84/pi_common/str_conv"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"syscall"
)

var (
	_gInnerIp   string
	_gInnerOnce sync.Once

	_gExternalIp   string
	_gExternalOnce sync.Once

	_gHostName     string
	_gHostNameOnce sync.Once

	netErr              net.Error
	dNSError            *net.DNSError
	parseError          *net.ParseError
	addrError           *net.AddrError
	unknownNetworkError *net.UnknownNetworkError
	dNSConfigError      *net.DNSConfigError
	invalidAddrError    *net.InvalidAddrError
	syscallError        *os.SyscallError
)

func isTimeOutError(err error) bool {
	ok := errors.As(err, &netErr)
	if !ok {
		return false
	}
	if netErr.Timeout() {
		return true
	}
	return false
}
func isCaredError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, io.EOF) {
		return false
	}

	ok := errors.As(err, &netErr)
	if !ok {
		return true
	}

	if netErr.Timeout() {
		return false
	}

	var opErr *net.OpError
	ok = errors.As(netErr, &opErr)
	if !ok {
		return false
	}

	switch {
	case errors.As(opErr.Err, &dNSError):
		return true
	case errors.As(opErr.Err, &parseError):
		return true
	case errors.As(opErr.Err, &addrError):
		return true
	case errors.As(opErr.Err, &unknownNetworkError):
		return true
	case errors.As(opErr.Err, &dNSConfigError):
		return true
	case errors.As(opErr.Err, &invalidAddrError):
		return true
	case errors.As(opErr.Err, &syscallError):
		var errno syscall.Errno
		if errors.As(opErr.Err, &errno) {
			switch {
			case errors.Is(errno, syscall.EINVAL):
				return true
			case errors.Is(errno, syscall.ECONNREFUSED):
				return true
			case errors.Is(errno, syscall.ETIMEDOUT):
				return false
			default:
			}
		}
	default:
	}
	return false
}

func GetInternalIP() string {
	_gInnerOnce.Do(func() {
		var err error
		_gInnerIp, err = getInternalIP()
		if err != nil {
			_gLogger.Fatal(err)
		}
	})
	return _gInnerIp
}

func GetExternalIP() string {
	_gExternalOnce.Do(func() {
		var err error
		_gExternalIp, err = getExternalIp()
		if err != nil {
			_gExternalIp = GetInternalIP()
		}
	})
	return _gExternalIp
}

func GetHostName() string {
	_gHostNameOnce.Do(func() {
		var err error
		_gHostName, err = os.Hostname()
		if err != nil {
			_gHostName = GetExternalIP()
		}
	})
	return _gHostName
}

func getInternalIP() (string, error) {
	// 获取本机的网络接口列表
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", errors.New("no internal ip found")
	}

	// 遍历接口列表，查找IPv4地址
	for _, iface := range interfaces {
		// 排除无效和回环接口
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err2 := iface.Addrs()
			if err2 != nil {
				continue
			}

			// 遍历接口地址，查找IPv4地址
			for _, addr := range addrs {
				// 检查地址是否为IPv4地址
				if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
					return ip.IP.String(), nil
				}
			}
		}
	}
	return "", errors.New("no internal ip found")
}

func getExternalIp() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("no external ip found")
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("no external ip found")
	}
	return str_conv.ToString(ip), nil
}

func getIPV6() (string, error) {
	resp, err := http.Get("https://ipv6.netarm.com")
	if err != nil {
		return "", errors.New("no ipv6 found")
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("no ipv6 found")
	}
	return str_conv.ToString(ip), nil
}

func SplitHostPort(addr string) (string, string, error) {
	return net.SplitHostPort(addr)
}

func GetHost(addr string) (string, error) {
	a, _, err := net.SplitHostPort(addr)
	return a, err
}

func GetPort(addr string) (string, error) {
	_, p, err := net.SplitHostPort(addr)
	if err != nil {
		return "0", err
	}
	return p, nil
}
