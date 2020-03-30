package helper

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// The ResourceData interface represents a schema.ResourceData and defines some
// of its methods that are commonly used by the helper packages.
type ResourceData interface {

	// IsNewResource reports whether or not the resource is seen for the first
	// time.
	IsNewResource() bool

	// HasChange reports whether or not the given key has been changed.
	HasChange(key string) bool

	// GetChange returns the old and new value for a given key.
	GetChange(key string) (interface{}, interface{})

	// Get returns the data for the given key, or nil if the key doesn't exist
	// in the schema.
	Get(key string) interface{}

	// GetOk returns the data for the given key and whether or not the key
	// has been set to a non-zero value at some point.
	GetOk(key string) (interface{}, bool)

	// GetOkExists returns the data for a given key and whether or not the key
	// has been set to a non-zero value. This is only useful for determining
	// if boolean attributes have been set, if they are Optional but do not
	// have a Default value.
	GetOkExists(key string) (interface{}, bool)

	// Set sets the value for the given key.
	//
	// If the key is invalid or the value is not a correct type, an error
	// will be returned.
	Set(key string, value interface{}) error
}

var _ ResourceData = (*schema.ResourceData)(nil)

// MapData wraps a map satisfying the Data interface, so it can be used in the
// accessor methods defined below.
//
// It is not possible to fully mirror the functionality of Data as some
// information available to schema.ResourceData is lost when dealing with maps.
type MapData map[string]interface{}

// IsNewResource always reports false.
func (md MapData) IsNewResource() bool {
	return false
}

// HasChange reports whether the key exists in the map.
func (md MapData) HasChange(key string) bool {
	_, ok := md[key]
	return ok
}

// GetChange returns the old and new value for a given key. The old and new
// values will always be the same.
func (md MapData) GetChange(key string) (interface{}, interface{}) {
	return md[key], md[key]
}

// Get returns the data for the given key, or nil if the key doesn't exist in
// the map.
func (md MapData) Get(key string) interface{} {
	return md[key]
}

// GetOk returns the data for the given key and whether or not the key has been
// set to a non-zero value at some point.
func (md MapData) GetOk(key string) (interface{}, bool) {
	v, ok := md[key]
	return v, ok
}

// GetOkExists returns the data for a given key and whether or not the key has
// been set to a non-nil and non-zero value.
func (md MapData) GetOkExists(key string) (interface{}, bool) {
	v, ok := md[key]
	return v, ok && !isNil(v) && !isZero(v)
}

// Set sets the value for the given key.
func (md MapData) Set(key string, value interface{}) error {
	md[key] = value
	return nil
}

func isNil(v interface{}) bool {
	return v == nil
}

func isZero(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}
