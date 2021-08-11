package service

import (
	"esnd/src/util"
	"net"
	"strconv"
	"sync"
)

var HID_INDEX int32 = 0
var Handlers = make(map[int32]*Handler)

var HandlersLock sync.Mutex

type NetService struct {
	Port int
	Lsn  net.Listener
}

func MakeNS(port int) (*NetService, error) {
	var ns NetService
	ns.Port = port
	lsn, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	ns.Lsn = lsn
	return &ns, nil
}

func (ns *NetService) Accept() {
	for {
		c, err := ns.Lsn.Accept()
		HandlersLock.Lock()
		util.DebugMsg("Listener", "New conn:"+c.RemoteAddr().Network())
		if err != nil {
			util.SaySub("Listener", "err:While accepting:"+err.Error())
			continue
		}

		h := MakeHandler(c)
		Handlers[h.HID] = h
		go h.Handle()

		HandlersLock.Unlock()

	}
}
