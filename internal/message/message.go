package message

import "errors"

// MessageType is an internal "enum" for message types
type MessageType int8

const (
	// Quit tells the server to close the connection
	Quit MessageType = iota
	// FileTransfer contains the header data for a new file transfer
	FileTransfer
	// FileDelete tells the server to remove a specific file
	FileDelete
	// FileSearch requests file info mathching a name pattern
	FileSearch
	// HashDelete tells the server to remove a file with a specific hash
	HashDelete
	// HashSearch requests file info matching a specific hash
	HashSearch
	// List requests a list of files available on the server
	List
	// Text is basically a bunch of bytes, usually a text response from the server
	Text
)

// MessageSize is the minimum size of header data for a message
const MessageSize = 8 + 8 + 64

var ErrInvalidMessageType = errors.New("invalid MessageType")

func NameFromType(t MessageType) (string, error) {
	switch t {
	case 0:
		return "Quit", nil
	case 1:
		return "FileTransfer", nil
	case 2:
		return "FileDelete", nil
	case 3:
		return "FileSearch", nil
	case 4:
		return "HashDelete", nil
	case 5:
		return "HashSearch", nil
	case 6:
		return "List", nil
	case 7:
		return "Text", nil
	default:
		return "", ErrInvalidMessageType
	}
}

// Message is the base information shared by all messages
type Message struct {
	Command    MessageType
	Version    uint8
	DataLength uint64
	Name       [64]byte
}

func NewText(payload []byte) Message {
	return Message{
		Command:    Text,
		DataLength: uint64(len(payload)),
	}
}

func NewList() Message {
	return Message{
		Command: List,
	}
}

func NewFileTransfer() Message {
	return Message{
		Command: FileTransfer,
	}
}

func QuitMessage() Message {
	return Message{
		Command: Quit,
	}
}
