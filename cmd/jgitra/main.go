package main

import (
	"log"

	"jgitra/internal/cli"
	"jgitra/internal/config"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Validate cli.ValidateCmd `cmd:"" help:"Validate configuration."`
}

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatal("error loading config: ", err)
	}
	ctx := kong.Parse(&CLI{},
		kong.Name(config.App),
		kong.Description("jira + git tool"),
		kong.UsageOnError(),
	)
	ctx.FatalIfErrorf(ctx.Run(c))
}
