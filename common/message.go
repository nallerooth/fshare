package common

import (
	"fmt"
	"os"
)

type MessageType int8

type messageFileInfo struct {
	length uint
	hash   string
	handle *os.File
}

const (
	Quit MessageType = iota
	Ping
	File
	DeleteFile
	DeleteHash
	List
	Text
	Search
)

type Message struct {
	Type   MessageType
	Length uint64
	Target [32]byte
}

func NewFileMessage(filename string) *Message {
	//info, err := fileInfo(filename)
	//if err != nil {
	//fmt.Fprintln(os.Stderr, err)
	//}

	return &Message{
		Type: File,
		//Length: int(info.length),
		//Checksum: info.hash,
	}
}

func fileInfo(filename string) (*messageFileInfo, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", stat)

	return &messageFileInfo{}, nil
}
