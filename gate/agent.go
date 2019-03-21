package gate

import (
	"comm/network"
)

type agent struct {
	id 		int64
	conn     network.TCPConn
	gate     *Gate
	backend  *Backend
}

func (a *agent) run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		
	}
}

func (a *agent) OnClose() {
}

func (a *agent) WriteMsg(msg interface{}) {
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}


