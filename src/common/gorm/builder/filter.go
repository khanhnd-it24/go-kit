package gormbuilder

import (
	"go-kit/src/common/utility"
	"gorm.io/gorm/clause"
)

func Eq[T any](columnName string, p *T) clause.Expression {
	if p == nil {
		return nil
	}
	return clause.Eq{Column: columnName, Value: *p}
}

func Neq[T any](columnName string, p *T) clause.Expression {
	if p == nil {
		return nil
	}
	return clause.Neq{Column: columnName, Value: *p}
}

func Gt[T any](columnName string, p *T) clause.Expression {
	if p == nil {
		return nil
	}
	return clause.Gt{Column: columnName, Value: *p}
}

func Gte[T any](columnName string, p *T) clause.Expression {
	if p == nil {
		return nil
	}
	return clause.Gte{Column: columnName, Value: *p}
}

func Lt[T any](columnName string, p *T) clause.Expression {
	if p == nil {
		return nil
	}
	return clause.Lt{Column: columnName, Value: *p}
}

func Lte[T any](columnName string, p *T) clause.Expression {
	if p == nil {
		return nil
	}
	return clause.Lte{Column: columnName, Value: *p}
}

func In[T any](columnName string, p []T) clause.Expression {
	if p == nil {
		return nil
	}
	list := utility.CvtToListInterface(p)
	return clause.IN{Column: columnName, Values: list}
}

func NIn[T any](columnName string, p []T) clause.Expression {
	if p == nil {
		return nil
	}
	list := utility.CvtToListInterface(p)
	return clause.Not(clause.IN{Column: columnName, Values: list})
}

func Like[T any](columnName string, p *T) clause.Expression {
	if p == nil {
		return nil
	}
	return clause.Like{Column: columnName, Value: p}
}

func filterNilConditions(conditions []clause.Expression) []clause.Expression {
	notNilConditions := make([]clause.Expression, 0, len(conditions))

	for _, c := range conditions {
		if c != nil {
			notNilConditions = append(notNilConditions, c)
		}
	}
	return notNilConditions
}

type logicOp int64

const (
	lopAnd logicOp = iota
	lopOr
	lopNor
	lopNot
)

type builder struct {
	conditionGroups map[logicOp][]clause.Expression
}

func Filter() *builder {
	b := &builder{}
	b.conditionGroups = make(map[logicOp][]clause.Expression)
	return b
}

func (b *builder) And(c ...clause.Expression) *builder {
	b.conditionGroups[lopAnd] = filterNilConditions(c)
	return b
}

func (b *builder) Or(c ...clause.Expression) *builder {
	notNilConditions := filterNilConditions(c)

	if len(notNilConditions) < 2 {
		panic("OR logical operator require minimum two conditions")
	}
	b.conditionGroups[lopOr] = notNilConditions
	return b
}

func (b *builder) Nor(c ...clause.Expression) *builder {
	notNilConditions := filterNilConditions(c)

	if len(notNilConditions) < 2 {
		panic("NOR logical operator require minimum two conditions")
	}
	b.conditionGroups[lopNor] = notNilConditions
	return b
}

func (b *builder) Not(c ...clause.Expression) *builder {
	b.conditionGroups[lopNot] = filterNilConditions(c)
	return b
}

func (b *builder) Build() []clause.Expression {
	m := make([]clause.Expression, 0)
	for lOp, conditions := range b.conditionGroups {
		switch lOp {
		case lopAnd:
			if len(conditions) > 0 {
				m = append(m, clause.And(conditions...))
			}

		case lopOr:
			if len(conditions) > 0 {
				m = append(m, clause.Or(conditions...))
			}

		case lopNor:
			if len(conditions) > 0 {
				m = append(m, clause.Not(clause.Or(conditions...)))
			}

		case lopNot:
			if len(conditions) > 0 {
				m = append(m, clause.Not(conditions...))
			}
		}
	}
	return m
}
