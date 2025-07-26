package reflect

import "reflect"

func IsNil(v interface{}) bool {
	if v == nil { // fast path
		return true
	}
	r := reflect.ValueOf(v)
	return r.IsNil() || !r.IsValid()
}
