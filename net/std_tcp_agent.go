package net

import (
	"errors"
	"github.com/ameise84/go_pool"
	"github.com/ameise84/logger"
	"github.com/ameise84/pi_common/common"
	"github.com/ameise84/pi_net/net/helper"
	"net"
	"sync"
)

var errOnAgentConnect = errors.New("hand OnAgentConnect return false")

func newStdTcpAgent(hand TcpAgentHandler, opts TcpAgentOptions) (*tcpAgent, error) {
	if hand == nil {
		return nil, helper.NetErrNoHandler
	}
	s := &tcpAgent{
		hand:     hand,
		opts:     opts,
		connPool: newConnPool(opts.ConnHoldSize(), opts.PacketSize(), opts.PackHandler()),
	}
	s.goRunner = go_pool.NewGoRunner(s, "tcp agent", go_pool.DefaultOptions().SetSimCount(0).SetBlock(false))
	return s, nil
}

type tcpAgent struct {
	common.Service
	hand     TcpAgentHandler
	opts     TcpAgentOptions
	connPool *stdTcpConnPool
	connMap  sync.Map
	goRunner go_pool.GoRunner
}

func (s *tcpAgent) LogFmt() string {
	return "net tcp agent"
}

func (s *tcpAgent) OnPanic(err error) {
	_gLogger.ErrorBeans([]logger.Bean{s}, err)
}

func (s *tcpAgent) Start() error {
	return s.Service.Start(nil)
}

func (s *tcpAgent) Stop() {
	s.Service.Stop(func() {
		s.DisConnectAll()
	})
}

func (s *tcpAgent) ConnectSync(tag Tag, ctx ConnContext, addr string) (Conn, error) {
	if !s.IsRunning() {
		return nil, errors.New("tcp agent is not running")
	}
	if ctx == nil {
		return nil, errors.New("stdTcpConn context is nil")
	}
	c, err := s.goRunner.SyncRun(dialSync, s, tag, ctx, addr)
	if err != nil {
		return nil, err
	}
	return c.(Conn), err
}

func (s *tcpAgent) ConnectAsync(tag Tag, ctx ConnContext, addr string) error {
	if !s.IsRunning() {
		return errors.New("tcp agent is not running")
	}
	if ctx == nil {
		return errors.New("stdTcpConn context is nil")
	}
	return s.goRunner.AsyncRun(dialAsync, s, tag, ctx, addr)
}

func (s *tcpAgent) DisConnectAll() {
	mp := make(map[uint64]*stdTcpConn)
	s.connMap.Range(func(key, value any) bool {
		c := value.(*stdTcpConn)
		mp[c.instId] = c
		return true
	})
	for _, conn := range mp {
		_ = conn.CloseAndWaitRecv()
	}
}

func (s *tcpAgent) handleCloseConn(c *stdTcpConn) {
	s.hand.OnAgentClosedConn(s, c.tag, c.ctx)
	s.connMap.Delete(c.instId)
	c.ctx = nil
	c.inActive()
	s.connPool.free(c)
}

func (s *tcpAgent) handleNewConnection(tag Tag, fd net.Conn, ctx ConnContext) *stdTcpConn {
	c := s.connPool.take()
	c.active(s, tag, fd, 0, WriteTimeout)
	c.ctx = ctx
	if c.run() {
		s.hand.OnAgentConnect(s, ctx, c)
		s.connMap.Store(c.instId, c)
		return c
	} else {
		s.hand.OnAgentConnectFailed(s, tag, ctx, errors.New("agent connect run failed"))
		return nil
	}
}

func (s *tcpAgent) connectTo(addr string) (net.Conn, error) {
	network := string(TCP)
	fd, err := net.Dial(network, addr)

	if err != nil {
		return nil, err
	}
	return fd, nil
}

func dialAsync(args ...any) {
	_, _ = dialSync(args...)
}

func dialSync(args ...any) (any, error) {
	s := args[0].(*tcpAgent)
	tag := args[1].(Tag)
	ctx := args[2].(ConnContext)
	addr := args[3].(string)
	fd, err := s.connectTo(addr)
	if err != nil {
		s.hand.OnAgentConnectFailed(s, tag, ctx, err)
		return nil, err
	}

	c := s.handleNewConnection(tag, fd, ctx)
	return c, nil
}
