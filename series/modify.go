package series

import (
	"fmt"
	"log"
	"sort"

	"github.com/ptiger10/pd/options"
)

// [START Sort interface]
func (s *Series) Swap(i, j int) {
	s.values.Swap(i, j)
	for lvl := 0; lvl < s.index.Len(); lvl++ {
		s.index.Levels[lvl].Labels.Swap(i, j)
		s.index.Levels[lvl].Refresh()
	}
}

func (s *Series) Less(i, j int) bool {
	return s.values.Less(i, j)
}

// [END Sort interface]

// [START return new Series]

// Sort sorts the series by its values and returns a new Series.
func (s *Series) Sort(asc bool) *Series {
	s = s.Copy()
	s.InPlace.Sort(asc)
	return s
}

// Insert inserts a new row into the Series immediately before the specified integer position and returns a new Series.
func (s *Series) Insert(pos int, val interface{}, idx []interface{}) (*Series, error) {
	s = s.Copy()
	s.InPlace.Insert(pos, val, idx)
	return s, nil
}

// dropRows drops multiple rows and returns a new Series
func (s *Series) dropRows(positions []int) (*Series, error) {
	s = s.Copy()
	s.InPlace.dropRows(positions)
	return s, nil
}

// Drop drops the row at the specified integer position and returns a new Series.
func (s *Series) Drop(pos int) (*Series, error) {
	s = s.Copy()
	s.InPlace.Drop(pos)
	return s, nil
}

// DropNull drops all null values and modifies the Series in place.
func (s *Series) DropNull() *Series {
	s = s.Copy()
	s.InPlace.DropNull()
	return s
}

// Append adds a row at the end and returns a new Series.
func (s *Series) Append(val interface{}, idx []interface{}) *Series {
	s, _ = s.Insert(s.Len(), val, idx)
	return s
}

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and returns a new Series.
func (s *Series) Join(s2 *Series) *Series {
	s = s.Copy()
	s.InPlace.Join(s2)
	return s
}

// [END return new Series]

// [START modify in place]

// InPlace contains methods for modifying a Series in place.
type InPlace struct {
	s *Series
}

// Sort sorts the series by its values and modifies the Series in place.
func (ip InPlace) Sort(asc bool) {
	if asc {
		sort.Stable(ip.s)
	} else {
		sort.Stable(sort.Reverse(ip.s))
	}
}

// Insert inserts a new row into the Series immediately before the specified integer position and modifies the Series in place.
func (ip InPlace) Insert(pos int, val interface{}, idx []interface{}) error {
	if len(idx) != ip.s.index.Len() {
		return fmt.Errorf("Series.Insert() len(idx) must equal number of index levels: supplied %v want %v",
			len(idx), ip.s.index.Len())
	}
	for i := 0; i < ip.s.index.Len(); i++ {
		err := ip.s.index.Levels[i].Labels.Insert(pos, idx[i])
		if err != nil {
			return fmt.Errorf("Series.Insert() with idx val %v at idx level %v: %v", val, i, err)
		}
		ip.s.index.Levels[i].Refresh()
	}
	if err := ip.s.values.Insert(pos, val); err != nil {
		return fmt.Errorf("Series.Insert() with val %v: %v", val, err)
	}
	return nil
}

