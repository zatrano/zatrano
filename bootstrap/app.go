package bootstrap

import (
	"github.com/zatrano/framework/app/providers"
	"github.com/zatrano/framework/core"
)

// App creates and configures the ZATRANO application.
func App() *core.Application {
	basePath, _ := findBasePath()
	application := core.NewApplication(basePath)

	application.RegisterProviders(
		&providers.AppServiceProvider{},
		&providers.DatabaseServiceProvider{},
		&providers.AuthServiceProvider{},
		&providers.EventServiceProvider{},
		&providers.ScheduleServiceProvider{},
		&providers.RouteServiceProvider{},
	)

	return application
}

func findBasePath() (string, error) {
	wd, err := lookWorkingDirectory()
	if err != nil {
		return ".", err
	}
	return wd, nil
}
