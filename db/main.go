package db

import (
	"crypto/tls"
	"database/sql"
	"os"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
)

func DBConnect(config map[string]interface{}) *pg.DB {
	conn := pg.Connect(&pg.Options{
		Database:  config["database_name"].(string),
		User:      config["database_user"].(string),
		TLSConfig: sslMode(),
	})

	return conn
}

func PgConnect() *sql.DB {
	conn, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable")

	if err != nil {
		panic(err)
	}

	return conn
}

func sslMode() *tls.Config {
	switch os.Getenv("SSL_MODE") {
	case "verify-ca", "verify-full":
		return &tls.Config{}
	case "allow", "prefer", "require":
		return &tls.Config{InsecureSkipVerify: true}
	default:
		return nil
	}
}
