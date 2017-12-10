package lib

import (
	"github.com/pkg/errors"
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

	evals = make([]*Evaluation, 0)

	// if currentState has something that is
	// not in the desiredState: delete
	for cHash, c := range currentMap {
		_, present = desiredMap[cHash]
		if present {
			continue
		}

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

		evals = append(evals, &Evaluation{
			Type:   EvaluationAddRecord,
			Record: d,
		})
	}

	return
}
