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

func TestCreateState(t *testing.T) {
	var testCases = []struct {
		desc        string
		asgs        []*AutoScalingGroup
		rules       []*FormattingRule
		expected    []*Record
		shouldError bool
	}{
		{
			desc:        "nil should fail",
			shouldError: true,
		},
		{
			desc: "single instance without formatting",
			asgs: []*AutoScalingGroup{
				{
					Name: "asg1",
					Instances: []*Instance{
						{
							Id:       "inst1",
							PublicIp: "1.1.1.1",
						},
					},
				},
			},
			rules: []*FormattingRule{
				{
					AutoScalingGroup: "asg1",
					Zone:             "apex1",
					Record:           "aaa",
				},
			},
			expected: []*Record{
				{
					Zone: "apex1",
					Name: "aaa",
					IPs: []string{
						"1.1.1.1",
					},
				},
			},
			shouldError: false,
		},
		{
			desc: "multiple instances without formatting",
			asgs: []*AutoScalingGroup{
				{
					Name: "asg1",
					Instances: []*Instance{
						{
							Id:       "inst1",
							PublicIp: "1.1.1.1",
						},
						{
							Id:       "inst1",
							PublicIp: "1.1.1.2",
						},
					},
				},
			},
			rules: []*FormattingRule{
				{
					AutoScalingGroup: "asg1",
					Zone:             "apex1",
					Record:           "aaa",
				},
			},
			expected: []*Record{
				{
					Zone: "apex1",
					Name: "aaa",
					IPs: []string{
						"1.1.1.1",
						"1.1.1.2",
					},
				},
			},
			shouldError: false,
		},
		{
			desc: "with formatting and multiple instances",
			asgs: []*AutoScalingGroup{
				{
					Name: "asg1",
					Instances: []*Instance{
						{
							Id:       "inst1",
							PublicIp: "1.1.1.1",
						},
						{
							Id:       "inst2",
							PublicIp: "1.1.1.2",
						},
					},
				},
			},
			rules: []*FormattingRule{
				{
					AutoScalingGroup: "asg1",
					Zone:             "apex1",
					Record:           "{{ .Id }}-asg1",
				},
			},
			expected: []*Record{
				{
					Zone: "apex1",
					Name: "inst1-asg1",
					IPs: []string{
						"1.1.1.1",
					},
				},
				{
					Zone: "apex1",
					Name: "inst2-asg1",
					IPs: []string{
						"1.1.1.2",
					},
				},
			},
			shouldError: false,
		},
	}

	var (
		records        []*Record
		expectedRecord *Record
		expectedIP     string
		err            error
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			records, err = CreateRecords(tc.asgs, tc.rules)
			if tc.shouldError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(records))

			for i, actualRecord := range records {
				expectedRecord = tc.expected[i]

				assert.Equal(t, expectedRecord.Name, actualRecord.Name)
				assert.Equal(t, expectedRecord.Zone, actualRecord.Zone)
				assert.Equal(t, len(expectedRecord.IPs), len(actualRecord.IPs))

				for k, actualIP := range actualRecord.IPs {
					expectedIP = expectedRecord.IPs[k]

					assert.Equal(t, expectedIP, actualIP)
				}
			}

		})
	}
}
