package service

import (
	"encoding/binary"
	"encoding/json"
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
	Json string
	Size int
	Code int
}

func ReadPackage(conn net.Conn) (*Package, error) {
	var p Package
	p.Code = ReadInt(conn)
	p.Size = ReadInt(conn)
	jsonBytes := make([]byte, p.Size)
	_, err := conn.Read(jsonBytes)
	if err != nil {
		return nil, err
	}
	p.Json = string(jsonBytes)
	return &p, nil
}

func WritePackage(conn net.Conn, obj interface{}, code int) (*Package, error) {
	jsonb, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	json := string(jsonb)
	var p Package
	p.Json = json
	p.Code = code
	dataByte := []byte(json)
	p.Size = len(dataByte)
	err = WriteInt(p.Code, conn)
	if err != nil {
		return nil, err
	}
	err = WriteInt(p.Size, conn)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(dataByte)
	if err != nil {
		return nil, err
	}
	return &p, err
}