// dropRows drops multiple rows
func (ip InPlace) dropRows(positions []int) error {
	sort.IntSlice(positions).Sort()
	for i, position := range positions {
		err := ip.s.InPlace.Drop(position - i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Drop drops a row at a specified integer position and modifies the Series in place.
func (ip InPlace) Drop(pos int) error {
	for i := 0; i < ip.s.index.Len(); i++ {
		err := ip.s.index.Levels[i].Labels.Drop(pos)
		if err != nil {
			return fmt.Errorf("Series.Drop(): %v", err)
		}
		ip.s.index.Levels[i].Refresh()
	}
	if err := ip.s.values.Drop(pos); err != nil {
		return fmt.Errorf("Series.Drop(): %v", err)
	}
	return nil
}

// DropNull drops all null values and modifies the Series in place.
func (ip InPlace) DropNull() {
	ip.dropRows(ip.s.null())
	return
}

// Append adds a row at a specified integer position and modifies the Series in place.
func (ip InPlace) Append(val interface{}, idx []interface{}) {
	_ = ip.s.InPlace.Insert(ip.s.Len(), val, idx)
	return
}

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and modifies s in place.
func (ip InPlace) Join(s2 *Series) {
	if ip.s.datatype == options.None {
		ip.s.replace(s2)
		return
	}
	for i := 0; i < s2.Len(); i++ {
		elem := s2.Element(i)
		ip.s.InPlace.Append(elem.Value, elem.Labels)
	}
}

// [END modify in place]

// [START type conversion methods]

// ToFloat64 converts Series values to float64 and returns a new Series.
func (s *Series) ToFloat64() *Series {
	s = s.Copy()
	s.values = s.values.ToFloat64()
	s.datatype = options.Float64
	return s
}

// ToInt64 converts Series values to int64 and returns a new Series.
func (s *Series) ToInt64() *Series {
	s = s.Copy()
	s.values = s.values.ToInt64()
	s.datatype = options.Int64

	return s
}

// ToString converts Series values to string and returns a new Series.
func (s *Series) ToString() *Series {
	s = s.Copy()
	s.values = s.values.ToString()
	s.datatype = options.String
	return s
}

// ToBool converts Series values to bool and returns a new Series.
func (s *Series) ToBool() *Series {
	s = s.Copy()
	s.values = s.values.ToBool()
	s.datatype = options.Bool
	return s
}

// ToDateTime converts Series values to time.Time and returns a new Series.
func (s *Series) ToDateTime() *Series {
	s = s.Copy()
	s.values = s.values.ToDateTime()
	s.datatype = options.DateTime
	return s
}

// ToInterface converts Series values to interface and returns a new Series.
func (s *Series) ToInterface() *Series {
	s = s.Copy()
	s.values = s.values.ToInterface()
	s.datatype = options.Interface
	return s
}

// [END type conversion methods]

// [START Index modifications]

// Index contains index selection and conversion
type Index struct {
	s *Series
}

// Levels returns the number of levels in the index
func (idx Index) Levels() int {
	return idx.s.index.Len()
}

// Len returns the number of items in each level of the index.
func (idx Index) Len() int {
	if len(idx.s.index.Levels) == 0 {
		return 0
	}
	return idx.s.index.Levels[0].Len()
}

// Swap swaps two labels at index level 0 and modifies the index in place.
func (idx Index) Swap(i, j int) {
	idx.s.Swap(i, j)
}

// Less compares two elements and returns true if the first is less than the second.
func (idx Index) Less(i, j int) bool {
	return idx.s.index.Levels[0].Labels.Less(i, j)
}

// Sort sorts the index by index level 0 and modifies the Series in place.
func (idx Index) Sort(asc bool) {
	if asc {
		sort.Stable(idx)
	} else {
		sort.Stable(sort.Reverse(idx))
	}
}

// At returns the index value at a specified index level and integer position.
func (idx Index) At(position int, level int) (interface{}, error) {
	if level >= idx.s.index.Len() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.Len())
	}
	if position >= idx.s.Len() {
		return nil, fmt.Errorf("invalid position: %d (len: %v)", position, idx.s.Len())
	}
	elem := idx.s.Element(position)
	return elem.Labels[level], nil
}

func (s *Series) rename(name string) {
	s = s.Copy()
	s.index.Levels[0].Name = name
}

// LevelToFloat64 converts the labels at a specified index level to float64 and returns a new Series.
func (idx Index) LevelToFloat64(level int) (*Series, error) {
	if level > idx.s.index.Len() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.Len())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToFloat64()
	return s, nil
}

// ToFloat64 converts the labels at index level 0 to float64 and returns a new Series.
func (idx Index) ToFloat64() *Series {
	if idx.s.index.Len() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToFloat64(0)
	return s
}

// LevelToInt64 converts the labels at a specified index level to int64 and returns a new Series.
func (idx Index) LevelToInt64(level int) (*Series, error) {
	if level > idx.s.index.Len() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.Len())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInt64()
	return s, nil
}

// ToInt64 converts the labels at index level 0 to int64 and returns a new Series.
func (idx Index) ToInt64() *Series {
	if idx.s.index.Len() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToInt64(0)
	return s
}

// LevelToString converts the labels at a specified index level to string and returns a new Series.
func (idx Index) LevelToString(level int) (*Series, error) {
	if level > idx.s.index.Len() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.Len())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToString()
	return s, nil
}

// ToString converts the labels at index level 0 to string and returns a new Series.
func (idx Index) ToString() *Series {
	if idx.s.index.Len() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToString(0)
	return s
}

// LevelToBool converts the labels at a specified index level to bool and returns a new Series.
func (idx Index) LevelToBool(level int) (*Series, error) {
	if level > idx.s.index.Len() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.Len())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToBool()
	return s, nil
}

// ToBool converts the labels at index level 0 to bool and returns a new Series.
func (idx Index) ToBool() *Series {
	if idx.s.index.Len() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToBool(0)
	return s
}

// LevelToDateTime converts the labels at a specified index level to DateTime and returns a new Series.
func (idx Index) LevelToDateTime(level int) (*Series, error) {
	if level > idx.s.index.Len() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.Len())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToDateTime()
	return s, nil
}

// ToDateTime converts the labels at index level 0 to DateTime and returns a new Series.
func (idx Index) ToDateTime() *Series {
	if idx.s.index.Len() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToDateTime(0)
	return s
}

// LevelToInterface converts the labels at a specified index level to interface and returns a new Series.
func (idx Index) LevelToInterface(level int) (*Series, error) {
	if level > idx.s.index.Len() {
		return nil, fmt.Errorf("invalid index level: %d (len: %v)", level, idx.s.index.Len())
	}
	s := idx.s.Copy()
	s.index.Levels[level] = s.index.Levels[level].ToInterface()
	return s, nil
}

// ToInterface converts the labels at index level 0 to interface and returns a new Series.
func (idx Index) ToInterface() *Series {
	if idx.s.index.Len() == 0 {
		log.Println("Cannot convert empty index")
		return nil
	}
	s, _ := idx.LevelToInterface(0)
	return s
}