package mgobuilder

import "strings"

type FindOption func(*FindOptionHolder)

type Pagination struct {
	Offset int64
	Limit  int64
}

type Sort struct {
	Field     string
	Direction int64
}

type FindOptionHolder struct {
	Pagination *Pagination
	Sorts      []*Sort
}

func NewFindOptionHolder(opts ...FindOption) *FindOptionHolder {
	holder := &FindOptionHolder{}

	for _, opt := range opts {
		opt(holder)
	}

	return holder
}

func WithPagination(page, limit *int64) FindOption {
	return func(o *FindOptionHolder) {
		if page == nil || limit == nil {
			return
		}

		offset := (*page - 1) * *limit

		o.Pagination = &Pagination{
			Offset: offset,
			Limit:  *limit,
		}
	}
}

func WithSort(sorts *[]string) FindOption {
	return func(o *FindOptionHolder) {
		if sorts == nil {
			return
		}
		sortOpts := make([]*Sort, 0, len(*sorts))
		for _, sort := range *sorts {
			pos := strings.LastIndex(sort, "_")
			if pos < 0 {
				continue
			}

			field, dir := sort[:pos], sort[pos+1:]

			if field == "" {
				continue
			}

			if dir == "asc" {
				sortOpts = append(sortOpts, &Sort{
					Field:     field,
					Direction: 1,
				})
			}

			if dir == "desc" {
				sortOpts = append(sortOpts, &Sort{
					Field:     field,
					Direction: -1,
				})
			}
		}
		if len(sortOpts) == 0 {
			return
		}

		if len(o.Sorts) > 0 {
			newSorts := append(o.Sorts, sortOpts...)
			o.Sorts = newSorts
		} else {
			o.Sorts = sortOpts
		}
	}
}

func WithSortNaturalAsc() FindOption {
	return WithSort(&[]string{"id_asc"})
}

func WithSortNaturalDesc() FindOption {
	return WithSort(&[]string{"id_desc"})
}
