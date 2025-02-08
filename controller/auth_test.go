package controller

import "database/sql"

type mockDB struct {
	*sql.DB
	execFunc     func(query string, args ...interface{}) (sql.Result, error)
	queryRowFunc func(query string, args ...interface{}) *sql.Row
}

type mockResult struct {
	lastInsertedId int64
	rowsAffected   int64
}


