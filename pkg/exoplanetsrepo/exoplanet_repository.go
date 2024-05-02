package exoplanetsrepo

import (
	"database/sql"
	"fmt"
	"space-api/pkg/models"
	"space-api/pkg/querybuilder"
	"space-api/pkg/querymodifiers"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

type ExoplanetsRepository struct {
	queryBuilder *querybuilder.QueryBuilder
	baseQuery    *goqu.SelectDataset
}

func New(queryBuilder *querybuilder.QueryBuilder) *ExoplanetsRepository {
	query := goqu.Select(
		"id",
		"planet_name",
		"host_name",
		"system_number",
		"discovery_method",
		"year_discovered",
	).From("exoplanet")

	return &ExoplanetsRepository{
		queryBuilder: queryBuilder,
		baseQuery:    query,
	}
}

func (e *ExoplanetsRepository) Read(queryModifiers *querymodifiers.QueryModifiers) ([]models.Exoplanet, error) {
	modifiedSQL, args, err := e.queryBuilder.Build(e.baseQuery, queryModifiers)
	if err != nil {
		return nil, fmt.Errorf("build sql: %w", err)
	}

	connInfo := "host=localhost port=5432 user=test password=test dbname=space sslmode=disable"
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	rows, err := db.Query(modifiedSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var records []models.Exoplanet

	for rows.Next() {
		var record models.Exoplanet
		if err = rows.Scan(
			&record.Id,
			&record.PlanetName,
			&record.HostName,
			&record.SystemNumber,
			&record.DiscoveryMethod,
			&record.YearDiscovered); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
