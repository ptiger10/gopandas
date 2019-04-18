package series

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestFloat32(t *testing.T) {
	s, _ := New([]float32{1, -2, 3.5, 0})
	got, _ := s.Sum()
	want := 2.5
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestFloat64(t *testing.T) {
	s, _ := New([]float64{1, -2, 3.5, 0})
	got, _ := s.Sum()
	want := 2.5
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestInt(t *testing.T) {
	s, _ := New([]int{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestInt32(t *testing.T) {
	s, _ := New([]int32{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestInt64(t *testing.T) {
	s, _ := New([]int64{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestUInt(t *testing.T) {
	s, err := New([]uint{1, 2, 3, 0})
	if err != nil {
		t.Errorf("Unable to create new series to %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestUInt32(t *testing.T) {
	s, err := New([]uint32{1, 2, 3, 0})
	if err != nil {
		t.Errorf("Unable to create new series to %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestUInt64(t *testing.T) {
	s, err := New([]uint64{1, 2, 3, 0})
	if err != nil {
		t.Errorf("Unable to create new series to %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestUnsupported(t *testing.T) {
	_, err := New([]complex64{10})
	if err == nil {
		t.Errorf("Returned nil error when constructing unsupported series type, want error")
	}
}

func TestString(t *testing.T) {
	s, err := New([]string{"low", "", "high"})
	if err != nil {
		t.Errorf("Unable to create new series to %v: %v", t.Name(), err)
	}
	_, err = s.Sum()
	if err == nil {
		t.Errorf("Returned nil error when when summing string series, want error")
	}
}

func TestBool(t *testing.T) {
	s, err := New([]bool{true, true, false, true})
	if err != nil {
		t.Errorf("Unable to create new series to %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 3.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestTime(t *testing.T) {
	s, err := New([]time.Time{
		time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC),
		time.Date(2019, 4, 19, 15, 0, 0, 0, time.UTC),
		time.Time{}})
	if err != nil {
		t.Errorf("Unable to create new series to %v: %v", t.Name(), err)
	}
	fmt.Println(s)

}

func TestInterfaceFloat(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(Float))
	if err != nil {
		t.Errorf("%v returned err, nil expected: %v", t.Name(), err)
		return
	}

	gotSum, _ := s.Sum()
	wantSum := 8.5
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}

func TestInterfaceInt(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(Int))
	if err != nil {
		t.Errorf("%v returned err, nil expected: %v", t.Name(), err)
		return
	}
	gotSum, _ := s.Sum()
	wantSum := 7.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}

func TestInterfaceString(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(String))
	if err != nil {
		t.Errorf("%v returned err, nil expected: %v", t.Name(), err)
		return
	}
	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}

func TestInterfaceBool(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(Bool))
	if err != nil {
		t.Errorf("%v returned err, nil expected: %v", t.Name(), err)
		return
	}
	gotSum, _ := s.Sum()
	wantSum := 8.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}
