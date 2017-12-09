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
