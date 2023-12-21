package floats

import "math"

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}
