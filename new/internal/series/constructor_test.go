package series

import (
	"log"
	"reflect"
	"testing"

	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// Calls New and panics if error. For use in testing
func mustNew(data interface{}, options ...Option) Series {
	s, err := New(data, options...)
	if err != nil {
		log.Panic(err)
	}
	return s
}

// Float Tests
// ------------------------------------------------
func TestConstructor_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5})
	if err != nil {
		t.Error(err)
	}

	_, err = New([]float64{-1.5, 0, 1.5})
	if err != nil {
		t.Error(err)
	}
}

func TestConstructor_Name_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, Name("float32"))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Name("float64"))
	if err != nil {
		t.Error(err)
	}
}

func TestConstructor_Index_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, Index([]float32{1, 2, 3}))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Index([]float64{1, 2, 3}))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Index([]float64{1}))
	if err == nil {
		t.Errorf("Returned nil error, want error due to mismatched value/index lengths")
	}
}

func TestConstructor_Kind_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Float))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Int))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.String))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Bool))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.DateTime))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Interface))
	if err == nil {
		t.Error("Returned nil error, want error due to unsupported conversion type")
	}
}

// func TestBaseConstructor_Index_Float(t *testing.T) {
// 	var err error
// 	var s Series

// 	wantName = "float32"
// 	s, err = New([]float32{-1.5, 0, 1.5}, Name(wantName))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if s.Name != wantName {
// 		t.Errorf("Returned name %v, want %v", s.Name, wantName)
// 	}

// 	wantName = "float64"
// 	s, err = New([]float64{-1.5, 0, 1.5}, Name(wantName))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if s.Name != wantName {
// 		t.Errorf("Returned name %v, want %v", s.Name, wantName)
// 	}
// }

// func TestConstructor_Float(t *testing.T) {
// 	// var tests = []struct{
// 	// 	input Series

// 	// }
// 	s, err := New([]float32{1, -2, 3.5, 0}, Name("float32"))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if s.Name != wantName {
// 		t.Errorf("Returned name %v, want %v", s.Name, wantName)
// 	}
// 	if s.Kind != wantKind {
// 		t.Errorf("Returned kind %v, want %v", s.Kind, wantKind)
// 	}

// 	// got, _ := s.Sum()
// 	// want := 2.5
// 	// if got != want {
// 	// 	t.Errorf("Sum() returned %v, want %v", got, want)
// 	// }
// 	gotName := s.Name
// 	wantName := "float32"
// 	if gotName != wantName {
// 		t.Errorf("Constructor did not read Name correctly: returned %s, want %s", gotName, wantName)
// 	}
// }

// func TestConstructor_Float64(t *testing.T) {
// 	s, _ := New([]float64{1, -2, 3.5, 0})
// 	got, _ := s.Sum()
// 	want := 2.5
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// 	gotName := s.Name
// 	wantName := ""
// 	if gotName != wantName {
// 		t.Errorf("Constructor did not create default Name correctly: returned %s, want %s", gotName, wantName)
// 	}
// }

// func TestConstructor_InterfaceFloat(t *testing.T) {
// 	s, err := New(
// 		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
// 		Type(Float))
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 		return
// 	}

// 	gotSum, _ := s.Sum()
// 	wantSum := 8.5
// 	if gotSum != wantSum {
// 		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
// 	}

// 	gotCount := s.Count()
// 	wantCount := 8
// 	if gotCount != wantCount {
// 		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
// 	}
// }

// func TestConstructor_Float_IntIndex_Default(t *testing.T) {
// 	s, _ := New([]float64{4, 5, 6})
// 	fmt.Println(s.Index.Levels[0].Labels.At(2))
// }

// // Int
// // ------------------------------------------------

// func TestConstructor_Int(t *testing.T) {
// 	s, _ := New([]int{1, -2, 3, 0})
// 	got, _ := s.Sum()
// 	want := 2.0
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// }
// func TestConstructor_Int32(t *testing.T) {
// 	s, _ := New([]int32{1, -2, 3, 0})
// 	got, _ := s.Sum()
// 	want := 2.0
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// }
// func TestConstructor_Int64(t *testing.T) {
// 	s, _ := New([]int64{1, -2, 3, 0})
// 	got, _ := s.Sum()
// 	want := 2.0
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// }

