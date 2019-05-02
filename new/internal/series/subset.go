package series

import "github.com/ptiger10/pd/new/internal/values"

// At subsets a Series by integer position
func (s Series) At(position int) Series {
	s = s.at(position)
	s.index.Refresh()
	return s
}

func (s Series) at(position int) Series {
	positions := []int{position}
	s.values = s.values.In(positions).(values.Values)
	for i, level := range s.index.Levels {
		s.index.Levels[i].Labels = level.Labels.In(positions).(values.Values)
	}
	return s
}