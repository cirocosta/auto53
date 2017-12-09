package lib

type EvaluationType int

const (
	EvaluationUnknown EvaluationType = iota
	EvaluationAddRecord
	EvaluationRemoveRecord
)

// Evaluation wraps an action that must be taken
// by the evaluator which acts as the system that
// mutates Route53.
type Evaluation struct {

	// Record is the record that we either add or remove
	// into/of a zone.
	Record string

	// Type is the type of evaluation to perform:
	// add or remove.
	Type EvaluationType
}

// AutoScalingGroup represents an AWS AutoScalingGroup
// filled with instances retrieved from EC2.
type AutoScalingGroup struct {

	// Name is the name of the ASG
	Name string

	// Instances contains the thin representation
	// of the set of EC2 instances that belong to
	// this ASG.
	Instances []*Instance
}

// Instance is a thin representation of
// an EC2 instance containing the values
// that can be used for formatting records.
type Instance struct {
	Id        string
	PublicIp  string
	PrivateIp string
	Tags      map[string]string

	// Running indicates whether the machine is
	// in "running" state of not.
	Running bool
}

// FormattingRule wraps an autoscaling group
// with a given template that is used to create records
// in a zone.
//
// A given AWS ASG can have more than one AutoScalingGroupFormat
// if desired, which would lead to the generation of multiple
// records for the machines in that ASG.
type FormattingRule struct {

	// AutoScalingGroup corresponds to the name of
	// the autoscaling group as registered
	// in AWS.
	// This is used to match an instance's autoscaling
	// group with records to be created.
	AutoScalingGroup string

	// Zone is the private or public zone created
	// in Route53 to use as the domain for the
	// record.
	Zone string

	// Record is a template that is used
	// as the name for the entry in the zone.
	// ps.: It can use the properties of the Instance
	// type.
	//
	// For instance:
	//	instance-{{ .Id }}-asg1
	// would be formatted as:
	//	instance-i-012931-asg1 for a machine
	// with the id `i-012931`.
	Record string
}
