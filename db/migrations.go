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

func MigrationVersion(conf map[string]interface{}, migrationPath string) error {
	databaseName := conf["database_name"].(string)

	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}

	v, dirty, err := version(driver, migrationPath)
	if err != nil {
		return err
	}

	fmt.Printf("%v database is currently migrated to version %v. ", databaseName, v)
	if dirty == true {
		fmt.Printf("The database is dirty.\n")
	} else {
		fmt.Printf("The database is clean.\n")
	}

	return nil
}

func MigrateStep(conf map[string]interface{}, migrationPath string, step int) error {
	databaseName := conf["database_name"].(string)

	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)
	if err != nil {
		return err
	}

	m.Steps(step)
	fmt.Printf("%v database migrated to step %v\n", databaseName, step)

	return nil
}

func MigrateUp(conf map[string]interface{}, migrationPath string) error {
	databaseName := conf["database_name"].(string)

	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)
	if err != nil {
		return err
	}

	m.Up()
	fmt.Printf("%v database migrated\n", databaseName)

	return nil
}

func MigrateDown(conf map[string]interface{}, migrationPath string) error {
	databaseName := conf["database_name"].(string)

	dbConn := DBConn(conf)
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)
	if err != nil {
		return err
	}

	m.Down()
	fmt.Printf("%v database schema rolled back\n", databaseName)

	return nil
}

func version(driver database.Driver, migrationPath string) (uint, bool, error) {
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("%v", uri.File(migrationPath)),
		"postgres", driver)

	version, dirty, err := m.Version()

	if err == migrate.ErrNilVersion {
		version = 0
	}

	return version, dirty, err
}
