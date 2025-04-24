

package main

import (
    "log"
    "os"
    "github.com/astropatty/gh-test/cmds"

    "github.com/urfave/cli/v2"
)

func main() {
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
}
