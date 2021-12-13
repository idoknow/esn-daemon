package services

type IAdapter interface {
	HandShake() (bool, error)                           //to check if this connection if an ESN connection
	Write(p interface{}, code int) (*NetPackage, error) //Write a NetPackage to peer
	Read() (*NetPackage, error)                         //Read a NetPackage from peer
	Dispose()                                           //Dispose connection
}
