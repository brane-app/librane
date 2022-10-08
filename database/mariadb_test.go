package database

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"os"
	"testing"
)

var (
	CONNECTION string = os.Getenv("DATABASE_CONNECTION")
)

func mapMod(source map[string]interface{}, mods ...map[string]interface{}) (modified map[string]interface{}) {
	modified = mapCopy(source)

	var key string
	var value interface{}

	var mod map[string]interface{}
	for _, mod = range mods {
		for key, value = range mod {
			modified[key] = value
		}
	}

	return
}

func TestMain(main *testing.M) {
	Connect(CONNECTION)
	if database_handle == nil {
		panic("database_handle nil after being set!")
	}

	database_handle.Exec("SET FOREIGN_KEY_CHECKS=OFF")

	var err error
	var table string
	for _, table = range listStringReverse(tableOrdered) {
		if _, err = database_handle.Query("DROP TABLE IF EXISTS " + table); err != nil {
			database_handle.Exec("SET FOREIGN_KEY_CHECKS=ON")
			panic(err)
		}
	}

	database_handle.Exec("SET FOREIGN_KEY_CHECKS=ON")
	Create()

	var result int = main.Run()

	for table, _ = range tables {
		if err = EmptyTable(table); err != nil {
			panic(err)
		}
	}

	os.Exit(result)
}

func Test_Connect_malformedAddress(test *testing.T) {
	defer func(test *testing.T) {
		if recover() == nil {
			test.Errorf("recover recovered nil!")
		}
	}(test)

	var existing *sqlx.DB = database_handle
	defer func(existing *sqlx.DB) { database_handle = existing }(existing)

	Connect("foobar")
}

func Test_Connect_unreachableAddress(test *testing.T) {
	defer func(test *testing.T) {
		if recover() == nil {
			test.Errorf("recover recovered nil!")
		}
	}(test)

	var existing *sqlx.DB = database_handle
	defer func(existing *sqlx.DB) { database_handle = existing }(existing)

	Connect("foo:bar@tcp(nothing)/table")
}

func Test_create_badTable(test *testing.T) {
	var backup map[string]string = tables

	defer func(test *testing.T, backup map[string]string) {
		tables = backup

		if recover() == nil {
			test.Errorf("recover recovered nil!")
		}
	}(test, backup)

	tables = map[string]string{"foobar": "id CHAR CHAR CHAR CHAR CHAR GAS GAS GAS DEJA VU"}
	Create()
}

func Test_EmptyTable(test *testing.T) {
	var writable map[string]interface{} = mapCopy(writableUser)
	writable["id"] = uuid.New().String()

	var err error
	if err = WriteUser(writable); err != nil {
		test.Fatal(err)
	}

	if err = EmptyTable(USER_TABLE); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if _, exists, err = ReadSingleUser(writable["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("table %s was not emptied!", USER_TABLE)
	}
}
