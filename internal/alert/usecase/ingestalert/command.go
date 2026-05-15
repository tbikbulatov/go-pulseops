package ingestalert

type IngestAlertCommand struct {
	IntegrationKey string
	ExternalID     string
	Service        string
	Environment    string
	Severity       string
	Name           string
	Message        string
	DedupKey       string
}
