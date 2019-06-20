package series

import (
	"fmt"
	"strings"
	"time"

	"github.com/ptiger10/pd/options"
)

func (s *Series) String() string {
	if Equal(s, newEmptySeries()) {
		return "Series{}"
	}
	return s.print()
}

// expects to receive a slice of typed value structs (eg values.float64Values)
func (s *Series) print() string {
	numLevels := len(s.index.Levels)
	var header string
	var printer string
	// [START header row]
	maxIndexWidths := s.index.MaxWidths()
	for j := 0; j < numLevels; j++ {
		name := s.index.Levels[j].Name
		padding := maxIndexWidths[j]
		header += fmt.Sprintf("%*v", padding, name)
		if j != numLevels-1 {
			// add buffer to all index levels except the last
			header += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
		}
	}
	// omit header line if empty
	if strings.TrimSpace((header)) != "" {
		printer += header + "\n"
	}

	// [END header row]

	// [START rows]
	prior := make(map[int]string)
	for i := 0; i < s.Len(); i++ {
		elem := s.Element(i)
		var newLine string

		// [START index printer]
		for j := 0; j < numLevels; j++ {
			var skip bool
			var buffer string
			padding := maxIndexWidths[j]
			idx := fmt.Sprint(elem.Labels[j])
			if j != numLevels-1 {
				// add buffer to all index levels except the last
				buffer = strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
				// skip repeated label values if this is not the last index level
				if prior[j] == idx {
					skip = true
					idx = ""
				}
			}

			printStr := fmt.Sprintf("%*v", padding, idx)
			// elide index string if longer than the max allowable width
			if padding == options.GetDisplayMaxWidth() {
				printStr = printStr[:len(printStr)-4] + "..."
			}

			newLine += printStr + buffer

			// set prior row value for each index level except the last
			if j != numLevels-1 && !skip {
				prior[j] = idx
			}
		}

		// [END index printer]

		// [START value printer]
		var valStr string
		if s.datatype == options.DateTime {
			valStr = elem.Value.(time.Time).Format(options.GetDisplayTimeFormat())
		} else {
			valStr = fmt.Sprint(elem.Value)
		}

		// add buffer at beginning
		val := strings.Repeat(" ", options.GetDisplayValuesWhitespaceBuffer()) + valStr
		// null string values must not return any trailing whitespace
		if valStr == "" {
			val = strings.TrimSpace(val)
		}
		newLine += val
		// Concatenate line onto printer string
		printer += fmt.Sprintln(newLine)
	}
	if s.datatype != options.None {
		printer += fmt.Sprintf("datatype: %s\n", s.datatype)
	}
	// [END rows]

	if s.name != "" {
		printer += fmt.Sprintf("name: %s\n", s.name)
	}
	return printer
}
