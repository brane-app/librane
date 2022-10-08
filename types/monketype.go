package types

type MonkeType interface {
	Map() map[string]interface{}
	JSON() ([]byte, error)
}
