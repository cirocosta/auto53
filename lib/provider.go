package lib

type Provider interface {
	GetAutoscalingGroups() (map[string]*AutoScalingGroup, error)
	GetZonesRecords() (map[string][]*Record, error)
	ExecuteEvaluations([]*Evaluation) error
}
