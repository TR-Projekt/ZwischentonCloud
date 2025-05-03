package model

import "database/sql"

type Zwischenton struct {
	ID          int    `json:"zwischenton_id"`
	Version     string `json:"zwischenton_version"`
	Valid       bool   `json:"zwischenton_is_valid" db:"zwischenton_is_valid"`
	Name        string `json:"zwischenton_name" db:"zwischenton_name"`
	Description string `json:"zwischenton_description" db:"zwischenton_description"`
	Include     any    `json:"include,omitempty"`
}

func ZwischentonScan(rs *sql.Rows) (Zwischenton, error) {
	var z Zwischenton
	return z, rs.Scan(&z.ID, &z.Version, &z.Valid, &z.Name, &z.Description)
}

type Situation struct {
	ID          int     `json:"situation_id"`
	Version     string  `json:"situation_version"`
	Latitude    float32 `json:"situation_lat" db:"situation_lat"`
	Longitude   float32 `json:"situation_lon" db:"situation_lon"`
	Radius      float32 `json:"situation_radius" db:"situation_radius"`
	Description string  `json:"situation_description" db:"situation_description"`
}

func SituationsScan(rs *sql.Rows) (Situation, error) {
	var p Situation
	return p, rs.Scan(&p.ID, &p.Version, &p.Latitude, &p.Longitude, &p.Radius, &p.Description)
}
