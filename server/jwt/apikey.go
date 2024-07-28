package token

import "database/sql"

type APIKey struct {
	ID      int    `json:"api_key_id" sql:"api_key_id"`
	Key     string `json:"api_key" sql:"api_key"`
	Comment string `json:"api_key_comment" sql:"api_key_comment"`
}

func APIKeyScan(rs *sql.Rows) (APIKey, error) {
	var u APIKey
	return u, rs.Scan(&u.ID, &u.Key, &u.Comment)
}
