package utils

import (
	"reflect"
)

func Clone(inter interface{}) interface{} {
	nInter := reflect.New(reflect.TypeOf(inter).Elem())
	nVal := nInter.Elem()

	source := reflect.ValueOf(inter).Elem()
	for i := 0; i < source.NumField(); i++ {
		dest := nVal.Field(i)
		if dest.CanSet() {
			copyValue(source.Field(i), dest)
		}
	}
	x := nInter.Interface()
	return x
}

func copyValue(source reflect.Value, dest reflect.Value) {
	if dest.Kind() == reflect.Ptr {
		if !source.IsZero() {
			field := reflect.New(source.Elem().Type())
			field.Elem().Set(source.Elem())
			dest.Set(field)
		}
	} else if dest.Kind() == reflect.Slice {
		for j := 0; j < source.Len(); j++ {
			if !source.Index(j).IsZero() {
				v := source.Index(j)
				var field reflect.Value
				if v.Kind() == reflect.Ptr {
					field = reflect.New(v.Elem().Type())
					field.Elem().Set(v.Elem())
				} else {
					field = reflect.New(v.Type())
					field.Set(v)
				}
				dest.Set(reflect.Append(dest, field))
			}
		}
	} else {
		dest.Set(source)
	}
}