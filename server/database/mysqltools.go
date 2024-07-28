package database

import "database/sql"

func ExecuteRowQuery(db *sql.DB, query string, args []interface{}) (*sql.Rows, error) {

	rows, err := db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rows, nil
}

func ExecuteQuery(db *sql.DB, query string, args []interface{}) (sql.Result, error) {

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	err = stmt.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}
