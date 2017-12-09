package lib

import (
	"github.com/pkg/errors"
)

func GetEvaluations(current, desired *State) (evals []*Evaluation, err error) {
	if current == nil || desired == nil {
		err = errors.Errorf("current and desired must be non-nil")
		return
	}

	return
}

// CreateState takes autoscalinggroup state and
// a set of formatting rules to produce a desired
// records state.
func CreateState(asgs []*AutoScalingGroup, rules []*FormattingRule) (state State, err error) {
	if asgs == nil || rules == nil {
		err = errors.Errorf("asgs and rules must be non-nil")
		return
	}

	return
}
