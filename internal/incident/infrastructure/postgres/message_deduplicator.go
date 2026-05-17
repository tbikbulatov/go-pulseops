package postgres

import (
	"context"
	"time"

	"github.com/tbikbulatov/go-pulseops/internal/platform/transaction"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MessageDeduplicator struct {
	db *gorm.DB
}

func NewMessageDeduplicator(db *gorm.DB) *MessageDeduplicator {
	return &MessageDeduplicator{db: db}
}

func (d *MessageDeduplicator) TryStartProcessing(ctx context.Context, consumerName string, messageID string) (bool, error) {
	model := ProcessedMessageModel{
		ConsumerName: consumerName,
		MessageID:    messageID,
		ProcessedAt:  time.Now(),
	}

	db := transaction.FromContext(ctx, d.db)
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "consumer_name"}, {Name: "message_id"}},
		DoNothing: true,
	}).Create(&model)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 1, nil
}
