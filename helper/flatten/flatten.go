// Package flatten contains helper functions to deal with arbitrary data
// structures with terraform providers.
package flatten

import (
	"github.com/alexkappa/terraform-plugin-helper/helper"
)

// A Flattener is used to flatten data into Terraform's internal representation.
type Flattener interface {
	Flatten(helper.ResourceData)
}

// The FlattenerFunc type is an adapter to allow the use of an ordinary function
// as a Flattener. If f is a function with the appropriate signature,
// FlattenerFunc(f) is a Flattener that calls f.
type FlattenerFunc func(helper.ResourceData)

// Flatten calls f(m).
func (fn FlattenerFunc) Flatten(d helper.ResourceData) {
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
func FlattenFunc(fn func(helper.ResourceData)) []interface{} {
	return Flatten(FlattenerFunc(fn))
}

// List is used when flattening list or set types into Terraform's internal
// representation.
//
// The methods require that the elements of the collection be enumerated by an
// integer index.
type List interface {
	// Len returns the number of elements in the collection.
	Len() int
	// Flatten flattens the element at index i into data d.
	Flatten(i int, d helper.ResourceData)
}

// FlattenerList is an implementation of List used to flatten []Flattener.
type FlattenerList []Flattener

// Len returns the number of elements in the collection.
func (f FlattenerList) Len() int { return len(f) }

// Flatten flattens the element at index i into data d.
func (f FlattenerList) Flatten(i int, d helper.ResourceData) { f[i].Flatten(d) }

// FlattenList flattens a List by iterating the List's elements and executing
// their Flatten method.
func FlattenList(l List) []interface{} {
	out := make([]interface{}, 0, l.Len())
	for i := 0; i < l.Len(); i++ {
		d := make(helper.MapData)
		l.Flatten(i, d)
		out = append(out, map[string]interface{}(d))
	}
	return out
}
