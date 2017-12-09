package lib

type Zone struct {
	Name    string
	Formats []*AutoScalingGroupFormat
}

type Machine struct {
	Id        string
	PublicIp  string
	PrivateIp string
}

// AutoScalingGroupFormat wraps an autoscaling group
// with a given template that is used to create records
// in a zone.
//
// A given AWS ASG can have more than one AutoScalingGroupFormat
// if desired, which would lead to the generation of multiple
// records for the machines in that ASG.
type AutoScalingGroupFormat struct {
	// Name corresponds to the name of
	// the autoscaling group as registered
	// in AWS.
	// This is used to match an instance's autoscaling
	// group with records to be created.
	Name string

	// Record is a template that is used
	// as the name for the entry in the zone.
	// ps.: It can use the properties of the Machine
	// type.
	//
	// For instance:
	//	instance-{{ .Id }}-asg1
	// would be formatted as:
	//	instance-i-012931-asg1 for a machine
	// with the id `i-012931`.
	Record string
}
