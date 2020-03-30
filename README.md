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
func resourceServerRead(d *schema.ResourceData, m interface{}) error {
  client := m.(*MyClient)
  server, ok := client.Get(d.Id())

  if !ok {
    log.Printf("[WARN] No Server found: %s", d.Id())
    d.SetId("")
    return nil
  }

  d.Set("address", server.Address)

  d.Set("task_spec", flatten.FlattenFunc(func(d helper.Data) {
  
  if taskTemplate := spec.TaskTemplate; taskTemplate != nil {
      if containerSpec := taskTemplate.ContainerSpec; containerSpec != nil {
        
        d.Set("container_spec", flatten.FlattenFunc(func (d helper.Data) {
  
          if mounts := containerSpec.Mounts; mounts != nil {
            d.Set("mounts", flatten.FlattenList(mountList(mounts)))
          }
        }))
      }
    }
  }))

  return nil
}
```

`mountList` wraps to the `[]*Mount` structure from the server response and implements the `List` interface so we can use it with the  `FlattenList` function.

```go
type mountList []*Mount

func (m mountList) Len() int { return len(m) }

func (m mountList) Flatten(i int, d helper.Data) {
	d.Set("target", m[i].Target)
	d.Set("source", m[i].Source)
	d.Set("type", m[i].Type)
}
```

In this example, we've used `flatten.FlattenFunc` mainly to flatten a nested data structure. If the data structures themselves implement `flatten.Flattener` it can flatten itself. Then the example becomes much easier, using `flatten.Flatten`.

## Expand

Using the same schema as an example, here's how we would write the `resourceServerCreate` function using the `expand` package.

```go
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

  api := &API{Spec: &Spec{}}

  expand.List(d, "task_spec").Elem(func(d helper.Data) {

    containerSpec := &ContainerSpec{}

    expand.List(d, "container_spec").Elem(func(d helper.Data) {

      mounts := make([]*Mount, 0)

      expand.Set(d, "mounts").Elem(func(d helper.Data) {

        mounts = append(mounts, &Mount{
          Target: expand.String(d, "target"),
          Source: expand.String(d, "source"),
          Type:   expand.String(d, "type"),
        })
      })

      containerSpec.Mounts = make([]*Mount, 0)
    })

    api.Spec.TaskTemplate = &TaskTemplate{containerSpec}
  })
}
```