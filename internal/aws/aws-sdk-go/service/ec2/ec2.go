package ec2

const (
	VolumeTypeStandard = "standard"
	VolumeTypeIo1      = "io1"
	VolumeTypeGp2      = "gp2"
	VolumeTypeSc1      = "sc1"
	VolumeTypeSt1      = "st1"
)

type BlockDeviceMapping struct {
	DeviceName *string
	Ebs        *EbsBlockDevice
}

type EbsBlockDevice struct {
	DeleteOnTermination *bool
	SnapshotId          *string
	Encrypted           *bool
	KmsKeyId            *string
	VolumeSize          *int64
	VolumeType          *string
	Iops                *int64
}
