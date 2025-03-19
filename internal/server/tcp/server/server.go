package server

import (
	"context"
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
	return &TCPServer{listener: listener, logger: logger}, nil
}

func (s *TCPServer) Run(ctx context.Context) {
	go s.acceptConnections(ctx)

	<-ctx.Done()
	s.stop()
}

// Остановка сервера, это прекращение прослушивания сокета,
// Завершение всех активных соединений
func (s *TCPServer) stop() {
	s.logger.Info("Server stopping",
		"address", s.listener.Addr(),
	)
	// Прекращение прослушивания сокета
	s.listener.Close()

	// Ждем завершения tcp соединений
	s.wg.Wait()

	s.logger.Info("Server stopped", "address", s.listener.Addr())
}

func (s *TCPServer) acceptConnections(ctx context.Context) {
	s.logger.Info("Server started accepting connections",
		"address", s.listener.Addr(),
		"network", s.listener.Addr().Network(),
	)
	for {
		select {
		case <-ctx.Done():
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

			go s.handleConnection(ctx, conn)
		}
	}
}

func (s *TCPServer) handleConnection(ctx context.Context, conn net.Conn) {
	connLogger := s.logger.With(
		"remote_addr", conn.RemoteAddr(),
		"local_addr", conn.LocalAddr(),
		"network", conn.RemoteAddr().Network(),
	)

	s.wg.Add(1)
	defer func() {
		s.wg.Done()
		connLogger.Info("Connection closed")
		if err := conn.Close(); err != nil {
			connLogger.Error("Error closing connection", "error", err)
		}
	}()

	connLogger.Info("New connection established")

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			connLogger.Info("Connection closing to server shutdown")
			return
		default:
			// чтобы I/O операция не держала цикл, ведь может прийти сигнал о закрытии - устанавливаем dedline
			conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					connLogger.Debug("Read timeout")
					continue
				} else if errors.Is(err, io.EOF) {
					connLogger.Debug("Client closed connection: EOF")
					return
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
