package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEvaluations(t *testing.T) {
	var testCases = []struct {
		desc        string
		current     []*Record
		desired     []*Record
		expected    []*Evaluation
		shouldError bool
	}{
		{
			desc:        "fail if nil current",
			current:     nil,
			shouldError: true,
		},
		{
			desc:        "fail if nil desired",
			current:     []*Record{},
			desired:     nil,
			shouldError: true,
		},
		{
			desc:        "success w/ no evaluations if both empty",
			current:     []*Record{},
			desired:     []*Record{},
			expected:    []*Evaluation{},
			shouldError: false,
		},
		{
			desc:    "additions if current is empty and one in desired",
			current: []*Record{},
			desired: []*Record{
				{
					Zone: "apex1",
					Name: "record1",
					IPs:  []string{"1.1.1.1"},
				},
			},
			expected: []*Evaluation{
				{
					Type: EvaluationAddRecord,
					Record: &Record{
						Zone: "apex1",
						Name: "record1",
						IPs:  []string{"1.1.1.1"},
					},
				},
			},
			shouldError: false,
		},
		{
			desc: "nothing if both equal",
			current: []*Record{
				{
					Zone: "apex1",
					Name: "record1",
					IPs:  []string{"1.1.1.1"},
				},
			},
			desired: []*Record{
				{
					Zone: "apex1",
					Name: "record1",
					IPs:  []string{"1.1.1.1"},
				},
			},
			expected:    []*Evaluation{},
			shouldError: false,
		},
		{
			desc: "update if ip-set changes with removal and addition",
			current: []*Record{
				{
					Zone: "apex1",
					Name: "record1",
					IPs:  []string{"1.1.1.1"},
				},
			},
			desired: []*Record{
				{
					Zone: "apex1",
					Name: "record1",
					IPs:  []string{"2.2.2.2"},
				},
			},
			expected: []*Evaluation{
				{
					Type: EvaluationRemoveRecord,
					Record: &Record{
						Zone: "apex1",
						Name: "record1",
						IPs:  []string{"1.1.1.1"},
					},
				},
				{
					Type: EvaluationAddRecord,
					Record: &Record{
						Zone: "apex1",
						Name: "record1",
						IPs:  []string{"2.2.2.2"},
					},
				},
			},
			shouldError: false,
		},
		{
			desc: "update if ip-set changes with addition",
			current: []*Record{
				{
					Zone: "apex1",
					Name: "record1",
					IPs:  []string{"1.1.1.1"},
				},
			},
			desired: []*Record{
				{
					Zone: "apex1",
					Name: "record1",
					IPs:  []string{"1.1.1.1", "2.2.2.2"},
				},
			},
			expected: []*Evaluation{
				{
					Type: EvaluationRemoveRecord,
					Record: &Record{
						Zone: "apex1",
						Name: "record1",
						IPs:  []string{"1.1.1.1"},
					},
				},
				{
					Type: EvaluationAddRecord,
					Record: &Record{
						Zone: "apex1",
						Name: "record1",
						IPs:  []string{"1.1.1.1", "2.2.2.2"},
					},
				},
			},
			shouldError: false,
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
