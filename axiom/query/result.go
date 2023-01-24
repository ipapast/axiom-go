package query

import (
	"encoding/json"
	"time"
)

// Result is the result of an APL query.
type Result struct {
	// Tables in the query result.
	Tables []Table `json:"tables"`
	// Status of the query result.
	Status Status `json:"status"`
	// TraceID is the ID of the trace that was generated by the server for this
	// results query request.
	TraceID string `json:"-"`
}

// Table in the [Result] of an APL query.
type Table struct {
	// Name of the table. Default name for unnamed results is "0", "1", "2", ...
	// etc.
	Name string `json:"name"`
	// Sources are the datasets that were consulted in order to create the
	// table.
	Sources []Source `json:"sources"`
	// Fields in the table matching the order of the [Columns] (e.g. the
	// [Column] at index 0 has the values for the [Field] at index 0).
	Fields []Field `json:"fields"`
	// Order of the fields in the table.
	Order []Order `json:"order"`
	// Groups are the groups of the table.
	Groups []Group `json:"groups"`
	// Range specifies the window the query was restricted to. Nil if the query
	// was not restricted to a time window.
	Range *RangeInfo `json:"range"`
	// Buckets defines if the query is bucketed (usually on the "_time" field).
	// Nil if the query returns a non-bucketed result.
	Buckets *BucketInfo `json:"buckets"`
	// Columns in the table matching the order of the [Fields] (e.g. the
	// [Column] at index 0 has the values for the [Field] at index 0). In case
	// of sub-groups, rows will repeat the group value.
	Columns []Column `json:"columns"`
}

// Field in a [Table].
type Field struct {
	// Name of the field.
	Name string `json:"name"`
	// Type of the field. Can also be composite types which are types separated
	// by a horizontal line "|".
	Type string `json:"type"`
	// Aggregation is the aggregation applied to the field.
	Aggregation Aggregation `json:"agg"`
}

// Aggregation that is applied to a [Field] in a [Table].
type Aggregation struct {
	// Name of the aggregation.
	Name string `json:"name"`
	// Args are the arguments of the aggregation.
	Args []any `json:"args"`
}

// Source that was consulted in order to create a [Table].
type Source struct {
	// Name of the source.
	Name string `json:"name"`
}

// Order of a [Field] in a [Table].
type Order struct {
	// Field is the name of the field to order by.
	Field string `json:"field"`
	// Desc is true if the order is descending. Otherwise the order is
	// ascending.
	Desc bool `json:"desc"`
}

// Group in a [Table].
type Group struct {
	// Name of the group.
	Name string `json:"name"`
}

// RangeInfo specifies the window a query was restricted to.
type RangeInfo struct {
	// Field specifies the field name on which the query range was restricted.
	// Usually "_time":
	Field string
	// Start is the starting time the query is limited by. Usually the start of
	// the time window. Queries are restricted to the interval [start,end).
	Start time.Time
	// End is the ending time the query is limited by. Usually the end of the
	// time window. Queries are restricted to the interval [start,end).
	End time.Time
}

// BucketInfo captures information about how a grouped query is sorted into
// buckets. Usually buckets are created on the "_time" column,
type BucketInfo struct {
	// Field specifies the field used to create buckets on. Usually this would
	// be "_time".
	Field string
	// An integer or float representing the fixed bucket size.
	// When the bucket field is "_time" this value is in nanoseconds.
	Size any
}

// Column in a [Table] containing the raw values of a [Field].
type Column []any

// Status of an APL query [Result].
type Status struct {
	// MinCursor is the id of the oldest row, as seen server side. May be lower
	// than what the results include if the server scanned more data than
	// included in the results. Can be used to efficiently resume time-sorted
	// non-aggregating queries (i.e. filtering only).
	MinCursor string `json:"minCursor"`
	// MaxCursor is the id of the newest row, as seen server side. May be higher
	// than what the results include if the server scanned more data than
	// included in the results. Can be used to efficiently resume time-sorted
	// non-aggregating queries (i.e. filtering only).
	MaxCursor string `json:"maxCursor"`
	// ElapsedTime is the duration it took the query to execute.
	ElapsedTime time.Duration `json:"elapsedTime"`
	// RowsExamined is the amount of rows that have been examined by the query.
	RowsExamined uint64 `json:"rowsExamined"`
	// RowsMatched is the amount of rows that matched the query.
	RowsMatched uint64 `json:"rowsMatched"`
}

// UnmarshalJSON implements [json.Unmarshaler]. It is in place to unmarshal the
// elapsed time into a proper [time.Duration] value because the server returns
// it in microseconds.
func (s *Status) UnmarshalJSON(b []byte) error {
	type localStatus *Status

	if err := json.Unmarshal(b, localStatus(s)); err != nil {
		return err
	}

	// Set to a proper [time.Duration] value by interpreting the server response
	// value in microseconds.
	s.ElapsedTime *= time.Microsecond

	return nil
}
