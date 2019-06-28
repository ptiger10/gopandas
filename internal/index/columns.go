package index

import (
	"fmt"

	"github.com/ptiger10/pd/internal/values"
)

// Columns is a collection of column levels, plus name mappings.
type Columns struct {
	NameMap LabelMap
	Levels  []ColLevel
}

// NewColumns returns a new Columns collection from a slice of column levels.
func NewColumns(levels ...ColLevel) Columns {
	if levels == nil {
		emptyLevel := NewColLevel(nil, "")
		levels = append(levels, emptyLevel)
	}
	cols := Columns{
		Levels: levels,
	}
	cols.updateNameMap()
	return cols
}

// NewColumnsFromConfig returns new Columns with default length n using a config struct.
func NewColumnsFromConfig(config Config, n int) (Columns, error) {
	var columns Columns

	// both nil: return default index
	if config.Cols == nil && config.MultiCol == nil {
		cols := NewDefaultColLevel(n, config.ColsName)
		return NewColumns(cols), nil

	}
	// both not nil: return error
	if config.Cols != nil && config.MultiCol != nil {
		return Columns{}, fmt.Errorf("columnFactory(): supplying both config.Cols and config.MultiCol is ambiguous; supply one or the other")
	}
	// single-level Columns
	if config.Cols != nil {
		newLevel := NewColLevel(config.Cols, config.ColsName)
		columns = NewColumns(newLevel)
	}

	// multi-level Columns
	if config.MultiCol != nil {
		if config.MultiColNames != nil && len(config.MultiColNames) != len(config.MultiCol) {
			return Columns{}, fmt.Errorf(
				"columnFactory(): if MultiColNames is not nil, it must must have same length as MultiCol: %d != %d",
				len(config.MultiColNames), len(config.MultiCol))
		}
		var newLevels []ColLevel
		for i := 0; i < len(config.MultiCol); i++ {
			var levelName string
			if i < len(config.MultiColNames) {
				levelName = config.MultiColNames[i]
			}
			newLevel := NewColLevel(config.MultiCol[i], levelName)
			newLevels = append(newLevels, newLevel)
		}
		columns = NewColumns(newLevels...)
	}
	return columns, nil
}

// NewDefaultColumns returns a new Columns collection with default range labels (0, 1, 2, ... n).
func NewDefaultColumns(n int) Columns {
	return NewColumns(NewDefaultColLevel(n, ""))
}

// Len returns the number of labels in every level of the column.
func (cols Columns) Len() int {
	if cols.NumLevels() == 0 {
		return 0
	}
	return cols.Levels[0].Len()
}

// NumLevels returns the number of column levels.
func (cols Columns) NumLevels() int {
	return len(cols.Levels)
}

// MaxNameWidth returns the number of characters in the column name with the most characters.
func (cols Columns) MaxNameWidth() int {
	var max int
	for k := range cols.NameMap {
		if length := len(fmt.Sprint(k)); length > max {
			max = length
		}
	}
	return max
}

// UpdateNameMap updates the holistic index map of {index level names: [index level positions]}
func (cols *Columns) updateNameMap() {
	nameMap := make(LabelMap)
	for i, lvl := range cols.Levels {
		nameMap[lvl.Name] = append(nameMap[lvl.Name], i)
	}
	cols.NameMap = nameMap
}

// Refresh updates the global name map and the label mappings at every level.
// Should be called after Series selection or index modification
func (cols *Columns) Refresh() {
	cols.updateNameMap()
	for i := 0; i < len(cols.Levels); i++ {
		cols.Levels[i].Refresh()
	}
}

// A ColLevel is a single collection of column labels within a Columns collection, plus label mappings and metadata.
// It is identical to an index Level except for the Labels, which are a simple []interface{} that do not satisfy the values.Values interface.
type ColLevel struct {
	Name     string
	Labels   []interface{}
	LabelMap LabelMap
}

// NewDefaultColLevel creates a column level with range labels (0, 1, 2, ...n) and optional name.
func NewDefaultColLevel(n int, name string) ColLevel {
	colsInt := values.MakeInterfaceRange(0, n)
	return NewColLevel(colsInt, name)
}

// NewColLevel returns a Columns level with updated label map.
func NewColLevel(labels []interface{}, name string) ColLevel {
	lvl := ColLevel{
		Labels: labels,
		Name:   name,
	}
	lvl.Refresh()
	return lvl
}

// Len returns the number of labels in the column level.
func (lvl ColLevel) Len() int {
	return len(lvl.Labels)
}

// Refresh updates all the label mappings value within a column level.
func (lvl *ColLevel) Refresh() {
	if lvl.Labels == nil {
		return
	}
	lvl.updateLabelMap()
}

// updateLabelMap updates a single level's map of {label values: [label positions]}.
// A level's label map is agnostic of the actual values in those positions.
func (lvl *ColLevel) updateLabelMap() {
	labelMap := make(LabelMap, lvl.Len())
	for i := 0; i < lvl.Len(); i++ {
		key := fmt.Sprint(lvl.Labels[i])
		labelMap[key] = append(labelMap[key], i)
	}
	lvl.LabelMap = labelMap
}

// Copy copies a Column Level
func (lvl ColLevel) Copy() ColLevel {
	lvlCopy := ColLevel{}
	lvlCopy = lvl
	lvlCopy.Labels = make([]interface{}, lvl.Len())
	for i := 0; i < lvl.Len(); i++ {
		lvlCopy.Labels[i] = lvl.Labels[i]
	}
	lvlCopy.LabelMap = make(LabelMap)
	for k, v := range lvl.LabelMap {
		lvlCopy.LabelMap[k] = v
	}
	return lvlCopy
}

// Copy returns a deep copy of each column level.
func (cols Columns) Copy() Columns {
	colsCopy := Columns{NameMap: LabelMap{}}
	for k, v := range cols.NameMap {
		colsCopy.NameMap[k] = v
	}
	for i := 0; i < len(cols.Levels); i++ {
		colsCopy.Levels = append(colsCopy.Levels, cols.Levels[i].Copy())
	}
	return colsCopy
}

// Subset returns a new Columns with all the column levels located at the specified integer positions
func (cols Columns) Subset(colPositions []int) (Columns, error) {
	cols = cols.Copy()
	var lvls []ColLevel
	for _, lvl := range cols.Levels {
		lvl, err := lvl.Subset(colPositions)
		if err != nil {
			return Columns{}, fmt.Errorf("internal columns.Subset(): %v", err)
		}
		lvls = append(lvls, lvl)
	}
	cols.Levels = lvls
	cols.updateNameMap()
	return cols, nil
}

// Subset returns the label values in a column level at specified integer positions.
func (lvl ColLevel) Subset(positions []int) (ColLevel, error) {
	var labels []interface{}
	for _, pos := range positions {
		if pos >= lvl.Len() {
			return ColLevel{}, fmt.Errorf("internal colLevel.Subset(): invalid integer position: %d (max %d)", pos, lvl.Len()-1)
		}
		labels = append(labels, lvl.Labels[pos])
	}
	lvl.Labels = labels

	lvl.Refresh()
	return lvl, nil
}