package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	conf := testInit(confYaml)
	forceDBDrop(conf)
	err := Create(conf)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDrop(t *testing.T) {
	conf := testInit(confYaml)
	Create(conf)
	err := Drop(conf)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSeed(t *testing.T) {
	conf := testInit(confYaml)
	Create(conf)
	migrationHelper(conf)
	err := Seed(conf)

	if err != nil {
		t.Fatal(err)
	}

	forceDBDrop(conf)
}

func TestClean(t *testing.T) {
	conf := testInit(confYaml)
	Create(conf)
	migrationHelper(conf)
	Seed(conf)

	err := Clean(conf)

	if err != nil {
		t.Fatal(err)
	}

	forceDBDrop(conf)
}

func TestInsertConfSql(t *testing.T) {
	conf := testInit(confYaml)
	assert.Regexp(
		t,
		`INSERT INTO cats \(.*\) VALUES\(.*\);`,
		insertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"name",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"age",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"colour",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		cat1.Name,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		strconv.Itoa(cat1.Age),
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		cat1.Colour,
		"Should return an SQL query with the correct values based on the configuration given",
	)
}

func TestInsertConfSqlWithMultipleRows(t *testing.T) {
	confYaml2 := []byte(
		fmt.Sprintf(`
%v
  data:
    cats:
      - name: %v
        age: %v
        colour: %v
      - name: Q-ee
        age: 3
        colour: Grey`,
			sharedTestConfig(),
			cat1.Name,
			cat1.Age,
			cat1.Colour,
		))

	conf := testInit(confYaml2)

	assert.Regexp(
		t,
		`INSERT INTO cats \(.*\) VALUES\(.*\);`,
		insertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"name",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"age",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"colour",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		cat1.Name,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		strconv.Itoa(cat1.Age),
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		cat1.Colour,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"Q-ee",
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"3",
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"Grey",
		"Should return an SQL query with the correct values based on the configuration given",
	)
}

func TestInsertConfSqlWithMultipleTables(t *testing.T) {
	confYaml2 := []byte(
		fmt.Sprintf(`
%v
  data:
    cats:
      - name: %v
        age: %v
        colour: %v
    birds:
      - name: Maggie
        colour: Black and white`,
			sharedTestConfig(),
			cat1.Name,
			cat1.Age,
			cat1.Colour,
		))

	conf := testInit(confYaml2)

	assert.Regexp(
		t,
		`INSERT INTO cats \(.*\) VALUES\(.*\);`,
		insertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"name",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"age",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"colour",
		"Should return an SQL query with the correct columns based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		cat1.Name,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		strconv.Itoa(cat1.Age),
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		cat1.Colour,
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Regexp(
		t,
		`INSERT INTO birds \(.*\) VALUES\(.*\);`,
		insertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
		"formatted",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"Maggie",
		"Should return an SQL query with the correct values based on the configuration given",
	)

	assert.Contains(
		t,
		insertConfSql(conf),
		"Black and white",
		"Should return an SQL query with the correct values based on the configuration given",
	)
}

func TestInsertConfSqlWithNoTables(t *testing.T) {
	confYaml2 := []byte(
		fmt.Sprintf(`
%v
  data:`,
			sharedTestConfig(),
		))

	conf := testInit(confYaml2)

	assert.Equal(
		t,
		"",
		insertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
	)
}

func TestInsertConfSqlWithNoData(t *testing.T) {
	confYaml2 := []byte(
		fmt.Sprintf(`
default: &default
  database_user: pgcli_user
  host: localhost
  port: 5432
test:
  <<: *default
  database_name: pgcli_test`))

	conf := testInit(confYaml2)

	assert.Equal(
		t,
		"",
		insertConfSql(conf),
		"Should return the correct SQL query based on the configuration given",
	)
}
