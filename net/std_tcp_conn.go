package net

import (
	"fmt"
	"github.com/ameise84/go_pool"
	"github.com/ameise84/lock"
	"github.com/ameise84/logger"
	"github.com/ameise84/pi_common/bytes_buffer"
	"github.com/ameise84/pi_common/errors"
	"github.com/ameise84/pi_net/net/packet"
	"net"
	"runtime"
	"sync/atomic"
	"time"
)

type iConnFactory interface {
	handleCloseConn(*stdTcpConn)
}

func newStdTcpConn(bufferPool *packetBufferPool, packHandle packet.PackHandler) *stdTcpConn {
	c := &stdTcpConn{
		bufferPool:  bufferPool,
		packHandle:  packHandle,
		packMsgBuff: bufferPool.Get(),
	}
	sendWorker := go_pool.NewGoRunner(c, "tcp conn send", go_pool.DefaultOptions().SetSimCount(1).SetBlock(false).SetCacheMode(true, 256))
	recvWorker := go_pool.NewGoRunner(c, "tcp conn recv", go_pool.DefaultOptions().SetSimCount(1))
	c.sendWorker = sendWorker
	c.recvWorker = recvWorker
	c.isSendWake.Store(false)
	return c
}

type stdTcpConn struct {
	instId      uint64
	factory     iConnFactory
	tag         Tag
	fd          net.Conn
	isConnected bool        //是否处于连接状态
	ctx         ConnContext //应用层扩展数据
	packHandle  packet.PackHandler
	readTime    time.Duration
	writeTime   time.Duration
	invalidTime int64 //失效时间
	closeOnce   lock.Once
	closeLock   lock.SpinLock
	sendWorker  go_pool.GoRunner
	recvWorker  go_pool.GoRunner
	isSendWake  atomic.Bool
	bufferPool  *packetBufferPool
	packMsgBuff bytes_buffer.ShiftBuffer
	sendCount   atomic.Uint64
	sendSize    atomic.Uint64
	recvCount   atomic.Uint64
	recvSize    atomic.Uint64
}

func (c *stdTcpConn) LogFmt() string {
	return fmt.Sprintf("net conn[<%s> <%s>][%v %d]", c.LocalAddr(), c.RemoteAddr(), c.tag, c.instId)
}

func (c *stdTcpConn) OnPanic(err error) {
	_gLogger.ErrorBeans([]logger.Bean{c}, err)
}

func (c *stdTcpConn) ID() uint64 {
	return c.instId
}

func (c *stdTcpConn) Tag() Tag {
	return c.tag
}

func (c *stdTcpConn) SetContext(ctx ConnContext) {
	c.closeLock.Lock()
	if c.isConnected {
		c.ctx = ctx
	}
	c.closeLock.Unlock()
}

func (c *stdTcpConn) Context() ConnContext {
	return c.ctx
}

func (c *stdTcpConn) LocalAddr() string {
	return c.fd.LocalAddr().String()
}

func (c *stdTcpConn) RemoteAddr() string {
	return c.fd.RemoteAddr().String()
}

func (c *stdTcpConn) SendSync(msg []byte) error {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()
	if !c.isConnected {
		return ErrorConnClosed
	}
	bf := c.bufferPool.Get()
	bf.Clean()
	wf, _ := bf.GetTailEmptyBytes()
	wn, err := c.packHandle.Pack(msg, wf)
	if err != nil {
		return err
	}
	bf.AddLen(wn)
	_, err = c.sendWorker.SyncRun(c.doSendSync, bf)
	return err
}

func (c *stdTcpConn) SendAsync(msg []byte) (err error) {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()

	if !c.isConnected {
		return ErrorConnClosed
	}

	bf := c.bufferPool.Get()
	bf.Clean()
	wf, _ := bf.GetTailEmptyBytes()
	wn, err := c.packHandle.Pack(msg, wf)
	if err != nil {
		return err
	}
	bf.AddLen(wn)
	err = c.sendWorker.AsyncRun(c.doSendAsync, bf)
	return err
}

func (c *stdTcpConn) IsConnected() bool {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()
	return c.isConnected
}

func (c *stdTcpConn) Close() (err error) {
	c.closeOnce.Do(func() {
		c.closeLock.Lock()
		defer c.closeLock.Unlock()
		if c.isConnected {
			c.sendWorker.Wait() //等待所有发送线程关闭
			c.isConnected = false
			if c.fd != nil {
				_ = c.fd.Close()
			}
			//启动一个临时协程来等待conn完全退出
			go_pool.NewGoFuncDo(c, "tcp conn close", func() {
				c.recvWorker.Wait() //
				logger.TracePrintf("%v close", c.TrafficData())
				c.factory.handleCloseConn(c)
			})
		}
	})
	return
}

func (c *stdTcpConn) CloseAndWaitRecv() error {
	err := c.Close()
	if err == nil {
		c.recvWorker.Wait()
	}
	return err
}

