package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/options"
)

func TestDataFrame_Describe(t *testing.T) {
	type want struct {
		len          int
		numIdxLevels int
		maxWidth     int
		values       []interface{}
		datatype     string
		name         string
	}
	tests := []struct {
		name  string
		input *Series
		want  want
	}{
		{name: "default index",
			input: MustNew("foo"),
			want: want{len: 1, numIdxLevels: 1, maxWidth: 3,
				values: []interface{}{"foo"}, datatype: "string", name: ""}},
		{"multi index",
			MustNew(1.0, Config{MultiIndex: []interface{}{"baz", "qux"}, Name: "foo"}),
			want{len: 1, numIdxLevels: 2, maxWidth: 1,
				values: []interface{}{1.0}, datatype: "float64", name: "foo"}},
		{"empty",
			newEmptySeries(),
			want{len: 0, numIdxLevels: 0, maxWidth: 0, datatype: "none"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input.Copy()
			gotLen := s.Len()
			if gotLen != tt.want.len {
				t.Errorf("s.Len(): got %v, want %v", gotLen, tt.want.len)
			}
			gotNumIdxLevels := s.Index.NumLevels()
			if gotNumIdxLevels != tt.want.numIdxLevels {
				t.Errorf("s.Index.NumLevels(): got %v, want %v", gotNumIdxLevels, tt.want.numIdxLevels)
			}
			gotMaxWidth := s.MaxWidth()
			if gotMaxWidth != tt.want.maxWidth {
				t.Errorf("s.MaxWidth(): got %v, want %v", gotMaxWidth, tt.want.maxWidth)
			}
			gotValues := s.Values()
			if !reflect.DeepEqual(gotValues, tt.want.values) {
				t.Errorf("s.Values(): got %v, want %v", gotMaxWidth, tt.want.values)
			}
			gotDatatype := s.DataType()
			if gotDatatype != tt.want.datatype {
				t.Errorf("s.Datatype(): got %v, want %v", gotDatatype, tt.want.datatype)
			}
			gotName := s.Name()
			if gotName != tt.want.name {
				t.Errorf("s.Name(): got %v, want %v", gotName, tt.want.name)
			}

		})
	}
}

func TestSeries_DataType(t *testing.T) {
	var tests = []struct {
		datatype options.DataType
		expected string
	}{

		{options.None, "none"},
		{options.Float64, "float64"},
		{options.Int64, "int64"},
		{options.String, "string"},
		{options.Bool, "bool"},
		{options.DateTime, "dateTime"},
		{options.Interface, "interface"},
		{options.Unsupported, "unsupported"},
		{-1, "unknown"},
		{100, "unknown"},
	}
	for _, test := range tests {
		s, _ := New([]int{1, 2, 3})
		s.datatype = test.datatype
		if s.DataType() != test.expected {
			t.Errorf("s.Datatype() for datatype %v returned %v, want %v", test.datatype, test.datatype.String(), test.expected)
		}
	}
}

func TestSeries_Equal(t *testing.T) {
	s, err := New("foo", Config{Index: "bar", Name: "baz"})
	if err != nil {
		t.Error(err)
	}
	s2, _ := New("foo", Config{Index: "bar", Name: "baz"})
	if !Equal(s, s2) {
		t.Errorf("Equal() returned false, want true")
	}
	s2.datatype = options.Bool
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different kind, want false")
	}

	s2, _ = New("quux", Config{Index: "bar", Name: "baz"})
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different values, want false")
	}
	s2, _ = New("foo", Config{Index: "corge", Name: "baz"})
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different index, want false")
	}
	s2, _ = New("foo", Config{Index: "bar", Name: "qux"})
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different name, want false")
	}
}

func TestSeries_ReplaceNil(t *testing.T) {
	s := MustNew(nil)
	s2 := MustNew([]int{1, 2})
	s.replace(s2)
	if !Equal(s, s2) {
		t.Errorf("Series.replace() returned %v, want %v", s, s2)
	}
}

func TestSeries_Describe_unsupported(t *testing.T) {
	s := MustNew([]float64{1, 2, 3})
	tm := s.Earliest()
	if (time.Time{}) != tm {
		t.Errorf("Earliest() got %v, want time.Time{} for unsupported type", tm)
	}
	tm = s.Latest()
	if (time.Time{}) != tm {
		t.Errorf("Latest() got %v, want time.Time{} for unsupported type", tm)
	}
}

// [START ensure tests]
func TestSeries_EnsureTypes_fail(t *testing.T) {
	defer log.SetOutput(os.Stderr)
	vals := []interface{}{1, 2, 3}

	var buf bytes.Buffer
	log.SetOutput(&buf)
	ensureFloatFromNumerics(vals)
	if buf.String() == "" {
		t.Errorf("ensureNumerics() returned no log message, want log due to fail")
	}
	buf.Reset()

	ensureDateTime(vals)
	if buf.String() == "" {
		t.Errorf("ensureDateTime() returned no log message, want log due to fail")
	}
	buf.Reset()

	ensureBools(vals)
	if buf.String() == "" {
		t.Errorf("ensureBools() returned no log message, want log due to fail")
	}
	buf.Reset()
}

// [END ensure tests]
