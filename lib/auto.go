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

// TODO paginate over all results
func (a *Auto) ListAutoscalingGroups(names []string) (asgsMap map[string]*AutoScalingGroup, err error) {
	if len(names) == 0 {
		err = errors.Errorf("names can't be empty")
		return
	}

	tagsFilter := &ec2.Filter{
		Name:   aws.String("tag:" + autoscalingGroupTag),
		Values: []*string{},
	}

	asgsMap = map[string]*AutoScalingGroup{}

	for _, name := range names {
		tagsFilter.Values = append(
			tagsFilter.Values,
			aws.String(name))
		asgsMap[name] = &AutoScalingGroup{Name: name}
	}

	var (
		input = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				tagsFilter,
			},
		}
		result *ec2.DescribeInstancesOutput
		asg    *AutoScalingGroup
		tags   map[string]string
	)

	result, err = a.ec2.DescribeInstances(input)
	if err != nil {
		err = errors.Wrapf(err, "failed to describe instances")
		return
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			tags = map[string]string{}
			asg = nil

			for _, tag := range instance.Tags {
				tags[*tag.Key] = *tag.Value
			}

			asg, _ = asgsMap[tags[autoscalingGroupTag]]
			if asg == nil {
				err = errors.Errorf(
					"couldn't find asg for instance %+v",
					instance)
				return
			}

			asg.Instances = append(asg.Instances, &Instance{
				Id:        *instance.InstanceId,
				PublicIp:  *instance.PublicIpAddress,
				PrivateIp: *instance.PrivateIpAddress,
				Tags:      tags,
				Running:   *instance.State.Name == "running",
			})
		}
	}

	return
}

// ListZones iterates over all the zones that can be
// fetched by the user. It then only shows those that
// are described by rules provided.
// TODO paginate over all results
func (a *Auto) ListZones() (zones []*Zone, err error) {
	var (
		input  = &route53.ListHostedZonesInput{}
		result *route53.ListHostedZonesOutput
	)

	result, err = a.route53.ListHostedZones(input)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to list hosted zone of account")
		return
	}

	zones = make([]*Zone, 0)
	for _, zone := range result.HostedZones {
		zones = append(zones, &Zone{
			ID:      *zone.Id,
			Name:    *zone.Name,
			Private: *zone.Config.PrivateZone,
		})
	}

	return
}

// ListZoneRecords lists the A records of a given zone
// identified by a ZoneID.
// TODO paginate over all results
func (a *Auto) ListZoneRecords(zone string) (records []*Record, err error) {
	var (
		input = &route53.ListResourceRecordSetsInput{
			HostedZoneId: aws.String(zone),
		}
		result *route53.ListResourceRecordSetsOutput
	)

	result, err = a.route53.ListResourceRecordSets(input)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to list resource records of zone %s",
			zone)
		return
	}

	records = make([]*Record, 0)
	for _, recordSet := range result.ResourceRecordSets {
		if *recordSet.Type != "A" {
			continue
		}

		record := &Record{
			Zone: zone,
			Name: *recordSet.Name,
			IPs:  []string{},
		}

		for _, resourceRecord := range recordSet.ResourceRecords {
			record.IPs = append(record.IPs, *resourceRecord.Value)
		}

		records = append(records, record)
	}

	return
}
