package exoplanetsrepo

import (
	"database/sql"
	"fmt"
	"space-api/pkg/models"
	"space-api/pkg/sqlutil"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

type ExoplanetsRepository struct {
	queryBuilder    *sqlutil.QueryBuilder
	exoplanetsQuery *goqu.SelectDataset
	countQuery      *goqu.SelectDataset
	connInfo        string
}

func New(queryBuilder *sqlutil.QueryBuilder, connInfo string) *ExoplanetsRepository {
	exoplanetsQuery := goqu.Select(
		"id",
		"planet_name",
		"host_name",
		"system_number",
		"discovery_method",
		"year_discovered",
	).From("exoplanet")

	countQuery := goqu.Select(
		goqu.L("COUNT(*) total"),
	).From("exoplanet")

	return &ExoplanetsRepository{
		queryBuilder:    queryBuilder,
		exoplanetsQuery: exoplanetsQuery,
		countQuery:      countQuery,
		connInfo:        connInfo,
	}
}

func (e *ExoplanetsRepository) ReadExoplanets(queryModifiers *sqlutil.QueryModifiers) ([]models.Exoplanet, error) {
	modifiedSQL, args, err := e.queryBuilder.BuildQuery(e.exoplanetsQuery, queryModifiers)
	if err != nil {
		return nil, fmt.Errorf("build sql: %w", err)
	}

	db, err := sql.Open("postgres", e.connInfo)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	rows, err := db.Query(modifiedSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	records := make([]models.Exoplanet, 0)

	for rows.Next() {
		var record models.Exoplanet
		if err = rows.Scan(
			&record.Id,
			&record.PlanetName,
			&record.HostName,
			&record.SystemNumber,
			&record.DiscoveryMethod,
			&record.YearDiscovered); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (e *ExoplanetsRepository) ReadCount(queryModifiers *sqlutil.QueryModifiers) (int, error) {
	modifiedSQL, args, err := e.queryBuilder.BuildCountQuery(e.countQuery, queryModifiers)
	if err != nil {
		return 0, fmt.Errorf("build sql: %w", err)
	}

	db, err := sql.Open("postgres", e.connInfo)
	if err != nil {
		return 0, fmt.Errorf("connect to postgres: %w", err)
	}

	var count int
	row := db.QueryRow(modifiedSQL, args...)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("scan row: %w", err)
	}

	return count, nil
}
