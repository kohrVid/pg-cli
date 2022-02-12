package db

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func Create(conf map[string]interface{}) error {
	databaseUser := conf["database_user"].(string)
	databaseName := conf["database_name"].(string)

	createRole := fmt.Sprintf("CREATE ROLE %v", databaseUser)

	alterRole := fmt.Sprintf(
		"ALTER ROLE %v WITH SUPERUSER LOGIN CREATEDB;",
		databaseUser,
	)

	createDB := fmt.Sprintf(
		"CREATE DATABASE %v WITH OWNER %v ENCODING 'UTF8';",
		databaseName,
		databaseUser,
	)

	conn := PgConnect(conf)

	_, err := conn.Exec(createRole)
	if err != nil {
		fmt.Println(err)
	}

	_, err = conn.Exec(alterRole)
	if err != nil {
		return err
	}

	_, err = conn.Exec(createDB)
	if err != nil {
		return err
	}

	fmt.Printf("%v database created\n", databaseName)
	return nil
}

func Drop(conf map[string]interface{}) error {
	databaseUser := conf["database_user"].(string)
	databaseName := conf["database_name"].(string)

	dropRole := fmt.Sprintf("DROP ROLE %v", databaseUser)
	dropDB := fmt.Sprintf("DROP DATABASE %v;", databaseName)

	conn := PgConnect(conf)
	defer conn.Close()

	_, err := conn.Exec(dropDB)
	if err != nil {
		return err
	}

	_, err = conn.Exec(dropRole)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v database deleted\n", databaseName)
	return nil
}

func Clean(conf map[string]interface{}) error {
	databaseUser := conf["database_user"].(string)

	truncateTables := fmt.Sprintf(`
CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
DECLARE
  statements CURSOR FOR
      SELECT
	tablename
      FROM pg_tables
      WHERE tableowner = username
	AND schemaname = 'public'
	AND tablename != 'gopg_migrations';

BEGIN
  FOR stmt IN statements LOOP
    EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) ||
    ' CASCADE; ALTER TABLE ' || quote_ident(stmt.tablename) ||
      ' ALTER COLUMN id RESTART WITH 1;';
  END LOOP;
END;
$$ LANGUAGE plpgsql;
		`)

	cleanDB := fmt.Sprintf(
		"%v SELECT truncate_tables('%v');",
		truncateTables,
		databaseUser,
	)

	conn := DBConnect(conf)
	defer conn.Close()

	_, err := conn.Exec(cleanDB)
	if err != nil {
		return err
	}

	return nil
}

func Seed(conf map[string]interface{}) error {
	seedDB := insertConfSql(conf)

	if len(seedDB) < 1 {
		log.Fatal("No data to seed")
	}

	conn := DBConnect(conf)
	defer conn.Close()

	_, err := conn.Exec(seedDB)
	if err != nil {
		return err
	}

	return nil
}

func insertConfSql(conf map[string]interface{}) string {
	var sql string
	d := conf["data"]

	if d == nil {
		return sql
	}

	data := make(map[string]interface{})

	if reflect.TypeOf(d).String() == "map[interface {}]interface {}" {
		for k, v := range d.(map[interface{}]interface{}) {
			data[k.(string)] = v
		}
	} else {
		data = d.(map[string]interface{})
	}

	for table, rows := range data {
		for _, row := range rows.([]interface{}) {

			cols := []string{}
			vals := []interface{}{}

			for col, val := range row.(map[interface{}]interface{}) {
				cols = append(cols, col.(string))
				vals = append(vals, val)
			}

			sql += fmt.Sprintf(
				"INSERT INTO %v (%v) VALUES(%v);\n",
				table,
				strings.Join(cols, ", "),
				mixedSqlSlice(vals),
			)
		}
	}

	return sql
}

func mixedSqlSlice(vals []interface{}) string {
	sql := ""
	for _, val := range vals {
		switch reflect.TypeOf(val).String() {
		case "string":
			sql += fmt.Sprintf(
				"'%v', ",
				val,
			)

		default:
			sql += fmt.Sprintf(
				"%v, ",
				val,
			)
		}
	}

	return sql[:len(sql)-2]
}
