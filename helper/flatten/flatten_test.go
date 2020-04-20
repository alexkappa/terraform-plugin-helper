package flatten

import (
	"testing"

	"github.com/alexkappa/terraform-plugin-helper/helper"
)

type flattener struct {
	foo string
}

func (f flattener) Flatten(d helper.ResourceData) {
	d.Set("foo", f.foo)
}

func TestFlatten(t *testing.T) {
	flat := Flatten(flattener{"bar"})
	t.Logf("%v", flat) // [map[foo:bar]]
}

func TestFlattenFunc(t *testing.T) {
	flat := Func(func(d helper.ResourceData) {
		d.Set("foo", "bar")
	})
	t.Logf("%v", flat) // [map[foo:bar]]
}

type flattenerList []flattener

func (f flattenerList) Len() int                             { return len(f) }
func (f flattenerList) Flatten(i int, d helper.ResourceData) { f[i].Flatten(d) }

func TestList(t *testing.T) {
	flatteners := []flattener{
		{"bar"},
		{"baz"},
	}
	flat := List(flattenerList(flatteners))
	t.Logf("%v", flat) // [map[foo:bar] map[foo:baz]]
}

type item struct{ name string }

type itemFlattener item

func (i itemFlattener) Flatten(d helper.ResourceData) { d.Set("name", i.name) }

func TestListWrap(t *testing.T) {
	items := []Flattener{
		itemFlattener(item{"bar"}),
		itemFlattener(item{"baz"}),
	}
	flat := List(Flatteners(items))
	t.Logf("%v", flat) // [map[name:bar] map[name:baz]]
}

func itemFlattenerFunc(i item) Flattener {
	return FlattenerFunc(func(d helper.ResourceData) {
		d.Set("name", i.name)
	})
}

func TestListWrapFunc(t *testing.T) {
	items := []Flattener{
		itemFlattenerFunc(item{"bar"}),
		itemFlattenerFunc(item{"baz"}),
	}
	flat := List(Flatteners(items))
	t.Logf("%v", flat) // [map[name:bar] map[name:baz]]
}
