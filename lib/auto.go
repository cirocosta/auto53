package lib

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Auto struct {
	logger          zerolog.Logger
	route53         *route53.Route53
	ec2             *ec2.EC2
	formattingRules []*FormattingRule
}

type AutoConfig struct {
	FormattingRules []*FormattingRule
}

func NewAuto(cfg AutoConfig) (a Auto, err error) {
	a.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "auto").
		Logger()

	sess, err := session.NewSession()
	if err != nil {
		err = errors.Wrapf(err, "failed to create aws session")
		return
	}

	a.route53 = route53.New(sess)
	a.ec2 = ec2.New(sess)

	return
}

const autoscalingGroupTag = "aws:autoscaling:groupName"

func (a *Auto) GetAutoscalingGroup(name string) (asg *AutoScalingGroup, err error) {
	if name == "" {
		err = errors.Errorf("name can't be nil")
		return
	}

	var (
		input = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				&ec2.Filter{
					Name: aws.String("tag:" + autoscalingGroupTag),
					Values: []*string{
						aws.String(name),
					},
				},
			},
		}
		result *ec2.DescribeInstancesOutput
	)

	result, err = a.ec2.DescribeInstances(input)
	if err != nil {
		err = errors.Wrapf(err, "failed to describe instances")
		return
	}

	asg = &AutoScalingGroup{
		Name:      name,
		Instances: []*Instance{},
	}
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			asg.Instances = append(asg.Instances, &Instance{
				Id:        *instance.InstanceId,
				PublicIp:  *instance.PublicIpAddress,
				PrivateIp: *instance.PrivateIpAddress,
			})
		}
	}

	return
}

func (a *Auto) ListZoneRecords(zone string) (records []*Record, err error) {
	return
}
