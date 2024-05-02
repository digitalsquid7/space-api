package sqlutil

import (
	"database/sql"
	"fmt"
)

type Selector[T any] struct {
	db         *sql.DB
	scanStruct func(rows *sql.Rows) (T, error)
}

func NewSelector[T any](db *sql.DB, scanStruct func(rows *sql.Rows) (T, error)) *Selector[T] {
	return &Selector[T]{
		db:         db,
		scanStruct: scanStruct,
	}
}

func (s *Selector[T]) Select(db *sql.DB, sql string, args ...any) ([]T, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var records []T

	for rows.Next() {
		record, err := s.scanStruct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan struct: %w", err)
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
