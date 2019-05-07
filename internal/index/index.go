package index

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// An Index is a collection of levels, plus label mappings
type Index struct {
	Levels  []Level
	NameMap LabelMap
}

// A Level is a single collection of labels within an index, plus label mappings and metadata
type Level struct {
	Kind     kinds.Kind
	Labels   values.Values
	LabelMap LabelMap
	Name     string
	Longest  int
}

// A LabelMap records the position of labels, in the form {label name: [label position(s)]}
type LabelMap map[string][]int

// In returns an index with only those levels located at specified integer positions
func (idx Index) In(positions []int) (Index, error) {
	idx = idx.Copy()
	var lvls []Level
	for _, pos := range positions {
		if pos >= len(idx.Levels) {
			return Index{}, fmt.Errorf("error indexing index levels: level %d is out of range", pos)
		}
		lvls = append(lvls, idx.Levels[pos])
	}
	newIdx := New(lvls...)
	return newIdx, nil
}

// Copy returns a deep copy of each index level
func (idx Index) Copy() Index {
	idxCopy := Index{}
	copier.Copy(&idxCopy, &idx)
	idxCopy.Levels = make([]Level, len(idx.Levels))
	for i := 0; i < len(idx.Levels); i++ {
		copier.Copy(&idxCopy.Levels[i], &idx.Levels[i])
	}
	return idxCopy
}

// Len returns the length of the first level of the index
func (idx Index) Len() int {
	return idx.Levels[0].Labels.Len()
}

// Aligned ensures that all index levels have the same length
func (idx Index) Aligned() bool {
	lvl0 := idx.Levels[0].Labels.Len()
	for i := 1; i < len(idx.Levels); i++ {
		if lvl0 != idx.Levels[i].Labels.Len() {
			return false
		}
	}
	return true
}
