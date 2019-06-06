package series

import (
	"github.com/ptiger10/pd/kinds"
)

// Index contains index selection and conversion
type Index struct {
	s  *Series
	To To
}

// To contains conversion methods
type To struct {
	s   *Series
	idx bool
}

// Float converts Series values to float64.
func (t To) Float() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToFloat()
	} else {
		t.s.values = t.s.values.ToFloat()
		t.s.kind = kinds.Float
	}
	return *t.s
}

// Int converts Series values to int64.
func (t To) Int() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToInt()
	} else {
		t.s.values = t.s.values.ToInt()
		t.s.kind = kinds.Int
	}
	return *t.s
}

// String converts Series values to string.
func (t To) String() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToString()
	} else {
		t.s.values = t.s.values.ToString()
		t.s.kind = kinds.String
	}
	return *t.s
}

// Bool converts Series values to bool.
func (t To) Bool() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToBool()
	} else {
		t.s.values = t.s.values.ToBool()
		t.s.kind = kinds.Bool
	}
	return *t.s
}

// DateTime converts Series values to time.Time.
func (t To) DateTime() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToDateTime()
	} else {
		t.s.values = t.s.values.ToDateTime()
		t.s.kind = kinds.DateTime
	}
	return *t.s
}

// Interface converts Series values to interface.
func (t To) Interface() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToInterface()
	} else {
		t.s.values = t.s.values.ToInterface()
		t.s.kind = kinds.Interface
	}
	return *t.s
}

// // IndexTo contains conversion methods
// type IndexTo struct {
// 	s *Series
// }

// // Float copies the index then converts the first level of index values to float64.
// func (t IndexTo) Float() Series {
// 	t.s.index = t.s.index.Copy()
// 	t.s.index.Levels[0] = t.s.index.Levels[0].ToFloat()
// 	return *t.s
// }

// // Int copies the index then converts the first level of index values to int64.
// func (t IndexTo) Int() Series {
// 	t.s.index = t.s.index.Copy()
// 	t.s.index.Levels[0] = t.s.index.Levels[0].ToInt()
// 	return *t.s
// }

// // String copies the index then converts the first level of index values to string.
// func (t IndexTo) String() Series {
// 	t.s.index = t.s.index.Copy()
// 	t.s.index.Levels[0] = t.s.index.Levels[0].ToString()
// 	return *t.s
// }

// // Bool copies the index then converts the first level of index values to bool.
// func (t IndexTo) Bool() Series {
// 	t.s.index = t.s.index.Copy()
// 	t.s.index.Levels[0] = t.s.index.Levels[0].ToBool()
// 	return *t.s
// }

// // DateTime copies the index then converts the first level of index values to DateTime.
// func (t IndexTo) DateTime() Series {
// 	t.s.index = t.s.index.Copy()
// 	t.s.index.Levels[0] = t.s.index.Levels[0].ToDateTime()
// 	return *t.s
// }

// // Interface copies the index then converts the first level of index values to Interface.
// func (t IndexTo) Interface() Series {
// 	t.s.index = t.s.index.Copy()
// 	t.s.index.Levels[0] = t.s.index.Levels[0].ToInterface()
// 	return *t.s
// }

// // // IndexLevelTo converts the specific integer level of the series index to the kind supplied
// // //
// // // Applies to All. If unsupported Kind or invalid level value is supplied, returns error.
// // func (s Series) IndexLevelTo(position int, kind kinds.Kind) (Series, error) {
// // 	copyS := s.copy()
// // 	if position >= len(s.index.Levels) {
// // 		return Series{}, fmt.Errorf("unable to convert index at level %d: index out of range (Series has %d levels)", position, len(s.index.Levels))
// // 	}
// // 	lvl, err := copyS.index.Levels[position].Convert(kind)
// // 	if err != nil {
// // 		return Series{}, fmt.Errorf("unable to convert index to kind %v: unsupported kind", kind)
// // 	}
// // 	copyS.index.Levels[position] = lvl
// // 	return copyS, nil
// // }
