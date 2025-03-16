package server

import (
	"log"
	"net"
)

type TCPServer struct {
	listener net.Listener
	quit     chan struct{}
}

func NewTCPServer(network, addres string) (*TCPServer, error) {
	listener, err := net.Listen(network, addres)
	if err != nil {
		return nil, err
	}
	log.Printf("Server unnonced %s", addres)
	return &TCPServer{listener: listener, quit: make(chan struct{})}, nil
}

func (s *TCPServer) Run() {
	log.Printf("Server started on %s", s.listener.Addr())
	go s.acceptConnections()

	<-s.quit

	log.Println("Server stopped")
}

func (s *TCPServer) Stop() {
	close(s.quit)
}

func (s *TCPServer) acceptConnections() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go s.handleConnection(conn)

	}

}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection: %s", conn.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("handle: ", err)
			log.Println("Connection closed ", conn.RemoteAddr())

			return
		}
		log.Printf("Recived from %s:\n--->%s", conn.RemoteAddr(), buf[:n])
		conn.Write(buf[:n])

	}

}
