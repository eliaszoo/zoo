package network

import (
	"net"
	"sync"
	"runtime"

	"github.com/eliaszoo/zoo/comm/log"
	"github.com/eliaszoo/zoo/comm/util"
)

type TCPServer struct {
	listener 	net.Listener
	conns   	map[net.Conn]struct{}
	mutexConns  sync.Mutex
	opts 		*TCPOptions
	wg 			util.WaitGroupWrapper
}

func NewTCPServer(options *TCPOptions) *TCPServer {
	s := &TCPServer {
	}

	return s
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

			s.handle(clientConn)
		}
	})
}

func (s *TCPServer) handle(conn net.Conn) {
	s.mutexConns.Lock()
	defer s.mutexConns.Unlock()

	if len(s.conns) >= s.getOpts().MaxConnNum {
		conn.Close()
		log.Logf(log.WARN, "too many connections")
		return
	}

	
}