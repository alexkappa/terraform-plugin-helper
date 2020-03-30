// Package expand contains helper functions used to map terraform configuration
// to an API object.
package expand

//go:generate go run gen.go > expand.gen.go

import (
	"strconv"

	"github.com/alexkappa/terraform-plugin-helper/helper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
)

type data struct {
	prefix string
	helper.ResourceData
}

func dataAtKey(key string, d helper.ResourceData) helper.ResourceData {
	return &data{key, d}
}

func dataAtIndex(i int, d helper.ResourceData) helper.ResourceData {
	return &data{strconv.Itoa(i), d}
}

func (d *data) IsNewResource() bool {
	return d.ResourceData.IsNewResource()
}

func (d *data) HasChange(key string) bool {
	return d.ResourceData.HasChange(d.prefix + "." + key)
}

func (d *data) GetChange(key string) (interface{}, interface{}) {
	return d.ResourceData.GetChange(d.prefix + "." + key)
}

func (d *data) Get(key string) interface{} {
	return d.ResourceData.Get(d.prefix + "." + key)
}

func (d *data) GetOk(key string) (interface{}, bool) {
	return d.ResourceData.GetOk(d.prefix + "." + key)
}

func (d *data) GetOkExists(key string) (interface{}, bool) {
	return d.ResourceData.GetOkExists(d.prefix + "." + key)
}

var _ helper.ResourceData = (*data)(nil)

func get(d helper.ResourceData, key string) (v interface{}, ok bool) {
	if d.IsNewResource() || d.HasChange(key) {
		v, ok = d.GetOkExists(key)
	}
	return
}

// Slice accesses the value held by key and type asserts it to a slice.
func Slice(d helper.ResourceData, key string) (s []interface{}) {
	if d.IsNewResource() || d.HasChange(key) {
		v, ok := d.GetOkExists(key)
		if ok {
			s = v.([]interface{})
		}
	}
	return
}

// Map accesses the value held by key and type asserts it to a map.
func Map(d helper.ResourceData, key string) (m map[string]interface{}) {
	if d.IsNewResource() || d.HasChange(key) {
		v, ok := d.GetOkExists(key)
		if ok {
			m = v.(map[string]interface{})
		}
	}
	return
}

// List accesses the value held by key and returns an iterator able to go over
// its elements.
func List(d helper.ResourceData, key string) Iterator {
	if d.IsNewResource() || d.HasChange(key) {
		v, ok := d.GetOkExists(key)
		if ok {
			return &list{dataAtKey(key, d), v.([]interface{})}
		}
	}
	return &list{}
}

// Set accesses the value held by key, type asserts it to a set and returns an
// iterator able to go over its elements.
func Set(d helper.ResourceData, key string) Iterator {
	if d.IsNewResource() || d.HasChange(key) {
		v, ok := d.GetOkExists(key)
		if ok {
			if s, ok := v.(*schema.Set); ok {
				return &set{dataAtKey(key, d), s}
			}
		}
	}
	return &set{nil, &schema.Set{}}
}

// Iterator enables access to the elements of a list or set.
type Iterator interface {

	// Elem iterates over all elements of the list or set, calling fn with each
	// iteration.
	//
	// The callback takes a Data interface as argument which is prefixed with
	// its parents key, making nested data access more convenient.
	//
	// The operation
	//
	// 	bar = d.Get("foo.0.bar").(string)
	//
	// can be expressed as
	//
	// 	List(d, "foo").Elem(func (d Data) {
	//		bar = String(d, "bar")
	// 	})
	//
	// making data access more intuitive for nested structures.
	Elem(func(d helper.ResourceData))

	// Range iterates over all elements of the list, calling fn in each iteration.
	Range(func(k int, v interface{}))

	// List returns the underlying list as a Go slice.
	List() []interface{}
}

type list struct {
	d helper.ResourceData
	v []interface{}
}

func (l *list) Range(fn func(key int, value interface{})) {
	for key, value := range l.v {
		fn(key, value)
	}
}

func (l *list) Elem(fn func(helper.ResourceData)) {
	for idx := range l.v {
		fn(dataAtIndex(idx, l.d))
	}
}

func (l *list) List() []interface{} {
	return l.v
}

type set struct {
	d helper.ResourceData
	s *schema.Set
}

func (s *set) hash(item interface{}) string {
	code := s.s.F(item)
	if code < 0 {
		code = -code
	}
	return strconv.Itoa(code)
}

func (s *set) Range(fn func(key int, value interface{})) {
	for key, value := range s.s.List() {
		fn(key, value)
	}
}

func (s *set) Elem(fn func(helper.ResourceData)) {
	for _, v := range s.s.List() {
		fn(dataAtKey(s.hash(v), s.d))
	}
}

func (s *set) List() []interface{} {
	return s.s.List()
}

// Diff accesses the value held by key and type asserts it to a set. It then
// compares it's changes if any and returns what needs to be added and what
// needs to be removed.
func Diff(d helper.ResourceData, key string) (add []interface{}, rm []interface{}) {
	if d.IsNewResource() {
		add = Set(d, key).List()
	}
	if d.HasChange(key) {
		o, n := d.GetChange(key)
		add = n.(*schema.Set).Difference(o.(*schema.Set)).List()
		rm = o.(*schema.Set).Difference(n.(*schema.Set)).List()
	}
	return
}

// JSON accesses the value held by key and unmarshals it into a map.
func JSON(d helper.ResourceData, key string) (m map[string]interface{}, err error) {
	if d.IsNewResource() || d.HasChange(key) {
		v, ok := d.GetOkExists(key)
		if ok {
			m, err = structure.ExpandJsonFromString(v.(string))
			if err != nil {
				return
			}
		}
	}
	return
}
