package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRecords(t *testing.T) {
	var testCases = []struct {
		desc        string
		asgs        map[string]*AutoScalingGroup
		rules       []*FormattingRule
		expected    []*Record
		shouldError bool
	}{
		{
			desc:        "nil should fail",
			shouldError: true,
		},
		{
			desc: "single instance single asg without formatting",
			asgs: map[string]*AutoScalingGroup{
				"asg1": {
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
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Record: "aaa",
				},
			},
			expected: []*Record{
				{
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Name: "aaa",
					IPs: []string{
						"1.1.1.1",
					},
				},
			},
			shouldError: false,
		},
		{
			desc: "multiple instances single asg without formatting",
			asgs: map[string]*AutoScalingGroup{
				"asg1": {
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
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Record: "aaa",
				},
			},
			expected: []*Record{
				{
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
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
			desc: "multiple instances multiple asgs without formatting",
			asgs: map[string]*AutoScalingGroup{
				"asg1": {
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
				"asg2": {
					Name: "asg2",
					Instances: []*Instance{
						{
							Id:       "vvvv1",
							PublicIp: "2.2.2.1",
						},
						{
							Id:       "vvvv2",
							PublicIp: "2.2.2.2",
						},
					},
				},
			},
			rules: []*FormattingRule{
				{
					AutoScalingGroup: "asg1",
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Record: "aaa",
				},
				{
					AutoScalingGroup: "asg2",
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Record: "aaa",
				},
			},
			expected: []*Record{
				{
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Name: "aaa",
					IPs: []string{
						"1.1.1.1",
						"1.1.1.2",
						"2.2.2.1",
						"2.2.2.2",
					},
				},
			},
			shouldError: false,
		},
		{
			desc: "with formatting and multiple instances",
			asgs: map[string]*AutoScalingGroup{
				"asg1": {
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
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Record: "{{ .Id }}-asg1",
				},
			},
			expected: []*Record{
				{
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
					Name: "inst1-asg1",
					IPs: []string{
						"1.1.1.1",
					},
				},
				{
					Zone: Zone{
						Name: "apex1",
						ID:   "zone123",
					},
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
