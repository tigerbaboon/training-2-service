package collections

import "reflect"

// KeyBy Create map[K]V from struct by use key is struct
func KeyBy[K comparable, V any](keyBy string, src []V, dest map[K]V) {
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dest)
	switch {
	case s.Kind() != reflect.Slice:
		panic("KeyBy() given a non-slice type")
	case d.Kind() != reflect.Map:
		panic("KeyBy() given a non-map type")
	}
	for i := 0; i < s.Len(); i++ {
		value := s.Index(i)
		key := s.Index(i)
		if key.Kind() == reflect.Ptr {
			key = value.Elem()
		}
		d.SetMapIndex(key.FieldByName(keyBy), value)
	}
}

// KeyBy2Array Create map[K][]V from struct by use key is struct
func KeyBy2Array[K comparable, V any](keyBy string, src []V, dest map[K][]V) {
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dest)
	switch {
	case s.Kind() != reflect.Slice:
		panic("KeyBy() src given a non-slice type")
	case d.Kind() != reflect.Map:
		panic("KeyBy() dest given a non-map type")
	}
	for i := 0; i < s.Len(); i++ {
		value := s.Index(i)
		beforeD := d.MapIndex(value.FieldByName(keyBy))
		typeD := d.Type().Elem()
		switch {
		case beforeD.Kind() == reflect.Invalid && typeD.Kind() == reflect.Slice:
			beforeD = reflect.Zero(typeD)
		case typeD.Kind() != reflect.Slice:
			panic("KeyBy() int dest given a non-slice type")
		}
		d.SetMapIndex(value.FieldByName(keyBy), reflect.Append(beforeD, value))
	}
}
