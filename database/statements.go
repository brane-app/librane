package database

const (
	CONTENT_FIELDS = `
id,
file_url,
author,
mime,
like_count,
dislike_count,
repub_count,
view_count,
comment_count,
created,
featured,
featurable,
removed,
nsfw`
	USER_FIELDS = `
id,
email,
nick,
bio,
subscriber_count,
subscription_count,
post_count,
created,
moderator,
admin`
	BAN_FIELDS = `
id,
banner,
banned,
reason,
created,
expires,
forever`
	REPORT_FIELDS = `
id,
reporter,
reported,
type,
reason,
created,
resolved,
resolution`

	READ_INDEX_OF_CONTENT                = "SELECT order_index FROM " + CONTENT_TABLE + " WHERE id=? LIMIT 1"
	READ_CONTENT_ID                      = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " WHERE id=? LIMIT 1"
	READ_MANY_CONTENT_AFTER_ID           = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " WHERE order_index<(" + READ_INDEX_OF_CONTENT + ") ORDER BY order_index DESC LIMIT ?"
	READ_MANY_CONTENT                    = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " ORDER BY order_index DESC LIMIT ?"
	READ_MANY_CONTENT_OF_AUTHOR          = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " WHERE author=? ORDER BY order_index DESC LIMIT ?"
	READ_MANY_CONTENT_OF_AUTHOR_AFTER_ID = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " WHERE author=? AND order_index<(" + READ_INDEX_OF_CONTENT + ") ORDER BY order_index DESC LIMIT ?"
	DELETE_CONTENT_ID                    = "DELETE FROM " + CONTENT_TABLE + " WHERE id=? LIMIT 1"

	READ_TAGS_OF_ID       = "SELECT tag FROM " + TAG_TABLE + " WHERE id=?"
	READ_TAGS_OF_MANY_ID  = "SELECT id, tag FROM " + TAG_TABLE + " WHERE id IN "
	WRITE_TAGS_OF_MANY_ID = "REPLACE INTO " + TAG_TABLE + " (id, tag, created) VALUES "
	DELETE_TAGS_OF_ID     = "DELETE FROM " + TAG_TABLE + " WHERE id=?"

	READ_USER_OF_ID                 = "SELECT " + USER_FIELDS + " FROM " + USER_TABLE + " WHERE id=? LIMIT 1"
	READ_USER_OF_EMAIL              = "SELECT " + USER_FIELDS + " FROM " + USER_TABLE + " WHERE email=? LIMIT 1"
	READ_USER_OF_NICK               = "SELECT " + USER_FIELDS + " FROM " + USER_TABLE + " WHERE nick=? LIMIT 1"
	DELETE_USER_OF_ID               = "DELETE FROM " + USER_TABLE + " WHERE id=? LIMIT 1"
	INCREMENT_USER_POST_COUNT_OF_ID = "UPDATE " + USER_TABLE + " SET post_count=post_count+1 WHERE id=?"
	READ_ANY_PRIVILEGE_OF_ID        = "SELECT admin, moderator FROM " + USER_TABLE + " WHERE id=?"
	READ_MODERATOR_OF_ID            = "SELECT moderator FROM " + USER_TABLE + " WHERE id=?"
	READ_ADMIN_OF_ID                = "SELECT admin FROM " + USER_TABLE + " WHERE id=?"
	WRITE_MODERATOR_OF_ID           = "UPDATE " + USER_TABLE + " SET moderator=? WHERE id=?"
	WRITE_ADMIN_OF_ID               = "UPDATE " + USER_TABLE + " SET admin=? WHERE id=?"

	READ_INDEX_OF_BAN          = "SELECT order_index FROM " + BAN_TABLE + " WHERE id=? LIMIT 1"
	READ_BAN_OF_ID             = "SELECT " + BAN_FIELDS + " FROM " + BAN_TABLE + " WHERE id=? LIMIT 1"
	READ_BANS_OF_USER          = "SELECT " + BAN_FIELDS + " FROM " + BAN_TABLE + " WHERE banned=? ORDER BY order_index DESC LIMIT ?"
	READ_BANS_OF_USER_AFTER_ID = "SELECT " + BAN_FIELDS + " FROM " + BAN_TABLE + " WHERE banned=? AND order_index<(" + READ_INDEX_OF_BAN + ") ORDER BY order_index DESC LIMIT ?"
	READ_BANS_OF_USER_COUNT    = "SELECT COUNT(id) FROM " + BAN_TABLE + " WHERE (banned=? AND forever) OR (banned=? AND expires>?) LIMIT 1"

	READ_REPORT_OF_ID                = "SELECT " + REPORT_FIELDS + " FROM " + REPORT_TABLE + " WHERE id=?"
	READ_INDEX_OF_REPORT             = "SELECT order_index FROM " + REPORT_TABLE + " WHERE id=? LIMIT 1"
	READ_REPORTS_UNRESOLVED          = "SELECT " + REPORT_FIELDS + " FROM " + REPORT_TABLE + " WHERE resolved=0 ORDER BY order_index DESC LIMIT ?"
	READ_REPORTS_UNRESOLVED_AFTER_ID = "SELECT " + REPORT_FIELDS + " FROM " + REPORT_TABLE + " WHERE resolved=0 AND order_index<(" + READ_INDEX_OF_REPORT + ") ORDER BY order_index DESC LIMIT ?"

	WRITE_SECRET_OF_ID  = "REPLACE INTO " + SECRET_TABLE + " (id, secret) VALUES (?, ?)"
	READ_SECRET_OF_ID   = "SELECT secret FROM " + SECRET_TABLE + " WHERE id=? LIMIT 1"
	DELETE_SECRET_OF_ID = "DELETE FROM " + SECRET_TABLE + " WHERE id=? LIMIT 1"

	WRITE_TOKEN_OF_ID  = "REPLACE INTO " + TOKEN_TABLE + " (id, token, created) VALUES (?, ?, ?)"
	READ_TOKEN_STAT    = "SELECT id, created FROM " + TOKEN_TABLE + " WHERE token=? LIMIT 1"
	DELETE_TOKEN       = "DELETE FROM " + TOKEN_TABLE + " WHERE token=?"
	DELETE_TOKEN_OF_ID = "DELETE FROM " + TOKEN_TABLE + " WHERE id=?"

	READ_HASH_OF_ID  = "SELECT hash FROM " + AUTH_TABLE + " WHERE id=? LIMIT 1"
	WRITE_HASH_OF_ID = "REPLACE INTO " + AUTH_TABLE + " (id, hash) VALUES (?, ?)"
)
