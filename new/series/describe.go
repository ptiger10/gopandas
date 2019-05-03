package series

import (
	"fmt"
	"log"
	"time"

	"github.com/ptiger10/pd/new/kinds"
)

// Len returns the length of the Series (including null values)
//
// Applies to: All
func (s Series) Len() int {
	all := s.values.All()
	return len(all)
}

// Describe the key details of the Series
//
// Applies to: All
func (s Series) Describe() {
	var err error
	// shared data
	origKind := s.Kind
	l := s.Len()
	valids := len(s.values.Valid())
	nulls := len(s.values.Null())
	length := fmt.Sprint(l)
	valid := fmt.Sprint(valids)
	null := fmt.Sprint(nulls)
	// type-specific data
	switch s.Kind {
	case kinds.Float, kinds.Int:
		precision := 4
		mean := fmt.Sprintf("%.*f", precision, s.Mean())
		min := fmt.Sprintf("%.*f", precision, s.Min())
		q1 := fmt.Sprintf("%.*f", precision, s.Quartile(1))
		q2 := fmt.Sprintf("%.*f", precision, s.Quartile(2))
		q3 := fmt.Sprintf("%.*f", precision, s.Quartile(3))
		max := fmt.Sprintf("%.*f", precision, s.Max())

		values := []string{length, valid, null, mean, min, q1, q2, q3, max}
		idx := Index([]string{"len", "valid", "null", "mean", "min", "25%", "50%", "75%", "max"})
		s, err = New(values, idx, Name("description"))

	case kinds.String:
		// value counts
		unique := fmt.Sprint(len(s.Unique()))
		values := []string{length, valid, null, unique}
		idx := Index([]string{"len", "valid", "null", "unique"})
		s, err = New(values, idx, Name("description"))
	case kinds.DateTime:
		unique := fmt.Sprint(len(s.Unique()))
		earliest := fmt.Sprint(s.Earliest())
		latest := fmt.Sprint(s.Latest())
		values := []string{length, valid, null, unique, earliest, latest}
		idx := Index([]string{"len", "valid", "null", "unique", "earliest", "latest"})
		s, err = New(values, idx, Name("description"))
	default:
		values := []string{length, valid, null}
		idx := Index([]string{"len", "valid", "null"})
		s, err = New(values, idx, Name("description"))
	}
	if err != nil {
		log.Printf("Internal error: s.Describe() could not construct Series: %v\nPlease open a Github issue.\n", err)
		return
	}
	// reset to pre-transformation Kind
	s.Kind = origKind
	fmt.Println(s)
	return
}

// ValueCounts returns a map of non-null value labels to number of occurrences in the Series.
//
// Applies to: All
func (s Series) ValueCounts() map[string]int {
	vals := s.validAll()
	counter := make(map[string]int)
	for _, val := range vals {
		counter[fmt.Sprint(val)]++
	}
	return counter
}

// Unique returns a de-duplicated list all element values (as strings) that appear in the Series.
//
// Applies to: All
func (s Series) Unique() []string {
	var ret []string
	counter := s.ValueCounts()
	for k := range counter {
		ret = append(ret, k)
	}
	return ret
}

func ensureDateTime(vals interface{}) []time.Time {
	if datetime, ok := vals.([]time.Time); ok {
		return datetime
	}
	log.Printf("Internal error: ensureDateTime has received an unallowable value: %v", vals)
	return nil
}

// Earliest returns the earliest non-null time.Time{} in the Series
//
// Applies to: time.Time. If inapplicable, defaults to time.Time{}.
func (s Series) Earliest() time.Time {
	earliest := time.Time{}
	vals := s.validVals()
	switch s.Kind {
	case kinds.DateTime:
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
	switch s.Kind {
	case kinds.DateTime:
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
