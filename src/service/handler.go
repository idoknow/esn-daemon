package service

import (
	"encoding/json"
	"errors"
	"esnd/src/cry"
	"esnd/src/db"
	"esnd/src/users"
	"esnd/src/util"
	"io/ioutil"
	"net"
	"strconv"
)

type Handler struct {
	HID    int32
	Conn   net.Conn
	Status int
	User   *users.User

	PrivateKey string
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
		case 3: //push
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
		case 4: //pull req
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
			err = SendNoti(*pack, h, pa.Crypto)
			if err != nil {
				util.DebugMsg("Handler-req", err.Error())
				WriteErr(err.Error(), h.Conn)
				continue
			}
			util.DebugMsg("Handler-req", "Response succ")
			continue
		case 6: //request priv list
			if h.Status != LOGINED {
				WriteErr("Not logined", h.Conn)
				continue
			}
			var resp PackReqPrivList
			resp.Priv = h.User.Priv
			rsakey := ""
			if pa.Crypto {
				rsakey = h.PrivateKey
			}
			WritePackage(h.Conn, resp, 6, rsakey)
			continue
		case 7:
			if h.Status != LOGINED {
				WriteErr("Not logined", h.Conn)
				continue
			}
			if !h.User.Can("account") {
				util.DebugMsg("Handler-account", "permission denied")
				WriteErr("You do not have account operation priv", h.Conn)
				continue
			}
			pack := &PackAccountOperation{}
			err := json.Unmarshal([]byte(pa.Json), &pack)
			if err != nil {
				h.CheckJSONSyntaxErr(err)
				continue
			}
			err = AccountOperation(*pack)
			if err != nil {
				util.DebugMsg("Handler-account", "err:"+err.Error())
				WriteErr(err.Error(), h.Conn)
				continue
			}
			util.DebugMsg("Handler-account", "Account operation succ")
			continue
		case 8: //request public key
			err := cry.Getkeys(strconv.Itoa(int(h.HID)))
			if err != nil {
				WriteErr("Cannot generate keypair", h.Conn)
				continue
			}

			publicKey, err := ioutil.ReadFile(".esnd/crypto/public/" + strconv.Itoa(int(h.HID)) + ".pem")
			if err != nil {
				WriteErr("Cannot read public key", h.Conn)
				continue
			}
			privateKey, err := ioutil.ReadFile(".esnd/crypto/private/" + strconv.Itoa(int(h.HID)) + ".pem")
			if err != nil {
				WriteErr("Cannot read private key", h.Conn)
				continue
			}
			h.PrivateKey = string(privateKey)

			var p0 PackRSAPublicKey
			p0.PublicKey = string(publicKey)

			WritePackage(h.Conn, p0, 9, "")
			util.DebugMsg("Handler-respPublicKey", "Send succ")
			continue
		default:
			WriteErr("Protocol Err", h.Conn)
			continue
		}
	}
}

func (h *Handler) Dispose() {
	HandlersLock.Lock()
	h.Conn.Close()
	delete(Handlers, h.HID)
	HandlersLock.Unlock()
}

func WriteErr(err string, c net.Conn) {
	var errp PackError
	errp.Err = err
	WritePackage(c, errp, 2, "")
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

func SendNoti(req PackRequest, h *Handler, crypto bool) error {
	rows, err := db.DB.Query("SELECT id,target,time,title,content,source FROM notis WHERE id>=" + strconv.Itoa(req.From) + " AND (target='" + h.User.Name + "' OR target='_global_') LIMIT 0," + strconv.Itoa(req.Limit))
	if err != nil {
		return err
	}
	for rows.Next() {
		var resp PackRespNotification
		err := rows.Scan(&resp.Id, &resp.Target, &resp.Time, &resp.Title, &resp.Content, &resp.Source)
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
func AccountOperation(req PackAccountOperation) error {
	if req.Name == "root" {
		return errors.New("cannot operate root account")
	}
	switch req.Oper {
	case "add":
		count := db.Count("SELECT count(*) FROM users WHERE name='" + req.Name + "'")
		if count >= 1 {
			return errors.New("account already exist")
		}
		_, err := db.DB.Exec("INSERT INTO users (name,mask,priv) VALUES ('" + req.Name + "','" + req.Pass + "','" + req.Priv + "')")
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
			WriteErr("Account info changed", h.Conn)
			h.Dispose()
		}
	}
}
