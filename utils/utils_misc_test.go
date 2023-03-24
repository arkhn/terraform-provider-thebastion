package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// FindStringIndex
func TestFindStringIndex(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input_array  []string
		input_target string
		expected     int
	}
	tests := map[string]testCase{
		"three elements": {
			input_array:  []string{"0", "1", "2"},
			input_target: "0",
			expected:     0,
		},
		"one element": {
			input_array:  []string{"0"},
			input_target: "0",
			expected:     0,
		},
		"zero elements": {
			input_array:  []string{},
			input_target: "0",
			expected:     -1,
		},
		"not in list": {
			input_array:  []string{"0", "1", "2"},
			input_target: "3",
			expected:     -1,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := FindStringIndex(test.input_array, test.input_target)

			if diff := cmp.Diff(got, test.expected); diff != "" {
				t.Errorf("unexpected diff (+wanted, -got): %s", diff)
			}
		})
	}
}

// ConvertInterfacesToStrings
func TestConvertInterfacesToStrings(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input    []interface{}
		expected []string
	}
	tests := map[string]testCase{
		"three elements": {
			input:    append(make([]interface{}, 0), "0", "1", "2"),
			expected: []string{"0", "1", "2"},
		},
		"one element": {
			input:    append(make([]interface{}, 0), "0"),
			expected: []string{"0"},
		},
		"zero elements": {
			input:    make([]interface{}, 0),
			expected: []string{},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := ConvertInterfacesToStrings(test.input)

			if diff := cmp.Diff(got, test.expected); diff != "" {
				t.Errorf("unexpected diff (+wanted, -got): %s", diff)
			}
		})
	}
}

// CompareLists
func TestCompareLists(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input_left          []string
		input_right         []string
		expected_left_only  []string
		expected_right_only []string
	}
	tests := map[string]testCase{
		"Left only": {
			input_left:          []string{"0", "1", "2"},
			input_right:         []string{},
			expected_left_only:  []string{"0", "1", "2"},
			expected_right_only: []string{},
		},
		"Right only": {
			input_left:          []string{},
			input_right:         []string{"0", "1", "2"},
			expected_left_only:  []string{},
			expected_right_only: []string{"0", "1", "2"},
		},
		"Empty arrays": {
			input_left:          []string{},
			input_right:         []string{},
			expected_left_only:  []string{},
			expected_right_only: []string{},
		},
		"Mixed": {
			input_left:          []string{"0", "1", "2"},
			input_right:         []string{"2", "3", "4"},
			expected_left_only:  []string{"0", "1"},
			expected_right_only: []string{"3", "4"},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got_left, got_right := CompareLists(test.input_left, test.input_right)

			if diff := cmp.Diff(got_left, test.expected_left_only); diff != "" {
				t.Errorf("unexpected diff (+wanted, -got): %s", diff)
			}
			if diff := cmp.Diff(got_right, test.expected_right_only); diff != "" {
				t.Errorf("unexpected diff (+wanted, -got): %s", diff)
			}
		})
	}
}
