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

// Func executes the provided function and wraps the result in a []interface{}
// which is used by Terraform list or set types.
func Func(fn func(helper.ResourceData)) []interface{} {
	return Flatten(FlattenerFunc(fn))
}

// FlattenerList is used when flattening list or set types into Terraform's
// internal representation.
//
// The methods require that the elements of the collection be enumerated by an
// integer index.
type FlattenerList interface {
	// Len returns the number of elements in the collection.
	Len() int
	// Flatten flattens the element at index i into data d.
	Flatten(i int, d helper.ResourceData)
}

// List flattens a ListFlattener by iterating the List's elements and calling
// their Flatten method.
func List(l FlattenerList) []interface{} {
	out := make([]interface{}, 0, l.Len())
	for i := 0; i < l.Len(); i++ {
		d := make(helper.MapData)
		l.Flatten(i, d)
		out = append(out, map[string]interface{}(d))
	}
	return out
}

// Flatteners is a type alias for []Flattener which implmements FlattenerList.
type Flatteners []Flattener

func (f Flatteners) Len() int                             { return len(f) }
func (f Flatteners) Flatten(i int, d helper.ResourceData) { f[i].Flatten(d) }
