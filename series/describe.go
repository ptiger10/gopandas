package series

import (
	"fmt"
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/opt"

	"github.com/ptiger10/pd/datatypes"
)

// Describe the key details of the Series.
func (s Series) Describe() {
	var err error
	// shared data
	origKind := s.datatype
	l := s.Len()
	valids := len(s.valid())
	nulls := len(s.null())
	length := fmt.Sprint(l)
	valid := fmt.Sprint(valids)
	null := fmt.Sprint(nulls)
	// type-specific data
	switch s.datatype {
	case datatypes.Float64, datatypes.Int64:
		precision := opt.GetDisplayFloatPrecision()
		mean := fmt.Sprintf("%.*f", precision, s.Math.Mean())
		min := fmt.Sprintf("%.*f", precision, s.Math.Min())
		q1 := fmt.Sprintf("%.*f", precision, s.Math.Quartile(1))
		q2 := fmt.Sprintf("%.*f", precision, s.Math.Quartile(2))
		q3 := fmt.Sprintf("%.*f", precision, s.Math.Quartile(3))
		max := fmt.Sprintf("%.*f", precision, s.Math.Max())

		values := []string{length, valid, null, mean, min, q1, q2, q3, max}
		idx := Idx([]string{"len", "valid", "null", "mean", "min", "25%", "50%", "75%", "max"})
		s, err = New(values, idx, opt.Name(s.Name))

	case datatypes.String:
		unique := fmt.Sprint(len(s.UniqueVals()))
		values := []string{length, valid, null, unique}
		idx := Idx([]string{"len", "valid", "null", "unique"})
		s, err = New(values, idx, opt.Name(s.Name))
	case datatypes.Bool:
		precision := opt.GetDisplayFloatPrecision()
		sum := fmt.Sprintf("%.*f", precision, s.Math.Sum())
		mean := fmt.Sprintf("%.*f", precision, s.Math.Mean())
		values := []string{length, valid, null, sum, mean}
		idx := Idx([]string{"len", "valid", "null", "sum", "mean"})
		s, err = New(values, idx, opt.Name(s.Name))
	case datatypes.DateTime:
		unique := fmt.Sprint(len(s.UniqueVals()))
		earliest := fmt.Sprint(s.Earliest())
		latest := fmt.Sprint(s.Latest())
		values := []string{length, valid, null, unique, earliest, latest}
		idx := Idx([]string{"len", "valid", "null", "unique", "earliest", "latest"})
		s, err = New(values, idx, opt.Name(s.Name))
	default:
		values := []string{length, valid, null}
		idx := Idx([]string{"len", "valid", "null"})
		s, err = New(values, idx, opt.Name(s.Name))
	}
	if err != nil {
		values.Warn(err, "nil (internal error)")
		return
	}
	// reset to pre-transformation Kind
	s.datatype = origKind
	fmt.Println(s)
	return
}

// ValueCounts returns a map of non-null value labels to number of occurrences in the Series.
//
// Applies to: All
func (s Series) ValueCounts() map[string]int {
	valid, _ := s.in(s.valid())
	vals := valid.all()
	counter := make(map[string]int)
	for _, val := range vals {
		counter[fmt.Sprint(val)]++
	}
	return counter
}

// UniqueVals returns a de-duplicated list all element values (as strings) that appear in the Series.
//
// Applies to: All
func (s Series) UniqueVals() []string {
	var ret []string
	counter := s.ValueCounts()
	for k := range counter {
		ret = append(ret, k)
	}
	return ret
}

// Earliest returns the earliest non-null time.Time{} in the Series
//
// Applies to: time.Time. If inapplicable, defaults to time.Time{}.
func (s Series) Earliest() time.Time {
	earliest := time.Time{}
	vals := s.validVals()
	switch s.datatype {
	case datatypes.DateTime:
		data := ensureDateTime(vals)
		for _, t := range data {
			if earliest == (time.Time{}) || t.Before(earliest) {
				earliest = t
			}
		}
		return earliest
	default:
		return earliest

	}
}

// Latest returns the latest non-null time.Time{} in the Series
//
// Applies to: time.Time. If inapplicable, defaults to time.Time{}.
func (s Series) Latest() time.Time {
	latest := time.Time{}
	vals := s.validVals()
	switch s.datatype {
	case datatypes.DateTime:
		data := ensureDateTime(vals)
		for _, t := range data {
			if latest == (time.Time{}) || t.After(latest) {
				latest = t
			}
		}
		return latest
	default:
		return latest

	}
}
