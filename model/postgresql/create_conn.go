package postgresql

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"discordBot/pkg/logger"
)

var (
	// 資料庫連線物件
	db *sql.DB
	mu sync.Mutex
)

// GetConn : 取得資料庫連線
func GetConn() (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		return db, nil
	}

	conn, err := createConn()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	db = conn
	return db, nil
}

// createConn : 建立資料庫連線
func createConn() (*sql.DB, error) {
	host := os.Getenv("DATABASE_Host")
	port := os.Getenv("DATABASE_Port")
	user := os.Getenv("DATABASE_User")
	password := os.Getenv("DATABASE_Password")
	database := os.Getenv("DATABASE_Name")

	if host == "" || port == "" || user == "" || database == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, database,
	)

	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 配置連接池
	conn.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 25))
	conn.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 10))
	conn.SetConnMaxLifetime(getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))
	conn.SetConnMaxIdleTime(getEnvDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute))

	// 驗證連線
	if err := conn.Ping(); err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			logger.Error("關閉資料庫連線失敗", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}

// Close : 關閉資料庫連線
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		if err := db.Close(); err != nil {
			logger.Error("關閉資料庫連線失敗", "error", err)
			return fmt.Errorf("failed to close database connection: %w", err)
		}
		db = nil
	}
	return nil
}

// getEnvInt : 從環境變量獲取整數值
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvDuration : 從環境變量獲取時間值
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultVal
}
