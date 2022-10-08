package database

import (
	"github.com/brane-app/types-library"
	"github.com/jmoiron/sqlx"

	"database/sql"
	"time"
)

func WriteBan(ban map[string]interface{}) (err error) {
	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(BAN_TABLE, ban)

	_, err = database_handle.Exec(statement, values...)
	return
}

/**
 * Read a single ban of id `ID`
 * Done in one query
 */
func ReadSingleBan(ID string) (ban types.Ban, exists bool, err error) {
	if err = database_handle.QueryRowx(READ_BAN_OF_ID, ID).StructScan(&ban); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	return
}

/**
 * Read a slice of bans of a user
 * Done in one query
 */
func ReadBansOfUser(ID, before string, count int) (bans []types.Ban, size int, err error) {
	var rows *sqlx.Rows
	if before == "" {
		rows, err = database_handle.Queryx(READ_BANS_OF_USER, ID, count)
	} else {
		rows, err = database_handle.Queryx(READ_BANS_OF_USER_AFTER_ID, ID, before, count)
	}

	if err != nil {
		return
	}

	defer rows.Close()

	bans = make([]types.Ban, count)
	size = 0
	for rows.Next() {
		rows.StructScan(&bans[size])
		size++
	}

	bans = bans[:size]
	return
}

/**
 * Get whether or not a user is banned, either by a permanent ban, or an expirable ban
 * Done in one query
 */
func IsBanned(ID string) (banned bool, err error) {
	var count int
	var now int64 = time.Now().Unix()
	if err = database_handle.QueryRowx(READ_BANS_OF_USER_COUNT, ID, ID, now).Scan(&count); err != nil {
		return
	}

	banned = count != 0
	return
}

/**
 * Create or update a report for some user
 * Done in one query
 */
func WriteReport(report map[string]interface{}) (err error) {
	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(REPORT_TABLE, report)

	_, err = database_handle.Exec(statement, values...)
	return
}

/**
 * Read a slice of unresolved reports (ie, the mod queue) by order of most recent
 * Done in one query
 */
func ReadManyUnresolvedReport(before string, count int) (reports []types.Report, size int, err error) {
	var rows *sqlx.Rows
	if before == "" {
		rows, err = database_handle.Queryx(READ_REPORTS_UNRESOLVED, count)
	} else {
		rows, err = database_handle.Queryx(READ_REPORTS_UNRESOLVED_AFTER_ID, before, count)
	}

	if err != nil {
		return
	}

	defer rows.Close()

	reports = make([]types.Report, count)
	size = 0
	for rows.Next() {
		rows.StructScan(&reports[size])
		size++
	}

	reports = reports[0:size]
	return
}

/**
 * Lookup single report by it's ID
 * Done in one query
 */
func ReadSingleReport(ID string) (report types.Report, exists bool, err error) {
	if err = database_handle.QueryRowx(READ_REPORT_OF_ID, ID).StructScan(&report); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}

	exists = true
	return
}
