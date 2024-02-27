package gormbuilder

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
	"testing"
)

func TestAnd(t *testing.T) {
	e := []clause.Expression{
		clause.And(clause.Eq{Column: "username", Value: "test"}),
	}

	testStr := "test"
	c := Filter().And(Eq[string]("username", &testStr)).Build()

	assert.Equal(t, e, c)
}

func TestAndNil(t *testing.T) {
	e := make([]clause.Expression, 0)
	filter := struct {
		UserName *string
	}{nil}
	c := Filter().And(Eq[string]("username", filter.UserName)).Build()

	assert.Equal(t, e, c)
}

func TestAndNoGeneric(t *testing.T) {
	e := []clause.Expression{
		clause.And(clause.Eq{Column: "username", Value: "test"}),
	}

	testStr := "test"
	c := Filter().And(Eq("username", &testStr)).Build()

	assert.Equal(t, e, c)
}

func TestAndIn(t *testing.T) {
	e := []clause.Expression{
		clause.And(clause.IN{Column: "username", Values: []interface{}{"a", "b"}}),
	}

	c := Filter().And(In("username", []string{"a", "b"})).Build()

	assert.Equal(t, e, c)
}

func TestOr(t *testing.T) {
	e := []clause.Expression{
		clause.Or(
			clause.Eq{Column: "username", Value: "a"},
			clause.Eq{Column: "password", Value: "b"},
		),
	}

	testStr := "a"
	testStr2 := "b"
	c := Filter().Or(
		Eq("username", &testStr),
		Eq("password", &testStr2),
	).Build()

	assert.Equal(t, e, c)
}
