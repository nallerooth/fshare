package common

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

// Message is the base information shared by all messages
type Message struct {
	Command    MessageType
	Version    uint8
	DataLength uint64
	Name       [64]byte
}
