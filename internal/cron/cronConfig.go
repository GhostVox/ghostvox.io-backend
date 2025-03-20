package cron

import (
	"bytes"
	"context"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/utils"
	"github.com/robfig/cron/v3"
)

type CronConfig struct {
	CheckForExpiredPolls string
	Jobs                 map[string]cron.EntryID
	Scheduler            *cron.Cron
	logger               *utils.Logger
}

func NewCronConfig(checkForExpiredPolls string) *CronConfig {
	buffer := bytes.NewBuffer([]byte{})

	return &CronConfig{
		CheckForExpiredPolls: checkForExpiredPolls,
		Jobs:                 make(map[string]cron.EntryID),
		Scheduler: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
			cron.SkipIfStillRunning(cron.DefaultLogger),
		)),
		logger: utils.NewLogger(buffer),
	}
}

func (c *CronConfig) StartCronJobs(ctx context.Context, cfg *config.APIConfig) {
	if c.Scheduler == nil {
		c.Scheduler = cron.New()
	}
	if c.Jobs == nil {
		c.Jobs = make(map[string]cron.EntryID)
	}
	updatePollJobID, err := c.Scheduler.AddFunc(c.CheckForExpiredPolls, func() {
		jobCtx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()
		UpdateExpiredPolls(jobCtx, cfg.Queries, c.logger)
	})
	if err != nil {
		c.logger.LogError(err)
		return
	}
	c.Jobs["updatePolls"] = updatePollJobID

	c.Scheduler.Start()

}

func (c *CronConfig) StopJobs() {
	if c.Scheduler != nil {
		c.Scheduler.Stop()
	}
}
