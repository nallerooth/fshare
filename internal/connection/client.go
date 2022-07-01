package connection

import (
	"fmt"
	"net"

	"github.com/nallerooth/fshare/internal/config"
)

// NewClientConnection creates a new TCP connection to the host specified in config
// TODO: Add support for remote alias
func NewClientConnection(conf config.ClientConfig) (net.Conn, error) {
	connStr := fmt.Sprintf("%s:%d", conf.ServerURL, conf.ServerPort)
	return net.Dial("tcp", connStr)
}
