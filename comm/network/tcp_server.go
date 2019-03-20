package network

import (
	"net"
	"sync"
	"runtime"

	"comm/log"
	"comm/util"
)

type TCPServer struct {
	listener 	net.Listener
	opts 		*TCPOptions
	newAgent 	func(*TCPConn) Agent
	conns 		map[net.Conn]struct{}
	wg 			util.WaitGroupWrapper
	wgConns 	util.WaitGroupWrapper

	sync.Mutex
}

func NewTCPServer(options *TCPOptions, newAgentFun func(*TCPConn) Agent) *TCPServer {
	return &TCPServer {
		opts: options,
		newAgent: newAgentFun,
		conns: make(map[net.Conn]struct{}),
	}
}

func (s *TCPServer) getOpts() *TCPOptions {
	return s.opts
}

func (s *TCPServer) Run() {
	var err error
	s.listener, err = net.Listen("tcp", s.getOpts().TCPAddr)
	if nil != err {
		log.Logf(log.ERROR, "listen (%s) failed - %s", s.getOpts().TCPAddr, err)
		return
	}

	s.wg.Wrap(func() {
		for {
			clientConn, err := s.listener.Accept()
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
					log.Logf(log.WARN, "temporary Accept() failure - %s", err)
					runtime.Gosched()
					continue
				}
				break
			}

			s.Lock()
			if len(s.conns) >= s.getOpts().MaxConnNum {
				s.Unlock()
				log.Logf(log.WARN, "too many connections")
				clientConn.Close()
				continue
			}

			s.conns[clientConn] = struct{}{}
			s.Unlock()

			tcpConn := newTCPConn(clientConn, s.getOpts().PendingWriteNum, nil)
			agent := s.newAgent(tcpConn)
			s.wgConns.Wrap(func() {
				agent.Run()

				tcpConn.Close()
				s.Lock()
				delete(s.conns, clientConn)
				s.Unlock()

				agent.OnClose()
			})
		}	
	})
}

func (s *TCPServer) Close() {
	s.listener.Close()
	s.wg.Wait()
	
	s.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.Unlock()

	s.wgConns.Wait()
}