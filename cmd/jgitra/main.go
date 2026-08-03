package main

import (
	"fmt"
	"log"

	"jgitra/internal/cli"
	"jgitra/internal/client"
	"jgitra/internal/config"

	"github.com/alecthomas/kong"
)

type VersionFlag bool

func (v VersionFlag) BeforeReset(app *kong.Kong) error {
	fmt.Printf("Version %s\n", config.Version)
	app.Exit(0)
	return nil
}

type Globals struct {
	Version VersionFlag `name:"version" short:"v" help:"Show version."`
}

type CLI struct {
	Globals

	Validate cli.ValidateCmd `cmd:"" help:"Validate configuration."`
}

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal("error loading config: ", err)
	}
	jira := client.NewJiraClient(c)
	ctx := kong.Parse(&CLI{},
		kong.Name(config.App),
		kong.Description("jira + git tool"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Tree: true,
		}),
	)
	ctx.FatalIfErrorf(ctx.Run(jira))
}
