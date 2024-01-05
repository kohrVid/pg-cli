package db

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

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
	err := viper.ReadConfig(bytes.NewBuffer(confYaml))
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return viper.GetStringMap("test")
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

// This method is used to test functions that print to stdout.
func CaptureOutput(f func()) string {
	out := make(chan string)
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, err := io.Copy(&buf, reader)
		if err != nil {
			panic(err)
		}

		out <- buf.String()
	}()
	wg.Wait()
	f()
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	return <-out
}
