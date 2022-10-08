package database

import (
	"github.com/jmoiron/sqlx"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Longest allowed mimetype is 255 ( {127}/{127} ) per RFC 4288
// Longest allowed email is 254 per RFC 5321
// Max CHAR size is 255

var (
	database_handle *sqlx.DB
	tables          map[string]string = map[string]string{
		CONTENT_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			file_url CHAR(64) NOT NULL,
			author CHAR(36) NOT NULL,
			mime CHAR(255) NOT NULL,
			like_count BIGINT UNSIGNED NOT NULL,
			dislike_count BIGINT UNSIGNED NOT NULL,
			repub_count BIGINT UNSIGNED NOT NULL,
			view_count BIGINT UNSIGNED NOT NULL,
			comment_count BIGINT UNSIGNED NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			featured BOOLEAN,
			featurable BOOLEAN,
			removed BOOLEAN,
			nsfw BOOLEAN,
			order_index BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT`,
		USER_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			email CHAR(254) UNIQUE NOT NULL,
			nick CHAR(16) UNIQUE NOT NULL,
			bio CHAR(255) NOT NULL,
			subscriber_count BIGINT UNSIGNED NOT NULL,
			subscription_count BIGINT UNSIGNED NOT NULL,
			post_count BIGINT UNSIGNED NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			moderator BOOLEAN NOT NULL,
			admin BOOLEAN NOT NULL,
			order_index BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT`,
		AUTH_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			hash BINARY(60) NOT NULL`,
		TOKEN_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			token BINARY(24) UNIQUE,
			created BIGINT UNSIGNED NOT NULL`,
		SECRET_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			secret BINARY(128) UNIQUE`,
		TAG_TABLE: `
			id CHAR(36) NOT NULL,
			tag CHAR(64) NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			order_index BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT,
			CONSTRAINT no_dupe_tags UNIQUE(id, tag),
			CONSTRAINT content_bound_tags FOREIGN KEY (id) REFERENCES ` + CONTENT_TABLE + `(id) ON DELETE CASCADE`,
		SUBSCRIPTION_TABLE: `
			subscriber CHAR(36) NOT NULL,
			subscription CHAR(36) NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			order_index BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT,
			CONSTRAINT no_dupe_subscriptions UNIQUE(subscriber, subscription)`,
		BAN_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			banner CHAR(36) NOT NULL,
			banned CHAR(36) NOT NULL,
			reason CHAR(255),
			expires BIGINT UNSIGNED NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			forever BOOLEAN,
			order_index BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT`,
		REPORT_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			reporter CHAR(36) NOT NULL,
			reported CHAR(36) NOT NULL,
			type CHAR(31) NOT NULL,
			reason CHAR(255) NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			resolved BOOLEAN NOT NULL,
			resolution CHAR(255) NOT NULL,
			order_index BIGINT UNSIGNED UNIQUE NOT NULL AUTO_INCREMENT`,
	}

	tableOrdered []string = []string{
		USER_TABLE,
		CONTENT_TABLE,
		AUTH_TABLE,
		TOKEN_TABLE,
		SECRET_TABLE,
		SUBSCRIPTION_TABLE,
		BAN_TABLE,
		REPORT_TABLE,
		TAG_TABLE,
	}
)

const (
	USER_TABLE         = "users"
	CONTENT_TABLE      = "content"
	AUTH_TABLE         = "auth"
	TOKEN_TABLE        = "token"
	SECRET_TABLE       = "secret"
	TAG_TABLE          = "tags"
	SUBSCRIPTION_TABLE = "subs"
	BAN_TABLE          = "bans"
	REPORT_TABLE       = "reports"
)

func listStringReverse(source []string) (reversed []string) {
	var size int = len(source)
	reversed = make([]string, size)

	var index int
	var it string
	for index, it = range source {
		reversed[size-index-1] = it
	}

	return
}

/**
 * Ping the database, and return any error
 * useful for health checks
 */
func Health() (err error) {
	err = database_handle.Ping()
	return
}

/**
 * Connect to a database, given a connection string
 * If the connection fails a ping, this function wil panic with the err
 * The conenction string should look something like
 * user:pass@tcp(addr)/table
 */
func Connect(address string) {
	var err error
	if database_handle, err = sqlx.Open("mysql", address); err != nil {
		panic(err)
	}

	if err = Health(); err != nil {
		panic(err)
	}
}

func Create() {
	var err error
	var table string
	for _, table = range tableOrdered {
		if _, err = database_handle.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", table, tables[table])); err != nil {
			panic(err)
		}
	}
}

func EmptyTable(table string) (err error) {
	_, err = database_handle.Exec("DELETE FROM " + table)
	return
}
