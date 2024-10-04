package sliceutil

import (
	"reflect"
	"testing"
)

func TestAppendUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		elems    interface{}
		expected interface{}
	}{
		{"AppendUnique with strings", []string{"apple", "banana"}, []string{"banana", "cherry", "apple"}, []string{"apple", "banana", "cherry"}},
		{"AppendUnique with integers", []int{1, 2, 3}, []int{3, 4, 1, 5}, []int{1, 2, 3, 4, 5}},
		{"AppendUnique with empty slice", []float64{}, []float64{1.1, 2.2, 1.1}, []float64{1.1, 2.2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Type switch based on the input type
			switch v := tt.input.(type) {
			case []string:
				result := AppendUnique(v, tt.elems.([]string)...)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("AppendUnique() = %v, expected %v", result, tt.expected)
				}
			case []int:
				result := AppendUnique(v, tt.elems.([]int)...)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("AppendUnique() = %v, expected %v", result, tt.expected)
				}
			case []float64:
				result := AppendUnique(v, tt.elems.([]float64)...)
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("AppendUnique() = %v, expected %v", result, tt.expected)
				}
			default:
				t.Errorf("Unhandled type %T", tt.input)
			}
		})
	}
}

// Test Contains with different data types
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		item     interface{}
		expected bool
	}{
		{"Contains with strings", []string{"apple", "banana", "cherry"}, "banana", true},
		{"Contains with strings (not present)", []string{"apple", "banana", "cherry"}, "mango", false},
		{"Contains with integers", []int{1, 2, 3, 4}, 3, true},
		{"Contains with integers (not present)", []int{1, 2, 3, 4}, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch v := tt.input.(type) {
			case []string:
				result = Contains(v, tt.item.(string))
			case []int:
				result = Contains(v, tt.item.(int))
			}
			if result != tt.expected {
				t.Errorf("Contains() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test Remove with different data types
func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		item     interface{}
		expected interface{}
	}{
		{"Remove with strings", []string{"apple", "banana", "cherry"}, "banana", []string{"apple", "cherry"}},
		{"Remove with integers", []int{1, 2, 3, 4}, 3, []int{1, 2, 4}},
		{"Remove with float64", []float64{1.1, 2.2, 3.3, 4.4}, 2.2, []float64{1.1, 3.3, 4.4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			switch v := tt.input.(type) {
			case []string:
				result = Remove(v, tt.item.(string))
			case []int:
				result = Remove(v, tt.item.(int))
			case []float64:
				result = Remove(v, tt.item.(float64))
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Remove() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test DeleteEmpty with different data types
func TestDeleteEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"DeleteEmpty with strings", []string{"apple", "", "banana", "", "cherry"}, []string{"apple", "banana", "cherry"}},
		{"DeleteEmpty with integers", []int{0, 1, 0, 2, 0, 3}, []int{1, 2, 3}},           // zero-value for int is 0
		{"DeleteEmpty with float64", []float64{0.0, 1.1, 0.0, 2.2}, []float64{1.1, 2.2}}, // zero-value for float64 is 0.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			switch v := tt.input.(type) {
			case []string:
				result = DeleteEmpty(v)
			case []int:
				result = DeleteEmpty(v)
			case []float64:
				result = DeleteEmpty(v)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeleteEmpty() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test Unique with different data types
func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"Unique with strings", []string{"apple", "banana", "apple", "cherry"}, []string{"apple", "banana", "cherry"}},
		{"Unique with integers", []int{1, 2, 2, 3, 1, 4}, []int{1, 2, 3, 4}},
		{"Unique with float64", []float64{1.1, 2.2, 1.1, 3.3}, []float64{1.1, 2.2, 3.3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			switch v := tt.input.(type) {
			case []string:
				result = Unique(v)
			case []int:
				result = Unique(v)
			case []float64:
				result = Unique(v)
			default:
				t.Fatalf("Unhandled type %T", tt.input)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unique() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
