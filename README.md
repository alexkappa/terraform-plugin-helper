# terraform-plugin-helper

Re-usable helpers for building your next terraform provider.

## What is this?

If you've been down the path of developing a terraform provider before, you might have had to write a lot of code to translate data from a terraform `schema.ResourceData` to your own API's format.

Here's a hypothetical `schema.CreateFunc` we've defined to create an example resource using a hypothetical API. 

```go
func resourceExampleCreate(d *schema.ResourceData, m interface{}) error {
    e := expandExample(d)

    e, err := m.(*ExampleAPI).Create(e)
    if err != nil {
		return err
	}
    
    d.SetId(e.ID)

	return resourceExampleRead(d, m)
}
```

We would define our `expandExample` function similar to the following.

```go
func expandExample(p []interface{}) (*Example, error) {
	e := &Example{}
	if len(p) == 0 || p[0] == nil {
		return e, nil
	}
    in := p[0].(map[string]interface{})
    
    if v, ok := in["name"].(string); ok && v != "" {
		e.Name = ptrToString(v)
	}

	if v, ok := in["age"].(int); ok && v > 0 {
		e.Age = ptrToInt(v)
	}

	if v, ok := in["address"].([]interface{}); ok && len(v) > 0 {
		a, err := expandExampleAddress(v)
		if err != nil {
			return e, err
		}
		e.Address = a
    }
    
    return e
}

func expandExampleAddress(p []interface{}) (*ExampleAddress, error) {
    e := &ExampleAddress{}
	if len(p) == 0 || p[0] == nil {
		return e, nil
	}
    in := p[0].(map[string]interface{})

    if v, ok := in["country"].(string); ok && v != "" {
		e.Country = ptrToString(v)
    }
    
    if v, ok := in["street"].(string); ok && v != "" {
		e.Street = ptrToString(v)
    }

    if v, ok := in["house_number"].(int); ok && v > 0 {
		e.HouseNumber = ptrToInt(v)
	}
    
    if v, ok := in["postal_code"].(string); ok && v != "" {
		e.PostalCode = ptrToString(v)
    }

    return e
}
```

... but it is a little repetitive and hard n the eyes. Using this package, the above can be rewritten to the following.

```go
func expandExample(d *schema.ResourceData) *Example {
    
    e := &Example{
        Name: String(d, "name"),
        Age: Int(d, "age"),
    }

    expand.List(d, "address").Elem(func(d Data) {
        e.Address = &Address{
            Country: expand.StringPtr(d, "country"),
            Street: expand.StringPtr(d, "street_address"),
            HouseNumber: expand.IntPtr(d, "house_number"),
            PostalCode: expand.StringPtr(d, "postal_code"),
        }
    })

    return e
}
```

## Is that it?

Yes! Well, no. There is also a `flattener` package to help with the inverse operation of flattening from your API's format to a Terraform `schema.ResourceData` type.

```go
```