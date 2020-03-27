package expand

import (
	"reflect"
	"testing"

	helper "github.com/alexkappa/terraform-plugin-helper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var s = map[string]*schema.Schema{
	"string": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"int": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"bool": {
		Type:     schema.TypeBool,
		Optional: true,
	},
	"map": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"list": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"bar": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	},
	"set": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"foo": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"bar": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	},
}

func TestExpand(t *testing.T) {

	r := map[string]interface{}{
		"string": "hello!",
		"int":    123,
		"bool":   true,
		"map":    map[string]interface{}{"foo": "bar"},
		"list": []interface{}{
			map[string]interface{}{
				"foo": "bar",
				"bar": 123,
			},
		},
		"set": []interface{}{
			map[string]interface{}{
				"foo": "bar",
				"bar": 123,
			},
		},
	}
	d := schema.TestResourceDataRaw(t, s, r)

	Expect(t, String(d, "string"), "hello!")
	Expect(t, Int(d, "int"), 123)
	Expect(t, Bool(d, "bool"), true)
	Expect(t, Map(d, "map"), map[string]interface{}{"foo": "bar"})

	Expect(t, String(d, "list.0.foo"), "bar")
	Expect(t, Int(d, "list.0.bar"), 123)

	Expect(t, String(d, "set.1122208398.foo"), "bar")
	Expect(t, Int(d, "set.1122208398.bar"), 123)

	var it Iterator

	it = List(d, "list")
	it.Elem(func(d helper.Data) {
		Expect(t, String(d, "foo"), "bar")
		Expect(t, Int(d, "bar"), 123)
	})
	Expect(t, it.List(), []interface{}{
		map[string]interface{}{
			"foo": "bar",
			"bar": 123,
		},
	})

	it = Set(d, "set")
	it.Elem(func(d helper.Data) {
		Expect(t, String(d, "foo"), "bar")
		Expect(t, Int(d, "bar"), 123)
	})
	Expect(t, it.List(), []interface{}{
		map[string]interface{}{
			"foo": "bar",
			"bar": 123,
		},
	})
}

func TestJSON(t *testing.T) {
	d := helper.MapData{"json": `{"foo": 123}`}
	v, err := JSON(d, "json")
	if err != nil {
		t.Error(err)
	}
	j, ok := v["foo"]
	if !ok {
		t.Errorf("Expected result to be a int, instead it was %T\n", j)
	}
}

func Expect(t *testing.T, x, y interface{}) bool {
	xv := reflect.ValueOf(x)
	if xv.Kind() == reflect.Ptr {
		xv = xv.Elem()
	}
	if !reflect.DeepEqual(xv.Interface(), y) {
		t.Errorf("Expected %v to equal %v\n", xv.Interface(), y)
		return false
	}
	return true
}
