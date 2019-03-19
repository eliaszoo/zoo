package network

import (
	"net"
	"sync"
	""
	"runtime"

	"github.com/eliaszoo/zoo/comm/log"
	"github.com/eliaszoo/zoo/comm/util"
)

type TCPServer struct {
	listener 	net.Listener
	connNum 	int32
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

func (s *TCPServer) getConnNum() int {
	return s.connNum
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

			if s.getConnNum() >= s.getOpts().MaxConnNum {
				log.Logf(log.WARN, "too many connections")
				clientConn.Close()
				continue
			}

			tcpConn := newTCPConn(clientConn, 100, nil)
		}	
	})
}
