package expand_test

import (
	"strings"

	"github.com/alexkappa/terraform-plugin-helper/helper"
	"github.com/alexkappa/terraform-plugin-helper/helper/expand"
	"github.com/alexkappa/terraform-plugin-helper/internal/aws/aws-sdk-go/service/ec2"
)

var d helper.ResourceData

func ExampleSet() {

	var blockDevices []*ec2.BlockDeviceMapping

	expand.Set(d, "ebs_block_device").Elem(func(d helper.ResourceData) {

		blockDevice := &ec2.EbsBlockDevice{
			DeleteOnTermination: expand.BoolPtr(d, "delete_on_termination"),
			SnapshotId:          expand.StringPtr(d, "snapshot_id"),
			Encrypted:           expand.BoolPtr(d, "encrypted"),
			KmsKeyId:            expand.StringPtr(d, "kms_key_id"),
			VolumeSize:          expand.Int64Ptr(d, "volume_size"),
			VolumeType:          expand.StringPtr(d, "volume_type"),
		}

		if strings.ToLower(expand.String(d, "volume_type")) == ec2.VolumeTypeIo1 {
			blockDevice.Iops = expand.Int64Ptr(d, "iops")
		}

		blockDevices = append(blockDevices, &ec2.BlockDeviceMapping{
			DeviceName: expand.StringPtr(d, "device_name"),
			Ebs:        blockDevice,
		})
	})
}
