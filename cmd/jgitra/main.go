package main

import (
	"fmt"
	"log"

	"jgitra/internal/cli"
	"jgitra/internal/config"
	"jgitra/internal/jira"

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
	Project  cli.ProjectCmd  `cmd:"" help:"Project configuration."`
}

func validateGit() error {
	return nil
}

func main() {
	// validating system
	if err := validateGit(); err != nil {
		log.Fatal("git validation error: ", err)
	}

	// parsing cli
	kCtx := kong.Parse(&CLI{},
		kong.Name(config.App),
		kong.Description("jira + git tool"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Tree: true,
		}),
	)

	// loading config
	c, err := config.Load()
	if err != nil {
		log.Fatal("error loading config: ", err)
	}
	kCtx.Bind(c)

	// jira initialisation
	jClient := jira.NewClient(c)
	kCtx.Bind(jClient)

	// run
	kCtx.FatalIfErrorf(kCtx.Run())
}
