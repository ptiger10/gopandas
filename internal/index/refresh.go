package index

import (
	"fmt"

	"github.com/ptiger10/pd/opt"
)

// updateLabelMap updates a single level's map of {label values: [label positions]}.
// A level's label map is agnostic of the actual values in those positions.
func (lvl *Level) updateLabelMap() {
	labelMap := make(LabelMap, lvl.Labels.Len())
	for i := 0; i < lvl.Labels.Len(); i++ {
		key := fmt.Sprint(lvl.Labels.Element(i).Value)
		labelMap[key] = append(labelMap[key], i)
	}
	lvl.LabelMap = labelMap
}

// UpdateNameMap updates the holistic index map of {index level names: [index level positions]}
func (idx *Index) UpdateNameMap() {
	nameMap := make(LabelMap)
	for i, lvl := range idx.Levels {
		nameMap[lvl.Name] = append(nameMap[lvl.Name], i)
	}
	idx.NameMap = nameMap
}

// updateLongest finds the max length of either the level name or the longest string in the LabelMap,
// for use in printing a Series or DataFrame
func (lvl *Level) updateLongest() {
	var max int
	for k := range lvl.LabelMap {
		if len(k) > max {
			max = len(k)
		}
	}
	if len(lvl.Name) > max {
		max = len(lvl.Name)
	}
	if max > opt.GetDisplayIndexMaxWidth() {
		max = opt.GetDisplayIndexMaxWidth()
	}
	lvl.Longest = max
}

// Refresh updates the global name map and the label mappings and longest value at every level.
// Should be called after Series selection or index modification
func (idx *Index) Refresh() {
	idx.UpdateNameMap()
	for i := 0; i < len(idx.Levels); i++ {
		idx.Levels[i].Refresh()
	}
}

// Refresh updates all the label mappings and longest value within a level.
func (lvl *Level) Refresh() {
	lvl.updateLabelMap()
	lvl.updateLongest()
}
