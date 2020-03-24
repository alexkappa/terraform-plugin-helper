package flatten

import "testing"

type flattener struct{}

func (f flattener) Flatten() map[string]interface{} {
	return map[string]interface{}{"foo": "bar"}
}

func TestFlatten(t *testing.T) {
	Flatten(flattener{})
}

func TestFlattenFunc(t *testing.T) {
	FlattenFunc(func() map[string]interface{} {
		return map[string]interface{}{"foo": "bar"}
	})
}
