package util

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

var EnableLog = false

var buffer bytes.Buffer

const BUFFER_MAX_SIZE = 16

func Say(msg string) {
	fmt.Print(msg)
	if EnableLog {
		buffer.WriteString(msg)
		if len(buffer.Bytes()) > BUFFER_MAX_SIZE {
			Flush()
		}
	}
}
func Sayln(msg string) {
	Say(msg + "\n")
}
func SaySub(sub string, msg string) {
	Sayln(getNowTimeStr() + "[" + sub + "]" + msg)
}

func Flush() {

	f, err := os.OpenFile("esnd.log", os.O_WRONLY, 0777)
	if err != nil {
		return
	} else {
		n, _ := f.Seek(0, 2)
		_, _ = f.WriteAt(buffer.Bytes(), n)
	}
	defer f.Close()
	defer buffer.Reset()
}

// Get timeStamp string as MM-DD,HH:mm:ss
func getNowTimeStr() string {
	t := time.Now()
	return strconv.Itoa(int(t.Month())) + "-" +
		strconv.Itoa(t.Day()) + "," +
		strconv.Itoa(t.Hour()) + ":" +
		strconv.Itoa(t.Minute()) + ":" +
		strconv.Itoa(t.Second())
}
func GetTimeStr(t time.Time) string {
	return strconv.Itoa(int(t.Month())) + "-" +
		strconv.Itoa(t.Day()) + "," +
		strconv.Itoa(t.Hour()) + ":" +
		strconv.Itoa(t.Minute()) + ":" +
		strconv.Itoa(t.Second())
}
