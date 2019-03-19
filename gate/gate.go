package gate

import (
	"net"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"encoding/json"

	"github.com/eliaszoo/zoo/comm/network"
	d "github.com/eliaszoo/zoo/comm/discovery"
	"github.com/eliaszoo/zoo/comm/util"
	"github.com/eliaszoo/zoo/comm/log"
)

type PeerInfo struct {
	Addr 		string  `json:"addr"`
	ConnNum 	int32	`json:"conn_num"`
	ConnNumMax 	int32	`json:"conn_num_max"`
}

type Peer struct {
	addr 		string
	conn 		net.Conn
	state 		int32
	PeerInfo 	
}

type Gate struct {
	tcpServer   network.TCPServer
	opts 		*Options
	apps 		[]Backend
	discovery	*d.Master

	util.WaitGroupWrapper
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
	g.Wrap(g.tcpServer.Run())
	g.Wrap(g.discovery.Watch())

	return nil
}

func (g *Gate) newAgent(conn net.conn) *agent {
	return &agent {
		conn,
		g,
		g.getBackend(),
	}
}

func (g *Gate) getBackend() *Backend {
	minIndex := 0
	minRate := 100
	for i, backend := range g.apps {
		if !backend.alive() {
			continue
		}

		idleRate := backend.idleRate()
		if idleRate < 25 {
			return &backend
		}

		if idleRate < minRate {
			minIndex = i
			minRate = idleRate
		}
	}

	return &g.apps[minIndex]
}

func (g *Gate) lookupLoop() {
	ticker := time.Tick(15 * time.Second)
	for {
		nodes := g.discovery.GetNodes()
		json.Unmarsh()
	}
}