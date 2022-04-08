package indexer

import (
	"time"

	"github.com/tendermint/tendermint/libs/pubsub/query/syntax"
)

// QueryRanges defines a mapping between a composite event key and a QueryRange.
//
// e.g.account.number => queryRange{lowerBound: 1, upperBound: 5}
type QueryRanges map[string]QueryRange

// QueryRange defines a range within a query condition.
type QueryRange struct {
	LowerBound        interface{} // int || time.Time
	UpperBound        interface{} // int || time.Time
	Key               string
	IncludeLowerBound bool
	IncludeUpperBound bool
}

// AnyBound returns either the lower bound if non-nil, otherwise the upper bound.
func (qr QueryRange) AnyBound() interface{} {
	if qr.LowerBound != nil {
		return qr.LowerBound
	}

	return qr.UpperBound
}

// LowerBoundValue returns the value for the lower bound. If the lower bound is
// nil, nil will be returned.
func (qr QueryRange) LowerBoundValue() interface{} {
	if qr.LowerBound == nil {
		return nil
	}

	if qr.IncludeLowerBound {
		return qr.LowerBound
	}

	switch t := qr.LowerBound.(type) {
	case int64:
		return t + 1

	case time.Time:
		return t.Unix() + 1

	default:
		panic("not implemented")
	}
}

// UpperBoundValue returns the value for the upper bound. If the upper bound is
// nil, nil will be returned.
func (qr QueryRange) UpperBoundValue() interface{} {
	if qr.UpperBound == nil {
		return nil
	}

	if qr.IncludeUpperBound {
		return qr.UpperBound
	}

	switch t := qr.UpperBound.(type) {
	case int64:
		return t - 1

	case time.Time:
		return t.Unix() - 1

	default:
		panic("not implemented")
	}
}

// LookForRanges returns a mapping of QueryRanges and the matching indexes in
// the provided query conditions.
func LookForRanges(conditions []syntax.Condition) (ranges QueryRanges, indexes []int) {
	ranges = make(QueryRanges)
	for i, c := range conditions {
		if IsRangeOperation(c.Op) {
			r, ok := ranges[c.Tag]
			if !ok {
				r = QueryRange{Key: c.Tag}
			}

			switch c.Op {
			case syntax.TGt:
				r.LowerBound = c.Op

			case syntax.TGeq:
				r.IncludeLowerBound = true
				r.LowerBound = c.Op

			case syntax.TLt:
				r.UpperBound = c.Op

			case syntax.TLeq:
				r.IncludeUpperBound = true
				r.UpperBound = c.Op
			}

			ranges[c.Tag] = r
			indexes = append(indexes, i)
		}
	}

	return ranges, indexes
}

// IsRangeOperation returns a boolean signifying if a query Operator is a range
// operation or not.
func IsRangeOperation(op syntax.Token) bool {
	switch op {
	case syntax.TGt, syntax.TGeq, syntax.TLt, syntax.TLeq:
		return true

	default:
		return false
	}
}
