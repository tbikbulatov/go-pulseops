package ingestalert

import (
	"context"
)

type IngestAlertUsecase interface {
	Handle(ctx context.Context, cmd IngestAlertCommand) (IngestAlertResult, error)
}
