package database

import (
	"github.com/google/uuid"
	"github.com/brane-app/types-library"

	"testing"
)

var (
	user         types.User         = types.NewUser("imonke", "mmm, monke", "me@imonke.io")
	writableUser map[string]interface{} = user.Map()
)

func userOK(test *testing.T, data map[string]interface{}, have types.User) {
	if data["id"].(string) != have.ID {
		test.Errorf("User ID mismatch! have: %s, want: %s", have.ID, data["id"])
	}

	if data["bio"].(string) != have.Bio {
		test.Errorf("User bio mismatch! have: %s, want: %s", have.Bio, data["bio"])
	}
}

func Test_WriteUser(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{},
		map[string]interface{}{
			"id":  uuid.New().String(),
			"bio": "' or 1=1; DROP TABLE user",
		},
	}

	var err error
	var mod map[string]interface{}
	for _, mod = range mods {
		mod = mapMod(writableUser, mod)
		if err = WriteUser(mod); err != nil {
			test.Fatal(err)
		}
	}
}

func Test_WriteUser_err(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{
			"id":  uuid.New().String(),
			"bio": nil,
		},
		map[string]interface{}{
			"id":     uuid.New().String(),
			"answer": 42,
		},
	}

	var mod map[string]interface{} = mapCopy(writableUser)
	delete(mod, "id")
	mod["bio"] = "foobar"

	var err error
	if err = WriteUser(mod); err == nil {
		test.Errorf("data %+v produced no error!", mod)
	}

	for _, mod = range mods {
		mod = mapMod(writableUser, mod)
		if err = WriteUser(mod); err == nil {
			test.Errorf("data %+v produced no error!", mod)
		}
	}

}

func Test_ReadSingleUser(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableUser)
	modified["id"] = uuid.New().String()

	WriteUser(modified)

	var user types.User
	var exists bool
	var err error
	if user, exists, err = ReadSingleUser(modified["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("user of id %s does not exist!", modified["id"])
	}

	userOK(test, modified, user)
}

func Test_ReadSingleUserNick(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableUser)
	modified["id"] = uuid.New().String()
	modified["nick"] = "readme"

	WriteUser(modified)

	var user types.User
	var exists bool
	var err error
	if user, exists, err = ReadSingleUserNick(modified["nick"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("user of id %s does not exist!", modified["id"])
	}

	userOK(test, modified, user)
}

func Test_ReadSingleUserEmail(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableUser)
	modified["id"] = uuid.New().String()
	modified["email"] = "read@me.io"

	WriteUser(modified)

	var user types.User
	var exists bool
	var err error
	if user, exists, err = ReadSingleUserEmail(modified["email"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("user of id %s does not exist!", modified["id"])
	}

	userOK(test, modified, user)
}

func Test_ReadSingleUser_NotExists(test *testing.T) {
	var id string = uuid.New().String()

	var user types.User
	var exists bool
	var err error
	if user, exists, err = ReadSingleUser(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("Query for nonexisting id got %+v", user)
	}
}

