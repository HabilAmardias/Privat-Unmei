package dbcommand

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func RunSeeder() error {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	seedFiles, err := os.ReadDir("./db/seeds")
	if err != nil {
		return err
	}
	var sqlString string
	for _, entry := range seedFiles {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join("./db/seeds", entry.Name())
		fmt.Println(filePath, "string")
		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		sqlString += string(sqlBytes)
		sqlString += "\n"
	}
	fmt.Println(sqlString, "sql string")

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(sqlString)
	if err != nil {
		return err
	}

	return nil
}
