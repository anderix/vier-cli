// Author: David M. Anderson
// Built with AI assistance (Claude, Anthropic)

package main

// SplitMix64 is the generator vier-desktop uses, so one seed gives both programs the same numbers.
type SplitMix64 struct {
	state uint64
}

func (r *SplitMix64) Next() uint64 {
	r.state += 0x9E3779B97F4A7C15
	z := r.state
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

// Unit returns a number in [0, 1).
func (r *SplitMix64) Unit() float64 {
	return float64(r.Next()>>11) / float64(uint64(1)<<53)
}

// Below returns a number in [0, bound).
func (r *SplitMix64) Below(bound int) int {
	return int(r.Next() % uint64(bound))
}
