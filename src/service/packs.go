package service

type PackTest struct { //0 both
	Integer int
	Msg     string
	Token   string
}

type PackLogin struct { //1 client
	User  string
	Pass  string
	Token string
}

type PackResult struct { //2 both
	Result string
	Error  string
	Token  string
}

type PackPush struct { //3 client
	Target  string
	Time    string
	Title   string
	Content string
	Token   string
}

type PackRequest struct { //4 client
	From  int
	Limit int
	Token string
}

type PackRespNotification struct { //5 server
	Id      int
	Target  string
	Time    string
	Title   string
	Content string
	Source  string
	Token   string
}

type PackReqPrivList struct { //6 both
	Priv  string //not nil when server response
	Token string
}

type PackAccountOperation struct { //7 client
	Oper  string //add/remove
	Name  string
	Pass  string
	Priv  string
	Kick  bool
	Token string
}

type PackReqRSAKey struct { //8 client
	Token string
}

type PackRSAPublicKey struct { //9 server
	PublicKey string
	Token     string
}
