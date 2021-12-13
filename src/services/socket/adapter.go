package socket

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"esnd/src/services"
	"esnd/src/util"
	"net"
	"strconv"
)

type SocketAdapter struct {
	Conn *net.Conn
}

//HandShake with income socket connection,return true if no err while handshaking
func (sa *SocketAdapter) HandShake() (bool, error) {
	identifier := ReadInt(*sa.Conn)
	if identifier != 119812525 {
		return false, errors.New("invalid identifier:" + strconv.Itoa(identifier))
	}
	err := WriteInt(util.ProtocolVersion, *sa.Conn)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (sa *SocketAdapter) Write(p interface{}, code int) (*services.NetPackage, error) {
	jsonb, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	json := string(jsonb)
	var np services.NetPackage
	np.JSON = json
	np.Code = code
	np.Crypto = false

	//encryption
	//unsupported

	np.Size = len(jsonb)
	err = WriteInt(np.Code, *sa.Conn) //code
	if err != nil {
		return nil, err
	}
	err = WriteInt(np.Size, *sa.Conn) //size
	if err != nil {
		return nil, err
	}
	err = WriteInt(0, *sa.Conn) //unencrypted
	if err != nil {
		return nil, err
	}
	_, err = (*sa.Conn).Write(jsonb) //json
	if err != nil {
		return nil, err
	}
	return &np, nil
}

func (sa *SocketAdapter) Read() (*services.NetPackage, error) {
	var np services.NetPackage
	np.Code = ReadInt(*sa.Conn)
	np.Size = ReadInt(*sa.Conn)
	if np.Size > 65536 {
		return nil, errors.New("package size(" + strconv.Itoa(np.Size) + ") is larger than limitation(65536)")
	}
	np.Crypto = ReadInt(*sa.Conn) == 1
	jsonBytes := make([]byte, np.Size)
	_, err := (*sa.Conn).Read(jsonBytes)
	if err != nil {
		return nil, err
	}
	if np.Crypto { //unsupported,return err when receive a encrypted package
		return nil, errors.New("encrypted package is now unsupported.")
	}
	np.JSON = string(jsonBytes)
	return &np, nil
}
func (sa *SocketAdapter) Dispose() {
	_ = (*sa.Conn).Close()
}

//Convert bytes to int
func ReadInt(conn net.Conn) int {
	// bytesBuffer := bytes.NewBuffer()
	var x int32
	binary.Read(conn, binary.BigEndian, &x)
	return int(x)
}

func WriteInt(n int, conn net.Conn) error {
	x := int32(n)
	err := binary.Write(conn, binary.BigEndian, x)
	return err
}
