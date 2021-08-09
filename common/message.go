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
