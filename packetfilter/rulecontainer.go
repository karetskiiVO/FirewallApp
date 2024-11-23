package packetfilter

type ruleContainer struct {
	DefaultAccept bool             `json:"default"`
	Rules         []map[string]any `json:"rules"`
}
