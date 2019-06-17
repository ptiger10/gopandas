package dataframe

import (
	"testing"

	"github.com/ptiger10/pd/series"
)

func Test_Sum(t *testing.T) {
	c := Config{Cols: []interface{}{"fooCol", "barCol"}, Index: []string{"foo", "bar", "baz"}}
	df, err := New([]interface{}{[]float64{1, 2, 3}, []float64{4, 5, 6}}, c)
	if err != nil {
		t.Error(err)
	}
	got := df.Sum()
	want, _ := series.New([]float64{6, 15}, series.Config{Index: []string{"fooCol", "barCol"}})
	if !series.Equal(got, want) {
		t.Errorf("df.Sum() returned %v, want %v", got, want)
	}
}

func Test_Mean(t *testing.T) {
	c := Config{Cols: []interface{}{"fooCol", "barCol"}, Index: []string{"foo", "bar", "baz"}}
	df, err := New([]interface{}{[]float64{1, 2, 3}, []float64{4, 5, 6}}, c)
	if err != nil {
		t.Error(err)
	}
	got := df.Mean()
	want, _ := series.New([]float64{2, 5}, series.Config{Index: []string{"fooCol", "barCol"}})
	if !series.Equal(got, want) {
		t.Errorf("df.Mean() returned %v, want %v", got, want)
	}
}
