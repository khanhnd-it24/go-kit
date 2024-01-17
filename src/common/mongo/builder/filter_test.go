package mgobuilder

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestAnd(t *testing.T) {
	e := bson.M{"username": bson.M{"$eq": "test"}}

	testStr := "test"
	c := Filter().And(Eq[string]("username", &testStr)).Build()

	assert.Equal(t, e, c)
}

func TestAndNil(t *testing.T) {
	e := bson.M{}
	filter := struct {
		UserName *string
	}{nil}
	c := Filter().And(Eq[string]("username", filter.UserName)).Build()

	assert.Equal(t, e, c)
}

func TestAndNoGeneric(t *testing.T) {
	e := bson.M{"username": bson.M{"$eq": "test"}}

	testStr := "test"
	c := Filter().And(Eq("username", &testStr)).Build()

	assert.Equal(t, e, c)
}
