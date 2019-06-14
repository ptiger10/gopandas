package index

import (
	"fmt"
)

// updateLabelMap updates a single level's map of {label values: [label positions]}.
// A level's label map is agnostic of the actual values in those positions.
func (lvl *Level) updateLabelMap() {
	labelMap := make(LabelMap, lvl.Len())
	for i := 0; i < lvl.Len(); i++ {
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

// Refresh updates the global name map and the label mappings at every level.
// Should be called after Series selection or index modification
func (idx *Index) Refresh() {
	if idx.Len() == 0 {
		return
	}
	idx.UpdateNameMap()
	for i := 0; i < len(idx.Levels); i++ {
		idx.Levels[i].Refresh()
	}
}

// Refresh updates all the label mappings value within a level.
func (lvl *Level) Refresh() {
	if lvl.Labels == nil {
		return
	}
	lvl.updateLabelMap()
}
