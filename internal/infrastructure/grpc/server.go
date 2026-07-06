package grpcserver

import (
	"fmt"
	"github.com/alimarzban99/notification-service/internal/application/notification"
	"net"

	"github.com/alimarzban99/notification-service/config"
	"github.com/alimarzban99/notification-service/internal/infrastructure/logger"

	notificationpb "github.com/alimarzban99/notification-service/proto/notification"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server *grpc.Server
	logger logger.Logger
	config *config.Config
}

func New(cfg *config.Config, log logger.Logger, notificationService notification.Service) *Server {

	s := grpc.NewServer()

	server := &Server{
		server: s,
		logger: log,
		config: cfg,
	}

	handler := NewHandler(notificationService)

	notificationpb.RegisterNotificationServiceServer(s, handler)

	if cfg.GRPC.Reflection {
		reflection.Register(s)
	}

	return server
}

func (s *Server) Start() error {

	address := fmt.Sprintf(
		"%s:%d",
		s.config.Server.Host,
		s.config.Server.Port,
	)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	return s.server.Serve(listener)
}

func (s *Server) Stop() {

	s.logger.Info("Stopping gRPC server")

	s.server.GracefulStop()
}
