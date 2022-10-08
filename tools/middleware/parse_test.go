package middleware

import (
	"net/http"
	"net/url"
	"testing"
)

type querySet struct {
	URL    string
	Size   int
	Before string
	Code   int
	OK     bool
}

func Test_PaginationParams(test *testing.T) {
	var request *http.Request = new(http.Request)
	var q_defaults map[string]interface{} = defaults()

	var set querySet
	var sets []querySet = []querySet{
		querySet{
			URL:    "http://imonke.io/?&size=10",
			Size:   10,
			Before: "",
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?before=foobar_baz&size=50",
			Size:   q_defaults["size"].(int),
			Before: "foobar_baz",
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?before=jol",
			Size:   q_defaults["size"].(int),
			Before: "jol",
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?&size=300",
			Size:   RANGE_SIZE_LIMIT,
			Before: q_defaults["before"].(string),
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?&size=lol",
			Size:   q_defaults["size"].(int),
			Before: q_defaults["before"].(string),
			Code:   400,
			OK:     false,
		},
		querySet{
			URL:    "http://imonke.io/?&size=-2",
			Size:   q_defaults["size"].(int),
			Before: q_defaults["before"].(string),
			Code:   400,
			OK:     false,
		},
	}

	var parsed map[string]interface{}
	var modified *http.Request
	var ok bool
	var code int
	var err error
	for _, set = range sets {
		if request.URL, err = url.Parse(set.URL); err != nil {
			test.Fatal(err)
		}

		if modified, ok, code, _, err = PaginationParams(request); err != nil {
			test.Fatal(err)
		}

		if code != set.Code {
			test.Errorf("got code %d", code)
		}

		if ok != set.OK {
			test.Errorf("got ok %t", ok)
		}

		if !ok {
			continue
		}

		parsed = modified.Context().Value("query").(map[string]interface{})

		if parsed["size"].(int) != set.Size {
			test.Errorf("size mismatch! have: %d, want: %d", parsed["size"], set.Size)
		}

		if parsed["before"].(string) != set.Before {
			test.Errorf("before mismatch! have: %s, want: %s", parsed["before"], set.Before)
		}
	}
}
