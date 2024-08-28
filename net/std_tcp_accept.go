package net

import (
	"github.com/ameise84/go_pool"
	"github.com/ameise84/pi_common/common"
	"github.com/libp2p/go-reuseport"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
)

type ListenerHandler interface {
	OnNewConnect(Tag, Conn) ConnContext //连接请求回调
	OnClosedConnect(ConnContext)
}

func newStdTcpAcceptor(tag Tag, addr string, opts TcpAceptorOptions, hand ListenerHandler, pool *stdTcpConnPool) (*stdTcpAcceptor, error) {
	network := string(TCP)

	ls, err := createStdTcpListener(network, addr, opts.ReusePort())
	if err != nil {
		return nil, err
	}

	ln := &stdTcpAcceptor{
		tag:      tag,
		network:  network,
		addr:     ls.Addr().String(),
		opts:     opts,
		hand:     hand,
		connPool: pool,
		ln:       ls,
	}
	ln.runner = go_pool.NewGoRunner(ln, "tcp acceptor", go_pool.DefaultOptions().SetSimCount(1))
	return ln, nil
}

func freeStdTcpAcceptor(ln *stdTcpAcceptor) {
	if ln.ln != nil {
		_ = ln.ln.Close()
		ln.ln = nil
	}
}

func createStdTcpListener(network, addr string, isReusePort bool) (net.Listener, error) {
	if isReusePort {
		return reuseport.Listen(network, addr)
	} else {
		return net.Listen(network, addr)
	}
}

type stdTcpAcceptor struct {
	s       common.Service
	tag     Tag
	network string
	addr    string

	opts     TcpAceptorOptions
	hand     ListenerHandler
	connPool *stdTcpConnPool
	ln       net.Listener

	pauseStat atomic.Bool
	isRunning atomic.Bool
	runner    go_pool.GoRunner
	connMap   sync.Map
}

func (a *stdTcpAcceptor) OnPanic(err error) {
	_gLogger.Error(err)
}

func (a *stdTcpAcceptor) where() string {
	return "std tcp acceptor " + a.tag
}

func (a *stdTcpAcceptor) start() error {
	return a.s.Start(func() error {
		if a.ln == nil {
			ls, err := createStdTcpListener(a.network, a.addr, a.opts.ReusePort())
			if err != nil {
				return err
			}
			a.ln = ls
		}
		return a.runner.AsyncRun(a.acceptLoop)
	})
}

func (a *stdTcpAcceptor) stop() {
	a.s.Stop(func() {
		a.isRunning.Store(false)
		_ = a.ln.Close()
		a.ln = nil
		a.runner.Wait()
	})
}

func (a *stdTcpAcceptor) getTag() Tag {
	return a.tag
}

func (a *stdTcpAcceptor) isPause() bool {
	return a.pauseStat.Load()
}

func (a *stdTcpAcceptor) pauseAccept() {
	a.pauseStat.Store(true)
}

func (a *stdTcpAcceptor) resumeAccept() {
	a.pauseStat.Store(false)
}

func (a *stdTcpAcceptor) closeAllConn() {
	mp := make(map[uint64]*stdTcpConn)
	a.connMap.Range(func(key, value any) bool {
		c := value.(*stdTcpConn)
		mp[c.instId] = c
		return true
	})
	for _, conn := range mp {
		_ = conn.CloseAndWaitRecv()
	}
}

func (a *stdTcpAcceptor) acceptLoop(...any) {
	a.isRunning.Store(true)
	for {
		fd, err := a.ln.Accept()
		if err != nil {
			if !a.isRunning.Load() {
				break
			}
			//logger.Error(errors.Wrap(err, "tcp server net accept"))
			runtime.Gosched()
			continue
		}
		if a.isPause() {
			_ = fd.Close()
			continue
		}
		c := a.connPool.take()
		c.active(a, a.tag, fd, a.opts.ReadTimeout(), a.opts.WriteTimeout())
		c.ctx = a.hand.OnNewConnect(a.tag, c)
		isRun := false
		if c.ctx != nil {
			isRun = c.run()
		}
		if isRun {
			a.connMap.Store(c.instId, c)
		} else {
			if c.ctx != nil {
				a.hand.OnClosedConnect(c.ctx)
				c.ctx = nil
			}
			c.inActive()
			_ = fd.Close()
			a.connPool.free(c)
			c.ctx = nil
		}
	}
	a.closeAllConn()
}

func (a *stdTcpAcceptor) handleCloseConn(c *stdTcpConn) {
	a.hand.OnClosedConnect(c.ctx)
	a.connMap.Delete(c.instId)
	c.ctx = nil
	c.inActive()
}
