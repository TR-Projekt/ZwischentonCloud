package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/TR-Projekt/zwischentoncloud/server/database"
	token "github.com/TR-Projekt/zwischentoncloud/server/jwt"
	"github.com/TR-Projekt/zwischentoncloud/server/model"
	"github.com/go-chi/chi/v5"
)

type HandlerEncapsulator3000 struct {
	ZwischentonDB *sql.DB
	IdentityDB    *sql.DB
	Auth          *token.AuthService
	Validator     *token.ValidationService
}

func GetObject(db *sql.DB, r *http.Request, entity string) ([]any, error) {

	objectID, err := ObjectID(r)
	if err != nil {
		return nil, err
	}
	values := r.URL.Query()

	return GetObjects(db, entity, []int{objectID}, values)
}

func GetObjects(db *sql.DB, entity string, objectIDs []int, values url.Values) ([]any, error) {

	var idValues []int
	var rels []string
	var err error
	idValues = append(idValues, objectIDs...)
	if len(values) != 0 {
		// search with name
		name := values.Get("name")
		if name != "" {
			return SearchObjects(db, entity, name)
		}
		// filter by ids
		ids := values.Get("ids")
		if ids != "" {
			var err error
			idValues, err = ObjectIDs(ids)
			if err != nil {
				return nil, err
			}
		}
		// handle include later
		include := values.Get("include")
		if include != "" {
			rels, err = RelationshipNames(include)
			if err != nil {
				return nil, err
			}
		}
	}
	rows, err := database.Select(db, entity, idValues)
	if err != nil {
		return nil, err
	}
	// no rows and no error indicate a successful query but an empty result
	if rows == nil {
		return []any{}, nil
	}
	var fetchedObjects []any
	// iterate over the rows an create
	for rows.Next() {
		// scan the object
		obj, err := AnonScan(entity, rows)
		if err != nil {
			return nil, err
		}
		if rels != nil {
			objId, err := AnonID(entity, obj)
			if err != nil {
				return nil, err
			}
			includedRels, err := GetRelationships(db, entity, objId, rels)
			if err != nil {
				return nil, err
			}
			obj, err = AnonInclude(entity, obj, includedRels)
			if err != nil {
				return nil, err
			}
		}
		// add object result slice
		fetchedObjects = append(fetchedObjects, obj)
	}
	return fetchedObjects, nil
}

func SearchObjects(db *sql.DB, entity string, name string) ([]any, error) {

	rows, err := database.Search(db, entity, name)
	if err != nil {
		return nil, err
	}
	// no rows and no error indicate a successful query but an empty result
	if rows == nil {
		return []any{}, nil
	}
	var fetchedObjects []any
	// iterate over the rows an create
	for rows.Next() {
		// scan the link
		obj, err := AnonScan(entity, rows)
		if err != nil {
			return nil, err
		}
		// add object result slice
		fetchedObjects = append(fetchedObjects, obj)
	}
	return fetchedObjects, nil
}

func GetRelationships(db *sql.DB, entity string, objectID int, relationships []string) (any, error) {

	// TODO check the relationship strings for sake of writing `include=links`instead of `include=link`?
	relsDict := make(map[string]any)
	for _, value := range relationships {
		objcts, err := GetAssociatedObjects(db, entity, objectID, value, nil)
		if err != nil {
			return nil, err
		}
		relsDict[value] = objcts
	}
	return relsDict, nil
}

func GetAssociation(db *sql.DB, r *http.Request, entity string, association string) ([]any, error) {

	objectID, err := ObjectID(r)
	if err != nil {
		return nil, err
	}
	includes := Includes(r)
	return GetAssociatedObjects(db, entity, objectID, association, includes)
}

func GetAssociatedObjects(db *sql.DB, entity string, objectID int, association string, includes []string) ([]any, error) {

	rows, err := database.Resource(db, entity, objectID, association)
	if err != nil {
		return nil, err
	}
	// no rows and no error indicate a successful query but an empty result
	if rows == nil {
		return []any{}, nil
	}
	var fetchedObjects []any
	// iterate over the rows an create
	for rows.Next() {
		// scan the link
		obj, err := AnonScan(association, rows)
		if err != nil {
			return nil, err
		}
		if includes != nil {
			objId, err := AnonID(association, obj)
			if err != nil {
				return nil, err
			}
			includedRels, err := GetRelationships(db, association, objId, includes)
			if err != nil {
				return nil, err
			}
			obj, err = AnonInclude(association, obj, includedRels)
			if err != nil {
				return nil, err
			}
		}
		// add object result slice
		fetchedObjects = append(fetchedObjects, obj)
	}
	if fetchedObjects == nil {
		fetchedObjects = []any{}
	}
	return fetchedObjects, nil
}

