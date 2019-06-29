package dataframe

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// A DataFrame is a 2D collection of one or more Series with a shared index and associated columns.
type DataFrame struct {
	name  string
	s     []*series.Series
	cols  index.Columns
	index index.Index
}

// Len returns the number of values in each Series of the DataFrame.
func (df *DataFrame) Len() int {
	if df.s == nil {
		return 0
	}
	return df.s[0].Len()
}

// Name returns the DataFrame's name.
func (df *DataFrame) Name() string {
	return df.name
}

// Rename the DataFrame.
func (df *DataFrame) Rename(name string) {
	df.name = name
}

// NumCols returns the number of columns in the DataFrame.
func (df *DataFrame) NumCols() int {
	if df.s == nil {
		return 0
	}
	return len(df.s)
}

// IndexLevels returns the number of index levels in the DataFrame.
func (df *DataFrame) IndexLevels() int {
	return df.index.NumLevels()
}

// ColLevels returns the number of column levels in the DataFrame.
func (df *DataFrame) ColLevels() int {
	return df.cols.NumLevels()
}

// Copy creates a new deep copy of a Series.
func (df *DataFrame) Copy() *DataFrame {
	var sCopy []*series.Series
	for i := 0; i < len(df.s); i++ {
		sCopy = append(sCopy, df.s[i].Copy())
	}
	idxCopy := df.index.Copy()
	colsCopy := df.cols.Copy()
	dfCopy := &DataFrame{
		s:     sCopy,
		index: idxCopy,
		cols:  colsCopy,
		name:  df.name,
	}
	// dfCopy.Apply = Apply{s: copyS}
	// dfCopy.Index = Index{s: copyS}
	// dfCopy.InPlace = InPlace{s: copyS}
	return dfCopy
}

func (df *DataFrame) selectByRows(rowPositions []int) (*DataFrame, error) {
	if err := df.ensureAlignment(); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe internal alignment error: %v", err)
	}
	idx := df.index.Subset(rowPositions)
	var seriesSlice []*series.Series
	for i := 0; i < df.NumCols(); i++ {
		s, err := df.s[i].Subset(rowPositions)
		if err != nil {
			return newEmptyDataFrame(), fmt.Errorf("dataframe.selectByRows(): %v", err)
		}
		seriesSlice = append(seriesSlice, s)
	}
	df = newFromComponents(seriesSlice, idx, df.cols, df.name)
	return df, nil
}

func (df *DataFrame) selectByCols(colPositions []int) (*DataFrame, error) {
	if err := df.ensureAlignment(); err != nil {
		return df, fmt.Errorf("dataframe internal alignment error: %v", err)
	}
	var seriesSlice []*series.Series
	for _, pos := range colPositions {
		if pos >= df.NumCols() {
			return nil, fmt.Errorf("dataframe.selectByCols(): invalid col position %d (max: %d)", pos, df.NumCols()-1)
		}
		seriesSlice = append(seriesSlice, df.s[pos])
	}
	columnsSlice, err := df.cols.Subset(colPositions)
	if err != nil {
		return nil, fmt.Errorf("dataframe.selectByCols(): %v", err)
	}

	df = newFromComponents(seriesSlice, df.index, columnsSlice, df.name)
	return df, nil
}

// Equal returns true if two dataframes contain equivalent values.
func Equal(df, df2 *DataFrame) bool {
	if df.NumCols() != df2.NumCols() {
		return false
	}
	for i := 0; i < df.NumCols(); i++ {
		if !series.Equal(df.s[i], df2.s[i]) {
			return false
		}
	}
	if !reflect.DeepEqual(df.index, df2.index) {
		return false
	}
	if !reflect.DeepEqual(df.cols, df2.cols) {
		return false
	}
	if df.name != df2.name {
		return false
	}
	return true
}

// Col returns the first Series with the specified column label at column level 0.
func (df *DataFrame) Col(label string) *series.Series {
	colPos, ok := df.cols.Levels[0].LabelMap[label]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("df.Col(): invalid column label: %v not in labels", label)
		}
		s, _ := series.New(nil)
		return s
	}
	df, _ = df.selectByCols(colPos)
	return df.s[0]
}

// DataTypes returns the DataTypes of the Series in the DataFrame.
func (df *DataFrame) DataTypes() *series.Series {
	var vals []interface{}
	var idx []interface{}
	for _, s := range df.s {
		vals = append(vals, s.DataType())
		idx = append(idx, s.Name())
	}
	s, err := newSingleIndexSeries(vals, idx, "datatypes")
	if err != nil {
		log.Printf("DataTypes(): %v", err)
		return nil
	}
	return s
}

// dataType is the data type of the DataFrame's values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (df *DataFrame) dataType() string {
	uniqueTypes := df.DataTypes().Unique()
	if len(uniqueTypes) == 1 {
		return df.s[0].DataType()
	}
	return "mixed"
}

// maxColWidths is the max characters in each column of a dataframe.
// exclusions should mimic the shape of the columns exactly
func (df *DataFrame) maxColWidths(exclusions [][]bool) []int {
	var maxColWidths []int
	if len(exclusions) != df.ColLevels() || len(exclusions) == 0 {
		return nil
	}
	if len(exclusions[0]) != df.NumCols() {
		return nil
	}
	for k := 0; k < df.NumCols(); k++ {
		max := df.s[k].MaxWidth()
		for j := 0; j < df.ColLevels(); j++ {
			if !exclusions[j][k] {
				if length := len(fmt.Sprint(df.cols.Levels[j].Labels[k])); length > max {
					max = length
				}
			}
		}
		maxColWidths = append(maxColWidths, max)
	}
	return maxColWidths
}

func (df *DataFrame) makeExclusionsTable() [][]bool {
	table := make([][]bool, df.ColLevels())
	for row := range table {
		table[row] = make([]bool, df.NumCols())
	}
	return table
}

// Subset a DataFrame to include only the rows at supplied integer positions.
func (df *DataFrame) Subset(rowPositions []int) (*DataFrame, error) {
	if len(rowPositions) == 0 {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.Subset(): no valid rows provided")
	}

	sub, err := df.selectByRows(rowPositions)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.Subset(): %v", err)
	}
	return sub, nil
}
