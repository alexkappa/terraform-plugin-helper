package helper_test

import (
	"encoding/json"
	"fmt"

	"github.com/alexkappa/terraform-plugin-helper/helper"
	"github.com/alexkappa/terraform-plugin-helper/helper/expand"
	"github.com/alexkappa/terraform-plugin-helper/helper/flatten"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type mountList []*Mount

func (m mountList) Len() int { return len(m) }

func (m mountList) Flatten(i int, d helper.ResourceData) {
	d.Set("target", m[i].Target)
	d.Set("source", m[i].Source)
	d.Set("type", m[i].Type)
}

func ExampleFlatten() {

	d := resourceData(nil)
	api := apiData()

	if spec := api.Spec; spec != nil {

		d.Set("task_spec", flatten.Func(func(d helper.ResourceData) {

			if taskTemplate := spec.TaskTemplate; taskTemplate != nil {
				if containerSpec := taskTemplate.ContainerSpec; containerSpec != nil {

					d.Set("container_spec", flatten.Func(func(d helper.ResourceData) {

						if mounts := containerSpec.Mounts; mounts != nil {

							d.Set("mounts", flatten.List(mountList(mounts)))
						}
					}))
				}
			}
		}))
	}

	fmt.Println(d.Get("task_spec.0.container_spec.0.mounts.1606541327.target"))
	fmt.Println(d.Get("task_spec.0.container_spec.0.mounts.1606541327.source"))
	fmt.Println(d.Get("task_spec.0.container_spec.0.mounts.1606541327.type"))
	// Output: /mount/test
	// tftest-volume
	// volume
}

func ExampleExpand() {

	d := resourceData(raw)

	api := &API{}
	api.Spec = &Spec{}

	expand.List(d, "task_spec").Elem(func(d helper.ResourceData) {

		api.Spec.TaskTemplate = &TaskTemplate{}
		api.Spec.TaskTemplate.ContainerSpec = &ContainerSpec{}

		expand.List(d, "container_spec").Elem(func(d helper.ResourceData) {

			api.Spec.TaskTemplate.ContainerSpec.Mounts = make([]*Mount, 0)

			expand.Set(d, "mounts").Elem(func(d helper.ResourceData) {
				api.Spec.TaskTemplate.ContainerSpec.Mounts = append(api.Spec.TaskTemplate.ContainerSpec.Mounts, &Mount{
					Target: expand.String(d, "target"),
					Source: expand.String(d, "source"),
					Type:   expand.String(d, "type"),
				})
			})
		})
	})

	fmt.Println(api.Spec.TaskTemplate.ContainerSpec.Mounts[0].Target)
	fmt.Println(api.Spec.TaskTemplate.ContainerSpec.Mounts[0].Source)
	fmt.Println(api.Spec.TaskTemplate.ContainerSpec.Mounts[0].Type)
	// Output: /mount/test
	// tftest-volume
	// volume
}

var raw = map[string]interface{}{
	"task_spec": []interface{}{
		map[string]interface{}{
			"container_spec": []interface{}{
				map[string]interface{}{
					"mounts": []interface{}{
						map[string]interface{}{
							"target": "/mount/test",
							"source": "tftest-volume",
							"type":   "volume",
						},
					},
				},
			},
		},
	},
}

var s = map[string]*schema.Schema{
	"task_spec": &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"container_spec": &schema.Schema{
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"mounts": &schema.Schema{
								Type:     schema.TypeSet,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"target": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"source": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"type": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"volume_options": &schema.Schema{
											Type:     schema.TypeList,
											Optional: true,
											MaxItems: 1,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"no_copy": &schema.Schema{
														Type:     schema.TypeBool,
														Optional: true,
													},
													"labels": &schema.Schema{
														Type:     schema.TypeMap,
														Optional: true,
														Elem:     &schema.Schema{Type: schema.TypeString},
													},
													"driver_name": &schema.Schema{
														Type:     schema.TypeString,
														Optional: true,
													},
													"driver_options": &schema.Schema{
														Type:     schema.TypeMap,
														Optional: true,
														Elem:     &schema.Schema{Type: schema.TypeString},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

type API struct {
	ID   string
	Spec *Spec
}

type Spec struct {
	Name         string
	Labels       map[string]string
	Address      string
	TaskTemplate *TaskTemplate
}

type TaskTemplate struct {
	ContainerSpec *ContainerSpec
}

type ContainerSpec struct {
	Mounts []*Mount
}

type Mount struct {
	Type          string
	Source        string
	Target        string
	VolumeOptions *VolumeOptions
}

type VolumeOptions struct {
	NoCopy       bool
	DriverConfig map[string]string
}

const apiRaw = `
{
	"ID": "ozfsuj7dblzwjo8zoguosr1l5",
	"Spec": {
		"Name": "tftest-service-basic",
		"Labels": {},
		"Address": "tf-test-address",
		"TaskTemplate": {
			"ContainerSpec": {
				"Mounts": [{
					"Type": "volume",
					"Source": "tftest-volume",
					"Target": "/mount/test",
					"VolumeOptions": {
						"NoCopy": true,
						"DriverConfig": {}
					}
				}]
			}
		}
	}
}
`

func resourceData(raw map[string]interface{}) *schema.ResourceData {

	c := terraform.NewResourceConfigRaw(raw)

	sm := schema.InternalMap(s)
	diff, err := sm.Diff(nil, c, nil, nil, true)
	if err != nil {
		panic(err)
	}

	result, err := sm.Data(nil, diff)
	if err != nil {
		panic(err)
	}

	return result
}

func apiData() (api *API) {
	err := json.Unmarshal([]byte(apiRaw), &api)
	if err != nil {
		panic(err)
	}
	return
}
