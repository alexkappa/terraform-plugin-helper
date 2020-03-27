package flatten

import (
	"testing"

	"github.com/alexkappa/terraform-plugin-helper/helper"
)

type flattener struct {
	foo string
}

func (f flattener) Flatten(d helper.Data) {
	d.Set("foo", f.foo)
}

func TestFlatten(t *testing.T) {
	flat := Flatten(flattener{"bar"})
	t.Logf("%v", flat) // [map[foo:bar]]
}

func TestFlattenFunc(t *testing.T) {
	flat := FlattenFunc(func(d helper.Data) {
		d.Set("foo", "bar")
	})
	t.Logf("%v", flat) // [map[foo:bar]]
}

func TestFlattenList(t *testing.T) {
	flatteners := []interface{}{
		flattener{"bar"},
		flattener{"baz"},
	}
	flat := FlattenList(flatteners)
	t.Logf("%v", flat) // [map[foo:bar] map[foo:baz]]
}

type item struct{ name string }

type itemFlattener item

func (i itemFlattener) Flatten(d helper.Data) {
	d.Set("name", i.name)
}

func TestFlattenListWrap(t *testing.T) {
	flatteners := []interface{}{
		itemFlattener(item{"bar"}),
		itemFlattener(item{"baz"}),
	}
	flat := FlattenList(flatteners)
	t.Logf("%v", flat) // [map[name:bar] map[name:baz]]
}

func itemFlattenerAlt(i item) Flattener {
	return FlattenerFunc(func(d helper.Data) {
		d.Set("name", i.name)
	})
}

func TestFlattenListWrapAlt(t *testing.T) {
	flatteners := []interface{}{
		itemFlattenerAlt(item{"bar"}),
		itemFlattenerAlt(item{"baz"}),
	}
	flat := FlattenList(flatteners)
	t.Logf("%v", flat) // [map[name:bar] map[name:baz]]
}

func TestFlattenListFunc(t *testing.T) {
	items := []item{{"foo"}, {"bar"}}
	flat := FlattenListFunc(items, func(i interface{}, d helper.Data) {
		d.Set("name", i.(item).name)
	})
	t.Logf("%v", flat)
}

func TestFlattenNested(t *testing.T) {
	type bar struct{ baz string }
	type foo struct {
		bar  *bar
		bars []*bar
	}

	v := &foo{
		bar: &bar{
			baz: "hey!",
		},
		bars: []*bar{
			{baz: "one"},
			{baz: "two"},
		},
	}

	flat := FlattenFunc(func(d helper.Data) {
		d.Set("bar", FlattenFunc(func(d helper.Data) {
			d.Set("baz", v.bar.baz)
		}))
		d.Set("bars", FlattenListFunc(v.bars, func(b interface{}, d helper.Data) {
			d.Set("baz", b.(*bar).baz)
		}))
	})
	t.Logf("%v", flat) // [map[bar:[map[baz:hey!]] bars:[map[baz:one] map[baz:two]]]]
}

func TestList(t *testing.T) {

}
