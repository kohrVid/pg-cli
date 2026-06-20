package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// DBConn is used to initialise connections to the postgres database.
// The following environment variable is required:
//   - DATABASE_PASSWORD
//
// The use of nested database fields is advised in the configuration file. The
// following keys are required for this format:
//   - DATABASE
//   - DATABASE.HOST
//   - DATABASE.PORT
//   - DATABASE.NAME
//   - DATABASE.USER
//
// The following keys are optional:
//   - DATABASE.SSL_MODE
func DBConn(conf map[string]interface{}) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		DBUser(conf),
		os.Getenv("DATABASE_PASSWORD"),
		host(conf),
		port(conf),
		DBName(conf),
		sslMode(conf),
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return dbConn
}

// PGUserDBConn is used to connect to a postgre instance as the default postgres user.
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

// DBName returns the name of the database in the configuration file.
func DBName(conf map[string]interface{}) string {
	if dbConf(conf)["name"] != nil {
		return dbConf(conf)["name"].(string)
	} else if conf["database_name"] != nil {
		return conf["database_name"].(string)
	}

	return ""
}

// DBUser returns the database user set in the configuration file.
func DBUser(conf map[string]interface{}) string {
	if dbConf(conf)["user"] != nil {
		return dbConf(conf)["user"].(string)
	} else if conf["database_user"] != nil {
		return conf["database_user"].(string)
	}

	return ""
}

func host(conf map[string]interface{}) string {
	if dbConf(conf)["host"] != nil {
		return dbConf(conf)["host"].(string)
	} else if conf["database_host"] != nil {
		return conf["database_host"].(string)
	}

	return "localhost"
}

func port(conf map[string]interface{}) int {
	if dbConf(conf)["port"] != nil {
		return dbConf(conf)["port"].(int)
	} else if conf["database_port"] != nil {
		return conf["database_port"].(int)
	}

	return 5432
}

func sslMode(conf map[string]interface{}) string {
	mode := "disable"

	if dbConf(conf)["ssl_mode"] != nil {
		mode = dbConf(conf)["ssl_mode"].(string)
	} else if conf["ssl_mode"] != nil {
		mode = conf["ssl_mode"].(string)
	}

	switch mode {
	case "verify-full", "verify-ca", "require", "prefer", "allow", "disable":
		return mode
	default:
		return "disable"
	}
}

func dbConf(conf map[string]interface{}) map[string]interface{} {
	if conf["database"] != nil {
		return conf["database"].(map[string]interface{})
	}

	return map[string]interface{}(nil)
}
