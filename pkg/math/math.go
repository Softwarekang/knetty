// Package math common mathematical methods
package math

func Max(a, b int) int {
	if a < b {
		return b
	}

	return a
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
