package exchange

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (e *ExchangeApp) StartSession() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		logger.LogError(err.Error(), "exchange-app")
		return "", domain.Failed("failed to generate room ID")
	}

	ctx := context.Background()

	err = e.rdb.Set(ctx, "exchange:session:"+uuid.String(), "2", 1*time.Minute).Err()
	if err != nil {
		logger.LogError(err.Error(), "exchange-app")
		return "", domain.Failed("failed to save room ID")
	}

	return uuid.String(), nil
}
