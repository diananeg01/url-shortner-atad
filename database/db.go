package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/diananeg01/url-shortner-atad/model"
	"github.com/google/uuid"
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

func CreateSession(
	sessionID uuid.UUID,
	userID uuid.UUID,
	expiresAt time.Time,
) error {

	query := `
        INSERT INTO sessions (session_id, user_id, expires_at)
        VALUES ($1, $2, $3)
    `

	conn := GetDB()
	_, err := conn.Exec(query, sessionID, userID, expiresAt)
	return err
}

func GetSession(sessionID string) (*model.Session, error) {
	s := &model.Session{}

	conn := GetDB()
	err := conn.QueryRow(
		"SELECT session_id, user_id, expires_at FROM sessions WHERE session_id = $1",
		sessionID,
	).Scan(&s.SessionId, &s.UserId, &s.ExpiresAt)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func DeleteSession(sessionID string) error {
	conn := GetDB()
	_, err := conn.Exec(
		"DELETE FROM sessions WHERE session_id = $1",
		sessionID,
	)
	return err
}

func GetUserByID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	conn := GetDB()
	err := conn.QueryRow(`SELECT user_id, email FROM user_data WHERE user_id = $1`, id).
		Scan(&user.UserId, &user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return user, nil
}
