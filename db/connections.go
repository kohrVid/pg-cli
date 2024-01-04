package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func DBConn(conf map[string]interface{}) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		conf["database_user"].(string),
		os.Getenv("DATABASE_PASSWORD"),
		host(conf),
		port(conf),
		conf["database_name"].(string),
		sslMode(conf),
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return dbConn
}

func PGUserDBConn(conf map[string]interface{}) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		"postgres",
		os.Getenv("DATABASE_PASSWORD"),
		host(conf),
		port(conf),
		"postgres",
		sslMode(conf),
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return dbConn
}

func host(conf map[string]interface{}) string {
	if conf["database_host"] != nil {
		return conf["database_host"].(string)
	}

	return "localhost"
}

func port(conf map[string]interface{}) int {
	if conf["database_port"] != nil {
		return conf["database_port"].(int)
	}

	return 5432
}

func sslMode(conf map[string]interface{}) string {
	if conf["ssl_mode"] != nil {
		mode := conf["ssl_mode"].(string)

		switch mode {
		case "verify-full", "verify-ca", "require", "prefer", "allow", "disable":
			return mode
		default:
			return "disable"
		}
	}

	return "disable"
}
