package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"go.lsp.dev/uri"
)

func MigrateStep(conf map[string]interface{}, migrationPath string, step int) {
	databaseName := conf["database_name"].(string)

	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)

	m.Steps(step)

	fmt.Printf("%v database migrated to step %v\n", databaseName, step)
}

func MigrationVersion(conf map[string]interface{}, migrationPath string) {
	databaseName := conf["database_name"].(string)
	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	v, dirty := version(driver, migrationPath)

	fmt.Printf("%v database is currently migrated to version %v. ", databaseName, v)
	if dirty == true {
		fmt.Printf("The database is dirty.\n")
	} else {
		fmt.Printf("The database is clean.\n")
	}
}

func MigrateUp(conf map[string]interface{}, migrationPath string) {
	databaseName := conf["database_name"].(string)

	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)

	m.Up()

	fmt.Printf("%v database migrated\n", databaseName)
}

func MigrateDown(conf map[string]interface{}, migrationPath string) {
	databaseName := conf["database_name"].(string)

	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)

	m.Down()

	fmt.Printf("%v database schema rolled back\n", databaseName)
}

func version(driver database.Driver, migrationPath string) (uint, bool) {
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)

	version, dirty, err := m.Version()

	if err == migrate.ErrNilVersion {
		version = 0
	} else if err != nil {
		panic(err)
	}

	return version, dirty
}
