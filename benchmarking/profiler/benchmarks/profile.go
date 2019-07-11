package benchmarks

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"
)

// Descriptions of the benchmarking tests
var Descriptions = map[string]string{
	"sum":    "Sum one column",
	"mean":   "Simple mean of one column",
	"median": "Median of one column",
}

// RunGoProfiler specifies all the benchmarks to profile and return in the benchmark table.
func RunGoProfiler() Results {
	fmt.Println("Profiling Go")
	Results := Results{
		"100k": {
			"sum":    ProfileGo(benchmarkSumFloat64_100000),
			"mean":   ProfileGo(benchmarkMeanFloat64_100000),
			"median": ProfileGo(benchmarkMedianFloat64_100000),
		},
	}
	return Results
}

// Results contains benchmarking results
// {"num of samples": {"test1": "10ms"...}}
type Results map[string]map[string]Result

// A Result of benchmarking data in the form [string, float64]
type Result []interface{}

// ProfileGo runs the normal Go benchmarking command but formats the result as a rounded string
// and raw ns float
func ProfileGo(f func(b *testing.B)) Result {
	benchmark := testing.Benchmark(f).NsPerOp()
	var speed string
	switch {
	case benchmark < int64(time.Microsecond):
		speed = fmt.Sprintf("%vns", benchmark)
	case benchmark < int64(time.Millisecond):
		speed = fmt.Sprintf("%.2fμs", float64(benchmark)/float64(time.Microsecond))
	case benchmark < int64(time.Second):
		speed = fmt.Sprintf("%.2fms", float64(benchmark)/float64(time.Millisecond))
	default:
		speed = fmt.Sprintf("%.2fs", float64(benchmark)/float64(time.Second))
	}
	return Result{speed, float64(benchmark)}
}

// RunPythonProfiler executes main.py in this directory, which is expected to return JSON
// in the form of Results. This command is expected to be initiated from the directory above.
func RunPythonProfiler() Results {
	fmt.Println("Profiling Python")
	cmd := exec.Command("python", "benchmarks/main.py")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var r Results
	err = json.Unmarshal(out, &r)
	if err != nil {
		log.Fatal(err)
	}
	return r
}
