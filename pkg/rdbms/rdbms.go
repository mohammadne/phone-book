package rdbms

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type MigrateDirection uint8

const (
	MigrateUp   MigrateDirection = 1
	MigrateDown MigrateDirection = 2
)

type RDBMS interface {
	// will be used for read and create
	QueryRow(query string, in []any, out []any) error

	// for cursor based queries
	Query(query string, in []any, out [][]any) error

	// will be used for update and delete
	Execute(query string, in []any) error
}

type rdbms struct {
	db *sql.DB
}

var (
	ErrPrepareStatement = "error when tying to prepare statement"
	ErrNotFound         = "there is no entry with provided arguments"
	ErrDuplicate        = "there is no entry with provided arguments"
	ErrorQueryRow       = "error when tying to read entry"
	ErrorQueryRows      = "error when tying to read entry"
	ErrExecute          = "error when tying to excute statement"
)

func (db *rdbms) QueryRow(query string, in []any, out []any) error {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s\n%v", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	if err = stmt.QueryRow(in...).Scan(out...); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New(ErrDuplicate)
		} else if err == sql.ErrNoRows {
			return errors.New(ErrNotFound)
		}
		return fmt.Errorf("%s\n%v", ErrorQueryRow, err)
	}

	return nil
}

func (db *rdbms) Query(query string, in []any, out [][]any) error {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s\n%v", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(in...)
	if err != nil {
		return fmt.Errorf("%s\n%v", ErrorQueryRows, err)
	}
	defer rows.Close()

	var index = 0
	for ; rows.Next(); index++ {
		if err = rows.Scan(out[index]...); err != nil {
			if err == sql.ErrNoRows {
				return errors.New(ErrNotFound)
			}
			return fmt.Errorf("%s\n%v", ErrorQueryRow, err)
		}
	}
	out = out[:index+1]

	if err := rows.Err(); err != nil {
		return fmt.Errorf("%s\n%v", ErrorQueryRows, err)
	}

	return nil
}

func (db *rdbms) Execute(query string, in []any) error {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s\n%v", ErrPrepareStatement, err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(in...); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New(ErrDuplicate)
		}
		return fmt.Errorf("%s\n%v", ErrExecute, err)
	}

	return nil
}
