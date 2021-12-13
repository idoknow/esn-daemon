package services

import (
	"errors"
	"esnd/src/cry"
	"esnd/src/db"
	"esnd/src/util"
	"regexp"
	"strconv"
	"strings"
)

//Store a received notification to DB
func StoreNoti(noti PackPush, source string) (int, error) {
	_, err := db.DB.Exec("INSERT INTO notis (target,time,title,content,source,token) values ('," + RawToEscape(noti.Target) + ",','" + RawToEscape(noti.Time) +
		"','" + RawToEscape(noti.Title) + "','" + RawToEscape(noti.Content) + "','" + RawToEscape(source) + "','" + RawToEscape(noti.Token) + "')")
	if err != nil {
		return -1, err
	}
	id := db.Count("SELECT id FROM notis WHERE token='" + noti.Token + "'")
	return id, nil
}

//Convert raw string to DB friendly string
func RawToEscape(raw string) string {
	return strings.ReplaceAll(raw, "'", "\\'")
}

//Send notification to specific handler
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
		h.Adapter.Write(resp, 5)
		util.DebugMsg("Handler-selectNoti", "Resp succ")
	}
	return nil
}

//Send recent notifications to specific handler
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
		h.Adapter.Write(desc[i], 5)
		util.DebugMsg("Handler-selectNoti", "Resp succ")
	}
	return nil
}

//Push incomed notification to target handlers
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
			h.Adapter.Write(send, 5)
		}
	}
}

//Execute account operation
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
			WriteErr("Account info changed", h, "LogoutPackage")
			h.Dispose()
		}
	}
}
