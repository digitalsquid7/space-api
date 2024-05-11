package sqlutil

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type QueryBuilder struct {
	dialect goqu.DialectWrapper
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		dialect: goqu.Dialect("postgres"),
	}
}

func (q *QueryBuilder) BuildQuery(query *goqu.SelectDataset, modifiers *QueryModifiers) (string, []any, error) {
	query = q.applyPaging(query, modifiers)
	query = q.applyFilters(query, modifiers)
	query = q.applySorting(query, modifiers)

	query.Prepared(true)
	return query.ToSQL()
}

func (q *QueryBuilder) BuildCountQuery(query *goqu.SelectDataset, modifiers *QueryModifiers) (string, []any, error) {
	query = q.applyFilters(query, modifiers)

	query.Prepared(true)
	return query.ToSQL()
}

func (q *QueryBuilder) applyPaging(query *goqu.SelectDataset, modifiers *QueryModifiers) *goqu.SelectDataset {
	if modifiers.Page == nil {
		return query
	}

	query = query.Limit(uint(modifiers.Page.Size))
	query = query.Offset(uint(modifiers.Page.Size) * (uint(modifiers.Page.Number) - 1))

	return query
}

func (q *QueryBuilder) applyFilters(query *goqu.SelectDataset, modifiers *QueryModifiers) *goqu.SelectDataset {
	expressions := make([]goqu.Expression, 0, len(modifiers.Filters))
	for _, filter := range modifiers.Filters {
		col := goqu.C(filter.Field.SQLName)

		var expression goqu.Expression
		switch filter.Operator {
		case Equals:
			expression = col.Eq(filter.Value)
		case GTE:
			expression = col.Gte(filter.Value)
		case LTE:
			expression = col.Lte(filter.Value)
		}

		expressions = append(expressions, expression)
	}

	return query.Where(goqu.And(expressions...))
}

func (q *QueryBuilder) applySorting(query *goqu.SelectDataset, modifiers *QueryModifiers) *goqu.SelectDataset {
	expressions := make([]exp.OrderedExpression, 0, len(modifiers.SortFields))

	for _, sortField := range modifiers.SortFields {
		if sortField.Direction == Ascending {
			expressions = append(expressions, goqu.C(sortField.Field.SQLName).Asc())
		} else {
			expressions = append(expressions, goqu.C(sortField.Field.SQLName).Desc())
		}
	}

	return query.Order(expressions...)
}
