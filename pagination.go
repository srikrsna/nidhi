package nidhi

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	sq "github.com/elgris/sqrl"
)

const (
	seperator = ":"
)

func addPagination(st *sq.SelectBuilder, scans []any, op *PaginationOptions, orderBy []OrderBy) (*sq.SelectBuilder, []any, error) {
	if op == nil {
		for _, ob := range orderBy {
			st = st.OrderBy(ob.Field.Name() + order(ob.Desc).Direction())
		}
		return st, scans, nil
	}
	if len(orderBy) == 0 {
		orderBy = append(orderBy, OrderBy{Field: orderById{}})
	}
	var (
		orderByIdAdded bool
	)
	for i, f := range orderBy {
		od := order(f.Desc)
		if op.Backward {
			od = !od
		}
		if i == 0 {
			if strings.TrimSpace(f.Field.Name()) == ColId {
				if op.Cursor != "" {
					st = st.Where(ColId+od.Cursor(), op.Cursor)
				}
			} else {
				scans = append(scans, f.Field.New())
				st = st.Column(f.Field.Name())
				if op.Cursor != "" {
					v, id, err := f.Field.Decode(op.Cursor)
					if err != nil {
						return nil, nil, err
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
			st = st.Column(ColId)
			var id string
			scans = append(scans, &id)
		}
		if strings.TrimSpace(f.Field.Name()) == ColId {
			orderByIdAdded = true
		}
		st = st.OrderBy(f.Field.Name() + od.Direction())
	}
	if !orderByIdAdded {
		st = st.OrderBy(ColId + order(op.Backward).Direction())
	}
	return st.Limit(op.Limit + 1), scans, nil
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
	Encode(v any, id string) string
	Decode(cursor string) (any, string, error)
	New() any
}

var (
	_ Orderer = OrderByFloat("")
	_ Orderer = OrderByInt("")
	_ Orderer = OrderByString("")
	_ Orderer = OrderByTime("")
)

type orderById struct{}

func (orderById) Name() string {
	return ColId
}

func (orderById) Encode(v any, id string) string {
	return ""
}
func (orderById) Decode(cursor string) (any, string, error) {
	return nil, "", nil
}
func (orderById) New() any { return nil }

type OrderByInt string

func (i OrderByInt) Name() string {
	return `COALESCE(` + string(i) + `, 0)`
}

func (OrderByInt) Encode(v any, id string) string {
	return base64.URLEncoding.EncodeToString([]byte(strconv.FormatInt(*v.(*int64), 10) + seperator + id))
}
func (OrderByInt) Decode(cursor string) (any, string, error) {
	dataBytes, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, "", InvalidCursor
	}

	data := string(dataBytes)
	splits := strings.SplitN(data, seperator, 2)
	if len(splits) != 2 {
		return nil, "", InvalidCursor
	}

	v, err := strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return nil, "", InvalidCursor
	}

	return v, splits[1], nil
}
func (OrderByInt) New() any { return Ptr[int64](0) }

type OrderByFloat string

func (i OrderByFloat) Name() string {
	return `COALESCE(` + string(i) + `, 0)`
}

func (OrderByFloat) Encode(v any, id string) string {
	return base64.URLEncoding.EncodeToString([]byte(strconv.FormatFloat(*v.(*float64), 'g', 2, 64) + seperator + id))
}
func (OrderByFloat) Decode(cursor string) (any, string, error) {
	dataBytes, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, "", InvalidCursor
	}

	data := string(dataBytes)
	splits := strings.SplitN(data, seperator, 2)
	if len(splits) != 2 {
		return nil, "", InvalidCursor
	}

	v, err := strconv.ParseFloat(splits[0], 64)
	if err != nil {
		return nil, "", InvalidCursor
	}

	return v, splits[1], nil
}
func (OrderByFloat) New() any { return Ptr[float64](0) }

type OrderByString string

func (i OrderByString) Name() string {
	return `COALESCE(` + string(i) + `, '')`
}

func (OrderByString) Encode(v any, id string) string {
	return base64.URLEncoding.EncodeToString([]byte(base64.URLEncoding.EncodeToString([]byte(*v.(*string))) + seperator + id))
}
func (OrderByString) Decode(cursor string) (any, string, error) {
	dataBytes, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, "", InvalidCursor
	}

	data := string(dataBytes)
	splits := strings.SplitN(data, seperator, 2)
	if len(splits) != 2 {
		return nil, "", InvalidCursor
	}

	v, err := base64.URLEncoding.DecodeString(splits[0])
	if err != nil {
		return nil, "", InvalidCursor
	}

	return string(v), splits[1], nil
}
func (OrderByString) New() any { return Ptr("") }

type OrderByTime string

func (i OrderByTime) Name() string {
	return `COALESCE(` + string(i) + `, '1970-01-01T00:00:00Z'::timestamp)`
}

func (OrderByTime) Encode(v any, id string) string {
	return base64.URLEncoding.EncodeToString([]byte(base64.URLEncoding.EncodeToString([]byte((*v.(*time.Time)).Format(time.RFC3339))) + seperator + id))
}
func (OrderByTime) Decode(cursor string) (any, string, error) {
	dataBytes, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, "", InvalidCursor
	}

	data := string(dataBytes)
	splits := strings.SplitN(data, seperator, 2)
	if len(splits) != 2 {
		return nil, "", InvalidCursor
	}

	v, err := base64.URLEncoding.DecodeString(splits[0])
	if err != nil {
		return nil, "", InvalidCursor
	}

	tv, err := time.Parse(time.RFC3339, string(v))
	if err != nil {
		return nil, "", InvalidCursor
	}

	return tv, splits[1], nil
}
func (OrderByTime) New() any { return Ptr(time.Time{}) }
