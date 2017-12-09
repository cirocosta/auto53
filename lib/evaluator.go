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
func CreateRecords(asgs map[string]*AutoScalingGroup, rules []*FormattingRule) (records []*Record, err error) {
	if asgs == nil || rules == nil {
		err = errors.Errorf("asgs and rules must be non-nil")
		return
	}

	var (
		recordsMap      = map[string]*Record{}
		ruleAsg         string
		asg             *AutoScalingGroup
		present         bool
		fqdn            string
		templatedRecord string
	)

	records = make([]*Record, 0)

	for _, rule := range rules {
		ruleAsg = rule.AutoScalingGroup

		err = rule.ParseRecordTemplate()
		if err != nil {
			err = errors.Wrapf(err, "failed to initialize record template")
			return
		}

		asg, present = asgs[ruleAsg]
		if !present {
			err = errors.Errorf("couldn't find asg %s for rule", ruleAsg)
			return
		}

		for _, instance := range asg.Instances {
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
