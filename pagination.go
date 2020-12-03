package nidhi

import (
	sq "github.com/elgris/sqrl"
)

func addPagination(st *sq.SelectBuilder, op *PaginationOptions) (*sq.SelectBuilder, []interface{}, error) {
	if op == nil {
		return st, nil, nil
	}

	if len(op.OrderBy) == 0 {
		op.OrderBy = append(op.OrderBy, OrderBy{Field: idOrderer{}})
	}

	var (
		selections     []interface{}
		orderByIdAdded bool
	)
	for i, f := range op.OrderBy {
		od := order(f.Desc)
		if op.Backward {
			od = !od
		}

		if i == 0 {
			if f.Field.Name() == "id" {
				st = st.Column(ColId)
				var id string
				selections = append(selections, &id)
				if op.Cursor != "" {
					st = st.Where(ColId+od.Cursor(), op.Cursor)
				}
			} else {
				selections = append(selections, f.Field.New())
				st = st.Column(f.Field.Name())
				if op.Cursor != "" {
					v, id, err := f.Field.Decode(op.Cursor)
					if err != nil {
						if err != InvalidCursor {
							// TODO: Log Warning here
						}
						return nil, nil, InvalidCursor
					}

					st = st.Where(
						sq.Or{
							sq.Expr(f.Field.Name()+od.Cursor(), v),
							sq.And{
								sq.Eq{f.Field.Name(): v},
								sq.Expr(ColId+order(op.Backward).Cursor(), id),
							},
						},
					)
				}
			}
		}

		if f.Field.Name() == "id" {
			orderByIdAdded = true
		}

		st = st.OrderBy(f.Field.Name() + od.Direction())
	}

	if !orderByIdAdded {
		st = st.OrderBy(ColId + order(op.Backward).Direction())
	}

	op.HasMore = false
	return st.Limit(op.Limit + 1), selections, nil
}

// true is descending
type order bool

func (o order) Direction() string {
	if o {
		return ` DESC`
	} else {
		return ` ASC`
	}
}

func (o order) Cursor() string {
	if o {
		return ` < ?`
	} else {
		return ` > ?`
	}
}

type Orderer interface {
	Name() string
	Encode(v interface{}, id string) string
	Decode(cursor string) (interface{}, string, error)
	New() interface{}
}

type idOrderer struct{}

func (idOrderer) Name() string {
	return "id"
}

func (idOrderer) Encode(v interface{}, id string) string {
	return ""
}
func (idOrderer) Decode(cursor string) (interface{}, string, error) {
	return nil, "", nil
}
func (idOrderer) New() interface{} { return nil }
