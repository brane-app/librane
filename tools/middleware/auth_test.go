package middleware

import (
	"github.com/brane-app/database-library"
	"github.com/brane-app/types-library"
	"github.com/google/uuid"

	"context"
	"net/http"
	"testing"
)

func Test_MustAuth_blank(test *testing.T) {
	var request *http.Request = new(http.Request)

	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = MustAuth(request); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("blank request got through")
	}

	if code != 401 {
		test.Errorf("got code %d", code)
	}
}

func Test_MustAuth_invalid(test *testing.T) {
	var request *http.Request = new(http.Request)
	request.Header = make(http.Header)
	request.Header.Add("Authorization", "foobar")

	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = MustAuth(request); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("invalid auth request got through")
	}

	if code != 401 {
		test.Errorf("got code %d", code)
	}
}

func Test_MustAuth_badauth(test *testing.T) {
	var request *http.Request = new(http.Request)
	request.Header = make(http.Header)
	request.Header.Add("Authorization", "Bearer foobar")

	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = MustAuth(request); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("bad auth request got through")
	}

	if code != 401 {
		test.Errorf("got code %d", code)
	}
}

func Test_MustAuth(test *testing.T) {
	var request *http.Request = new(http.Request)
	request.Header = make(http.Header)
	request.Header.Add("Authorization", "Bearer "+token)

	var modified *http.Request
	var ok bool
	var err error
	if modified, ok, _, _, err = MustAuth(request); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Errorf("ok request did not get through")
	}

	var owner string = modified.Context().Value("requester").(string)
	if owner != user.ID {
		test.Errorf("modified is not owned by %s, but by %s", user.ID, owner)
	}
}

func Test_RejectBanned(test *testing.T) {
	var banned string = uuid.New().String()
	var ban types.Ban = types.NewBan(uuid.New().String(), banned, "", 0, true)
	database.WriteBan(ban.Map())

	var request *http.Request = new(http.Request).WithContext(context.WithValue(
		context.TODO(),
		"requester",
		banned,
	))

	var ok bool
	var code int
	var r_map map[string]interface{}
	var err error
	if _, ok, code, r_map, err = RejectBanned(request); err != nil {
		test.Fatal(err)
	}

	if code != 403 {
		test.Errorf("got code %d", code)
	}

	if ok {
		test.Errorf("banned not rejected")
	}

	var error string
	if error, ok = r_map["error"].(string); !ok || error != "banned" {
		test.Errorf("%#v", r_map)
	}
}
