package main

import (
	"fmt"
	"os"

	appconsole "github.com/zatrano/framework/app/console"
	"github.com/zatrano/framework/bootstrap"
	"github.com/zatrano/framework/core/console"
)

func main() {
	app := bootstrap.App()
	cli := console.New(app)
	appconsole.Register(cli, app)

	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
