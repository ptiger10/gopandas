package series

import (
	"log"
	"time"
)

// appropriate for numeric data only ([]float64 or []int64)
func ensureFloatFromNumerics(vals interface{}) []float64 {
	var data []float64
	if ints, ok := vals.([]int64); ok {
		data = convertIntToFloat(ints)
	} else if floats, ok := vals.([]float64); ok {
		data = floats
	} else {
		log.Printf("Internal error: ensureFloatFromNumerics has received an unallowable value: %v", vals)
		return nil
	}
	return data
}

func convertIntToFloat(vals []int64) []float64 {
	var ret []float64
	for _, val := range vals {
		ret = append(ret, float64(val))
	}
	return ret
}

func ensureStrings(vals interface{}) []string {
	if strings, ok := vals.([]string); ok {
		return strings
	}
	log.Printf("Internal error: ensureStrings has received an unallowable value: %v", vals)
	return nil
}

func ensureBools(vals interface{}) []bool {
	if bools, ok := vals.([]bool); ok {
		return bools
	}
	log.Printf("Internal error: ensureBools has received an unallowable value: %v", vals)
	return nil
}

func ensureDateTime(vals interface{}) []time.Time {
	if datetime, ok := vals.([]time.Time); ok {
		return datetime
	}
	log.Printf("Internal error: ensureDateTime has received an unallowable value: %v", vals)
	return nil
}