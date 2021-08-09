package postgresql

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	//資料庫連線物件
	pool map[string]*sql.DB

	//同步鎖
	mu *sync.Mutex
)

func init() {
	pool = make(map[string]*sql.DB)
	mu = &sync.Mutex{}
}

// GetConn : 依照資料庫名稱取得DB連線
func GetConn(dbName string) *sql.DB {
	mu.Lock()
	defer mu.Unlock()

	if conn, ok := pool[dbName]; ok {
		if err := conn.Ping(); err == nil {
			return conn
		}
		conn.Close()
	}

	pool[dbName] = createConn(dbName)
	return pool[dbName]
}

// createConn : 建立資料庫連線
func createConn(dbName string) *sql.DB {
	host := os.Getenv("DATABASE_Host")
	port := os.Getenv("DATABASE_Port")
	user := os.Getenv("DATABASE_User")
	password := os.Getenv("DATABASE_Password")
	database := os.Getenv("DATABASE_Name")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port, user, password, database)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return db
}
