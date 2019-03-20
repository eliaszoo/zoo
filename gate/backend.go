package gate

import (
	"sync/atomic"

	"comm/network"
)

type Backend struct {
	connNumMax 		int32
	connNum 		int32
	client 			*network.TCPClient
}

func (b *Backend) getConnNum() int32 {
	return atomic.LoadInt32(&b.connNum)
}

func (b *Backend) idle() bool {
	return getConnNum() <= b.connNumMax >> 2
}

func (b *Backend) idleRate() int {
	return getConnNum() * 100 / b.connNumMax
}

func (b *Backend) alive() bool {
	return true
}

func (b *Backend) addConnNum() int32 {
	return atomic.AddInt32(&b.ConnNum, 1)
}
