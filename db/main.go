package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DBConnect(conf map[string]interface{}) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%v:@%v:%v/%v?sslmode=%v",
		conf["database_user"].(string),
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

func PgConnect(conf map[string]interface{}) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%v:@%v:%v/%v?sslmode=%v",
		"postgres",
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
	if conf["host"] != nil {
		return conf["host"].(string)
	}

	return "localhost"
}

func port(conf map[string]interface{}) int {
	if conf["port"] != nil {
		return conf["port"].(int)
	}

	return 5432
}

func sslMode(conf map[string]interface{}) string {
	if conf["ssl_mode"] != nil {
		mode := conf["string"].(string)

		switch mode {
		case "verify-full", "verify-ca", "require", "prefer", "allow", "disable":
			return mode
		default:
			return "disable"
		}
	}

	return "disable"
}
