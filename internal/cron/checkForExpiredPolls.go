package cron

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/GhostVox/ghostvox.io-backend/internal/utils"
)

const jobName = "checkForExpiredPolls"

func UpdateExpiredPolls(ctx context.Context, q *database.Queries, logger *utils.Logger) {
	successCount := 0
	failureCount := 0
	expiredPolls, err := q.GetPollsThatHaveExpired(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.LogJob(jobName, "No expired polls found")
			return
		}
		logger.LogError(err)
		return

	}

	for _, poll := range expiredPolls {
		_, err := q.UpdatePoll(ctx, database.UpdatePollParams{
			ID:     poll.ID,
			Status: database.PollStatus(database.PollStatusArchived),
		})
		if err != nil {
			logger.LogError(fmt.Errorf("poll %s failed to update: %v", poll.ID.String(), err))
			failureCount++
			continue
		}
		successCount++
	}
	logger.LogJob(jobName, fmt.Sprintf("Processed %d polls: %d updated successfully, %d failed",
		len(expiredPolls), successCount, failureCount))

}
