package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Cannot find .env file skipping ")
	}

	dir := ""
	if len(os.Args) != 2 {
		fmt.Println("Migration directory not found, attempting default of ./migrations")
		dir = "./migrations"
	} else {
		dir = os.Args[1]
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal("Path not valid", err)
	}
	_, err = os.Stat(abs)
	if os.IsNotExist(err) {
		log.Fatal("path to migrations is invalid ", dir)
	}
	log.Println("Migrations directory found ", abs)
	databaseurl := os.Getenv("DATABASE_URL")
	log.Println("Database url ", databaseurl)
	m, err := migrate.New("file://"+abs, databaseurl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running migrations file://" + abs)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("Migrations complete")
}
