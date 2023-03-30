package utils_test

import (
	"fmt"
	"terraform-provider-thebastion/utils"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

// FindStringIndex
func TestFindStringIndex(t *testing.T) {
	require := require.New(t)
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

			got := utils.FindStringIndex(test.input_array, test.input_target)
			diff, err := fmt.Printf("unexpected diff (+wanted, -got): %s", cmp.Diff(got, test.expected))
			require.Equal(err, nil)
			require.Equal(got, test.expected, diff)
		})
	}
}

// ConvertInterfacesToStrings
func TestConvertInterfacesToStrings(t *testing.T) {
	require := require.New(t)
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

			got := utils.ConvertInterfacesToStrings(test.input)
			require.ElementsMatchf(got, test.expected, "unexpected diff (+wanted, -got): %s", cmp.Diff(got, test.expected))
		})
	}
}

// CompareLists
func TestCompareLists(t *testing.T) {
	require := require.New(t)
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

			got_left, got_right := utils.CompareLists(test.input_left, test.input_right)
			require.ElementsMatchf(got_left, test.expected_left_only, "unexpected diff (+wanted, -got): %s", cmp.Diff(got_left, test.expected_left_only))
			require.ElementsMatchf(got_right, test.expected_right_only, "unexpected diff (+wanted, -got): %s", cmp.Diff(got_right, test.expected_right_only))
		})
	}
}
