package index

import (
	"fmt"

	"github.com/ptiger10/pd/kinds"
)

// Convert an index level from one kind to another, then refresh the LabelMap
func (lvl Level) Convert(kind kinds.Kind) (Level, error) {
	var convertedLvl Level
	switch kind {
	case kinds.None:
		return Level{}, fmt.Errorf("unable to convert index level: must supply a valid Kind")
	case kinds.Float:
		convertedLvl = lvl.toFloat()
	case kinds.Int:
		convertedLvl = lvl.toInt()
	case kinds.String:
		convertedLvl = lvl.toString()
	case kinds.Bool:
		convertedLvl = lvl.toBool()
	case kinds.DateTime:
		convertedLvl = lvl.toDateTime()
	case kinds.Interface:
		convertedLvl = lvl.toInterface()
	default:
		return Level{}, fmt.Errorf("unable to convert level: kind not supported: %v", kind)
	}
	convertedLvl.Refresh()
	return convertedLvl, nil
}

func (lvl Level) toFloat() Level {
	lvl.Labels = lvl.Labels.ToFloat()
	lvl.Kind = kinds.Float
	return lvl
}

func (lvl Level) toInt() Level {
	lvl.Labels = lvl.Labels.ToInt()
	lvl.Kind = kinds.Int
	return lvl
}

func (lvl Level) toString() Level {
	lvl.Labels = lvl.Labels.ToString()
	lvl.Kind = kinds.String
	return lvl
}

func (lvl Level) toBool() Level {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.Kind = kinds.Bool
	return lvl
}

func (lvl Level) toDateTime() Level {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.Kind = kinds.DateTime
	return lvl
}

func (lvl Level) toInterface() Level {
	lvl.Labels = lvl.Labels.ToInterface()
	lvl.Kind = kinds.Interface
	return lvl
}
