package tools

import (
	"strings"
)

func splitter(it rune) (ok bool) {
	ok = it == '/'
	return
}

func SplitPath(path string) (parts []string) {
	parts = strings.FieldsFunc(path, splitter)
	return
}
