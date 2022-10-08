package database

import (
	"github.com/google/uuid"

	"encoding/base64"
	"testing"
	"time"
)

func Test_SetPassword(test *testing.T) {
	var password string
	var err error
	if password, err = randomString(64); err != nil {
		test.Fatal(err)
	}

	var id string = uuid.New().String()

	if err = SetPassword(id, password); err != nil {
		test.Fatal(err)
	}

	var ok bool
	if ok, err = CheckPassword(id, password); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Errorf("Set password %s does not match retrieved!", password)
	}
}

func Test_SetPassword_length(test *testing.T) {
	var id string = uuid.New().String()
	var err error
	var password string
	var index int = 1
	for index != 4*64 {
		index = index * 4
		if password, err = randomString(index); err != nil {
			test.Fatal(err)
		}
		if err = SetPassword(id, password); err != nil {
			test.Fatal(err)
		}
	}
}

func Test_CheckPassword_wrong(test *testing.T) {
	var sets []string = []string{
		"password",
		"",
	}

	var id string = uuid.New().String()

	var err error
	var password string
	if password, err = randomString(64); err != nil {
		test.Fatal(err)
	}

	if err = SetPassword(id, password); err != nil {
		test.Fatal(err)
	}

	var ok bool
	var set string
	for _, set = range sets {
		if ok, err = CheckPassword(id, set); err != nil {
			test.Fatal(err)
		}

		if ok {
			test.Errorf("password %s should not match, but does!", set)
		}
	}
}

func Test_CheckPassword_nobody(test *testing.T) {
	var ok bool
	var err error
	if ok, err = CheckPassword(uuid.New().String(), "password"); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("password for random uuid is ok!")
	}
}

func Test_CreateSecret(test *testing.T) {
	var id string = uuid.New().String()

	var err error
	var secret string
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal(err)
	}

	_ = secret
}

func Test_CheckSecret(test *testing.T) {
	var id string = uuid.New().String()

	var secret string
	var err error
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal(err)
	}

	var valid bool
	if valid, err = CheckSecret(id, secret); err != nil {
		test.Fatal(err)
	}

	if !valid {
		test.Errorf("Just set secret %s is invalid for %s!", secret, id)
	}
}

func Test_CheckSecret_invalid(test *testing.T) {
	var id string = uuid.New().String()

	var secret string
	var err error
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal(err)
	}

	var valid bool
	if valid, err = CheckSecret(uuid.New().String(), secret); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("Just set secret %s is valid for a random uuid!", secret)
	}

	if valid, err = CheckSecret(id, "not_a_secret"); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("A bad secret is valid for uuid %s!", id)
	}
}

func Test_RevokeSecretOf(test *testing.T) {
	var id string = uuid.New().String()

	var secret string
	var err error
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal()
	}

	var valid bool
	if valid, err = CheckSecret(id, secret); err != nil {
		test.Fatal(err)
	}

	if !valid {
		test.Errorf("Just set secret %s is invalid for %s!", secret, id)
	}

	if err = RevokeSecretOf(id); err != nil {
		test.Fatal(err)
	}

	if valid, err = CheckSecret(id, secret); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("revoked secret %s for %s is still valid!", secret, id)
	}
}

func Test_CreateToken(test *testing.T) {
	var id string = uuid.New().String()

	var expires int64
	var err error
	if _, expires, err = CreateToken(id); err != nil {
		test.Fatal(err)
	}

	var now int64 = time.Now().Unix()
	var delta int64 = expires - now
	if delta > TOKEN_TTL {
		test.Errorf("token expiry off! now: %d, expires: %d (delta %d, more than %d)", now, expires, delta, TOKEN_TTL)
	}
}

func Test_ReadTokenStat(test *testing.T) {
	var id string = uuid.New().String()

	var token string
	var err error
	if token, _, err = CreateToken(id); err != nil {
		test.Fatal(err)
	}

	var owner string
	var valid bool
	if owner, valid, err = ReadTokenStat(token); err != nil {
		test.Fatal(err)
	}

	if !valid {
		test.Errorf("token %s is not valid", token)
	}

	if owner != id {
		test.Errorf("owner mismatch! have: %s, want: %s", owner, id)
	}
}

func Test_ReadTokenStat_expired(test *testing.T) {
	var id string = uuid.New().String()

	var token string
	var err error
	if token, _, err = CreateToken(id); err != nil {
		test.Fatal(err)
	}

	var bytes []byte
	if bytes, err = base64.URLEncoding.DecodeString(token); err != nil {
		test.Fatal(err)
	}

	var statement string = "REPLACE INTO " + TOKEN_TABLE + " (id, token, created) VALUES (?, ?, ?)"
	if _, err = database_handle.Exec(statement, id, bytes, 1); err != nil {
		test.Fatal(err)
	}

	var valid bool
	if _, valid, err = ReadTokenStat(token); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("expired token %s is still valid!", token)
	}
}

func Test_ReadTokenStat_nobody(test *testing.T) {
	var token string
	var err error
	if token, err = randomString(TOKEN_LENGTH); err != nil {
		test.Fatal(err)
	}

	var owner string
	var valid bool
	if owner, valid, err = ReadTokenStat(token); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("Random token %s is valid for %s!", token, owner)
	}
}

func Test_ReadTokenStat_err(test *testing.T) {
	var valid bool
	var err error
	if _, _, err = ReadTokenStat("f"); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("invalid base64 is valid!")
	}
}

func Test_RevokeToken(test *testing.T) {
	var id string = uuid.New().String()

	var token string
	var err error
	if token, _, err = CreateToken(id); err != nil {
		test.Fatal(err)
	}

	if err = RevokeToken(token); err != nil {
		test.Fatal(err)
	}

	var valid bool
	if _, valid, err = ReadTokenStat(token); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("token %s is valid", token)
	}
}

func Test_RevokeTokenOf(test *testing.T) {
	var id string = uuid.New().String()

	var token string
	var err error
	if token, _, err = CreateToken(id); err != nil {
		test.Fatal(err)
	}

	if err = RevokeTokenOf(id); err != nil {
		test.Fatal(err)
	}

	var valid bool
	if _, valid, err = ReadTokenStat(token); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("token %s is valid", token)
	}
}
