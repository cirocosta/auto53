package lib

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	logger = zerolog.New(os.Stdout).
		With().
		Str("from", "evaluator").
		Logger()
)

// GetEvaluations retrieves a list of ordered
// evaluations to be performed on the current
// state to reach the desired state.
//
// It discovers which elements in a given array
// are missing and which ones are meant to be
// deleted.
//
// For each record, take the hash of it.
// Given the two arrays of hashed objects, calculate
// the additions and removals by looking at a map of current
// records and a map of desired records.
func GetEvaluations(current, desired []*Record) (evals []*Evaluation, err error) {
	if current == nil || desired == nil {
		err = errors.Errorf("current and desired must be non-nil")
		return
	}

	var (
		currentMap = map[uint64]*Record{}
		desiredMap = map[uint64]*Record{}
		present    bool
	)

	for _, c := range current {
		err = c.ComputeHash()
		if err != nil {
			err = errors.Wrapf(err, "failed to compute hash of record")
			return
		}

		currentMap[c.hash] = c
	}

	for _, d := range desired {
		err = d.ComputeHash()
		if err != nil {
			err = errors.Wrapf(err, "failed to compute hash of record")
			return
		}

		desiredMap[d.hash] = d
	}

	logger.Info().Interface("map", currentMap).Msg("--CURRENT")
	logger.Info().Interface("map", desiredMap).Msg("++DESIRED")

	evals = make([]*Evaluation, 0)

	// if currentState has something that is
	// not in the desiredState: delete
	for cHash, c := range currentMap {
		_, present = desiredMap[cHash]
		if present {
			continue
		}

		logger.Info().Interface("record", c).Msg("REMOVE")

		evals = append(evals, &Evaluation{
			Type:   EvaluationRemoveRecord,
			Record: c,
		})
	}

	// if desiredState has something that is
	// not in the currentState: add
	for dHash, d := range desiredMap {
		_, present = currentMap[dHash]
		if present {
			continue
		}

		logger.Info().Interface("record", d).Msg("ADD")

		evals = append(evals, &Evaluation{
			Type:   EvaluationAddRecord,
			Record: d,
		})
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
