package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Entity string

const (
	Zwischenton Entity = "zwischenton"
)

func Select(db *sql.DB, table string, objectIDs []int) (*sql.Rows, error) {

	var query string
	var vars []int
	if len(objectIDs) == 0 {
		query = "SELECT * FROM " + table + "s;"
		vars = []int{}
	} else {
		placeholder := DBPlaceholderForIDs(objectIDs)
		query = "SELECT * FROM " + table + "s WHERE " + table + "_id IN (" + placeholder + ");"
		vars = objectIDs
	}
	return ExecuteRowQuery(db, query, InterfaceInt(vars))
}

func Search(db *sql.DB, table string, name string) (*sql.Rows, error) {

	query := "SELECT * FROM " + table + "s WHERE " + table + "_name LIKE CONCAT('%', ?, '%');"
	vars := []interface{}{name}
	rows, err := ExecuteRowQuery(db, query, vars)
	if err != nil {
		err = errors.New("failed to search `" + table + "` for `" + name + "` with error: " + err.Error())
	}
	return rows, err
}

func Resource(db *sql.DB, object string, objectID int, resource string) (*sql.Rows, error) {

	var query string
	if object == "tag" || resource == "festival" {
		query = "SELECT * FROM " + resource + "s WHERE " + resource + "_id IN (SELECT `associated_" + resource + "` FROM `map_" + resource + "_" + object + "` WHERE `associated_" + object + "`=?);"
	} else {
		query = "SELECT * FROM " + resource + "s WHERE " + resource + "_id IN (SELECT `associated_" + resource + "` FROM `map_" + object + "_" + resource + "` WHERE `associated_" + object + "`=?);"
	}
	vars := []interface{}{objectID}

	rows, err := ExecuteRowQuery(db, query, vars)
	if err != nil {
		err = errors.New("failed to fetch `" + resource + "` for `" + object + "` with error: " + err.Error())
	}
	return rows, err
}

func SetResource(db *sql.DB, object string, objectID int, resource string, resourceID int) error {

	query := "SELECT `map_id` from `map_" + object + "_" + resource + "` WHERE associated_" + object + " =?;"

	vars := []interface{}{objectID}
	rows, err := ExecuteRowQuery(db, query, vars)
	if err != nil {
		return errors.New("failed to set `" + resource + "` for `" + object + "` with error: " + err.Error())
	}
	var mapID string
	for rows.Next() {
		err = rows.Scan(&mapID)
		if err != nil {
			return errors.New("failed to set `" + resource + "` for `" + object + "` with error: " + err.Error())
		}
	}
	// for to-many relationships we want to create a new map entry
	if mapID == "" || resource == "event" || resource == "tag" || resource == "link" || resource == "artist" {
		query = "INSERT INTO `map_" + object + "_" + resource + "` ( `associated_" + object + "` , `associated_" + resource + "` ) VALUES (?,?);"
		vars := []interface{}{objectID, resourceID}
		result, err := ExecuteQuery(db, query, vars)
		if err != nil {
			return errors.New("failed to set `" + resource + "` for `" + object + "` with error: " + err.Error())
		}
		_, err = result.LastInsertId()
		if err != nil {
			return errors.New("failed to set `" + resource + "` for `" + object + "` with error: " + err.Error())
		}
		return nil
	} else {
		query = "UPDATE `map_" + object + "_" + resource + "` SET associated_" + resource + "=? WHERE map_id=?;"
		vars := []interface{}{resourceID, mapID}
		_, err := ExecuteQuery(db, query, vars)
		if err != nil {
			return errors.New("failed to set `" + resource + "` for `" + object + "` with error: " + err.Error())
		}
		return nil
	}
}

