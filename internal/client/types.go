package client

import "github.com/nallerooth/fshare/internal/connection"

type InternalClientMessage struct {
	Type           connection.MessageType
	LocalFilename  string
	CleanFilename  string
	RemoteFilename string
	Sha256sum      string
}
