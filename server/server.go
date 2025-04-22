package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sms-tcpserver/handlers"
	"syscall"
)

// Server represents the TCP server
type Server struct {
	listener net.Listener
	port     string
}

// NewServer creates a new TCP server instance
func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

// Start initializes and starts the TCP server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	log.Printf("TCP Server started on %s", s.port)
	return nil
}

// Run starts accepting connections
func (s *Server) Run() {
	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start accepting connections in a goroutine
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				continue
			}

			// Handle each connection in a new goroutine
			go handlers.HandleConnection(conn)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down server...")
	s.Stop()
	handlers.StopHTTPSender()
}

// Stop closes the server listener
func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}
