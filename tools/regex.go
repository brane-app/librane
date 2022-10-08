package tools

import (
	"regexp"
)

const (
	NICK_PATTERN  = `[A-Za-z0-9.\-_]{3,16}`
	EMAIL_PATTERN = `[^@]+@[^@]+`
)

var (
	NickRegex  *regexp.Regexp = regexp.MustCompile("^" + NICK_PATTERN + "$")
	EmailRegex *regexp.Regexp = regexp.MustCompile("^" + EMAIL_PATTERN + "$")
)
