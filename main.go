package main

import (
	"github.com/astropatty/gbtest/cmds"
	"github.com/astropatty/gbtest/config"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	err := config.InitializeConfig()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "greet",
		Usage: "fight the loneliness!",
		Action: func(cCtx *cli.Context) error {
			cmds.Deploy()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	viper.WriteConfig()
}
