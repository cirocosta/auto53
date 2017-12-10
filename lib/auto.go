package lib

import (
	"os"
	"strings"

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
	Debug           bool
}

func NewAuto(cfg AutoConfig) (a Auto, err error) {
	if len(cfg.FormattingRules) == 0 {
		err = errors.Errorf("FormattingRules must be specified")
		return
	}

	a.formattingRules = cfg.FormattingRules
	a.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "auto").
		Logger()

	var awsConfig = &aws.Config{}

	if cfg.Debug {
		awsConfig.LogLevel =
			aws.LogLevel(aws.LogDebug | aws.LogDebugWithRequestErrors)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		err = errors.Wrapf(err, "failed to create aws session")
		return
	}

	a.route53 = route53.New(sess)
	a.ec2 = ec2.New(sess)

	return
}

const (
	autoscalingGroupTag = "aws:autoscaling:groupName"
	runningState        = "running"
)

// TODO paginate over all results
func (a *Auto) GetAutoScalingGroups() (asgsMap map[string]*AutoScalingGroup, err error) {
	var present bool

	tagsFilter := &ec2.Filter{
		Name:   aws.String("tag:" + autoscalingGroupTag),
		Values: []*string{},
	}

	asgsMap = map[string]*AutoScalingGroup{}

	for _, rule := range a.formattingRules {
		if rule.AutoScalingGroup == "" {
			err = errors.Errorf(
				"Rule %+v does not have an autoscalinggroup specified",
				rule)
			return
		}

		_, present = asgsMap[rule.AutoScalingGroup]
		if present {
			continue
		}

		asgsMap[rule.AutoScalingGroup] = &AutoScalingGroup{
			Name: rule.AutoScalingGroup,
		}

		tagsFilter.Values = append(
			tagsFilter.Values,
			aws.String(rule.AutoScalingGroup))
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
				Running:   *instance.State.Name == runningState,
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

func (a *Auto) GetZonesRecords() (recordsMap map[string][]*Record, err error) {
	var (
		present bool
		records []*Record
	)

	recordsMap = map[string][]*Record{}

	for _, rule := range a.formattingRules {
		if rule.Zone == "" {
			err = errors.Errorf(
				"Rule %+v does not have a zone specified",
				rule)
			return
		}

		_, present = recordsMap[rule.Zone]
		if present {
			continue
		}

		records, err = a.ListZoneRecords(rule.Zone)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to retrieve records from zone %s",
				rule.Zone)
			return
		}

		recordsMap[rule.Zone] = records
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
		result   *route53.ListResourceRecordSetsOutput
		zoneName string
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
		if *recordSet.Type == "SOA" {
			zoneName = "." + *recordSet.Name
			break
		}
	}

	if zoneName == "" {
		err = errors.Errorf(
			"couldn't find SOA record fone zone %s",
			zone)
		return
	}

	for _, recordSet := range result.ResourceRecordSets {
		if *recordSet.Type != "A" {
			continue
		}

		record := &Record{
			Zone: zone,
			Name: strings.TrimSuffix(*recordSet.Name, zoneName),
			IPs:  []string{},
		}

		for _, resourceRecord := range recordSet.ResourceRecords {
			record.IPs = append(record.IPs, *resourceRecord.Value)
		}

		records = append(records, record)
	}

	return
}
