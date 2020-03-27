# terraform-plugin-helper

Re-usable helpers for building your next terraform provider.

## What is this?

If you've been down the path of developing a terraform provider before, you might have had to write a lot of code translating data from a terraform [`schema.ResourceData`](https://godoc.org/github.com/hashicorp/terraform-plugin-sdk/helper/schema#ResourceData) to your own API's format and the other way around.

In the [Implementing Create](https://www.terraform.io/docs/extend/writing-custom-providers.html#implement-create) and [Implementing Read](https://www.terraform.io/docs/extend/writing-custom-providers.html#implementing-read) guides, explain how to write your own logic of creating and reading upstream resources using Terraform.

The guide to [Implementing a more complex Read](https://www.terraform.io/docs/extend/writing-custom-providers.html#implementing-a-more-complex-read) is where `flattening` is introduced, and how it is used to map a nested object from an API into the `terraform.state`. The inverse is called `expanding` and it is used to map the terraform configuration to an API object.

## Flatten

The `flatten` package contains several helper functions to deal with flattening.

Using the [more complex example](https://www.terraform.io/docs/extend/writing-custom-providers.html#implementing-a-more-complex-read) from the Terraform docs, we can rewrite the `flattenTaskSpec` function as follows using helper functions from this package.

```go
func flattenTaskSpec(in *server.TaskSpec) []interface{} {
    return flatten.FlattenFunc(func(m map[string]interface{}) {
        if in.ContainerSpec != nil {
            m["container_spec"] = flatten.FlattenFunc(func(m map[string]interface{}) {
                if in.ContainerSpec.Mounts != nil && len(in.ContainerSpec.Mounts) > 0 {
                    
                }
            })
        }
    })
    
    // NOTE: the top level structure to set is a map
    m := make(map[string]interface{})
    if in.ContainerSpec != nil {
      m["container_spec"] = flattenContainerSpec(in.ContainerSpec)
    }
    /// ...

    return []interface{}{m}
}
```

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