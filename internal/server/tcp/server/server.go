package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type TCPServer struct {
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewTCPServer(network, address string) (*TCPServer, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen %s: %w", address, err)
	}
	log.Printf("Server announced %s", address)
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

	log.Println("Server is stopped")

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
				log.Printf("Failed to accept connection: %s", err)
				continue
			}
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	// Закрытие соединения

	defer func() {
		s.wg.Done()
		log.Println("Connection closed ", conn.RemoteAddr())
		conn.Close()
	}()

	log.Printf("New connection: %s", conn.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		select {
		case <-s.quit:
			return
		default:
			// чтобы I/O операция не держала цикл, ведь может прийти сигнал о закрытии - устанавливаем dedline
			conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("handle: ", err)
				// Если ошибка timeout, то продолжаем
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}
				return
			}
			log.Printf("Received from %s:\n--->%s", conn.RemoteAddr(), buf[:n])
			if _, err := conn.Write(buf[:n]); err != nil {
				log.Println("handle error write: ", err)
				return
			}

		}
	}
}
