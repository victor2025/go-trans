package runner

import (
	"log"
	"time"
)

/**
@author: victor2022
@since: 2025/1/27
*/

// TickerRunner 周期执行器
type TickerRunner struct {
	Runner
	Period time.Duration
}

func NewTimedRunner(period time.Duration, runner Runner) *TickerRunner {
	return &TickerRunner{
		Runner: runner,
		Period: period,
	}
}

func (t *TickerRunner) Run() {
	ticker := time.NewTicker(t.Period)
	defer ticker.Stop()
	t.MarkStarted()
	for range ticker.C {
		if !t.IsOn() {
			break
		}
		t.Handle()
	}
	log.Println("processor stopped")
}
