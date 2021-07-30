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

	numBytes, err := conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("Error writing to connection: %s", err)
	}
	fmt.Printf("Wrote %d bytes to connection\n", numBytes)

	resp := common.Message{}
	responseBuffer := make([]byte, 41)

	running := true
	for running {

		numBytes, err = conn.Read(responseBuffer)
		if err != nil {
			return err
		}
		fmt.Printf("Read %d bytes: %v\n", numBytes, responseBuffer)
		err = binary.Read(bytes.NewReader(responseBuffer), binary.BigEndian, &resp)
		if err != nil {
			return fmt.Errorf("Error in binary.Read: %s", err)
		}

		fmt.Printf("Received response: %+v\n", resp)

		switch resp.Type {
		case common.Text:
			if resp.Length > 0 {
				payload := make([]byte, resp.Length)
				numBytes, err = conn.Read(payload)
				if err != nil {
					return fmt.Errorf("read common.Text: %s", err)
				}
				fmt.Printf("Received %d bytes\n", numBytes)
				fmt.Println(string(payload))
			}

		case common.Quit:
			fmt.Println("Received end of transmission")
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
