package main

import (
	"database/sql"
	"log"

	"github.com/ETNA10-7/ecom/cmd/api"
	"github.com/ETNA10-7/ecom/config"
	"github.com/ETNA10-7/ecom/db"
	_ "github.com/lib/pq"
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

	initStorage(db)
	//Create a server instance or call
	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: SuccessFUlly Connected")

}
