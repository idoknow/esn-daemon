package service

import (
	"encoding/json"
	"errors"
	"esnd/src/cry"
	"esnd/src/db"
	"esnd/src/users"
	"esnd/src/util"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type Handler struct {
	HID    int32
	Conn   net.Conn
	Status int
	User   *users.User

	PrivateKey string
}

//status of handler
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

	h.User = &users.User{}
	h.User.Name = ""
	h.User.Md5 = ""
	h.User.Priv = ""

	identifier := ReadInt(h.Conn)
	util.DebugMsg("READ", "ReadPack From:"+h.Conn.RemoteAddr().String())
	if identifier != 119812525 {
		h.Dispose()
		util.DebugMsg("Handler-flag", "invalid conn:"+strconv.Itoa(identifier)+" from:"+h.Conn.RemoteAddr().String())
		return
	}
	util.DebugMsg("Handler", "writeVersion")
	err := WriteInt(util.ProtocolVersion, h.Conn)
	if err != nil {
		h.Dispose()
		return
	}
	for {
		pa, err := ReadPackage(h.Conn, h.PrivateKey)
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
		util.DebugMsg("recvJson", "#####"+pa.Json)

		//Parse json
		//Check code
		var pack IDataPackage

		switch pa.Code {
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
			WriteErr("Protocol Err"+strconv.Itoa(pa.Code), h.Conn, "ErrorPackage")
			continue
		}
		//unmarshall
		err = json.Unmarshal([]byte(pa.Json), &pack)
		if err != nil {
			h.CheckJSONSyntaxErr(err)
			h.Dispose()
			continue
		}
		//process
		_ = pack.Process(h, pa) //TODO ignore error here because error msg will be sent to client in Process()
	}
}

func (h *Handler) Dispose() {
	HandlersLock.Lock()
	h.Conn.Close()
	delete(Handlers, h.HID)
	HandlersLock.Unlock()
}

func WriteResult(result string, c net.Conn, token string) {
	var resultp PackResult
	resultp.Error = ""
	resultp.Result = result
	resultp.Token = token
	WritePackage(c, resultp, 2, "")
}
func WriteErr(err string, c net.Conn, token string) {
	util.DebugMsg("ERR", "err:"+err)
	var errp PackResult
	errp.Error = err
	errp.Result = ""
	errp.Token = token
	WritePackage(c, errp, 2, "")
}

func (h *Handler) CheckJSONSyntaxErr(err error) {
	if err != nil {
		WriteErr("JSON syntax err", h.Conn, "ErrorPackage")
	}
}

func StoreNoti(noti PackPush, source string) (int, error) {
	_, err := db.DB.Exec("INSERT INTO notis (target,time,title,content,source,token) values ('," + RawToEscape(noti.Target) + ",','" + RawToEscape(noti.Time) +
		"','" + RawToEscape(noti.Title) + "','" + RawToEscape(noti.Content) + "','" + RawToEscape(source) + "','" + RawToEscape(noti.Token) + "')")
	if err != nil {
		return -1, err
	}
	id := db.Count("SELECT id FROM notis WHERE token='" + noti.Token + "'")
	return id, nil
}

func RawToEscape(raw string) string {
	return strings.ReplaceAll(raw, "'", "\\'")
}

func SendNoti(req PackRequest, h *Handler, crypto bool, token string) error {
	to := 2147483647

	if req.To != 0 {
		to = req.To
	}
	rows, err := db.DB.Query("SELECT id,target,time,title,content,source FROM notis WHERE id>=" + strconv.Itoa(req.From) +
		" AND id<=" + strconv.Itoa(to) + " AND (target like '%," + RawToEscape(h.User.Name) +
		",%' OR target like '%,_global_,%') LIMIT 0," + strconv.Itoa(req.Limit))
	if err != nil {
		return err
	}
	index := 1
	for rows.Next() {
		var resp PackRespNotification
		err := rows.Scan(&resp.Id, &resp.Target, &resp.Time, &resp.Title, &resp.Content, &resp.Source)
		resp.Token = token + "-" + strconv.Itoa(index) //*
		index++
		util.DebugMsg("Handler-select", "select:"+resp.Target+" "+resp.Content)
		if err != nil {
			return err
		}
		rsakey := ""
		if crypto {
			rsakey = h.PrivateKey
		}
		WritePackage(h.Conn, resp, 5, rsakey)
		util.DebugMsg("Handler-selectNoti", "Resp succ")
	}
	return nil
}

