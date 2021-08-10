package service

import (
	"encoding/json"
	"esnd/src/db"
	"esnd/src/users"
	"esnd/src/util"
	"net"
	"strconv"
)

type Handler struct {
	HID    int32
	Conn   net.Conn
	Status int
	User   *users.User
}

const (
	ESTABLISHED = iota
	LOGINED
	KILLED
)

func MakeHandler(conn net.Conn) *Handler {
	var h Handler
	h.HID = HID_INDEX
	HID_INDEX++
	h.Conn = conn
	h.Status = ESTABLISHED
	return &h
}

func (h *Handler) Handle() {
	for {
		pa, err := ReadPackage(h.Conn)
		if err != nil {
			util.DebugMsg("Handler", "err:While read pack:"+err.Error())
			h.Dispose()
			break
		}
		if pa == nil {
			util.DebugMsg("Handler", "err:Pack is nil")
			h.Dispose()
			break
		}
		switch pa.Code {
		case 0: //test
			pack := &PackTest{}
			err := json.Unmarshal([]byte(pa.Json), &pack)
			if err != nil {
				h.CheckJSONSyntaxErr(err)
				continue
			}
			util.DebugMsg("Handler", "PackTest:0:  int:"+strconv.Itoa(pack.Integer)+" msg:"+pack.Msg)
			continue
		case 1: //login
			if h.Status != ESTABLISHED {
				WriteErr("Cannot login", h.Conn)
				continue
			}
			pack := &PackLogin{}
			err := json.Unmarshal([]byte(pa.Json), &pack)
			if err != nil {
				h.CheckJSONSyntaxErr(err)
				continue
			}
			user, err := users.Auth(pack.User, pack.Pass)
			if err != nil {
				util.DebugMsg("Handler-auth", "err:"+err.Error())
				WriteErr("Login failed:"+err.Error(), h.Conn)
				h.Dispose()
				continue
			}
			h.User = user
			h.Status = LOGINED
			util.DebugMsg("Handler-auth", "Login succ:"+pack.User)
			continue
		case 3:
			if h.Status != LOGINED {
				WriteErr("Not logined", h.Conn)
				continue
			}
			if !h.User.Can("push") {
				WriteErr("You do not have push priv", h.Conn)
				continue
			}
			pack := &PackNotification{}
			err := json.Unmarshal([]byte(pa.Json), &pack)
			if err != nil {
				h.CheckJSONSyntaxErr(err)
				continue
			}
			err = StoreNoti(*pack, h.User.Name)
			if err != nil {
				util.DebugMsg("Handler-pushNoti", err.Error())
				WriteErr(err.Error(), h.Conn)
				continue
			}
			util.DebugMsg("Handler-pushNoti", "Push succ.")
			continue
		case 4:
			if h.Status != LOGINED {
				WriteErr("Not logined", h.Conn)
				continue
			}
			if !h.User.Can("pull") {
				WriteErr("You do not have pull priv", h.Conn)
				continue
			}
			pack := &PackRequest{}
			err := json.Unmarshal([]byte(pa.Json), &pack)
			if err != nil {
				h.CheckJSONSyntaxErr(err)
				continue
			}
			err = SendNoti(*pack, h)
			if err != nil {
				util.DebugMsg("Handler-req", err.Error())
				WriteErr(err.Error(), h.Conn)
				continue
			}
			util.DebugMsg("Handler-req", "Response succ")
		}
	}
}

func (h *Handler) Dispose() {
	HandlersLock.Lock()
	h.Conn.Close()
	delete(Handles, h.HID)
	HandlersLock.Unlock()
}

func WriteErr(err string, c net.Conn) {
	var errp PackError
	errp.Err = err
	WritePackage(c, errp, 2)
}

func (h *Handler) CheckJSONSyntaxErr(err error) {
	if err != nil {
		WriteErr("JSON syntax err", h.Conn)
	}
}

func StoreNoti(noti PackNotification, source string) error {
	_, err := db.DB.Exec("INSERT INTO notis (target,time,title,content,source) values ('" + noti.Target + "','" + noti.Time +
		"','" + noti.Title + "','" + noti.Content + "','" + source + "')")
	return err
}

func SendNoti(req PackRequest, h *Handler) error {
	rows, err := db.DB.Query("SELECT target,time,title,content,source FROM notis WHERE id>=" + strconv.Itoa(req.From) + " AND (target='" + h.User.Name + "' OR target='_global_') LIMIT 0," + strconv.Itoa(req.Limit))
	if err != nil {
		return err
	}
	for rows.Next() {
		var resp PackRespNotification
		err := rows.Scan(&resp.Target, &resp.Time, &resp.Title, &resp.Content, &resp.Source)
		util.DebugMsg("Handler-select", "select:"+resp.Target+" "+resp.Content)
		if err != nil {
			return err
		}
		WritePackage(h.Conn, resp, 5)
		util.DebugMsg("Handler-selectNoti", "Resp succ")
	}
	return nil
}
