package types

import (
	"github.com/google/uuid"

	"testing"
)

func acceptMonkeType(it MonkeType) {
	return
}

func Test_MonkeType(test *testing.T) {
	acceptMonkeType(Content{})
	acceptMonkeType(User{})
}

func Test_Ban(test *testing.T) {
	var banner string = uuid.New().String()
	var banned string = uuid.New().String()
	var reason string = "they smell like cheese"
	var duration int64 = 24 * 60 * 60
	var ban Ban = NewBan(banner, banned, reason, duration, false)

	if ban.Banner != banner {
		test.Errorf("ban properties not being set for banner! have: %s, want: %s", ban.Banner, banner)
	}

	if ban.Map()["banner"].(string) != banner {
		test.Errorf("bad ban map! %#v", ban.Map())
	}

	var err error
	if _, err = ban.JSON(); err != nil {
		test.Fatal(err)
	}
}

func Test_Report(test *testing.T) {
	var reporter string = uuid.New().String()
	var reported string = uuid.New().String()
	var reason string = "they smell like cheese"
	var report Report = NewReport(reporter, reported, "user", reason)

	if report.Reporter != reporter {
		test.Errorf("report properties not being set for reporter! have: %s, want: %s", report.Reporter, reporter)
	}

	if report.Map()["reporter"].(string) != reporter {
		test.Errorf("bad report map! %#v", report.Map())
	}

	var err error
	if _, err = report.JSON(); err != nil {
		test.Fatal(err)
	}
}

func Test_User(test *testing.T) {
	var nick string = "imonke"
	var user User = NewUser(nick, "", "")

	if user.Nick != nick {
		test.Errorf("user properties not being set for nick! have: %s, want: %s", user.Nick, nick)
	}

	if user.Map()["nick"].(string) != nick {
		test.Errorf("bad user map! %#v", user.Map())
	}

	var err error
	if _, err = user.JSON(); err != nil {
		test.Fatal(err)
	}

	var public map[string]interface{} = user.PublicMap()
	var key string
	for _, key = range userPrivateFields {
		if public[key] != nil {
			test.Errorf("private map has key %s\n%#v", key, public)
		}
	}

	var map_source User
	map_source.FromMap(user.Map())

	if map_source.Nick != nick {
		test.Errorf("nick %s not sourced from map %#v", nick, user.Map())
	}
}

func Test_Content(test *testing.T) {
	var author string = uuid.New().String()
	var content Content = NewContent("", author, "", nil, false, false)

	if content.Author != author {
		test.Errorf("content properties not being set for author! have: %s, want: %s", content.Author, author)
	}

	if content.Tags == nil {
		test.Errorf("Tags are nil")
	}

	if content.Map()["author"].(string) != author {
		test.Errorf("bad content map! %#v", content.Map())
	}

	var err error
	if _, err = content.JSON(); err != nil {
		test.Fatal(err)
	}

	var map_source Content
	map_source.FromMap(content.Map())

	if map_source.Author != author {
		test.Errorf("author %s not osurced from map %#v", author, content.Map())
	}
}
