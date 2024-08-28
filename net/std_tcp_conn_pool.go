package net

import (
	"github.com/ameise84/lock"
	"github.com/ameise84/pi_net/net/packet"
	"github.com/ameise84/queue"
	"sync"
	"sync/atomic"
	"time"
)

const (
	gcTime      = 5  //每5秒检测一次gc
	invalidTime = 20 //连接对象释放超过20秒可以复用
)

type stdTcpConnPool struct {
	bufferPool *packetBufferPool
	packHandle packet.PackHandler
	gcLock     lock.SpinLock
	gc         queue.IDeque[*stdTcpConn]
	pool       *sync.Pool
	lastGC     atomic.Int64
}

func newConnPool(ConnHoldSize uint32, packetSize uint32, packHandle packet.PackHandler) *stdTcpConnPool {
	bufferPool := newPacketBufferPool(int(packetSize))
	return &stdTcpConnPool{
		bufferPool: bufferPool,
		packHandle: packHandle,
		gc:         queue.NewRingQueue[*stdTcpConn](int32(ConnHoldSize)),
		pool: &sync.Pool{
			New: func() any {
				c := newStdTcpConn(bufferPool, packHandle)
				return c
			},
		},
	}
}

func (s *stdTcpConnPool) take() (c *stdTcpConn) {
	c = s.checkGC(true)
	if c == nil {
		c = s.pool.Get().(*stdTcpConn)
	}
	return
}

func (s *stdTcpConnPool) free(c *stdTcpConn) {
	s.gcLock.Lock()
	err := s.gc.PushBack(c)
	s.gcLock.Unlock()
	if err != nil {
		s.checkGC(false)
	}
}

func (s *stdTcpConnPool) checkGC(takeOne bool) (c *stdTcpConn) {
	now := time.Now().Unix()
	if now-s.lastGC.Load() > gcTime {
		s.lastGC.Store(now)
		s.gcLock.Lock()
		defer s.gcLock.Unlock()
		for {
			tmpConn, err := s.gc.PopFront()
			if err != nil {
				break
			}
			if now-tmpConn.invalidTime > invalidTime {
				if takeOne && c == nil {
					c = tmpConn
				} else {
					s.pool.Put(tmpConn)
				}
			} else {
				_ = s.gc.PushFront(tmpConn)
				break
			}
		}
	}
	return
}