func SetAssociation(db *sql.DB, r *http.Request, entity string, association string) error {

	objectID, err := ObjectID(r)
	if err != nil {
		return err
	}
	resourceID, err := ResourceID(r)
	if err != nil {
		return err
	}
	err = database.SetResource(db, entity, objectID, association, resourceID)
	if err != nil {
		return err
	}
	return nil
}

func RemoveAssociation(db *sql.DB, r *http.Request, entity string, association string) error {

	objectID, err := ObjectID(r)
	if err != nil {
		return err
	}
	resourceID, err := ResourceID(r)
	if err != nil {
		return err
	}
	err = database.RemoveResource(db, entity, objectID, association, resourceID)
	if err != nil {
		return err
	}
	return nil
}

func Create(db *sql.DB, r *http.Request, entity string) ([]any, error) {

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		return nil, readBodyErr
	}
	objectToCreate, err := AnonUnmarshal(entity, body)
	if err != nil {
		return nil, err
	}
	rows, err := database.Insert(db, entity, objectToCreate)
	if err != nil {
		return nil, err
	}
	// no rows and no error indicate a successful query but an empty result
	if rows == nil {
		return []any{}, nil
	}
	var fetchedObjects []any
	// iterate over the rows an create
	for rows.Next() {
		// scan the link
		obj, err := AnonScan(entity, rows)
		if err != nil {
			return nil, err
		}
		// add object result slice
		fetchedObjects = append(fetchedObjects, obj)
	}
	return fetchedObjects, nil
}

func Update(db *sql.DB, r *http.Request, entity string) ([]any, error) {

	objectID, err := ObjectID(r)
	if err != nil {
		return nil, err
	}
	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		return nil, readBodyErr
	}
	objectToUpdate, err := AnonUnmarshal(entity, body)
	if err != nil {
		return nil, err
	}
	rows, err := database.Update(db, entity, objectID, objectToUpdate)
	if err != nil {
		return nil, err
	}
	var fetchedObjects []any
	// iterate over the rows an create
	for rows.Next() {
		// scan the link
		obj, err := AnonScan(entity, rows)
		if err != nil {
			return nil, err
		}
		// add object result slice
		fetchedObjects = append(fetchedObjects, obj)
	}
	return fetchedObjects, nil
}

func Delete(db *sql.DB, r *http.Request, entity string) error {

	objectID, err := ObjectID(r)
	if err != nil {
		return err
	}
	err = database.Delete(db, entity, objectID)
	if err != nil {
		return err
	}
	return nil
}

func ObjectID(r *http.Request) (int, error) {

	objectID := chi.URLParam(r, "objectID")
	num, err := strconv.ParseUint(objectID, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(num), nil
}

func ResourceID(r *http.Request) (int, error) {

	objectID := chi.URLParam(r, "resourceID")
	num, err := strconv.ParseUint(objectID, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(num), nil
}

func Includes(r *http.Request) []string {

	include := r.URL.Query().Get("include")
	if include != "" {
		includes, err := RelationshipNames(include)
		if err == nil {
			return includes
		}
	}
	return nil
}

func ObjectIDs(idsString string) ([]int, error) {

	var ids []int
	for id := range strings.SplitSeq(idsString, ",") {

		idNum, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(idNum))
	}
	if len(ids) == 0 {
		return nil, errors.New("ids parsing: failed to provide an id")
	}
	return ids, nil
}

func RelationshipNames(includes string) ([]string, error) {

	return strings.Split(includes, ","), nil
}

func AnonScan(entity string, rs *sql.Rows) (any, error) {

	if entity == "zwischenton" {
		return model.ZwischentonScan(rs)
	}
	if entity == "situation" {
		return model.SituationsScan(rs)
	} else {
		return nil, errors.New("scan row: tried to scan an unknown entity")
	}
}

func AnonInclude(entity string, object any, includes any) (any, error) {

	if entity == "zwischenton" {
		realObject := object.(model.Zwischenton)
		realObject.Include = includes
		return realObject, nil
	}
	if entity == "situation" {
		realObject := object.(model.Situation)
		return realObject, nil
	} else {
		return nil, errors.New("include relationship: tried to add relationships to an unknown entity")
	}
}

func AnonID(entity string, object any) (int, error) {

	if entity == "zwischenton" {
		return object.(model.Zwischenton).ID, nil
	}
	if entity == "situation" {
		return object.(model.Situation).ID, nil
	} else {
		return -1, errors.New("get id: tried to retrieve the ID of an unknown entity")
	}
}

func AnonUnmarshal(entity string, body []byte) (any, error) {

	if entity == "zwischenton" {
		var objectToCreate model.Zwischenton
		err := json.Unmarshal(body, &objectToCreate)
		if err != nil {
			return nil, err
		}
		return objectToCreate, nil
	}
	if entity == "situation" {
		var objectToCreate model.Situation
		err := json.Unmarshal(body, &objectToCreate)
		if err != nil {
			return nil, err
		}
		return objectToCreate, nil
	} else {
		return nil, errors.New("unmarshal object: tried to unmarshal an unknown entity")
	}
}
