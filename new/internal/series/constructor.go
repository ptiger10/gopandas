package series

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/new/kinds"

	"github.com/ptiger10/pd/new/internal/index"
	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	"github.com/ptiger10/pd/new/internal/values"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
)

// An Option is an optional parameter in the Series constructor
type Option func(*seriesConfig)
type seriesConfig struct {
	indices []miniIndex
	kind    reflect.Kind
	name    string
}

// Kind will convert either values or an index level to the specified kind
func Kind(kind reflect.Kind) Option {
	return func(c *seriesConfig) {
		c.kind = kind
	}
}

// Name will name either values or an index level
func Name(n string) Option {
	return func(c *seriesConfig) {
		c.name = n
	}
}

// Index returns a Option for use in the Series constructor New(),
// and takes an optional Name.
func Index(data interface{}, options ...Option) Option {
	config := seriesConfig{}
	for _, option := range options {
		option(&config)
	}
	return func(c *seriesConfig) {
		idx := miniIndex{
			Data: data,
			Kind: config.kind,
			Name: config.name,
		}
		c.indices = append(c.indices, idx)
	}
}

// New Series constructor
// Optional
// - Index(): If no index is supplied, defaults to a single index of IntValues (0, 1, 2, ...n)
// - Name(): If no name is supplied, no name will appear when Series is printed
// - Kind(): Convert the Series values to the specified kind
// If passing []interface{}, must supply a type expectation for the Series.
// Options: Float, Int, String, Bool, DateTime
func New(data interface{}, options ...Option) (Series, error) {
	// Setup
	config := seriesConfig{}

	for _, option := range options {
		option(&config)
	}
	suppliedKind := config.kind
	var kind reflect.Kind
	name := config.name

	var v values.Values
	var idx index.Index
	var err error

	// Values
	switch reflect.ValueOf(data).Kind() {
	case reflect.Slice:
		v, kind, err = constructVal.ValuesFromSlice(data)

	default:
		return Series{}, fmt.Errorf("Unable to construct new Series: type not supported: %T", data)
	}

	// Optional kind conversion
	if suppliedKind != kinds.None {
		v, err = constructVal.Convert(v, suppliedKind)
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
	}
	// Index
	// Default case: no client-supplied Index
	if config.indices == nil {
		idx = constructIdx.Default(v.Len())
	} else {
		idx, err = indexFromMiniIndex(config.indices, v.Len())
		if err != nil {
			return Series{}, fmt.Errorf("Unable to construct new Series: %v", err)
		}
	}

	// Construct Series
	s := Series{
		index:  idx,
		values: v,
		Kind:   kind,
		Name:   name,
	}

	return s, err
}

// [START MiniIndex]

// an untyped representation of one index level.
// It is used for unpacking client-supplied index data and optional metadata
type miniIndex struct {
	Data interface{}
	Kind reflect.Kind
	Name string
}

// creates a full index from a mini client-supplied representation of an index level,
// but returns an error if every index level is not the same length as requiredLen

func indexFromMiniIndex(minis []miniIndex, requiredLen int) (index.Index, error) {
	var levels []index.Level
	for _, miniIdx := range minis {
		if reflect.ValueOf(miniIdx.Data).Kind() != reflect.Slice {
			return index.Index{}, fmt.Errorf("Unable to construct index: custom index must be a Slice: unsupported index type: %T", miniIdx.Data)
		}
		level, err := constructIdx.LevelFromSlice(miniIdx.Data, miniIdx.Name)
		if err != nil {
			return index.Index{}, fmt.Errorf("Unable to construct index: %v", err)
		}
		if level.Labels.Len() != requiredLen {
			return index.Index{}, fmt.Errorf("Unable to construct index %v:"+
				"mismatch between supplied index length (%v) and expected length (%v)",
				miniIdx.Data, level.Labels.Len(), requiredLen)
		}
		if miniIdx.Kind != kinds.None {
			level.Convert(miniIdx.Kind)
		}
		levels = append(levels, level)
	}
	idx := constructIdx.New(levels...)
	return idx, nil

}

// [END MiniIndex]