// func TestConstructor_UInt(t *testing.T) {
// 	s, err := New([]uint{1, 2, 3, 0})
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 	}
// 	got, _ := s.Sum()
// 	want := 6.0
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// }
// func TestConstructor_UInt32(t *testing.T) {
// 	s, err := New([]uint32{1, 2, 3, 0})
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 	}
// 	got, _ := s.Sum()
// 	want := 6.0
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// }
// func TestConstructor_UInt64(t *testing.T) {
// 	s, err := New([]uint64{1, 2, 3, 0})
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 	}
// 	got, _ := s.Sum()
// 	want := 6.0
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// }

// func TestConstructor_InterfaceInt(t *testing.T) {
// 	s, err := New(
// 		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
// 		Type(Int))
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 		return
// 	}
// 	gotSum, _ := s.Sum()
// 	wantSum := 7.0
// 	if gotSum != wantSum {
// 		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
// 	}

// 	gotCount := s.Count()
// 	wantCount := 8
// 	if gotCount != wantCount {
// 		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
// 	}
// }

// // String
// // ------------------------------------------------
// func TestConstructor_SliceString(t *testing.T) {
// 	_, err := New([]string{"low", "", "high"})
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 	}
// }

// func TestConstructor_InterfaceString(t *testing.T) {
// 	s, err := New(
// 		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
// 		Type(String))
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 		return
// 	}
// 	gotCount := s.Count()
// 	wantCount := 8
// 	if gotCount != wantCount {
// 		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
// 	}
// }

// // Bool
// // ------------------------------------------------
// func TestConstructor_Bool(t *testing.T) {
// 	s, err := New([]bool{true, true, false, true})
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 	}
// 	got, _ := s.Sum()
// 	want := 3.0
// 	if got != want {
// 		t.Errorf("Sum() returned %v, want %v", got, want)
// 	}
// }

// func TestConstructor_InterfaceBool(t *testing.T) {
// 	s, err := New(
// 		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
// 		Type(Bool))
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 		return
// 	}
// 	gotSum, _ := s.Sum()
// 	wantSum := 8.0
// 	if gotSum != wantSum {
// 		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
// 	}

// 	gotCount := s.Count()
// 	wantCount := 8
// 	if gotCount != wantCount {
// 		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
// 	}
// }

// // DateTime
// // ------------------------------------------------
// func TestConstructor_DateTime(t *testing.T) {
// 	s, err := New([]time.Time{
// 		time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC),
// 		time.Date(2019, 4, 19, 15, 0, 0, 0, time.UTC),
// 		time.Time{}})
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 	}
// 	_, err = s.Sum()
// 	if err == nil {
// 		t.Errorf("Returned nil error when when summing datetime series, want error")
// 	}

// 	gotCount := s.Count()
// 	wantCount := 2
// 	if gotCount != wantCount {
// 		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
// 	}

// }

// func TestConstructor_InterfaceDateTime(t *testing.T) {
// 	s, err := New(
// 		[]interface{}{
// 			time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC), "1/1/18", "Jan 1, 2018", "January 1 2018",
// 			"1pm", "1", // times cannot be parsed
// 			[]string{"1", "2"}, // slice cannot be parsed
// 			time.Location{},    // struct other than time.Time cannot be parsed
// 			nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
// 		Type(DateTime))
// 	if err != nil {
// 		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
// 		return
// 	}

// 	gotCount := s.Count()
// 	wantCount := 4
// 	if gotCount != wantCount {
// 		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
// 	}
// }

func TestMini_single(t *testing.T) {
	mini := miniIndex{
		Data: []int{1, 2, 3},
		Kind: kinds.Int,
		Name: "test",
	}
	got, err := indexFromMiniIndex([]miniIndex{mini}, 3)
	if err != nil {
		t.Error(err)
	}
	want := constructIdx.New(
		constructIdx.SliceInt([]int{1, 2, 3}, "test"),
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MiniIndex returned %v, want %v", got, want)
	}

}