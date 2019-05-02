package values

import (
	"fmt"
	"math"
	"time"

	"github.com/ptiger10/pd/new/options"
)

// [START Definitions]

// IntValues is a slice of Int64-typed Value/Null structs
type IntValues []IntValue

// An IntValue is one Int64-typed Value/Null struct
type IntValue struct {
	v    int64
	null bool
}

// Int constructs an IntValue
func Int(v int64, null bool) IntValue {
	return IntValue{
		v:    v,
		null: null,
	}
}

// [END Definitions]

// [START Converters]

// ToFloat converts IntValues to FloatValues
//
// 1: 1.0
func (vals IntValues) ToFloat() Values {
	var ret FloatValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Float(math.NaN(), true))
		} else {
			v := float64(val.v)
			ret = append(ret, Float(v, false))
		}
	}
	return ret
}

// ToInt returns itself
func (vals IntValues) ToInt() Values {
	return vals
}

// ToBool converts IntValues to BoolValues
//
// x != 0: true; x == 0: false; null: false
func (vals IntValues) ToBool() Values {
	var ret BoolValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, Bool(false, true))
		} else {
			if val.v == 0 {
				ret = append(ret, Bool(false, false))
			} else {
				ret = append(ret, Bool(true, false))
			}
		}
	}
	return ret
}

// ToDateTime converts IntValues to DateTimeValues.
// Tries to convert from Unix EPOCH timestamp.
// Defaults to 1970-01-01 00:00:00 +0000 UTC.
func (vals IntValues) ToDateTime() Values {
	var ret DateTimeValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, DateTime(time.Time{}, true))
		} else {
			ret = append(ret, intToDateTime(val.v))
		}
	}
	return ret
}

func intToDateTime(i int64) DateTimeValue {
	// convert from nanoseconds to seconds
	i /= 1000000000
	v := time.Unix(i, 0).UTC()
	return DateTime(v, false)
}

// [END Converters]

// [START Methods]

// Describe the values in the collection
func (vals IntValues) Describe() string {
	offset := options.DisplayValuesWhitespaceBuffer
	l := len(vals)
	len := fmt.Sprintf("%-*s%d\n", offset, "len", l)
	return fmt.Sprint(len)
}

// [END Methods]