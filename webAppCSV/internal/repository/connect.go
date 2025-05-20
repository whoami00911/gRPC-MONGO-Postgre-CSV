package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"webApp/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type Config struct {
	Db         Postgres
	MaxRetries int
	Delay      time.Duration
}

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
	Sslmode  string
}

func ConfigInicialize() *Config {
	return &Config{
		Db:         Postgres{},
		MaxRetries: 3,
		Delay:      3 * time.Second,
	}
}

func ConnectPostgres() (*sql.DB, error) {
	logger := logger.GetLogger()

	if err := godotenv.Load(".env"); err != nil {
		logger.Errorf("Can't load environment: %s", err)
		log.Fatalf("Can't load environment: %s", err)
	}

	cfg := ConfigInicialize()

	if err := envconfig.Process("db", &cfg.Db); err != nil {
		logger.Errorf("Can't read environment: %s", err)
		log.Fatalf("Can't read environment: %s", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s", cfg.Db.Host, cfg.Db.Port, cfg.Db.Username, cfg.Db.Password, cfg.Db.Dbname, cfg.Db.Sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Errorf("Can't open connection: %s", err)
		db, err = ConnectWithRetry(cfg)
		return db, err
	}

	if err = db.Ping(); err != nil {
		logger.Errorf("Ping pg database failed: %s", err)
		db.Close()
		db, err = ConnectWithRetry(cfg)
		return db, err
	}

	return db, err
}
func ConnectWithRetry(cfg *Config) (*sql.DB, error) {
	logger := logger.GetLogger()
	var err error
	var db *sql.DB

	for i := 0; i < cfg.MaxRetries; i++ {
		fmt.Printf("Попытка подключения к БД (%d/%d)...\n", i+1, cfg.MaxRetries)

		// Пытаемся подключиться к БД
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s", cfg.Db.Host, cfg.Db.Port, cfg.Db.Username, cfg.Db.Password, cfg.Db.Dbname, cfg.Db.Sslmode)
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			fmt.Printf("Ошибка при открытии соединения: %v", err)
			logger.Error(fmt.Sprintf("Retry connect to DB faild: %s", err))
		} else if err = db.Ping(); err == nil {
			// Успешное подключение
			fmt.Println("Успешное подключение к базе данных!")
			return db, err
		}

		// Если подключение не удалось, закрываем его
		if db != nil {
			_ = db.Close()
		}

		fmt.Printf("Не удалось подключиться, ожидаем %v перед повторной попыткой...\n", cfg.Delay)
		time.Sleep(cfg.Delay)
	}

	logger.Error(fmt.Sprintf("Retry connect to DB faild: %s", err))
	fmt.Printf("не удалось подключиться к базе данных после %d попыток: %s", cfg.MaxRetries, err)
	return nil, err
}
