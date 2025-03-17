package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewTCPServer(network, addres string) (*TCPServer, error) {
	listener, err := net.Listen(network, addres)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %s %w", addres, err)
	}
	log.Printf("Server unnonced %s", addres)
	return &TCPServer{listener: listener, quit: make(chan struct{})}, nil
}

func (s *TCPServer) Run() {
	log.Printf("Server started on %s", s.listener.Addr())
	go s.acceptConnections()

	<-s.quit

	log.Println("The server began to stop")
}

// Остановка сервера, это прекращение прослушивания сокета,
// Завершение всех активных соединений
func (s *TCPServer) Stop() {

	// Информирование tcp соединений о завершении
	close(s.quit)

	// Прекращение прослушивания сокета
	s.listener.Close()

	// Ждем завершения tcp соединений
	s.wg.Wait()

	log.Println("Server is stoped")

}

// План такой: выходим, а handle горутины когда нибудь сами завершаться
// Проблема: Нет точного времени, гарантирующего завершени их, мы сами не разрываем соединение

func (s *TCPServer) acceptConnections() {
	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				log.Printf("Connection dont accepting %s", err)
				continue
			}
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	// Закрытие соединения

	defer s.wg.Done()
	defer log.Println("Connection closed ", conn.RemoteAddr())
	defer conn.Close()

	log.Printf("New connection: %s", conn.RemoteAddr())
	buf := make([]byte, 1024)
	isReaded := make(chan struct{})
	isClosed := make(chan struct{})

	go func() { isReaded <- struct{}{} }()
	for {
		select {
		case <-s.quit:
			return
		case <-isClosed:
			return
		default:
			select {
			case <-isReaded:
				go func() {
					n, err := conn.Read(buf)
					if err != nil {
						log.Println("handle: ", err)
						isClosed <- struct{}{}
						return
					}
					log.Printf("Recived from %s:\n--->%s", conn.RemoteAddr(), buf[:n])
					conn.Write(buf[:n])
					isReaded <- struct{}{}
				}()
			case <-time.After(time.Second * 1):
				continue
			}

		}
	}
}
