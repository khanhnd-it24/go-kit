package mgobuilder

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Option func(*FindOptionBuilder)

type FindOptionBuilder struct {
	fieldMap map[string]string
}

func NewFindOptionBuilder() *FindOptionBuilder {
	mapper := map[string]string{
		"id": "_id",
	}

	return &FindOptionBuilder{fieldMap: mapper}
}

func WithMapper(mapper map[string]string) Option {
	return func(b *FindOptionBuilder) {
		mergeMap := make(map[string]string, len(b.fieldMap)+len(mapper))

		for k, v := range b.fieldMap {
			mergeMap[k] = v
		}
		for k, v := range mapper {
			mergeMap[k] = v
		}
		b.fieldMap = mergeMap
	}
}

func (b *FindOptionBuilder) getField(field string) string {
	if dbField, ok := b.fieldMap[field]; ok {
		return dbField
	}
	return field
}

func NewFindOption(option *FindOptionHolder, opts ...Option) *options.FindOptions {
	builder := NewFindOptionBuilder()

	for _, opt := range opts {
		opt(builder)
	}

	findOpts := options.Find()

	if option.Pagination != nil {
		p := option.Pagination

		findOpts.SetSkip(p.Offset).SetLimit(p.Limit)
	}

	if len(option.Sorts) > 0 {
		sortOps := bson.D{}

		for _, sortOpt := range option.Sorts {
			dbField := builder.getField(sortOpt.Field)
			sortOps = append(sortOps, bson.E{Key: dbField, Value: sortOpt.Direction})
		}

		findOpts.SetSort(sortOps)
	}

	return findOpts
}

func NewFindOneOption(option FindOptionHolder, opts ...Option) *options.FindOneOptions {
	builder := NewFindOptionBuilder()

	for _, opt := range opts {
		opt(builder)
	}

	findOpts := options.FindOne()

	if len(option.Sorts) > 0 {
		sortOps := bson.D{}

		for _, sortOpt := range option.Sorts {
			dbField := builder.getField(sortOpt.Field)
			sortOps = append(sortOps, bson.E{Key: dbField, Value: sortOpt.Direction})
		}

		findOpts.SetSort(sortOps)
	}

	return findOpts
}
