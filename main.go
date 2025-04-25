

package main

import (
    "log"
    "os"
    "github.com/astropatty/gbtest/config"
    "github.com/astropatty/gbtest/cmds"
    "github.com/urfave/cli/v2"
    "github.com/spf13/viper"
)

func main() {
    err  := config.InitializeConfig()
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
