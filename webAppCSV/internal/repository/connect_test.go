package repository

/*
import (
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	loadEnv       = godotenv.Load
	processEnv    = envconfig.Process
	sqlOpen       = sql.Open
	pingDB        = func(db *sql.DB) error { return db.Ping() }
	retryFunction = ConnectWithRetry
)


	func TestConnectPostgre(t *testing.T) {
		loadEnv = func(filename ...string) error {
			return nil
		}
		processEnv = func(prefix string, spec interface{}) error {
			// Прописать допустимые поля в spec.(*Config).Db
			db := spec.(*Postgres)
			db.Host, db.Port, db.Username, db.Password, db.Dbname, db.Sslmode =
				"h", "p", "u", "pw", "db", "disable"
			return nil
		}

		// 2. Подменяем sql.Open и ping
		sqlOpen = func(driver, dsn string) (*sql.DB, error) {
			// возвращаем фиктивный *sql.DB — главное, чтобы pingDB вызвал наш мок
			return &sql.DB{}, nil
		}
		pingDB = func(db *sql.DB) error {
			return nil
		}

		// 3. Гарантируем, что retry не вызовется
		called := false
		retryFunction = func(cfg *Config) (*sql.DB, error) {
			called = true
			return nil, errors.New("should-not-be-called")
		}

		mockBehavior = func(){

		}

		testTables := []struct {
			loadEnv       func(filename ...string) error
			processEnv    func(prefix string, spec interface{}) error
			retryFunction func(cfg *Config) (*sql.DB, error)
		}{
			{},
		}
	}

func TestConnectPostgres_TableDriven(t *testing.T) {
	// сохраняем оригинальные глобальные функции, чтобы вернуть их в конце
	origLoadEnv := loadEnv
	origProcessEnv := processEnv
	origSqlOpen := sqlOpen
	origPingDB := pingDB
	origRetryFunction := retryFunction
	defer func() {
		loadEnv = origLoadEnv
		processEnv = origProcessEnv
		sqlOpen = origSqlOpen
		pingDB = origPingDB
		retryFunction = origRetryFunction
	}()

	// стандартная processEnv, заполняющая любые поля, но без ошибки
	goodProcessEnv := func(string, interface{}) error { return nil }

	// заглушка для успешного *sql.DB
	dummyDB := &sql.DB{}

	tests := []struct {
		name       string
		loadEnvErr error
		openErr    error
		pingErr    error
		retryDB    *sql.DB
		retryErr   error

		wantDB          *sql.DB
		wantErr         bool
		wantRetryCalled bool
	}{
		{
			name:            "Success on first try",
			loadEnvErr:      nil,
			openErr:         nil,
			pingErr:         nil,
			retryDB:         nil,
			retryErr:        nil,
			wantDB:          dummyDB,
			wantErr:         false,
			wantRetryCalled: false,
		},
		{
			name:            "Open fails → retry succeeds",
			loadEnvErr:      nil,
			openErr:         errors.New("open-fail"),
			pingErr:         nil,
			retryDB:         dummyDB,
			retryErr:        nil,
			wantDB:          dummyDB,
			wantErr:         false,
			wantRetryCalled: true,
		},
		{
			name:            "Ping fails → retry fails",
			loadEnvErr:      nil,
			openErr:         nil,
			pingErr:         errors.New("ping-fail"),
			retryDB:         nil,
			retryErr:        errors.New("retry-fail"),
			wantDB:          nil,
			wantErr:         true,
			wantRetryCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1) подменяем loadEnv и processEnv
			loadEnv = func(filenames ...string) error { return tc.loadEnvErr }
			processEnv = goodProcessEnv

			// 2) подменяем sqlOpen
			sqlOpen = func(driver, dsn string) (*sql.DB, error) {
				if tc.openErr != nil {
					return nil, tc.openErr
				}
				return dummyDB, nil
			}

			// 3) подменяем pingDB
			pingDB = func(db *sql.DB) error {
				return tc.pingErr
			}

			// 4) отслеживаем вызов retryFunction
			retryCalled := false
			retryFunction = func(cfg *Config) (*sql.DB, error) {
				retryCalled = true
				return tc.retryDB, tc.retryErr
			}

			// 5) запускаем тестируемый метод
			db, err := ConnectPostgres()

			// 6) Assertions
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantDB, db)
			require.Equal(t, tc.wantRetryCalled, retryCalled)
		})
	}
}*/
