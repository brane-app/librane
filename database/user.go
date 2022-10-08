package database

import (
	"github.com/brane-app/types-library"

	"database/sql"
)

/**
 * Write some user `user` into USER_TABLE
 * Uses 1 query
 * 		write user: 	REPLACE INTO USER_TABLE (keys...) VALUES (values...)
 * Returns error, if any
 */
func WriteUser(user map[string]interface{}) (err error) {
	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(USER_TABLE, user)

	_, err = database_handle.Query(statement, values...)
	return
}

/**
 * Delete some user from USER_TABLE
 * Uses 1 query:
 * 		delete user: 	DELETE FROM USER_TABLE WHERE id=ID LIMIT 1
 */
func DeleteUser(ID string) (err error) {
	_, err = database_handle.Exec(DELETE_USER_OF_ID, ID)
	return
}

func readSingleUserKey(statement, query string) (user types.User, exists bool, err error) {
	if err = database_handle.QueryRowx(statement, query).StructScan(&user); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	return
}

/**
 * Read some user of id `ID` from USER_TABLE
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE id=ID LIMIT 1
 */
func ReadSingleUser(ID string) (user types.User, exists bool, err error) {
	user, exists, err = readSingleUserKey(READ_USER_OF_ID, ID)
	return
}

/**
 * Read some user of email `email` from USER_TABLE
 * Works in the same way as ReadSingleUser, but with email
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE email=email LIMIT 1
 */
func ReadSingleUserEmail(email string) (user types.User, exists bool, err error) {
	user, exists, err = readSingleUserKey(READ_USER_OF_EMAIL, email)
	return
}

/**
 * Read some user of email `email` from USER_TABLE
 * Works in the same way as ReadSingleUser, but with nick
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE nick=nick LIMIT 1
 */
func ReadSingleUserNick(nick string) (user types.User, exists bool, err error) {
	user, exists, err = readSingleUserKey(READ_USER_OF_NICK, nick)
	return
}

/**
 * Increment the post count of user of id `ID` by one
 * Done in one query
 * 		increment: UPDATE USER_TABLE SET post_count=post_count+1 WHERE id=ID
 */
func IncrementPostCount(ID string) (err error) {
	_, err = database_handle.Exec(INCREMENT_USER_POST_COUNT_OF_ID, ID)
	return
}

func IsModerator(ID string) (moderator bool, err error) {
	var admin bool
	if err = database_handle.QueryRowx(READ_ANY_PRIVILEGE_OF_ID, ID).Scan(&admin, &moderator); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	moderator = moderator || admin
	return
}

func IsAdmin(ID string) (admin bool, err error) {
	if err = database_handle.QueryRowx(READ_ADMIN_OF_ID, ID).Scan(&admin); err == sql.ErrNoRows {
		err = nil
	}

	return
}

func SetModerator(ID string, state bool) (err error) {
	_, err = database_handle.Exec(WRITE_MODERATOR_OF_ID, state, ID)
	return
}

func SetAdmin(ID string, state bool) (err error) {
	_, err = database_handle.Exec(WRITE_ADMIN_OF_ID, state, ID)
	return
}
