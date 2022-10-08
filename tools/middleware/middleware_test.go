package middleware

import (
	"github.com/brane-app/database-library"
	"github.com/brane-app/types-library"

	"os"
	"testing"
)

const (
	nick  = "foobar"
	email = "foo@bar.com"
)

var (
	user  types.User
	token string
)

func TestMain(main *testing.M) {
	database.Connect(os.Getenv("DATABASE_CONNECTION"))
	database.Create()
	user = types.NewUser(nick, "", email)

	var err error
	if err = database.WriteUser(user.Map()); err != nil {
		panic(err)
	}

	if token, _, err = database.CreateToken(user.ID); err != nil {
		panic(err)
	}

	var result int = main.Run()
	database.DeleteUser(user.ID)
	os.Exit(result)
}
