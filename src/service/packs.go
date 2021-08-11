package service

type PackTest struct { //0 both
	Integer int
	Msg     string
}

type PackLogin struct { //1 client
	User string
	Pass string
}

type PackError struct { //2 both
	Err string
}

type PackNotification struct { //3 client
	Target  string
	Time    string
	Title   string
	Content string
}

type PackRequest struct { //4 client
	From  int
	Limit int
}

type PackRespNotification struct { //5 server
	Id      int
	Target  string
	Time    string
	Title   string
	Content string
	Source  string
}

type PackReqPrivList struct { //6 both
	Priv string //not nil when server response
}

type PackAccountOperation struct { //7 client
	Oper string //add/remove
	Name string
	Pass string
	Priv string
	Kick bool
}

type PackReqRSAKey struct { //8 client
}

type PackRSAPublicKey struct { //9 server
	PublicKey string
}
