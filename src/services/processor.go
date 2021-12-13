package services

import (
	"esnd/src/cry"
	"esnd/src/db"
	"esnd/src/users"
	"esnd/src/util"
	"io/ioutil"
	"strconv"
)

//Process packages from peer client

//PackTest
func (pack *PackTest) Process(h *Handler, p *NetPackage) error {
	util.DebugMsg("Handler", "PackTest:0:  int:"+strconv.Itoa(pack.Integer)+" msg:"+pack.Msg)
	WriteResult("Done", h, pack.Token)
	return nil
}

//PackLogin
func (pack *PackLogin) Process(h *Handler, p *NetPackage) error {
	if h.Status != ESTABLISHED {
		WriteErr("Cannot login", h, pack.Token)
		return nil
	}
	user, err := users.Auth(pack.User, pack.Pass)
	if err != nil {
		util.DebugMsg("Handler-auth", "err:"+err.Error())
		WriteErr("Login failed:"+err.Error(), h, pack.Token)
		h.Dispose()
		return nil
	}
	h.User = user
	h.Status = LOGINED
	util.DebugMsg("Handler-auth", "Login succ:"+pack.User)
	WriteResult("Done", h, pack.Token)
	return nil
}

//PackPush
func (pack *PackPush) Process(h *Handler, p *NetPackage) error {

	if h.Status != LOGINED {
		WriteErr("Not logined", h, pack.Token)
		return nil
	}
	if !h.User.Can("push") {
		WriteErr("You do not have push priv", h, pack.Token)
		return nil
	}

	/*default(real-time) id is -1
	,if this is not a rel-time notif,the id will be determined by database*/
	id := -1

	util.DebugMsg("Recv", "Real-time:"+strconv.FormatBool(pack.Realtime))
	var err error
	if !pack.Realtime {
		id, err = StoreNoti(*pack, h.User.Name)
		if err != nil {
			util.DebugMsg("Handler-pushNoti", err.Error())
			WriteErr(err.Error(), h, pack.Token)
			return nil
		}
	}

	util.DebugMsg("Handler-pushNoti", "Push succ.")
	WriteResult("Done", h, pack.Token)
	util.SaySub("Handler", "Push:source:"+h.User.Name+" target:"+pack.Target+" title:"+pack.Title)

	go PushToTarget(*pack, id, h.User.Name)
	return nil
}

//PackRequest
func (pack *PackRequest) Process(h *Handler, p *NetPackage) error {

	if h.Status != LOGINED {
		WriteErr("Not logined", h, pack.Token)
		return nil
	}
	if !h.User.Can("pull") {
		WriteErr("You do not have pull priv", h, pack.Token)
		return nil
	}
	WriteResult("Done", h, pack.Token)

	err := SendNoti(*pack, h, p.Crypto, pack.Token)
	if err != nil {
		util.DebugMsg("Handler-req", err.Error())
		WriteErr(err.Error(), h, pack.Token+"-1") //*
		return nil
	}
	util.DebugMsg("Handler-req", "Response succ")
	return nil
}

//PackReqPrivList
func (pack *PackReqPrivList) Process(h *Handler, p *NetPackage) error {
	if h.Status != LOGINED {
		WriteErr("Not logined", h, pack.Token)
		return nil
	}
	WriteResult("Done", h, pack.Token)
	var resp PackReqPrivList
	resp.Priv = h.User.Priv
	resp.Token = pack.Token + "-1" //*
	h.Adapter.Write(resp, 6)
	return nil
}

//PackAccountOperation
func (pack *PackAccountOperation) Process(h *Handler, p *NetPackage) error {

	if h.Status != LOGINED {
		WriteErr("Not logined", h, pack.Token)
		return nil
	}
	if !h.User.Can("account") {
		util.DebugMsg("Handler-account", "permission denied")
		WriteErr("You do not have account operation priv", h, pack.Token)
		return nil
	}
	err := AccountOperation(*pack)
	if err != nil {
		util.DebugMsg("Handler-account", "err:"+err.Error())
		WriteErr(err.Error(), h, pack.Token)
		return nil
	}
	util.DebugMsg("Handler-account", "Account operation succ")
	WriteResult("Done", h, pack.Token)

	util.SaySub("Handler", "Account Oper:subject:"+h.User.Name+" object.name:"+pack.Name+" oper:"+pack.Oper)

	return nil
}

//PackReqRSAKey
func (pack *PackReqRSAKey) Process(h *Handler, p *NetPackage) error {
	err := cry.Getkeys(strconv.Itoa(int(h.HID)))
	if err != nil {
		WriteErr("Cannot generate keypair", h, pack.Token)
		return nil
	}

	publicKey, err := ioutil.ReadFile(".esnd/crypto/public/" + strconv.Itoa(int(h.HID)) + ".pem")
	if err != nil {
		WriteErr("Cannot read public key", h, pack.Token)
		return nil
	}
	privateKey, err := ioutil.ReadFile(".esnd/crypto/private/" + strconv.Itoa(int(h.HID)) + ".pem")
	if err != nil {
		WriteErr("Cannot read private key", h, pack.Token)
		return nil
	}
	h.PrivateKey = string(privateKey)
	WriteResult("Done", h, pack.Token)

	var p0 PackRSAPublicKey
	p0.PublicKey = string(publicKey)
	p0.Token = pack.Token + "-1" //*

	h.Adapter.Write(p0, 9)
	util.DebugMsg("Handler-respPublicKey", "Send succ")
	return nil
}

//PackReqRecent
func (pack *PackReqRecent) Process(h *Handler, p *NetPackage) error {

	if h.Status != LOGINED {
		WriteErr("Not logined", h, pack.Token)
		return nil
	}
	if !h.User.Can("pull") {
		WriteErr("You do not have pull priv", h, pack.Token)
		return nil
	}
	WriteResult("Done", h, pack.Token)

	err := SendRecent(*pack, h, false, pack.Token)
	if err != nil {
		util.DebugMsg("Handler-req-recent", err.Error())
		WriteErr(err.Error(), h, pack.Token+"-1") //*
		return nil
	}
	util.DebugMsg("Handler-req-recent", "Resp succ")
	return nil
}

//PackCount
func (pack *PackCount) Process(h *Handler, p *NetPackage) error {

	if h.Status != LOGINED {
		WriteErr("Not logined", h, pack.Token)
		return nil
	}
	if !h.User.Can("pull") {
		WriteErr("You do not have pull priv", h, pack.Token)
		return nil
	}

	WriteResult("Done", h, pack.Token)

	//pack resp package

	to := 2147483647
	if pack.To != 0 {
		to = pack.To
	}

	var p0 PackRespCount
	p0.Amount = db.Count("SELECT count(*) FROM notis WHERE id>=" + strconv.Itoa(pack.From) +
		" AND id<=" + strconv.Itoa(to) + " AND (target like '%," + RawToEscape(h.User.Name) +
		",%' OR target like '%,_global_,%')")
	p0.Token = pack.Token + "-1" //*
	h.Adapter.Write(p0, 12)
	util.DebugMsg("Handler-count", "SELECT count(*) FROM notis WHERE id>="+strconv.Itoa(pack.From)+
		" AND id<="+strconv.Itoa(to)+" AND (target like '%,"+RawToEscape(h.User.Name)+
		",%' OR target like '%,_global_,%')")
	util.DebugMsg("Handler-count", "Result:"+strconv.Itoa(p0.Amount))

	return nil
}
