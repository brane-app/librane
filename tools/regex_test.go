package tools

import (
	"testing"
)

type regexCase struct {
	Condition string
	Pass      bool
}

func Test_NickRegex(test *testing.T) {
	var cases []regexCase = []regexCase{
		regexCase{"gastrodon", true},
		regexCase{"jim", true},
		regexCase{"jumbounclesteve", true},
		regexCase{"___lllIIIlll___", true},
		regexCase{"zero", true},
		regexCase{"...", true},
		regexCase{"foobar.io", true},
		regexCase{"", false},
		regexCase{"foo bar", false},
		regexCase{"|hello|", false},
		regexCase{"thisnameiswaytoolongwhoallowedthis", false},
		regexCase{"foo@bar.io", false},
	}

	var testcase regexCase
	for _, testcase = range cases {
		if NickRegex.MatchString(testcase.Condition) != testcase.Pass {
			test.Errorf("%s -> %t", testcase.Condition, testcase.Pass)
		}
	}
}

func Test_EmailRegex(test *testing.T) {
	var cases []regexCase = []regexCase{
		regexCase{"gastrodon", false},
		regexCase{"gastrodon@", false},
		regexCase{"@gastrodon", false},
		regexCase{"@", false},
		regexCase{"@@@", false},
		regexCase{"gastrodon@localhost", true},
	}

	var testcase regexCase
	for _, testcase = range cases {
		if EmailRegex.MatchString(testcase.Condition) != testcase.Pass {
			test.Errorf("%s -> %t", testcase.Condition, testcase.Pass)
		}
	}
}
