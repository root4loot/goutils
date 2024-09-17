package sliceutil

import (
	"reflect"
	"testing"
)

func TestAppendUnique(t *testing.T) {
	initial := []string{"a", "b"}
	result := AppendUnique(initial, "b", "c", "d")
	expected := []string{"a", "b", "c", "d"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("AppendUnique failed, expected %v, got %v", expected, result)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !Contains(slice, "b") {
		t.Error("Contains failed, expected true, got false")
	}

	if Contains(slice, "d") {
		t.Error("Contains failed, expected false, got true")
	}
}

func TestEqual(t *testing.T) {
	if !Equal("a", "a", "b", "c") {
		t.Error("Equal failed, expected true, got false")
	}

	if Equal("d", "a", "b", "c") {
		t.Error("Equal failed, expected false, got true")
	}
}

func TestRemove(t *testing.T) {
	initial := []string{"a", "b", "c"}
	result := Remove(initial, "b")
	expected := []string{"a", "c"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Remove failed, expected %v, got %v", expected, result)
	}
}

func TestDeleteEmpty(t *testing.T) {
	initial := []string{"a", "", "b", "", "c"}
	result := DeleteEmpty(initial)
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("DeleteEmpty failed, expected %v, got %v", expected, result)
	}
}

func TestUnique(t *testing.T) {
	initial := []string{"a", "b", "a", "c", "b"}
	result := Unique(initial)
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unique failed, expected %v, got %v", expected, result)
	}
}
