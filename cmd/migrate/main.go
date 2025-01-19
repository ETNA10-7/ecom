package main

//go get -u github.com/golang-migrate/migrate/v4

import (
	"log"
	"os"

	"github.com/ETNA10-7/ecom/config"
	"github.com/ETNA10-7/ecom/db"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"

	//"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//"github.com/lib/pq"
)

func main() {
	cfg := db.Connc{
		//Host:     config.Envs.PublicHost,

		Port:     config.Envs.Port,
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		Address:  config.Envs.DBAddress,
		Name:     config.Envs.DBName,
	}
	db, err := db.PostGresSqlStorage(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new database driver for migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)

	// v, d, _ := m.Version()
	// log.Printf("Version: %d, dirty: %v", v, d)
	//Entry Point for Application
	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {

		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	}
}
