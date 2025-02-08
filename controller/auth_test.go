package controller

import (
	"database/sql"
)

type mockDB struct {
	*sql.DB
	execFunc     func(query string, args ...interface{}) (sql.Result, error)
	queryRowFunc func(query string, args ...interface{}) *sql.Row
}
type Result interface {
	lastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
type Database interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type mockResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (m mockResult) LastInsertId() (int64, error) {
	return m.lastInsertId, nil
}

func (m mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.execFunc != nil {
		return m.execFunc(query, args...)
	}
	return &mockResult{1, 1}, nil
}

func (m *mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.queryRowFunc != nil {
		return m.queryRowFunc(query, args...)
	}
	row := sql.Row{}
	return &row
}
