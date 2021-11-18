// Package options makes options from kong structs
package options

import (
	"fmt"
	"reflect"
)

// Transfer the values from source to target
// struct fields must have the same name and type (pointer to type is accepted)
// the source value must be set
// the target value must be valid and settable
func Transfer(source, target interface{}) error {
	src := reflect.ValueOf(source)
	tgt := reflect.ValueOf(target)

	if !isValid(src, tgt) {
		return fmt.Errorf("src and tgt must be a pointer to a struct")
	}

	src = src.Elem()
	tgt = tgt.Elem()

	for i := 0; i < src.NumField(); i++ {
		n := src.Type().Field(i).Name

		sv := src.Field(i) // source value

		if sv.IsZero() {
			continue // source value is not set
		}

		tf, ok := tgt.Type().FieldByName(n)
		if !ok {
			continue // field does not exist in target
		}

		if !isEqualType(sv.Type(), tf.Type) {
			continue // source and target type are not equal
		}

		tv := tgt.FieldByName(n) // target value

		// check if target value is valid and settable
		if !tv.IsValid() || !tv.CanSet() {
			continue // target field is not settable
		}

		// set target value
		if tf.Type.Kind() != reflect.Ptr {
			nv := reflect.New(tf.Type).Elem()
			nv.Set(sv) // set source value
			tv.Set(nv) // set value on target
		} else {
			nv := reflect.New(tf.Type.Elem()).Elem()
			nv.Set(sv)        // set source value
			tv.Set(nv.Addr()) // set value on target
		}
	}

	return nil
}

// isValid checks if input is a pointer to struct for source and target
func isValid(src, tgt reflect.Value) bool {
	return src.Kind() == reflect.Ptr &&
		tgt.Kind() == reflect.Ptr &&
		src.Elem().Kind() == reflect.Struct &&
		tgt.Elem().Kind() == reflect.Struct
}

// isEqualType checks if source and target type are equal
func isEqualType(st, tt reflect.Type) bool {
	if st.Kind() == reflect.Ptr {
		st = tt.Elem()
	}
	if tt.Kind() == reflect.Ptr {
		tt = tt.Elem()
	}

	return st == tt
}
