// Command work does a fixed amount of CPU and memory work, so that its run
// time can be measured on a CI runner.
//
// The amount is scaled by the percentage in slowdown.txt, which is embedded at
// build time: a value of 10 makes every run do 10% more of the same work. The
// experiments compare a build with 0 against a build with 0 (an A/A test) or
// with a known slowdown, to see how often a shared runner confuses the two.
package main

import (
	_ "embed"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

//go:embed slowdown.txt
var slowdownText string

// baseRounds is the number of rounds of the unchanged build. One round sorts
// 200,000 numbers and sums a strided walk over them, about 20ms on a desktop.
const baseRounds = 5

const size = 200_000

func main() {
	percent, err := strconv.ParseFloat(strings.TrimSpace(slowdownText), 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "slowdown.txt:", err)
		os.Exit(2)
	}
	// Work is counted in units of 1/100 round, so 5% of 5 rounds is exact.
	units := int(float64(baseRounds*100) * (1 + percent/100))

	data := make([]uint32, size)
	var state uint32 = 2463534242
	var sum uint64
	for u := 0; u < units; u += 100 {
		n := size
		if units-u < 100 {
			n = size * (units - u) / 100
		}
		for i := range n {
			// xorshift32: the same numbers in every run and every build.
			state ^= state << 13
			state ^= state >> 17
			state ^= state << 5
			data[i] = state
		}
		part := data[:n]
		slices.Sort(part)
		for i := 0; i < n; i += 7 {
			sum += uint64(part[(i*7919)%n])
		}
	}
	fmt.Println(sum)
}
