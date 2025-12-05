package main

import (
	"slices"
	"testing"
)

type testCase struct {
	name     string
	input    Ranges
	expected Ranges
}

func TestReduceNoOverlaps(t *testing.T) {
	tests := []testCase{
		{
			name: "No overlaps",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 12, End: 15},
			},
			expected: Ranges{
				{Start: 1, End: 10},
				{Start: 12, End: 15},
			},
		},
		{
			name: "Adjacent",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 11, End: 15},
			},
			expected: Ranges{
				{Start: 1, End: 10},
				{Start: 11, End: 15},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.Reduce()
			if !slices.Equal(tc.input, tc.expected) {
				t.Errorf(`%v != %v`, tc.input, tc.expected)
			}
		})
	}
}

func TestReduceSimpleOverlaps(t *testing.T) {
	tests := []testCase{
		{
			name: "Overlaps end",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 10, End: 15},
			},
			expected: Ranges{
				{Start: 1, End: 15},
			},
		},
		{
			name: "Overlaps start",
			input: Ranges{
				{Start: 5, End: 10},
				{Start: 0, End: 8},
			},
			expected: Ranges{
				{Start: 0, End: 10},
			},
		},
		{
			name: "Duplicate",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 1, End: 10},
			},
			expected: Ranges{
				{Start: 1, End: 10},
			},
		},
		{
			name: "Single ID at start",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 1, End: 1},
			},
			expected: Ranges{
				{Start: 1, End: 10},
			},
		},
		{
			name: "Single ID at end",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 10, End: 10},
			},
			expected: Ranges{
				{Start: 1, End: 10},
			},
		},
		{
			name: "Single ID in middle",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 5, End: 5},
			},
			expected: Ranges{
				{Start: 1, End: 10},
			},
		},
		{
			name: "Fully contained",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 3, End: 8},
			},
			expected: Ranges{
				{Start: 1, End: 10},
			},
		},
		{
			name: "Fully contained reverse",
			input: Ranges{
				{Start: 3, End: 8},
				{Start: 1, End: 10},
			},
			expected: Ranges{
				{Start: 1, End: 10},
			},
		},
		{
			name: "One non-overlapping",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 10, End: 15},
				{Start: 20, End: 25},
			},
			expected: Ranges{
				{Start: 1, End: 15},
				{Start: 20, End: 25},
			},
		},
		{
			name: "Multiple",
			input: Ranges{
				{Start: 1, End: 10},
				{Start: 10, End: 15},
				{Start: 5, End: 10},
				{Start: 0, End: 8},
				{Start: 1, End: 10},
				{Start: 3, End: 8},
			},
			expected: Ranges{
				{Start: 0, End: 15},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.Reduce()
			if !slices.Equal(tc.input, tc.expected) {
				t.Errorf(`%v != %v`, tc.input, tc.expected)
			}
		})
	}
}
