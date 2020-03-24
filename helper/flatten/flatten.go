package flatten

// A Flattener is used to flatten data into Terraform's internal representation.
type Flattener interface {
	Flatten() map[string]interface{}
}

// The FlattenerFunc type is an adapter to allow the use of an ordinary function
// as a Flattener. If f is a function with the appropriate signature,
// FlattenerFunc(f) is a Flattener that calls f.
type FlattenerFunc func() map[string]interface{}

// Flatten calls f().
func (fn FlattenerFunc) Flatten() map[string]interface{} {
	return fn()
}

// Flatten executes the provided flatteners Flatten method and wraps the result
// in a []interface{} which is used by Terraform list or set types.
func Flatten(f Flattener) []interface{} {
	return []interface{}{f.Flatten()}
}

// FlattenFunc executes the provided function and wraps the result in a
// []interface{} which is used by Terraform list or set types.
func FlattenFunc(fn func() map[string]interface{}) []interface{} {
	return []interface{}{Flatten(FlattenerFunc(fn))}
}
