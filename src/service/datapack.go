package service

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"esnd/src/cry"
	"esnd/src/util"
	"net"
)

func WriteInt(n int, conn net.Conn) error {
	x := int32(n)
	err := binary.Write(conn, binary.BigEndian, x)
	return err
}

//字节转换成整形
func ReadInt(conn net.Conn) int {
	// bytesBuffer := bytes.NewBuffer()

	var x int32
	binary.Read(conn, binary.BigEndian, &x)

	return int(x)
}

type Package struct {
	Json   string
	Size   int
	Code   int
	Crypto bool
}

func ReadPackage(conn net.Conn, privateKey string) (*Package, error) {
	var p Package
	p.Code = ReadInt(conn)
	p.Size = ReadInt(conn)
	p.Crypto = ReadInt(conn) == 1
	jsonBytes := make([]byte, p.Size)
	_, err := conn.Read(jsonBytes)
	if err != nil {
		return nil, err
	}
	//加密了
	if p.Crypto {
		util.DebugMsg("ReadPack", "decrypting:\n"+string(jsonBytes))
		if privateKey == "" {
			return nil, errors.New("no private key to decrypt json")
		}
		de, err := cry.RSA_decrypter(privateKey, jsonBytes)
		if err != nil {
			return nil, err
		}
		jsonBytes = de
	}

	p.Json = string(jsonBytes)
	return &p, nil
}

func WritePackage(conn net.Conn, obj interface{}, code int, rsakey string) (*Package, error) {
	jsonb, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	json := string(jsonb)
	var p Package
	p.Json = json
	p.Code = code
	p.Crypto = rsakey != ""
	dataByte := []byte(json)
	//加密
	if p.Crypto {
		en, err := cry.RSA_encrypter(rsakey, dataByte)
		if err != nil {
			return nil, err
		}
		dataByte = en
	}

	p.Size = len(dataByte)
	err = WriteInt(p.Code, conn)
	if err != nil {
		return nil, err
	}
	err = WriteInt(p.Size, conn)
	if err != nil {
		return nil, err
	}
	cryptoLabel := 0
	if p.Crypto {
		cryptoLabel = 1
	}
	err = WriteInt(cryptoLabel, conn)
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(dataByte)
	if err != nil {
		return nil, err
	}
	return &p, err
}
