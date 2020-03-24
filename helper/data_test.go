package helper

import "testing"

func TestMapData(t *testing.T) {
	d := MapData{
		"one":   1,
		"zero":  0,
		"foo":   "foo",
		"empty": "",
		"nil":   nil,
	}

	for key, expectOk := range map[string]bool{
		"one":   true,
		"zero":  false,
		"foo":   true,
		"empty": false,
		"nil":   false,
	} {
		if _, ok := d.GetOkExists(key); ok != expectOk {
			t.Errorf("d.GetOkExists(%s) should retport ok == %t", key, expectOk)
		}
	}
}