func Test_DeleteUser(test *testing.T) {
	var mod map[string]interface{} = map[string]interface{}{
		"id":    uuid.New().String(),
		"email": "delete@monke.io",
	}

	var err error
	if err = WriteUser(mapMod(writableUser, mod)); err != nil {
		test.Fatal(err)
	}

	if err = DeleteUser(mod["id"].(string)); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if _, exists, err = ReadSingleUser(mod["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("user %s exists after being deleted!", mod["id"])
	}
}

func Test_IncrementPostCount(test *testing.T) {
	var writable map[string]interface{} = types.NewUser("increment", "", "i@monke.io").Map()
	var unchanged map[string]interface{} = types.NewUser("unchanged", "", "u@monke.io").Map()

	var err error
	if err = WriteUser(writable); err != nil {
		test.Fatal(err)
	}

	if err = WriteUser(unchanged); err != nil {
		test.Fatal(err)
	}

	var id string = writable["id"].(string)
	if err = IncrementPostCount(id); err != nil {
		test.Fatal(err)
	}

	var fetched types.User
	if fetched, _, err = ReadSingleUser(id); err != nil {
		test.Fatal(err)
	}

	var count int = int(writable["post_count"].(float64))
	if fetched.PostCount != count+1 {
		test.Errorf("post count not incremented! %d -> %d", count, fetched.PostCount)
	}

	id = unchanged["id"].(string)
	if fetched, _, err = ReadSingleUser(id); err != nil {
		test.Fatal(err)
	}

	if fetched.PostCount != count {
		test.Errorf("post count also affected %s!", id)
	}
}

func Test_IsModerator(test *testing.T) {
	var moderator types.User = types.NewUser("mod", "", "mod@imonke.io")
	moderator.Moderator = true

	defer DeleteUser(moderator.ID)
	WriteUser(moderator.Map())

	var is_mod bool
	var err error
	if is_mod, err = IsModerator(moderator.ID); err != nil {
		test.Fatal(err)
	}

	if !is_mod {
		test.Errorf("%s is not a moderator", moderator.ID)
	}

}

func Test_IsModerator_nomod(test *testing.T) {
	var moderator types.User = types.NewUser("mod", "", "mod@imonke.io")

	defer DeleteUser(moderator.ID)
	WriteUser(moderator.Map())

	var is_mod bool
	var err error
	if is_mod, err = IsModerator(moderator.ID); err != nil {
		test.Fatal(err)
	}

	if is_mod {
		test.Errorf("%s is a moderator", moderator.ID)
	}
}

func Test_IsModerator_nobody(test *testing.T) {
	var is_mod bool
	var err error
	if is_mod, err = IsModerator(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if is_mod {
		test.Errorf("nobody is a moderator")
	}
}

func Test_IsAdmin(test *testing.T) {
	var moderator types.User = types.NewUser("admin", "", "admin@imonke.io")
	moderator.Admin = true

	defer DeleteUser(moderator.ID)
	WriteUser(moderator.Map())

	var is_mod bool
	var err error
	if is_mod, err = IsAdmin(moderator.ID); err != nil {
		test.Fatal(err)
	}

	if !is_mod {
		test.Errorf("%s is not a moderator", moderator.ID)
	}

}

func Test_IsAdmin_nomod(test *testing.T) {
	var moderator types.User = types.NewUser("admin", "", "admin@imonke.io")

	defer DeleteUser(moderator.ID)
	WriteUser(moderator.Map())

	var is_mod bool
	var err error
	if is_mod, err = IsAdmin(moderator.ID); err != nil {
		test.Fatal(err)
	}

	if is_mod {
		test.Errorf("%s is a moderator", moderator.ID)
	}
}

func Test_IsAdmin_nobody(test *testing.T) {
	var is_mod bool
	var err error
	if is_mod, err = IsAdmin(uuid.New().String()); err != nil {
		test.Fatal(err)
	}

	if is_mod {
		test.Errorf("nobody is a moderator")
	}
}

func Test_SetModerator(test *testing.T) {
	var moderator types.User = types.NewUser("mod", "", "mod@imonke.io")

	defer DeleteUser(moderator.ID)
	var err error
	if err = WriteUser(moderator.Map()); err != nil {
		test.Fatal(err)
	}

	if err = SetModerator(moderator.ID, true); err != nil {
		test.Fatal(err)
	}

	var is_mod bool
	if is_mod, err = IsModerator(moderator.ID); err != nil {
		test.Fatal(err)
	}

	if !is_mod {
		test.Errorf("%s was not made a moderator", moderator.ID)
	}

	if err = SetModerator(moderator.ID, false); err != nil {
		test.Fatal(err)
	}

	if is_mod, err = IsModerator(moderator.ID); err != nil {
		test.Fatal(err)
	}

	if is_mod {
		test.Errorf("%s is still a moderator", moderator.ID)
	}
}

func Test_SetAdmin(test *testing.T) {
	var admin types.User = types.NewUser("mod", "", "mod@imonke.io")

	defer DeleteUser(admin.ID)
	var err error
	if err = WriteUser(admin.Map()); err != nil {
		test.Fatal(err)
	}

	if err = SetAdmin(admin.ID, true); err != nil {
		test.Fatal(err)
	}

	var is_mod bool
	if is_mod, err = IsAdmin(admin.ID); err != nil {
		test.Fatal(err)
	}

	if !is_mod {
		test.Errorf("%s was not made a admin", admin.ID)
	}

	if err = SetAdmin(admin.ID, false); err != nil {
		test.Fatal(err)
	}

	if is_mod, err = IsAdmin(admin.ID); err != nil {
		test.Fatal(err)
	}

	if is_mod {
		test.Errorf("%s is still a admin", admin.ID)
	}
}
