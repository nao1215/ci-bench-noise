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

// baseRounds is the number of rounds of the unchanged build. Every round is
// the same: it sorts 20,000 numbers and sums a strided walk over them. A
// slowdown adds whole rounds, one round per percent, so the work grows by
// exactly the percentage (a partial round would sort fewer numbers, and
// sorting does not scale linearly with the count).
const baseRounds = 100

const size = 20_000

func main() {
	percent, err := strconv.Atoi(strings.TrimSpace(slowdownText))
	if err != nil || percent < 0 {
		fmt.Fprintln(os.Stderr, "slowdown.txt: want a non-negative integer percentage")
		os.Exit(2)
	}
	rounds := baseRounds + baseRounds*percent/100

	data := make([]uint32, size)
	var state uint32 = 2463534242
	var sum uint64
	for range rounds {
		n := size
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
