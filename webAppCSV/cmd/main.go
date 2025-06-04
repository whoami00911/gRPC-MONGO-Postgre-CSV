package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webApp/internal/repository"
	"webApp/internal/server"
	"webApp/internal/service"
	"webApp/internal/transport/handlers"
	"webApp/pkg/logger"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %v", err)
	}
}

func main() {
	logger := logger.GetLogger()

	db, err := repository.ConnectPostgres()
	if err != nil {
		log.Fatalf("Can't connect postgres: %s", err)
	}

	repo := repository.NewRepository(db, logger)
	service := service.NewService(repo)
	handlers := handlers.NewProductForHandlers(logger, service)

	srv := new(server.Server)

	go func() {
		if err := srv.ListenAndServer(handlers.InitRoutes()); err != nil {
			logger.Errorf("http-server can't start: %s", err)
			log.Fatalf("http-server can't start: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := srv.ShutDown(ctx); err != nil {
		logger.Error("Shutdown error: " + err.Error())
	}

	log.Println("timeout of 1 seconds.")
	<-ctx.Done()

	log.Println("Server exiting")
}
