package network

import (
	"net"
	"sync"
	"time"

	"comm/util"
)

type TCPClient struct {
	opts 			*TCPClientOptions
	newAgent        func(*TCPConn) Agent
	conns           map[net.Conn]struct{}
	wg              util.WaitGroupWrapper
	closeFlag       bool

	sync.Mutex
}

func (c *TCPClient) NewTCPClient(options *TCPClientOptions, newAgentFun func(*TCPConn) Agent) *TCPClient {
	return &TCPClient {
		opts: options,
		newAgent: newAgentFun,
		conns: make(map[net.Conn]struct{}),
	}
}

func (c *TCPClient) getOpts() *TCPClientOptions {
	return c.opts
}

func (c *TCPClient) Run() {
	for i := 0; i < c.getOpts().ConnNum; i ++ {
		c.wg.Wrap(c.connect)
	}
}

func (c *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", c.opts.ServerAddr)
		if err == nil || c.closeFlag {
			return conn
		}

		time.Sleep(c.opts.ConnectInterval)
		continue
	}
}

func (c *TCPClient) connect() {
reconnect:
	conn := c.dial()
	if conn == nil {
		return
	}

	c.Lock()
	if c.closeFlag {
		c.Unlock()
		conn.Close()
		return
	}
	c.conns[conn] = struct{}{}
	c.Unlock()

	tcpConn := newTCPConn(conn, c.opts.PendingWriteNum, nil)
	agent := c.newAgent(tcpConn)
	agent.Run()

	tcpConn.Close()
	c.Lock()
	delete(c.conns, conn)
	c.Unlock()
	agent.OnClose()

	if c.opts.AutoReconnect {
		time.Sleep(c.opts.ConnectInterval)
		goto reconnect
	}
}

func (c *TCPClient) Close() {
	c.Lock()
	c.closeFlag = true
	for conn := range c.conns {
		conn.Close()
	}
	c.conns = nil
	c.Unlock()

	c.wg.Wait()
}
