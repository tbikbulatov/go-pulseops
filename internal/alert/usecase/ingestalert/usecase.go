package ingestalert

import "context"

type Usecase struct {
}

func (u *Usecase) Handle(ctx context.Context, cmd IngestAlertCommand) (IngestAlertResult, error) {
	return IngestAlertResult{AlertID: "todo", Status: "accepted"}, nil
}
