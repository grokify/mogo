package sortutil

import (
	"reflect"
	"testing"
)

func TestIntegerSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "basic numeric suffix",
			input:    []string{"ABC-10", "ABC-2", "ABC-1"},
			expected: []string{"ABC-1", "ABC-2", "ABC-10"},
		},
		{
			name:     "mixed numeric and non-numeric",
			input:    []string{"ABC-10", "ABC-2", "XYZ", "FOO-bar", "ABC-abc", "ABC-1"},
			expected: []string{"ABC-1", "ABC-2", "ABC-10", "ABC-abc", "FOO-bar", "XYZ"},
		},
		{
			name:     "no numeric suffixes",
			input:    []string{"banana", "apple", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "numeric and empty suffix",
			input:    []string{"Task-3", "Task-A", "Task-2", "Task"},
			expected: []string{"Task-2", "Task-3", "Task", "Task-A"},
		},
		{
			name:     "stable ordering",
			input:    []string{"ABC-1", "ABC-1", "ABC-2"},
			expected: []string{"ABC-1", "ABC-1", "ABC-2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntegerSuffix(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("IntegerSuffix() = %v, want %v", got, tt.expected)
			}
		})
	}
}
