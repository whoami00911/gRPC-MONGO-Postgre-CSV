package repository

import (
	"context"
	"fmt"
	"gRPC-server/pkg/logger"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	maxRetries    int
	retryDelation time.Duration
	db            Mongodb
}

type Mongodb struct {
	URI      string
	User     string
	Password string
	Database string
}

func ConfigInicialize() *Config {
	return &Config{
		maxRetries:    3,
		retryDelation: 1 * time.Second,
		db:            Mongodb{},
	}
}

func NewMongoConnect() (*mongo.Database, error) { //добавить error
	logger := logger.GetLogger()
	cfg := ConfigInicialize()

	if err := godotenv.Load(".env"); err != nil {
		logger.Error(fmt.Sprintf("godotenv can't load env: %s", err))
		log.Panic("godotenv can't load env")
	}

	if err := envconfig.Process("db", &cfg.db); err != nil {
		logger.Error(fmt.Sprintf("envconfig cant parse to struct: %s", err))
		log.Panic("envconfig cant parse to struct")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Создать объект настроек
	opts := options.Client()

	//Передать строку для подключения
	opts.ApplyURI(cfg.db.URI)

	//Установить аутентификационные данные
	opts.SetAuth(options.Credential{
		Username: cfg.db.User,
		Password: cfg.db.Password,
	})

	//Подключиться к серверу базы данных
	server, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error(fmt.Sprintf("Open connection to db server failed: %s", err))
		db, err := connectWithRetry(ctx, cfg, logger, opts)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	if err := server.Ping(context.Background(), nil); err != nil {
		logger.Error(fmt.Sprintf("ping db server failed: %s", err))
		db, err := connectWithRetry(ctx, cfg, logger, opts)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	//Выбрать базу данных для передачи и эксплутации следующими методами
	db := server.Database(cfg.db.Database)

	return db, nil
}

func connectWithRetry(ctx context.Context, cfg *Config, logger *logger.Logger, opts *options.ClientOptions) (*mongo.Database, error) {
	var err error
	var server *mongo.Client

	for i := 0; i < cfg.maxRetries; i++ {

		server, err = mongo.Connect(ctx, opts)
		if err != nil {
			fmt.Printf("Ошибка при открытии соединения к MongoDb: %v", err)
			logger.Error(fmt.Sprintf("Retry connect to DB server failed: %s", err))

			return nil, err

		} else if err = server.Ping(context.Background(), nil); err == nil {
			fmt.Println("Успешное подключение к базе данных MongoDb!")

			return server.Database(cfg.db.Database), nil
		}

		server.Disconnect(ctx)
		time.Sleep(cfg.retryDelation)
	}

	logger.Error(fmt.Sprintf("Retry connect to DB faild: %s", err))
	fmt.Printf("не удалось подключиться к базе данных после %d попыток: %s", cfg.maxRetries, err)

	return nil, err
}
