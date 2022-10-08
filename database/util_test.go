package database

import (
	"testing"
)

func Test_getSQLParams(test *testing.T) {
	var data map[string]interface{} = map[string]interface{}{
		"foo":  "bar",
		"baz":  "buz",
		"this": "that",
		"some": 1,
		"none": 0,
	}

	var keys []string
	var values []interface{}
	keys, values = getSQLParams(data)

	if len(keys) != len(values) {
		test.Errorf("Slice length mismatch! keys: %d, values: %d", len(keys), len(values))
	}

	var index int = 0
	for index != len(keys) {
		if values[index] != data[keys[index]] {
			test.Errorf(
				"key value mismatch on %s (at %d)! have: %v, want: %v",
				keys[index], index, values[index], data[keys[index]],
			)
		}

		index++
	}
}

func Test_makeSQLInsertable(test *testing.T) {
	var data map[string]interface{} = map[string]interface{}{
		"foo":  "bar",
		"baz":  "buz",
		"this": "that",
		"some": 1,
		"none": 0,
	}

	_, _ = makeSQLInsertable("table", data)
}
