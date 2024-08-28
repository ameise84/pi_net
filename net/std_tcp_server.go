package net

import (
	"errors"
	"github.com/ameise84/lock"
	"github.com/ameise84/pi_common/common"
	"github.com/ameise84/pi_net/net/helper"
)

func newStdTcpServer(opts TcpServerOptions) (*stdTcpServer, error) {
	s := &stdTcpServer{
		connPool:    newConnPool(opts.ConnHoldSize(), opts.PacketSize(), opts.PackHandler()),
		acceptorMap: map[Tag]*stdTcpAcceptor{},
		opts:        opts,
	}
	return s, nil
}

type stdTcpServer struct {
	common.Service
	mu          lock.ReinLock
	connPool    *stdTcpConnPool
	acceptorMap map[Tag]*stdTcpAcceptor
	opts        TcpServerOptions
}

func (s *stdTcpServer) Start() error {
	if len(s.acceptorMap) == 0 {
		return helper.NetErrNoListener
	}
	return s.Service.Start(
		func() (err error) {
			s.mu.Lock()
			for _, ln := range s.acceptorMap {
				err = ln.start()
				if err != nil {
					break
				}
			}
			if err != nil {
				for _, l := range s.acceptorMap {
					l.stop()
				}
			}
			s.mu.Unlock()
			return
		},
	)
}

func (s *stdTcpServer) Stop() {
	s.Service.Stop(
		func() {
			s.mu.Lock()
			for _, ln := range s.acceptorMap {
				delete(s.acceptorMap, ln.tag)
				ln.stop()
				freeStdTcpAcceptor(ln)
			}
			s.mu.Unlock()
		})
}

func (s *stdTcpServer) AddAcceptor(hand ListenerHandler, tag Tag, addr string, pause bool, opts ...TcpAceptorOptions) (string, error) {
	var o TcpAceptorOptions
	if opts != nil && opts[0] != nil {
		o = opts[0]
	} else {
		o = DefaultTcpAcceptorOptions()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.acceptorMap[tag]; ok {
		return "", errors.New("the tag is added")
	}

	acceptor, err := newStdTcpAcceptor(tag, addr, o, hand, s.connPool)

	if err != nil {
		return "", err
	}

	s.acceptorMap[tag] = acceptor

	if pause {
		acceptor.pauseAccept()
	}
	if s.IsRunning() {
		err = acceptor.start()
		if err != nil {
			freeStdTcpAcceptor(acceptor)
			return "", err
		}
	}
	return acceptor.addr, nil
}

func (s *stdTcpServer) RemoveAcceptor(tag Tag) {
	s.mu.Lock()
	acceptor, ok := s.acceptorMap[tag]
	delete(s.acceptorMap, tag)
	s.mu.Unlock()
	if ok {
		acceptor.stop()
		freeStdTcpAcceptor(acceptor)
	}
}

func (s *stdTcpServer) DisConnectByTag(tag Tag, pause bool) {
	s.mu.Lock()
	acceptor, ok := s.acceptorMap[tag]
	s.mu.Unlock()

	if ok {
		isPause := acceptor.isPause()
		if !isPause {
			acceptor.pauseAccept()
		}
		acceptor.closeAllConn()
		if !isPause && !pause {
			acceptor.resumeAccept()
		}
	}
}

func (s *stdTcpServer) DisConnectAll(pause bool) {
	s.mu.Lock()
	for _, acceptor := range s.acceptorMap {
		isPause := acceptor.isPause()
		if !isPause {
			acceptor.pauseAccept()
		}
		acceptor.closeAllConn()
		if !isPause && !pause {
			acceptor.resumeAccept()
		}
	}
	s.mu.Unlock()
}

func (s *stdTcpServer) PauseAcceptByTag(tag Tag) {
	s.mu.Lock()
	if acceptor, ok := s.acceptorMap[tag]; ok {
		acceptor.pauseAccept()
	}
	s.mu.Unlock()
}

func (s *stdTcpServer) PauseAccept() {
	s.mu.Lock()
	for _, acceptor := range s.acceptorMap {
		acceptor.pauseAccept()
	}
	s.mu.Unlock()
}

func (s *stdTcpServer) ResumeAcceptByTag(tag Tag) {
	s.mu.Lock()
	if acceptor, ok := s.acceptorMap[tag]; ok {
		acceptor.resumeAccept()
	}
	s.mu.Unlock()
}

func (s *stdTcpServer) ResumeAccept() {
	s.mu.Lock()
	for _, acceptor := range s.acceptorMap {
		acceptor.resumeAccept()
	}
	s.mu.Unlock()
}
