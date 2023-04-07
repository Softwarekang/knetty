package utils

import (
	"testing"
)

func TestAdjustNToPowerOfTwo(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected int
	}{
		{"Test case 1", 9, 16},
		{"Test case 2", 15, 16},
		{"Test case 3", 64, 64},
		{"Test case 4", 128, 128},
		{"Test case 5", 255, 256},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := AdjustNToPowerOfTwo(tc.input)
			if result != tc.expected {
				t.Errorf("Unexpected result. Expected: %d, got: %d", tc.expected, result)
			}
		})
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Test case 1", 2, true},
		{"Test case 2", 3, false},
		{"Test case 3", 4, true},
		{"Test case 4", 7, false},
		{"Test case 5", 8, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsPowerOfTwo(tc.input)
			if result != tc.expected {
				t.Errorf("Unexpected result. Expected: %t, got: %t", tc.expected, result)
			}
		})
	}
}
