package db

import (
	"database/sql"
	"fmt"
	"os"
	"privat-unmei/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type CustomDB struct {
	db     *sql.DB
	logger logger.CustomLogger
}

func (cdb *CustomDB) Query(query string, args ...any) (*sql.Rows, error) {
	cdb.logger.Info(query)
	return cdb.db.Query(query, args...)
}

func (cdb *CustomDB) QueryRow(query string, args ...any) *sql.Row {
	cdb.logger.Info(query)
	return cdb.db.QueryRow(query, args...)
}

func (cdb *CustomDB) Exec(query string, args ...any) (sql.Result, error) {
	cdb.logger.Info(query)
	return cdb.db.Exec(query, args...)
}

func (cdb *CustomDB) Close() error {
	return cdb.db.Close()
}

func (cdb *CustomDB) Begin() (*sql.Tx, error) {
	return cdb.db.Begin()
}

func wrapDB(db *sql.DB, logger logger.CustomLogger) *CustomDB {
	return &CustomDB{db, logger}
}

func ConnectDB(logger logger.CustomLogger) (*CustomDB, error) {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")

	connString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}
	return wrapDB(conn, logger), nil
}
