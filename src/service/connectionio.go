package service

type IPackageIO interface {
	Write(p interface{}, code int) error
	Read() (Package, error)
	HandShake() bool
}

type SocketIO struct {
}

type WebSocketIO struct {
}