func RemoveResource(db *sql.DB, object string, objectID int, resource string, resourceID int) error {

	query := "SELECT `map_id` FROM `map_" + object + "_" + resource + "` WHERE associated_" + object + " =? AND associated_" + resource + "=?;"
	vars := []interface{}{objectID, resourceID}
	rows, err := ExecuteRowQuery(db, query, vars)
	if err != nil {
		return errors.New("failed to remove `" + resource + "` for `" + object + "` with error: " + err.Error())
	}
	if rows != nil {
		var mapID string
		for rows.Next() {
			err = rows.Scan(&mapID)
			if err != nil {
				return errors.New("failed to remove `" + resource + "` for `" + object + "` with error: " + err.Error())
			}
		}
		vars = []interface{}{mapID}
		query = "DELETE FROM `map_" + object + "_" + resource + "` WHERE map_id=?;"
		result, err := ExecuteQuery(db, query, vars)
		if err != nil {
			return errors.New("failed to remove `" + resource + "` for `" + object + "` with error: " + err.Error())
		}
		numOfAffectedRows, err := result.RowsAffected()
		if err != nil {
			return errors.New("failed to remove `" + resource + "` for `" + object + "` with error: " + err.Error())
		}
		if numOfAffectedRows != 1 {
			return errors.New("failed to remove `" + resource + "` for `" + object + "` with error: " + strconv.FormatInt(numOfAffectedRows, 10) + " rows where affected")
		}
		return nil
	} else {
		return errors.New("failed to remove `" + resource + "` for `" + object + "` with error: there is no resource to remove")
	}
}

func Insert(db *sql.DB, table string, object interface{}) (*sql.Rows, error) {

	fields, err := DBFields(object)
	if err != nil {
		return nil, errors.New("failed to insert `" + fmt.Sprintf("%s", object) + "` into `" + table + "` with error: " + err.Error())
	}

	placeholder, err := DBPlaceholder(object)
	if err != nil {
		return nil, errors.New("failed to insert `" + fmt.Sprintf("%s", object) + "` into `" + table + "` with error: " + err.Error())
	}

	vars, err := DBValues(object)
	if err != nil {
		return nil, errors.New("failed to insert `" + fmt.Sprintf("%s", object) + "` into `" + table + "` with error: " + err.Error())
	}

	query := "INSERT INTO " + table + "s(" + fields + ") VALUES (" + placeholder + ");"
	result, err := ExecuteQuery(db, query, vars)
	if err != nil {
		return nil, errors.New("failed to insert `" + fmt.Sprintf("%s", object) + "` into `" + table + "` with error: " + err.Error())
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("failed to insert `" + fmt.Sprintf("%s", object) + "` into `" + table + "` with error: " + err.Error())
	}

	rows, err := Select(db, table, []int{int(insertID)})
	if err != nil {
		err = errors.New("failed to insert `" + fmt.Sprintf("%s", object) + "` into `" + table + "` with error: " + err.Error())
	}
	return rows, err
}

func Update(db *sql.DB, table string, objectID int, object interface{}) (*sql.Rows, error) {

	keyValuePairs, err := DBKeyValuePairs(object)
	if err != nil {
		return nil, errors.New("failed to update `" + fmt.Sprintf("%s", object) + "` in `" + table + "` with error: " + err.Error())
	}

	vars, err := DBValues(object)
	if err != nil {
		return nil, errors.New("failed to update `" + fmt.Sprintf("%s", object) + "` in `" + table + "` with error: " + err.Error())
	}

	vars = append(vars, objectID) // for *table*_id value
	query := "UPDATE " + table + "s SET " + keyValuePairs + " WHERE `" + table + "_id`=?;"
	_, err = ExecuteQuery(db, query, vars)
	if err != nil {
		return nil, errors.New("failed to update `" + fmt.Sprintf("%s", object) + "` in `" + table + "` with error: " + err.Error())
	}

	rows, err := Select(db, table, []int{int(objectID)})
	if err != nil {
		err = errors.New("failed to update `" + fmt.Sprintf("%s", object) + "` in `" + table + "` with error: " + err.Error())
	}
	return rows, err
}

func Delete(db *sql.DB, table string, objectID int) error {

	query := "DELETE FROM " + table + "s WHERE " + table + "_id=?"
	vars := []interface{}{objectID}
	result, err := ExecuteQuery(db, query, vars)
	if err != nil {
		return errors.New("failed to delete `" + strconv.Itoa(objectID) + "` from `" + table + "` with error: " + err.Error())
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to delete `" + strconv.Itoa(objectID) + "` from `" + table + "` with error: " + err.Error())
	}
	if numOfAffectedRows != 1 {
		return errors.New("failed to delete `" + strconv.Itoa(objectID) + "` from `" + table + "` with error: " + strconv.FormatInt(numOfAffectedRows, 10) + " rows where affected")
	}
	return nil
}
