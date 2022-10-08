package database

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"
)

const (
	BCRYPT_ITERS  = 12
	TOKEN_LENGTH  = 24
	SECRET_LENGTH = 128
	TOKEN_TTL     = 60 * 60 * 24
)

func randomBytes(size int) (generated []byte, err error) {
	generated = make([]byte, size)
	_, err = rand.Read(generated)
	return
}

func randomString(size int) (generated string, err error) {
	var bytes []byte
	if bytes, err = randomBytes(size); err == nil {
		generated = base64.URLEncoding.EncodeToString(bytes)
	}

	return
}

/**
 * Create a secret for some user of id `ID`
 * Any existing secret for that user is destroyed
 * Done in one query:
 * 		update secret 	REPLACE INTO SECRET_TABLE (id, secret) VALUES ID, new_secret
 */
func CreateSecret(ID string) (secret string, err error) {
	var bytes []byte
	if bytes, err = randomBytes(SECRET_LENGTH); err != nil {
		return
	}

	secret = base64.URLEncoding.EncodeToString(bytes)
	_, err = database_handle.Exec(WRITE_SECRET_OF_ID, ID, bytes)
	return
}

/**
 * Check that a secret `secret` for some user of id `ID` matches
 * Done in one query:
 * 		read secret: 	SELECT secret FROM SECRET_TABLE WHERE id=ID LIMIT 1
 */
func CheckSecret(ID, secret string) (valid bool, err error) {
	var bytes []byte
	if err = database_handle.QueryRowx(READ_SECRET_OF_ID, ID).Scan(&bytes); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	valid = secret == base64.URLEncoding.EncodeToString(bytes)
	return
}

/**
 * Revoke the secret of some user of id `ID`
 * Done in one query:
 * 		delete row: 	DELETE FROM SECRET_TABLE WHERE id=ID LIMIT 1
 */
func RevokeSecretOf(ID string) (err error) {
	_, err = database_handle.Exec(DELETE_SECRET_OF_ID, ID)
	return
}

/**
 * Create a token for some user of id `ID` that expires in 24 hours
 * Any existing token for that user is destroyed
 * Done in one query:
 * 		update secret 	REPLACE INTO TOKEN_TABLE (id, token, created) VALUES ID, new_token, now
 */
func CreateToken(ID string) (token string, expires int64, err error) {
	var bytes []byte
	if bytes, err = randomBytes(TOKEN_LENGTH); err != nil {
		return
	}

	var now int64 = time.Now().Unix()
	expires = now + TOKEN_TTL
	token = base64.URLEncoding.EncodeToString(bytes)
	_, err = database_handle.Exec(WRITE_TOKEN_OF_ID, ID, bytes, now)
	return
}

/**
 * Read information about some token `token`
 * Returns who it belongs to, and whether or not it's valid
 * done in one query:
 * 		read token: SELECT id, created FROM TOKEN_TABLE WHERE token=? LIMIT 1
 */
func ReadTokenStat(token string) (owner string, valid bool, err error) {
	var bytes []byte
	if bytes, err = base64.URLEncoding.DecodeString(token); err != nil {
		err = nil
		return
	}

	var rows *sqlx.Rows
	if rows, err = database_handle.Queryx(READ_TOKEN_STAT, bytes); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	if !rows.Next() {
		return
	}

	var created int64
	if err = rows.Scan(&owner, &created); err == nil {
		valid = created <= time.Now().Unix() && created+TOKEN_TTL >= time.Now().Unix()
	}

	return
}

/**
 * Revoke some token `token`
 * Done in one query:
 * 		delete row: 	DELETE FROM TOKEN_TABLE WHERE token=token LIMIT 1
 */
func RevokeToken(token string) (err error) {
	var bytes []byte
	if bytes, err = base64.URLEncoding.DecodeString(token); err != nil {
		return
	}

	_, err = database_handle.Exec(DELETE_TOKEN, bytes)
	return
}

/**
 * Revoke the token of some user of id `ID`
 * Done in one query:
 * 		delete row: 	DELETE FROM TOKEN_TABLE WHERE id=ID LIMIT 1
 */
func RevokeTokenOf(ID string) (err error) {
	_, err = database_handle.Exec(DELETE_TOKEN_OF_ID, ID)
	return
}

/**
 * Check that password `password` matches the hash for user of id `ID`
 * Done in one query:
 *  		read hash: 		SELECT hash FROM AUTH_TABLE WHERE id=ID LIMIT 1
 */
func CheckPassword(ID, password string) (valid bool, err error) {
	var hash []byte
	if err = database_handle.QueryRowx(READ_HASH_OF_ID, ID).Scan(&hash); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	valid = bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
	return
}

/**
 * Set a password `password` for some user of id `ID`
 * Done in one query:
 * 		write row:		REPLACE INTO AUTH_TABLE (id, hash) VALUES (ID, hash(password))
 */
func SetPassword(ID, password string) (err error) {
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(password), BCRYPT_ITERS); err != nil {
		return
	}

	_, err = database_handle.Exec(WRITE_HASH_OF_ID, ID, hash)

	return
}
