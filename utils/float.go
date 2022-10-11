package utils

import "math"

func FloatEquals(a, b, delta float32) bool {
	diff := math.Abs(float64(a - b))

	return diff > float64(delta)
}
