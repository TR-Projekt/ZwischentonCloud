package model

import "database/sql"

type Zwischenton struct {
	ID          int         `json:"zwischenton_id"`
	Version     string      `json:"zwischenton_version"`
	Valid       bool        `json:"zwischenton_is_valid" db:"zwischenton_is_valid"`
	Name        string      `json:"zwischenton_name" db:"zwischenton_name"`
	Description string      `json:"zwischenton_description" db:"zwischenton_description"`
	Include     interface{} `json:"include,omitempty"`
}

func ZwischentonScan(rs *sql.Rows) (Zwischenton, error) {
	var z Zwischenton
	return z, rs.Scan(&z.ID, &z.Version, &z.Valid, &z.Name, &z.Description)
}
