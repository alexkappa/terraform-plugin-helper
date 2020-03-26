package flatten

import (
	"fmt"
	"reflect"
)

// A Flattener is used to flatten data into Terraform's internal representation.
type Flattener interface {
	Flatten(map[string]interface{})
}

// The FlattenerFunc type is an adapter to allow the use of an ordinary function
// as a Flattener. If f is a function with the appropriate signature,
// FlattenerFunc(f) is a Flattener that calls f.
type FlattenerFunc func(map[string]interface{})

// Flatten calls f(m).
func (fn FlattenerFunc) Flatten(m map[string]interface{}) {
	fn(m)
}

// Flatten executes the provided flatteners Flatten method and wraps the result
// in a []interface{} which is used by Terraform list or set types.
func Flatten(f Flattener) []interface{} {
	m := make(map[string]interface{})
	f.Flatten(m)
	return []interface{}{m}
}

// FlattenFunc executes the provided function and wraps the result in a
// []interface{} which is used by Terraform list or set types.
func FlattenFunc(fn func(map[string]interface{})) []interface{} {
	return Flatten(FlattenerFunc(fn))
}

func FlattenList(in interface{}) (out []interface{}) {

	switch reflect.TypeOf(in).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(in)

		for i := 0; i < s.Len(); i++ {
			m := make(map[string]interface{})
			f, ok := s.Index(i).Interface().(Flattener)
			if !ok {
				panic(fmt.Sprintf("FlattenList: in[%d] %s is not a Flattener", i, s.Index(i).Type()))
			}
			f.Flatten(m)
			out = append(out, m)
		}

	default:
		panic("FlattenList: input is not a slice or array type")
	}

	return
}

func FlattenListFunc(in interface{}, fn func(interface{}, map[string]interface{})) (out []interface{}) {

	switch reflect.TypeOf(in).Kind() {
	case reflect.Slice, reflect.Array:

		s := reflect.ValueOf(in)

		for i := 0; i < s.Len(); i++ {
			m := make(map[string]interface{})
			fn(s.Index(i).Interface(), m)
			out = append(out, m)
		}
	default:
		panic("FlattenListFunc: input is not a slice or array type")
	}

	return
}

type FlattenerList []Flattener

func List(in interface{}) FlattenerList {

	var fl []Flattener

	v := reflect.ValueOf(in)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:

		fl = make([]Flattener, v.Len())
		for i := 0; i < v.Len(); i++ {
			f, ok := v.Index(i).Interface().(Flattener)
			if !ok {
				panic(fmt.Sprintf("FlattenList: in[%d] %s is not a Flattener", i, v.Index(i).Type()))
			}
			fl[i] = f
		}

	default:
		panic("FlattenList: input is not a slice or array type")
	}

	return FlattenerList(fl)
}

// ---

// flatten.ListFunc(in.ContainerSpec.Mounts, func(v interface{}, m map[string]interface{} {
// 	m["target"] = v.(mount.Mount).Target
// 	m["source"] = v.(mount.Mount).Source
// 	m["type"] = v.(mount.Mount).Type
// }))

// flatten.List(in.ContainerSpec.Mounts).FlattenFunc(func(v interface{}, m map[string]interface{}){
// 		m["target"] = v.(mount.Mount).Target
// 		m["source"] = v.(mount.Mount).Source
// 		m["type"] = v.(mount.Mount).Type
// })

// func List() Iterator

// Iterator.Flatten()