func (c *stdTcpConn) TrafficData() string {
	return fmt.Sprintf("%s%s", c.LogFmt(), fmt.Sprintf("{send:{count:%v,size:%v},recv:{count:%v,size:%v}}", c.sendCount.Load(), c.sendSize.Load(), c.recvCount.Load(), c.recvSize.Load()))
}

func (c *stdTcpConn) loopRead(...any) {
	var err error
	var n int
	rdBuffer := bytes_buffer.NewShiftBuffer(c.bufferPool.size*4, 0, c.bufferPool.grown)
	for {
		if c.readTime > 0 {
			readOverTime := time.Now().Add(c.readTime)
			err = c.fd.SetReadDeadline(readOverTime)
			if err != nil {
				break
			}
		}
		bf, _ := rdBuffer.GetTailEmptyBytes()
		n, err = c.fd.Read(bf)
		if err != nil {
			break
		}
		if n <= 0 {
			runtime.Gosched()
			continue
		}
		c.recvSize.Add(uint64(n))
		rdBuffer.AddLen(n)
		err = c.handleUnPacketMsg(rdBuffer)
		if err != nil {
			break
		}
	}

	if isCaredError(err) { //排除正常关闭
		_gLogger.ErrorBeans([]logger.Bean{c}, errors.WrapNoStack(err, "handle recv"))
	} else if rdBuffer.GetDataSize() != 0 {
		_gLogger.ErrorBeans([]logger.Bean{c}, errors.NewOrWrapNoStack(err, "handle recv data!=0"))
	}
	_ = c.Close()
}

func (c *stdTcpConn) handleUnPacketMsg(read bytes_buffer.Reader) error {
	var rsp []byte
	w, _ := c.packMsgBuff.GetTailEmptyBytes()

	for {
		r, _ := read.Peek()
		if len(r) == 0 {
			return nil
		}

		rn, msg, err := c.packHandle.UnPack(r, w)
		if rn == 0 || err != nil {
			return err
		}

		_, _ = read.FetchLen(rn)
		c.recvCount.Add(1)
		rsp, err = c.ctx.OnRecv(c, msg)
		if len(rsp) > 0 {
			sendErr := c.SendAsync(rsp)
			if sendErr != nil {
				err = errors.NewOrWrap(err, sendErr.Error())
			}
		}
		if err != nil {
			return err
		}
	}
}

func (c *stdTcpConn) active(f iConnFactory, tag Tag, fd net.Conn, readTime int64, writeTime int64) {
	c.instId = _gConnInstId.Add(1) //正常情况下instId不可能重复, 除非遇见当前 conn 申请后,一直持续到后面有 MAXUINT64 个用户登录,当前conn还持续在线.实际环境中,几乎是不可能的
	c.factory = f
	c.tag = tag
	c.fd = fd
	c.readTime = time.Duration(readTime) * time.Second
	c.writeTime = time.Duration(writeTime) * time.Second
	c.packMsgBuff.Clean()
	c.sendCount.Store(0)
	c.sendSize.Store(0)
	c.recvCount.Store(0)
	c.recvSize.Store(0)
	c.closeOnce.Reset()
	return
}

func (c *stdTcpConn) inActive() (fd net.Conn) {
	c.factory = nil
	c.invalidTime = time.Now().Unix()
	fd = c.fd
	c.fd = nil
	return
}

func (c *stdTcpConn) run() bool {
	c.isConnected = true
	if err := c.recvWorker.AsyncRun(c.loopRead); err != nil {
		c.isConnected = false
		return false
	}
	return true
}

func (c *stdTcpConn) doSendAsync(args ...any) {
	bf := args[0].(bytes_buffer.ShiftBuffer)
	defer func() {
		c.bufferPool.Put(bf)
	}()
	msg, _, _ := bf.Fetch()
	err := c.doSend(msg)
	if err != nil {
		_gLogger.WarnBeans([]logger.Bean{c}, errors.WrapNoStack(err, fmt.Sprintf("do send:%v", msg)).Error())
	}
}

func (c *stdTcpConn) doSendSync(args ...any) (any, error) {
	bf := args[0].(bytes_buffer.ShiftBuffer)
	defer func() {
		c.bufferPool.Put(bf)
	}()
	msg, _, _ := bf.Fetch()
	return nil, c.doSend(msg)
}

func (c *stdTcpConn) doSend(msg []byte) (err error) {
	fd := c.fd
	var n int
	if c.writeTime > 0 {
		if err := fd.SetWriteDeadline(time.Now().Add(c.writeTime)); err != nil {
			return err
		}
	}

	writeLen := len(msg)
	idx := 0
	for {
		n, err = fd.Write(msg[idx:])
		if err != nil {
			go_pool.NewGoFuncDo(c, "tcp conn send close", func() {
				_ = c.Close()
			})
			return err
		}
		idx += n
		if idx == writeLen {
			break
		}
	}
	c.sendCount.Add(1)
	c.sendSize.Add(uint64(writeLen))
	return nil
}
