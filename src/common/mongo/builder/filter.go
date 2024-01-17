package mgobuilder

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type op string

const (
	eqOp     = op("$eq")
	neOp     = op("$ne")
	gtOp     = op("$gt")
	gteOp    = op("$gte")
	ltOp     = op("$lt")
	lteOp    = op("$lte")
	inOp     = op("$in")
	ninOp    = op("$nin")
	existsOp = op("$exists")
	regexOp  = op("$regex")
)

type Condition struct {
	op        op
	value     interface{}
	fieldName string
}

func Eq[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}

	return &Condition{op: eqOp, value: *p, fieldName: fieldName}
}

func Ne[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: neOp, value: *p, fieldName: fieldName}
}

func Gt[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: gtOp, value: *p, fieldName: fieldName}
}

func Gte[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: gteOp, value: *p, fieldName: fieldName}
}

func Lt[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: ltOp, value: *p, fieldName: fieldName}
}

func Lte[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: lteOp, value: *p, fieldName: fieldName}
}

func In[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: inOp, value: *p, fieldName: fieldName}
}

func NIn[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: ninOp, value: *p, fieldName: fieldName}
}

func Exist[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: existsOp, value: *p, fieldName: fieldName}
}

func Contains[T any](fieldName string, p *T) *Condition {
	if p == nil {
		return nil
	}
	return &Condition{op: regexOp, value: *p, fieldName: fieldName}
}

func getArrayOfM(cond []*Condition) []bson.M {
	ar := make([]bson.M, 0, len(cond))
	for _, c := range cond {
		ar = append(ar, bson.M{c.fieldName: getM(c)})
	}
	return ar
}

func getM(c *Condition) bson.M {
	if c.op == regexOp {
		m := bson.M{string(c.op): primitive.Regex{Pattern: c.value.(string), Options: "i"}}
		return m
	}
	m := bson.M{string(c.op): c.value}
	return m
}

func mergeM(m1, m2 bson.M) bson.M {
	m := bson.M{}

	for k, v := range m1 {
		m[k] = v
	}

	for k, v := range m2 {
		m[k] = v
	}
	return m
}

func filterNilConditions(conditions []*Condition) []*Condition {
	notNilConditions := make([]*Condition, 0, len(conditions))

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
	conditionGroups map[logicOp][]*Condition
}

func Filter() *builder {
	b := &builder{}
	b.conditionGroups = make(map[logicOp][]*Condition)
	return b
}

func (b *builder) And(c ...*Condition) *builder {
	b.conditionGroups[lopAnd] = filterNilConditions(c)
	return b
}

func (b *builder) Or(c ...*Condition) *builder {
	notNilConditions := filterNilConditions(c)

	if len(notNilConditions) < 2 {
		panic("OR logical operator require minimum two conditions")
	}
	b.conditionGroups[lopOr] = notNilConditions
	return b
}

func (b *builder) Nor(c ...*Condition) *builder {
	notNilConditions := filterNilConditions(c)

	if len(notNilConditions) < 2 {
		panic("NOR logical operator require minimum two conditions")
	}
	b.conditionGroups[lopNor] = notNilConditions
	return b
}

func (b *builder) Not(c ...*Condition) *builder {
	b.conditionGroups[lopNot] = filterNilConditions(c)
	return b
}

func (b *builder) Build() bson.M {
	m := bson.M{}
	for lOp, conditions := range b.conditionGroups {
		switch lOp {
		case lopAnd:
			for _, c := range conditions {
				v, _ := m[c.fieldName]
				if oldM, ok := v.(bson.M); ok {
					m[c.fieldName] = mergeM(oldM, getM(c))
				} else {
					m[c.fieldName] = getM(c)
				}
			}

		case lopOr:
			m["$or"] = getArrayOfM(conditions)

		case lopNor:
			m["$nor"] = getArrayOfM(conditions)

		case lopNot:
			for _, c := range conditions {
				v, _ := m[c.fieldName]
				if oldM, ok := v.(bson.M); ok {
					m[c.fieldName] = mergeM(oldM, bson.M{"$not": getM(c)})
				} else {
					m[c.fieldName] = bson.M{"$not": getM(c)}
				}
			}
		}
	}
	return m
}
