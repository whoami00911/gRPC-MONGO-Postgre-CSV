package server

import (
	"context"
	"fmt"
	"gRPC-server/pkg/logger"
	"gRPC-server/pkg/parseCSV/grpcPb"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type grpcServer struct {
	grpcServer *grpc.Server
	logger     *logger.Logger
	csvService grpcPb.SortServiceServer
	addr       string
}

func NewGrpcServer(csvService grpcPb.SortServiceServer, logger *logger.Logger) *grpcServer {
	return &grpcServer{
		grpcServer: grpc.NewServer(),
		logger:     logger,
		csvService: csvService,
		addr:       viper.GetString("grpc.addr"),
	}
}

func (s *grpcServer) ListenAndServer() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Errorf("Can't lisnet grpc address: %s", err)
		log.Fatalf("Can't lisnet grpc address: %s", err)
	}

	grpcPb.RegisterSortServiceServer(s.grpcServer, s.csvService)
	fmt.Println("gRPC server has been started")

	if err := s.grpcServer.Serve(listener); err != nil {
		s.logger.Errorf("Can't server grpc server: %s", err)
		log.Fatalf("Can't server grpc server: %s", err)
	}
}

func (s *grpcServer) GracefulShutDown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		s.logger.Errorf("Graceful shutdown timed out, forcing immediate stop: %s", ctx.Err())
		fmt.Println("Graceful shutdown timed out, forcing immediate stop")
		return ctx.Err()
	case <-done:
		fmt.Println("Server shutdown gracefully")
		return nil
	}
}
