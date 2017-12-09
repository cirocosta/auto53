package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEvaluations(t *testing.T) {
	var testCases = []struct {
		desc        string
		current     *State
		desired     *State
		expected    []*Evaluation
		shouldError bool
	}{
		{
			desc:        "nil current",
			shouldError: true,
		},
		{
			desc:        "nil desired",
			shouldError: true,
		},
	}

	var (
		evals []*Evaluation
		err   error
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			evals, err = GetEvaluations(tc.current, tc.desired)
			if tc.shouldError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(evals))
		})
	}
}
