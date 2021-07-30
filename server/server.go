package server

import (
	"fmt"
	"net"
)

type Config struct {
	Passphrase string
	Port       uint
	Salt       string
	URL        string
	Workdir    string
}

type Server struct {
	config Config
	files  HashFileMap
}

func New(config Config) (*Server, error) {
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
