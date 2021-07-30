package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/nallerooth/fshare/client"
	"github.com/nallerooth/fshare/common"
)

var clientConf = client.Config{
	Passphrase: "",
	RemoteURL:  "localhost",
	RemotePort: 32000,
}

func connectAndSend(msg *common.Message) error {
	connStr := fmt.Sprintf("%s:%d", clientConf.RemoteURL, clientConf.RemotePort)

	conn, err := net.Dial("tcp", connStr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, *msg)
	if err != nil {
		return fmt.Errorf("Error converting message to bytes: %s", err)
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Error writing to connection: %s", err)
	}

	resp := common.Message{}
	responseBuffer := make([]byte, 41)

	running := true
	for running {

		_, err = conn.Read(responseBuffer)
		if err != nil {
			return err
		}
		err = binary.Read(bytes.NewReader(responseBuffer), binary.BigEndian, &resp)
		if err != nil {
			return fmt.Errorf("Error in binary.Read: %s", err)
		}

		switch resp.Type {
		case common.Text:
			if resp.Length > 0 {
				payload := make([]byte, resp.Length)
				_, err = conn.Read(payload)
				if err != nil {
					return fmt.Errorf("read common.Text: %s", err)
				}
				fmt.Println(string(payload))
			}

		case common.Quit:
			fmt.Printf("\nReceived end of transmission\n")
			running = false
		}
	}

	return nil
}

func main() {
	msg := &common.Message{
		Type: common.List,
	}

	err := connectAndSend(msg)
	if err != nil {
		fmt.Println(err)
	}
}
