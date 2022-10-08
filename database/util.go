package database

import (
	"github.com/brane-app/types-library"
	"github.com/jmoiron/sqlx"

	"strings"
)

func getSQLParams(it map[string]interface{}) (keys []string, values []interface{}) {
	var size int = len(it)
	keys, values = make([]string, size), make([]interface{}, size)

	var index int = 0
	var key string
	var value interface{}
	for key, value = range it {
		keys[index] = key
		values[index] = value
		index++
	}

	return
}

func makeSQLInsertable(table string, it map[string]interface{}) (statement string, values []interface{}) {
	var keys []string
	keys, values = getSQLParams(it)
	statement = "REPLACE INTO " + table + " (" + strings.Join(keys, ", ") + ") VALUES " + "(" + manyParamString("?", len(keys)) + ")"

	return
}

func manyParamString(param string, size int) (param_string string) {
	var param_slice []string = make([]string, size)
	for size != 0 {
		size--
		param_slice[size] = param
	}

	param_string = strings.Join(param_slice, ", ")
	return
}

func interfaceStrings(them ...string) (faces []interface{}) {
	faces = make([]interface{}, len(them))

	var index int
	for index, _ = range them {
		faces[index] = them[index]
	}

	return
}

func mapCopy(source map[string]interface{}) (copy map[string]interface{}) {
	copy = map[string]interface{}{}

	var key string
	var value interface{}
	for key, value = range source {
		copy[key] = value
	}

	return
}

func scanManyContent(rows *sqlx.Rows, count int) (content []types.Content, size int, err error) {
	var ids []string = make([]string, count)
	var scanned []types.Content = make([]types.Content, count)
	size = 0

	for rows.Next() {
		rows.StructScan(&scanned[size])
		ids[size] = scanned[size].ID
		size++
	}

	content = make([]types.Content, size)
	copy(content, scanned)

	var tags map[string][]string
	if tags, err = getManyTags(ids); err != nil {
		return
	}

	var index int
	for index, _ = range content {
		content[index].Tags = tags[content[index].ID]
	}

	return
}
