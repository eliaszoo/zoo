package network

import (
	"time"
)

type TCPOptions struct {
	TCPAddr				string
	MaxConnNum 			int
	PendingWriteNum 	int32
}

type TCPClientOptions struct {
	ServerAddr 		string
	ConnNum 		int
	ConnectInterval time.Duration
	PendingWriteNum int32
	AutoReconnect 	bool
}