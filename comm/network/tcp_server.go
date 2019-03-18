package network

import (
	"net"
	"runtime"
	"strings"

	"comm/log"
)

func TCPServer(listener net.Listener, handler TCPHandler) {
	log.Logf(log.INFO, "TCP: listening on %s", listener.Addr())

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Logf(log.WARN, "temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			break
		}
		go handler.Handle(clientConn)
	}

	log.Logf(lg.INFO, "TCP: closing %s", listener.Addr())
}
