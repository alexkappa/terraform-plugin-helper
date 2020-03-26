package flatten

import (
	"fmt"
	"reflect"

	"github.com/alexkappa/terraform-plugin-helper/helper"
)

// A Flattener is used to flatten data into Terraform's internal representation.
type Flattener interface {
	Flatten(helper.Data)
}

// The FlattenerFunc type is an adapter to allow the use of an ordinary function
// as a Flattener. If f is a function with the appropriate signature,
// FlattenerFunc(f) is a Flattener that calls f.
type FlattenerFunc func(helper.Data)

// Flatten calls f(m).
func (fn FlattenerFunc) Flatten(d helper.Data) {
	fn(d)
}

// Flatten executes the provided flatteners Flatten method and wraps the result
// in a []interface{} which is used by Terraform list or set types.
func Flatten(f Flattener) []interface{} {
	d := make(helper.MapData)
	f.Flatten(d)
	return []interface{}{map[string]interface{}(d)}
}

// FlattenFunc executes the provided function and wraps the result in a
// []interface{} which is used by Terraform list or set types.
func FlattenFunc(fn func(helper.Data)) []interface{} {
	return Flatten(FlattenerFunc(fn))
}

// type FlattenerList []Flattener

// func (fls FlattenerList) Flatten() []interface{} {
// 	out := make([]interface{}, len(fls))
// 	for i, f := range fls {
// 		m := make(helper.MapData)
// 		f.Flatten(m)
// 		out[i] = map[string]interface{}(m)
// 	}
// 	return out
// }

func FlattenList(in interface{}) (out []interface{}) {

	v := reflect.ValueOf(in)

	switch v.Kind() {
	case reflect.Slice, reflect.Array:

		out = make([]interface{}, v.Len())

		for i := 0; i < v.Len(); i++ {
			f, ok := v.Index(i).Interface().(Flattener)
			if !ok {
				panic(fmt.Sprintf("flatten: in[%d] %s is not a Flattener", i, v.Index(i).Type()))
			}
			out[i] = f
		}

	default:
		panic("flatten: input is not a slice or array type")
	}

	return
}

func FlattenListFunc(in interface{}, fn func(interface{}, helper.Data)) (out []interface{}) {

	switch reflect.TypeOf(in).Kind() {
	case reflect.Slice, reflect.Array:

		s := reflect.ValueOf(in)

		for i := 0; i < s.Len(); i++ {
			d := make(helper.MapData)
			fn(s.Index(i).Interface(), d)
			out = append(out, d)
		}
	default:
		panic("flatten: input is not a slice or array type")
	}

	return
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
