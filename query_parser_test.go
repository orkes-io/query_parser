package query_parser

import (
	"log"
	"testing"

	"github.com/orkes-io/query_parser/util"
)

func TestQueryParser(t *testing.T) {
	whereClause := `name != "john" AND (age < 60) OR (salary > 0 AND color = "Grey")`
	mongoQueryMap := ConvertWhereClauseToMongoFilter(whereClause)

	jsonString, err := util.ConvertMapToJsonString(mongoQueryMap)
	if err != nil {
		log.Printf("Error encoding in JSON: %s", err)
		t.Error(err)
	}

	log.Printf("Mongo Query Map: %+v", jsonString)

	expected := `{
	"$or": [
		{
			"$and": [
				{
					"color": {
						"$eq": "\"Grey\""
					}
				},
				{
					"salary": {
						"$gt": "0"
					}
				}
			]
		},
		{
			"$and": [
				{
					"age": {
						"$lt": "60"
					}
				},
				{
					"name": {
						"$ne": "\"john\""
					}
				}
			]
		}
	]
}`

	if expected != jsonString {
		t.Errorf("\nExpected: %+v\nReceived: %+v", jsonString, expected)
	}
}

func TestQueryParser2(t *testing.T) {
	whereClause := `name != "john" OR (age < 60) OR (salary > 0 AND color = "Grey")`
	mongoQueryMap := ConvertWhereClauseToMongoFilter(whereClause)

	jsonString, err := util.ConvertMapToJsonString(mongoQueryMap)
	if err != nil {
		log.Printf("Error encoding in JSON: %s", err)
		t.Error(err)
	}

	log.Printf("Mongo Query Map: %+v", jsonString)

	expected := `{
	"$or": [
		{
			"$and": [
				{
					"color": {
						"$eq": "\"Grey\""
					}
				},
				{
					"salary": {
						"$gt": "0"
					}
				}
			]
		},
		{
			"$or": [
				{
					"age": {
						"$lt": "60"
					}
				},
				{
					"name": {
						"$ne": "\"john\""
					}
				}
			]
		}
	]
}`

	if expected != jsonString {
		t.Errorf("\nExpected: %+v\nReceived: %+v", jsonString, expected)
	}
}

func TestQueryParser3(t *testing.T) {
	whereClause := `name != "john" OR salary > 0 AND color = "Grey"`
	mongoQueryMap := ConvertWhereClauseToMongoFilter(whereClause)

	jsonString, err := util.ConvertMapToJsonString(mongoQueryMap)
	if err != nil {
		log.Printf("Error encoding in JSON: %s", err)
		t.Error(err)
	}

	log.Printf("Mongo Query Map: %+v", jsonString)

	expected := `{
	"$or": [
		{
			"$and": [
				{
					"color": {
						"$eq": "\"Grey\""
					}
				},
				{
					"salary": {
						"$gt": "0"
					}
				}
			]
		},
		{
			"name": {
				"$ne": "\"john\""
			}
		}
	]
}`

	if expected != jsonString {
		t.Errorf("\nExpected: %+v\nReceived: %+v", jsonString, expected)
	}
}
