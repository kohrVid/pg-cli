package db

import (
	"testing"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/assert"
)

var conf map[string]interface{} = testInit(confYaml)

func TestMigrationVersion(t *testing.T) {
	Create(conf)
	MigrateUp(conf, "../example/migrations")

	actual := CaptureOutput(func() {
		MigrationVersion(conf, "../example/migrations")
	})

	assert.Regexp(
		t,
		`pgcli_test database is currently migrated to version 2. The database is clean.`,
		actual,
		"Should return the correct version number",
		"formatted",
	)

	forceDBDrop(conf)
}

func TestMigrationVersionWithUnknownPath(t *testing.T) {
	Create(conf)

	assert.Error(
		t,
		MigrationVersion(conf, "./unknown"),
		`no such file or directory`,
		"Should be called with a valid migration path",
		"formatted",
	)

	forceDBDrop(conf)
}

func TestMigrationVersionWithoutMigrations(t *testing.T) {
	Create(conf)

	actual := CaptureOutput(func() {
		MigrationVersion(conf, "../example/migrations")
	})

	assert.Regexp(
		t,
		`pgcli_test database is currently migrated to version 0. The database is clean.`,
		actual,
		"Should return the correct version number",
		"formatted",
	)

	forceDBDrop(conf)
}

func TestMigrateStep(t *testing.T) {
	Create(conf)
	err := MigrateStep(conf, "../example/migrations", 1)

	if err != nil {
		t.Fatal(err)
	}

	checkDBVersion(t, conf, "../example/migrations", 1)

	forceDBDrop(conf)
}

func TestMigrateNegativeStep(t *testing.T) {
	Create(conf)
	err := MigrateStep(conf, "../example/migrations", 2)
	if err != nil {
		t.Fatal(err)
	}

	err = MigrateStep(conf, "../example/migrations", -1)
	if err != nil {
		t.Fatal(err)
	}

	checkDBVersion(t, conf, "../example/migrations", 1)

	forceDBDrop(conf)
}

func TestMigrateUp(t *testing.T) {
	Create(conf)
	err := MigrateUp(conf, "../example/migrations")

	if err != nil {
		t.Fatal(err)
	}

	checkDBVersion(t, conf, "../example/migrations", 2)

	forceDBDrop(conf)
}

func TestMigrateDown(t *testing.T) {
	Create(conf)
	MigrateUp(conf, "../example/migrations")
	err := MigrateDown(conf, "../example/migrations")

	if err != nil {
		t.Fatal(err)
	}

	checkDBVersion(t, conf, "../example/migrations", 0)

	forceDBDrop(conf)
}

func checkDBVersion(t *testing.T, conf map[string]interface{}, path string, expectation int) {
	driver, _ := postgres.WithInstance(DBConn(conf), &postgres.Config{})
	v, _, _ := version(driver, "../example/migrations")

	assert.Equal(
		t,
		uint(expectation),
		v,
		"Should return the correct version number",
		"formatted",
	)
}
