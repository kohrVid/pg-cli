package db

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

type Cat struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Colour string `json:"colour"`
}

var cat1 Cat = Cat{
	Id:     1,
	Name:   "QT",
	Age:    4,
	Colour: "Tabby and white",
}

var confYaml []byte = []byte(
	fmt.Sprintf(`
%v
  data:
    cats:
      - name: %v
        age: %v
        colour: %v`,
		sharedTestConfig(),
		cat1.Name,
		cat1.Age,
		cat1.Colour,
	))

func testInit(confYaml []byte) map[string]interface{} {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(confYaml))
	return viper.Get("test").(map[string]interface{})
}

func sharedTestConfig() string {
	return `
default: &default
  database_user: pgcli_user
  host: localhost
  port: 5432
test:
  <<: *default
  database_name: pgcli_test`
}

func forceDBDrop(conf map[string]interface{}) {
	databaseName := conf["database_name"]
	dbConn := PGUserDBConn(conf)
	defer dbConn.Close()

	/* Terminate all other sessions in the database */
	_, err := dbConn.Exec(
		fmt.Sprintf(`
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = '%v'
  AND pid <> pg_backend_pid();`,
			databaseName),
	)
	if err != nil {
		panic(err)
	}

	dropDB := fmt.Sprintf("DROP DATABASE IF EXISTS %v;", databaseName)
	_, err = dbConn.Exec(dropDB)
	if err != nil {
		panic(err)
	}
}
