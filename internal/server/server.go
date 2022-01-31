package server

import (
	"fmt"
	"net"
)

// Config is a temporary struct for keeping track of settings
// TODO: Move to a proper settings package
type Config struct {
	Passphrase string
	Port       uint
	Salt       string
	URL        string
	Workdir    string
}

// Server is the main data storage and socket owner
type Server struct {
	config Config
	files  HashFileMap
}

// NewServer creates a new server with the given config, but does not start it.
func NewServer(config Config) (*Server, error) {
	s := &Server{
		config: config,
		files:  HashFileMap{},
	}

	if s.validateConfig() != nil {
		return nil, fmt.Errorf("invalid server configuration")
	}

	return s, nil
}

func (s *Server) validateConfig() error {
	// TODO:: Add validation here
	return nil
}

// Start sets up a TCP socket in listening mode and waits for connections
func (s *Server) Start() error {
	err := s.LoadWorkdir()
	if err != nil {
		return fmt.Errorf("error loading workdir: %s", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.config.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	// TODO: verbose check
	fmt.Println("Listening on TCP port", s.config.Port)
	fmt.Printf("Currently serving %d files found in workdir\n", len(s.AvailableFiles()))

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.HandleConnection(conn)
	}
}
