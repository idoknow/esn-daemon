package service

type PackTest struct { //0
	Integer int
	Msg     string
}

type PackLogin struct { //1
	User string
	Pass string
}

type PackError struct { //2
	Err string
}

type PackNotification struct { //3
	Target  string
	Time    string
	Title   string
	Content string
}

type PackRequest struct { //4
	From  int
	Limit int
}

type PackRespNotification struct { //5
	Target  string
	Time    string
	Title   string
	Content string
	Source  string
}
