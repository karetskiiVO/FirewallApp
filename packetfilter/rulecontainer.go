package packetfilter

type ruleContainer struct {
	DefaultAccept string           `json:"default"`
	Rules         []map[string]any `json:"rules"`
}
