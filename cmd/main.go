package main

import (
	"context"
	"fmt"
	"gRPC-server/internal/repository"
	"gRPC-server/internal/server"
	"gRPC-server/internal/service"
	"gRPC-server/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %v", err)
	}
}

func main() {
	logger := logger.GetLogger()
	db := repository.MongoInit(repository.NewMongoConnect(), logger)
	repo := repository.NewRepo(db, logger)
	service := service.NewService(repo, logger)
	sortService := server.NewSortServerService(service, logger)
	server := server.NewGrpcServer(sortService, logger)

	go server.ListenAndServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.GracefulShutDown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Shutdown error: %s", err))
	}
}
