package helpers

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/kohrVid/pg-cli/db"
)

func Clean(conf map[string]interface{}) {
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

	db := db.DBConnect(conf)
	defer db.Close()

	_, err := db.Exec(cleanDB)
	if err != nil {
		log.Fatal(err)
	}
}

func Seed(conf map[string]interface{}) {
	seedDB := InsertConfSql(conf)

	if len(seedDB) < 1 {
		log.Fatal("No data to seed")
	}

	db := db.DBConnect(conf)
	defer db.Close()

	_, err := db.Exec(seedDB)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		`INSERT INTO candidate_time_slots
		  (candidate_id, time_slot_id) VALUES(1, 1);
		INSERT INTO interviewer_time_slots
		  (interviewer_id, time_slot_id) VALUES(1, 2)`,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func InsertConfSql(conf map[string]interface{}) string {
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

func SetSqlColumns(model *structs.Struct, params *structs.Struct) string {
	sql := "SET "
	for _, k := range model.Names() {
		if !params.Field(k).IsZero() {
			col := model.Field(k)

			sql += fmt.Sprintf(
				"%v = ",
				model.Field(k).Tag("json"),
			)

			switch col.Kind().String() {
			case "string":
				sql += fmt.Sprintf(
					"'%v', ",
					params.Field(k).Value(),
				)

			default:
				sql += fmt.Sprintf(
					"%v, ",
					params.Field(k).Value(),
				)
			}

			/*
			  This function exists because go-pg doesn't
			  support the structs library. If it did, the string
			  manipulation could be replaced with the line below.
			  Currently this line is required to ensure that the
			  controller returns the correct JSON object but isn't
			  used in the ORM command.
			*/
			col.Set(params.Field(k).Value())
		}
	}

	return sql[:len(sql)-2]
}
