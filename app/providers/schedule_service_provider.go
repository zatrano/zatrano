package providers

import (
	"github.com/zatrano/framework/core"
)

// ScheduleServiceProvider registers scheduled tasks.
type ScheduleServiceProvider struct{}

func (p *ScheduleServiceProvider) Register(app *core.Application) {}

func (p *ScheduleServiceProvider) Boot(app *core.Application) {
	// Register scheduled commands here, e.g.:
	// app.Scheduler().Command("reports", fn).Daily()
}
