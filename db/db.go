package db

import (
	"database/sql"
	"fmt"
	"strings"

	//"fmt"
	"log"

	//"github.com/ETNA10-7/ecom/config"
	_ "github.com/lib/pq"
	//"golang.org/x/tools/go/cfg"
	//"golang.org/x/tools/go/cfg"
)

// FormatDSN generates a PostgreSQL DSN string from a Config struct
// func  FormatDSN(cfg config.Config) string {
// 	return fmt.Sprintf(
// 		"host=%s port=%d user=%s password=%s dbname=%s dbaddress=%s",
// 		cfg.PublicHost, cfg.Port, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBAddress,
// 	)
// }

type Connc struct {
	//Host     string
	Port     string
	User     string
	Password string
	Address  string
	Name     string
}

func (cfg *Connc) FormatDSN() (string, error) {
	// Split the address into host and port
	parts := strings.Split(cfg.Address, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid address format: expected host:port, got %s", cfg.Address)
	}
	//sslmode is needed part of the connection
	host := parts[0] // Extract host
	port := parts[1] // Extract port
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, cfg.User, cfg.Password, cfg.Name,
	), nil

}

//	func FormatDSN(cfg *Connc) string {
//		return fmt.Sprintf(
//			"host=%s port=%s user=%s password=%s dbname=%s dbaddress=%s",
//			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.Address,
//		)
//	}
func PostGresSqlStorage(cfg *Connc) (*sql.DB, error) {

	dsn, err := cfg.FormatDSN()
	if err != nil {
		return nil, fmt.Errorf("failed to format DSN: %w", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
