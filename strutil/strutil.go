package strutil

import "reflect"

// Contains checks if an item is present in a list.
func Contains(list interface{}, item interface{}) bool {
	listValue := reflect.ValueOf(list)

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < listValue.Len(); i++ {
			if reflect.DeepEqual(listValue.Index(i).Interface(), item) {
				return true
			}
		}
	}

	return false
}
