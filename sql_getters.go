package main

import "database/sql"

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