package manageincident

type AcknowledgeCommand struct {
	IncidentID     string
	CommandID      string
	Actor          string
	ExpectedStatus string
}

func (c *AcknowledgeCommand) toMap() map[string]any {
	return map[string]any{
		"incident_id":     c.IncidentID,
		"command_id":      c.CommandID,
		"actor":           c.Actor,
		"expected_status": c.ExpectedStatus,
	}
}

type ResolveCommand struct {
	IncidentID     string
	CommandID      string
	Actor          string
	Reason         string
	ExpectedStatus string
}

func (c *ResolveCommand) toMap() map[string]any {
	return map[string]any{
		"incident_id":     c.IncidentID,
		"command_id":      c.CommandID,
		"actor":           c.Actor,
		"reason":          c.Reason,
		"expected_status": c.ExpectedStatus,
	}
}
