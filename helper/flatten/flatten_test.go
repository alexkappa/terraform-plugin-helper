package flatten

import "testing"

type flattener struct {
	foo string
}

func (f flattener) Flatten(m map[string]interface{}) {
	m["foo"] = f.foo
}

func TestFlatten(t *testing.T) {
	flat := Flatten(flattener{"bar"})
	t.Logf("%v", flat) // [map[foo:bar]]
}

func TestFlattenFunc(t *testing.T) {
	flat := FlattenFunc(func(m map[string]interface{}) {
		m["foo"] = "bar"
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

func (i itemFlattener) Flatten(m map[string]interface{}) {
	m["name"] = i.name
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
	return FlattenerFunc(func(m map[string]interface{}) {
		m["name"] = i.name
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
	flat := FlattenListFunc(items, func(i interface{}, m map[string]interface{}) {
		m["name"] = i.(item).name
	})
	t.Logf("%v", flat)
}

func TestFlattenNested(t *testing.T) {
	type bar struct{ baz string }
	type foo struct{ bar *bar }

	v := &foo{&bar{"hey!"}}

	flat := FlattenFunc(func(m map[string]interface{}) {
		m["bar"] = FlattenFunc(func(m map[string]interface{}) {
			m["baz"] = v.bar.baz
		})
	})
	t.Logf("%v", flat)
}

func TestList(t *testing.T) {

}
