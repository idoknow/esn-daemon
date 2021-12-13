package socket

import (
	"esnd/src/services"
	"esnd/src/util"
	"net"
	"strconv"
)

type SocketService struct {
	Port int
	Lsn  net.Listener
}

func MakeService(port int) (*SocketService, error) {
	var ss SocketService
	ss.Port = port
	lsn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	ss.Lsn = lsn
	return &ss, nil
}

/*
Accept socket connections.
Call makeSocketHandler() to process an new incomed connection
*/
func (ss *SocketService) Accept() {
	for {
		c, err := ss.Lsn.Accept()
		util.DebugMsg("Listener", "New socket conn:"+c.RemoteAddr().String())
		if err != nil {
			util.SaySub("Listener", "err:While socket accepting:"+err.Error())
			continue
		}
		go makeSocketHandler(c)
	}
}

//Check connection and do handshaking,if success:call handlerMgr to create handler
func makeSocketHandler(c net.Conn) {
	sa := &SocketAdapter{ //a socket connection
		&c,
	}
	hsResult, err := sa.HandShake()
	if err != nil {
		util.SaySub("Listener", "err:While socket handshaking:"+err.Error())
		return
	}
	if !hsResult {
		util.SaySub("Listener", "Failed to socket handshake.")
		return
	}
	services.CreateHandler(sa)

}
