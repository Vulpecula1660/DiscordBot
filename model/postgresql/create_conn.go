package postgresql

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"discordBot/pkg/config"
)

var (
	// 資料庫連線物件
	pool map[string]*sql.DB

	// 同步鎖
	mu *sync.Mutex
)

func init() {
	pool = make(map[string]*sql.DB)
	mu = &sync.Mutex{}
}

// GetConn : 依照資料庫名稱取得DB連線
func GetConn(dbName string) (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	if conn, ok := pool[dbName]; ok {
		if err := conn.Ping(); err == nil {
			return conn, nil
		}
		// 連線失效，關閉並移除
		conn.Close()
		delete(pool, dbName)
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
	db.SetMaxOpenConns(config.GetEnvInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(config.GetEnvInt("DB_MAX_IDLE_CONNS", 10))
	db.SetConnMaxLifetime(config.GetEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))
	db.SetConnMaxIdleTime(config.GetEnvDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute))

	// 驗證連線
	if err := db.Ping(); err != nil {
		db.Close()
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
