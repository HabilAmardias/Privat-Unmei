package dbcommand

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RunSeeder() error {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

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

	db, err := sql.Open("postgres", connString)
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
