package database

import (
	"database/sql"
	"errors"

	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
)

func GetAllAPIKeys(db *sql.DB) ([]string, error) {

	query := "SELECT * FROM api_keys;"
	vars := []interface{}{}

	rows, err := ExecuteRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}
	keys := []string{}
	for rows.Next() {
		key, err := token.APIKeyScan(rows)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key.Key)
	}
	return keys, nil
}

func AddAPIKey(db *sql.DB, key token.APIKey) error {

	query := "INSERT INTO api_keys(`api_key`, `api_key_comment`) VALUES (?, ?);"
	vars := []interface{}{key.Key, key.Comment}

	result, err := ExecuteQuery(db, query, vars)
	if err != nil {
		return err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	if insertID == 0 {
		return errors.New("failed to insert new API key without mysql error")
	}
	return nil
}

func UpdateAPIKey(db *sql.DB, key token.APIKey) error {

	query := "UPDATE api_keys SET `api_key`=?, `api_key_comment`=? WHERE `api_key_id`=?;"
	vars := []interface{}{key.Key, key.Comment, key.ID}

	result, err := ExecuteQuery(db, query, vars)
	if err != nil {
		return err
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numOfAffectedRows != 1 {
		return errors.New("failed to update API key without mysql error")
	}
	return nil
}

func RemoveAPIKey(db *sql.DB, keyID string) error {

	query := "DELETE FROM api_keys WHERE `api_key_id`=?;"
	vars := []interface{}{keyID}

	result, err := ExecuteQuery(db, query, vars)
	if err != nil {
		return err
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numOfAffectedRows != 1 {
		return errors.New("failed to remove API key without mysql error")
	}
	return nil
}
