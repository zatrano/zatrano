package providers

import (
	"github.com/zatrano/framework/core"
)

// EventServiceProvider registers application events and listeners.
type EventServiceProvider struct{}

func (p *EventServiceProvider) Register(app *core.Application) {}

func (p *EventServiceProvider) Boot(app *core.Application) {
	// Register listeners here as your application grows.
}
