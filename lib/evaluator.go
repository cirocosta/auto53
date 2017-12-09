package lib

import (
	"github.com/pkg/errors"
)

func GetEvaluations(current, desired []*Record) (evals []*Evaluation, err error) {
	if current == nil || desired == nil {
		err = errors.Errorf("current and desired must be non-nil")
		return
	}

	return
}

// CreateState takes autoscalinggroup state and
// a set of formatting rules to produce a desired
// records state.
// TODO possibly take privateIp instead of
//	publicIp
func CreateRecords(asgs []*AutoScalingGroup, rules []*FormattingRule) (records []*Record, err error) {
	if asgs == nil || rules == nil {
		err = errors.Errorf("asgs and rules must be non-nil")
		return
	}

	var (
		recordsMap      = map[string]*Record{}
		asgMap          = map[string][]*Instance{}
		fqdn            string
		templatedRecord string
	)

	records = make([]*Record, 0)

	for _, asg := range asgs {
		asgMap[asg.Name] = asg.Instances
	}

	for _, rule := range rules {
		asg := rule.AutoScalingGroup

		err = rule.ParseRecordTemplate()
		if err != nil {
			err = errors.Wrapf(err, "failed to initialize record template")
			return
		}

		instances, present := asgMap[asg]
		if !present {
			err = errors.Errorf("couldn't find asg %s for rule", asg)
			return
		}

		for _, instance := range instances {
			templatedRecord, err = rule.TemplateRecord(instance)
			if err != nil {
				err = errors.Wrapf(err, "failed to template record")
				return
			}

			fqdn = templatedRecord + "." + rule.Zone

			existingRecord, present := recordsMap[fqdn]
			if present {
				existingRecord.IPs = append(existingRecord.IPs,
					instance.PublicIp)
			} else {
				recordsMap[fqdn] = &Record{
					Zone: rule.Zone,
					Name: templatedRecord,
					IPs:  []string{instance.PublicIp},
				}
			}
		}
	}

	for _, record := range recordsMap {
		records = append(records, record)
	}

	return
}
