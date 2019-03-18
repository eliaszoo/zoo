package gate

import (
	"net"
	"fmt"
	"sync/atomic"

	"comm/network"
	"comm/util"
)

type PeerInfo struct {

}

type Peer struct {
	addr 		string
	conn 		net.Conn
	state 		int32
	PeerInfo 	
	util.WaitGroupWrapper
}

type Client interface {

}

type Gate struct {
	listener 	net.Listener
	clients 	map[int64]Client
	peers 		atomic.Value
	opts 		*Options

}

func New(options *Options) *Gate {
	g := &Gate {
		clients: 	make(map[int64]Client),
		opts: options,
	}

	return g
}

func (g *Gate) getOpts() *Options {
	return g.opts
}

func (g *Gate) Main() error {
	var err error
	g.listener, err = net.Listen("tcp", g.getOpts().TCPAddr)
	if nil != err {
		return fmt.Errorf("listen (%s) failed - %s", g.getOpts().TCPAddr, err)
	}

	
	g.Wrap(network.TCPServer(g.listener, ))
	

	return nil
}

func (g *Gate)