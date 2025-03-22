package main

import (
	"database/sql"
)

func getInt32Value (n sql.NullInt32) int {
	if n.Valid {
		return int(n.Int32)
	} else {
		return 0
	}
}

func getSqlStringValue (s sql.NullString) string {
	if s.Valid {
		return s.String
	} else {
		return ""
	}
}

func validateSqlNullString (s string) sql.NullString {
	if s != "" {
		return sql.NullString{Valid: true, String: s}
	}
	return sql.NullString{Valid: false}
}

func validateSqlNullInt32 (n int) sql.NullInt32 {
	if n != 0 {
		return sql.NullInt32{Valid: true, Int32: int32(n)}
	}
	return sql.NullInt32{Valid: false}
}