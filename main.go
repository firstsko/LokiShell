package main

import (
	"log"
	"lokishell/shell"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	username := ""
	app_id := ""
	app := &cli.App{
		Action: func(c *cli.Context) error {
			username = c.Args().Get(0)
			app_id = c.Args().Get(1)
			//fmt.Printf("Hello %q %q", user,app_id)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	shell.Run(username, app_id)

}
