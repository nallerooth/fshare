package common

type MessageType int8

const (
	Quit MessageType = iota
	FileTransfer
	FileDelete
	FileSearch
	HashDelete
	HashSearch
	List
	Text
)

const MessageSize = 8 + 8 + 64

type Message struct {
	Command    MessageType
	Version    uint8
	DataLength uint64
	Name       [64]byte
}
