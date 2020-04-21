package flatten

import (
	"testing"

	"github.com/alexkappa/terraform-plugin-helper/helper"
	"github.com/alexkappa/terraform-plugin-helper/internal/testing/expect"
)

// flattener satisfies the Flattener interface and can be used with the packages
// Flatten() function.
type flattener struct{ foo string }

func (f flattener) Flatten(d helper.ResourceData) { d.Set("foo", f.foo) }

var _ Flattener = flattener{}

func TestFlatten(t *testing.T) {
	flat := Flatten(flattener{"bar"})
	expect.Expect(t, len(flat), 1)
	expect.Expect(t, flat[0].(map[string]interface{})["foo"], "bar")
	t.Logf("%v", flat) // [map[foo:bar]]
}

func TestFlattenFunc(t *testing.T) {
	flat := Func(func(d helper.ResourceData) {
		d.Set("foo", "bar")
	})
	expect.Expect(t, len(flat), 1)
	expect.Expect(t, flat[0].(map[string]interface{})["foo"], "bar")
	t.Logf("%v", flat) // [map[foo:bar]]
}

// flattenerList satisfies the List interface and can be used with the packages
// FlattenList function.
type flattenerList []flattener

func (f flattenerList) Len() int                             { return len(f) }
func (f flattenerList) Flatten(i int, d helper.ResourceData) { f[i].Flatten(d) }

var _ List = flattenerList{}

func TestList(t *testing.T) {
	flatteners := []flattener{
		{"bar"},
		{"baz"},
	}
	flat := FlattenList(flattenerList(flatteners))
	expect.Expect(t, len(flat), 2)
	expect.Expect(t, flat[0].(map[string]interface{})["foo"], "bar")
	expect.Expect(t, flat[1].(map[string]interface{})["foo"], "baz")
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
	flat := FlattenList(Flatteners(items))
	expect.Expect(t, len(flat), 2)
	expect.Expect(t, flat[0].(map[string]interface{})["name"], "bar")
	expect.Expect(t, flat[1].(map[string]interface{})["name"], "baz")
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
	flat := FlattenList(Flatteners(items))
	expect.Expect(t, len(flat), 2)
	expect.Expect(t, flat[0].(map[string]interface{})["name"], "bar")
	expect.Expect(t, flat[1].(map[string]interface{})["name"], "baz")
	t.Logf("%v", flat) // [map[name:bar] map[name:baz]]
}
