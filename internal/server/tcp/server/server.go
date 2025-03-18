package server

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"
)

type TCPServer struct {
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup
	logger   *slog.Logger
}

func NewTCPServer(network, address string, logger *slog.Logger) (*TCPServer, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen %s: %w", address, err)
	}
	logger.Info("Server announced",
		"network", network,
		"address", address,
	)
	return &TCPServer{listener: listener, quit: make(chan struct{}), logger: logger}, nil
}

func (s *TCPServer) Run() {
	go s.acceptConnections()

	<-s.quit

}

// Остановка сервера, это прекращение прослушивания сокета,
// Завершение всех активных соединений
func (s *TCPServer) Stop() {
	s.logger.Info("Server stopping",
		"address", s.listener.Addr(),
	)
	close(s.quit)

	// Прекращение прослушивания сокета
	s.listener.Close()

	// Ждем завершения tcp соединений
	s.wg.Wait()

	s.logger.Info("Server stopped", "address", s.listener.Addr())
}

// План такой: выходим, а handle горутины когда нибудь сами завершаться
// Проблема: Нет точного времени, гарантирующего завершени их, мы сами не разрываем соединение

func (s *TCPServer) acceptConnections() {
	s.logger.Info("Server started accepting connections",
		"address", s.listener.Addr(),
		"network", s.listener.Addr().Network(),
	)
	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				// Проверяем, является ли ошибка результатом закрытия сервера
				if errors.Is(err, net.ErrClosed) {
					s.logger.Info("Server is shutting down", "address", s.listener.Addr())
					return
				}
				// Логируем неожиданные ошибки
				s.logger.Error("Failed to accept connection", "error", err)
				continue
			}

			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	connLogger := s.logger.With(
		"remote_addr", conn.RemoteAddr(),
		"local_addr", conn.LocalAddr(),
		"network", conn.RemoteAddr().Network(),
	)

	defer func() {
		s.wg.Done()
		connLogger.Info("Connection closed")
		if err := conn.Close(); err != nil {
			connLogger.Error("Error closing connection", "error", err)
		}
	}()

	s.wg.Add(1)
	connLogger.Info("New connection established")

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
				if errors.Is(err, os.ErrDeadlineExceeded) {
					connLogger.Debug("Read timeout")
					continue
				} else {
					connLogger.Error("Read error",
						"error", err,
						"is_eof", errors.Is(err, io.EOF),
					)
					return
				}
			}

			connLogger.Info("Data received",
				"bytes_read", n,
				"data", string(buf[:n]),
			)

			if _, err := conn.Write(buf[:n]); err != nil {
				connLogger.Error("Write error",
					"error", err,
					"bytes_to_write", n,
				)
				return
			}
		}
	}
}
