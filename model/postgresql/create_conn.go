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
	pool map[string]*sql.DB

	// 讀寫鎖，減少讀取路徑的競爭
	mu sync.RWMutex
)

func init() {
	pool = make(map[string]*sql.DB)
}

// GetConn : 依照資料庫名稱取得DB連線
func GetConn(dbName string) (*sql.DB, error) {
	// 快速路徑：讀鎖取得已存在的連線
	mu.RLock()
	conn, ok := pool[dbName]
	mu.RUnlock()

	if ok {
		// 在鎖外驗證連線狀態，避免阻塞其他 goroutine
		if err := conn.Ping(); err == nil {
			return conn, nil
		}

		// 連線失效，需要重建
		mu.Lock()
		// Double-check：其他 goroutine 可能已經重建
		if currentConn, exists := pool[dbName]; exists && currentConn == conn {
			if err := conn.Close(); err != nil {
				logger.Error("關閉資料庫連線失敗", "dbName", dbName, "error", err)
			}
			delete(pool, dbName)
		}
		mu.Unlock()
	}

	// 慢速路徑：寫鎖建立新連線
	mu.Lock()
	defer mu.Unlock()

	// Double-check：其他 goroutine 可能已經建立
	if conn, ok := pool[dbName]; ok {
		return conn, nil
	}

	conn, err := createConn(dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	pool[dbName] = conn
	return conn, nil
}

// createConn : 建立資料庫連線
func createConn(dbName string) (*sql.DB, error) {
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

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 配置連接池
	db.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 10))
	db.SetConnMaxLifetime(getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))
	db.SetConnMaxIdleTime(getEnvDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute))

	// 驗證連線
	if err := db.Ping(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("關閉資料庫連線失敗", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// CloseAll : 關閉所有數據庫連線
func CloseAll() error {
	mu.Lock()
	defer mu.Unlock()

	var errs []error
	for name, conn := range pool {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close database connection %s: %w", name, err))
		}
		delete(pool, name)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing database connections: %v", errs)
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
