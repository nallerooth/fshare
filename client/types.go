package client

import "github.com/nallerooth/fshare/common"

type InternalClientMessage struct {
	Type          common.MessageType
	LocalFilename string
	CleanFilename string
	Sha256sum     string
}
