package services

import (
	"encoding/json"
	"esnd/src/users"
	"esnd/src/util"
	"strconv"
	"sync"
)

var HID_INDEX int32 = 0
var Handlers = make(map[int32]*Handler)

var HandlersLock sync.Mutex

type Handler struct {
	Adapter IAdapter
	HID     int32
	Status  int
	User    *users.User

	PrivateKey string
}

// status of handler
const (
	ESTABLISHED = iota
	LOGINED
	KILLED
)

func CreateHandler(adapter IAdapter) {
	HandlersLock.Lock()
	//handshook successfully,make handler
	handler := &Handler{
		adapter,
		HID_INDEX,
		ESTABLISHED,
		&users.User{},
		"",
	}
	HID_INDEX++
	Handlers[handler.HID] = handler
	go handler.Handle()

	HandlersLock.Unlock()
}

func (h *Handler) Handle() {
	for {
		np, err := h.Adapter.Read()
		if err != nil {
			util.DebugMsg("Handler", "err:While read socket pack:"+err.Error())
			h.Dispose()
			break
		}
		if np == nil {
			util.DebugMsg("Handler", "err:Pack from socket is nil.")
			h.Dispose()
			break
		}
		util.DebugMsg("recvJson", "#####"+np.JSON)

		//Parse json
		//Check code
		var pack IDataPackage

		switch np.Code {
		case 0: //test
			pack = &PackTest{}
		case 1: //login
			pack = &PackLogin{}
		case 3: //push
			pack = &PackPush{}
		case 4: //pull req
			pack = &PackRequest{}
		case 6: //request priv list
			pack = &PackReqPrivList{}
		case 7: //account
			pack = &PackAccountOperation{}
		case 8: //request public key
			pack = &PackReqRSAKey{}
		case 10: //req recent
			pack = &PackReqRecent{}
		case 11: //count notifications amount
			pack = &PackCount{}
		default:
			util.SaySub("Error", "Protocol Err:unrecognizable code:"+strconv.Itoa(np.Code))
			WriteErr("Protocol Err:unrecognizable code:"+strconv.Itoa(np.Code), h, "ErrorPackage")
			continue
		}
		//unmarshall
		err = json.Unmarshal([]byte(np.JSON), &pack)
		if err != nil {
			h.CheckJSONSyntaxErr(err)
			h.Dispose()
			continue
		}
		//process
		_ = pack.Process(h, np)
	}
}

// Dispose this handler and unregister
func (h *Handler) Dispose() {
	HandlersLock.Lock()

	h.Adapter.Dispose()
	delete(Handlers, h.HID)

	HandlersLock.Unlock()
}

func WriteResult(result string, h *Handler, token string) error {
	var resultp PackResult
	resultp.Error = ""
	resultp.Result = result
	resultp.Token = token
	_, err := h.Adapter.Write(resultp, 2)
	return err
}
func WriteErr(err string, h *Handler, token string) error {
	util.DebugMsg("ERR", "err:"+err)
	var errp PackResult
	errp.Error = err
	errp.Result = ""
	errp.Token = token
	_, error := h.Adapter.Write(errp, 2)
	return error
}

func (h *Handler) CheckJSONSyntaxErr(err error) {
	if err != nil {
		WriteErr("JSON syntax err", h, "ErrorPackage")
	}
}
