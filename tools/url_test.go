package tools

import (
	"testing"
)

type SplitPathTest struct {
	Path  string
	Parts []string
}

func Test_SplitPath_many(test *testing.T) {
	var cases []SplitPathTest = []SplitPathTest{
		{Path: "foo/bar/baz", Parts: []string{"foo", "bar", "baz"}},
		{Path: "/Spam/eggs/yoink", Parts: []string{"Spam", "eggs", "yoink"}},
		{Path: "Jinkies/Zoinks/Jeepers/", Parts: []string{"Jinkies", "Zoinks", "Jeepers"}},
		{Path: "/Yikes/Ruh-roh/idk/", Parts: []string{"Yikes", "Ruh-roh", "idk"}},
		{Path: "/", Parts: []string{}},
		{Path: "", Parts: []string{}},
	}

	var testcase SplitPathTest
	for _, testcase = range cases {
		var parts []string = SplitPath(testcase.Path)
		if len(parts) != len(testcase.Parts) {
			test.Errorf("length mismatch, %d != %d", len(parts), len(testcase.Parts))
		}

		var index int
		for index = range testcase.Parts {
			if testcase.Parts[index] != parts[index] {
				test.Errorf("mismatch at position %d, %s != %s", index, testcase.Parts[index], parts[index])
			}
		}
	}
}