func SendRecent(req PackReqRecent, h *Handler, crypto bool, token string) error {

	count := db.Count("SELECT count(*) FROM notis ORDER BY id DESC")

	if count == 0 {
		return nil
	}
	if count > req.Limit {
		count = req.Limit
	}
	rows, err := db.DB.Query("SELECT id,target,time,title,content,source FROM notis ORDER BY id DESC LIMIT 0," + strconv.Itoa(req.Limit))
	if err != nil {
		return err
	}
	desc := make([]PackRespNotification, count)
	index := 0
	for rows.Next() {
		var resp PackRespNotification
		err := rows.Scan(&resp.Id, &resp.Target, &resp.Time, &resp.Title, &resp.Content, &resp.Source)
		resp.Token = token + "-" + strconv.Itoa(index+1)
		util.DebugMsg("Handler-select", "select:"+resp.Target+" "+resp.Content)
		if err != nil {
			return err
		}
		desc[index] = resp
		index++
	}
	for i := count - 1; i >= 0; i-- {
		rsakey := ""
		if crypto {
			rsakey = h.PrivateKey
		}
		WritePackage(h.Conn, desc[i], 5, rsakey)
		util.DebugMsg("Handler-selectNoti", "Resp succ")
	}
	return nil
}

func PushToTarget(pack PackPush, id int, source string) {
	var send PackRespNotification
	send.Id = id
	send.Target = "," + pack.Target + ","
	send.Time = pack.Time
	send.Title = pack.Title
	send.Content = pack.Content
	send.Source = source
	send.Token = pack.Token
	for _, h := range Handlers {
		if strings.Contains(send.Target, ",_global_,") ||
			strings.Contains(send.Target, ","+h.User.Name+",") {
			WritePackage(h.Conn, send, 5, "")
		}
	}
}

func AccountOperation(req PackAccountOperation) error {
	if req.Name == "root" {
		return errors.New("cannot operate root account")
	}
	reg, _ := regexp.Compile("^[0-9a-zA-Z_]{1,}$")
	if !reg.MatchString(req.Name) {
		return errors.New("invalid user name")
	}
	switch req.Oper {
	case "add":
		count := db.Count("SELECT count(*) FROM users WHERE name='" + req.Name + "'")
		if count >= 1 {
			return errors.New("account already exist")
		}
		_, err := db.DB.Exec("INSERT INTO users (name,mask,priv) VALUES ('" + req.Name + "','" + cry.MD5(req.Pass) + "','" + RawToEscape(req.Priv) + "')")
		return err
	case "remove":
		count := db.Count("SELECT count(*) FROM users WHERE name='" + req.Name + "'")
		if count < 1 {
			return errors.New("account not found")
		}

		row := db.DB.QueryRow("SELECT id FROM users WHERE name='" + req.Name + "'")
		var id int
		err := row.Scan(&id)
		if err != nil {
			return err
		}
		_, err = db.DB.Exec("DELETE FROM users WHERE id=" + strconv.Itoa(id))
		if err != nil {
			return err
		}
		Kick(req.Kick, req.Name)
		return nil
	default:
		return errors.New("no such operation")
	}
}
func Kick(kick bool, user string) {
	HandlersLock.Lock()
	toKick := make(map[int]*Handler)
	index := 0
	for _, h := range Handlers {
		if h.User.Name == user {
			toKick[index] = h
			index++
		}
	}
	HandlersLock.Unlock()

	for _, h := range toKick {
		if h != nil {
			WriteErr("Account info changed", h.Conn, "LogoutPackage")
			h.Dispose()
		}
	}
}
