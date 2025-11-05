package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var dbConnection *sql.DB

func computeLocalDSN() string {
	user := "postgres"
	password := "postgres"
	host := "localhost"
	port := "5432"
	dbname := "postgres"

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
}
func Init() {
	var err error
	dbConnection, err = sql.Open("pgx", computeLocalDSN())
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	dbConnection.SetMaxOpenConns(10)
	dbConnection.SetMaxIdleConns(5)
	dbConnection.SetConnMaxIdleTime(0)

	// Test dbConnection
	if err := dbConnection.Ping(); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	fmt.Println("✅ Connected to Postgres")
}

func GetDB() *sql.DB {
	if dbConnection == nil {
		panic("db connection is nil")
	}
	return dbConnection
}

func Close() {
	if dbConnection != nil {
		if err := dbConnection.Close(); err != nil {
			log.Printf("failed to close DB dbConnection: %v", err)
		}
	}
}
