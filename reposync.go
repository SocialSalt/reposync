package reposync

import (
	cli "github.com/urfave/cli/v2"
	"log"
	"os"
)

func App() {
	log.Println("starting")
	app := &cli.App{
		Name:  "reposync",
		Usage: "Sync current git repo with remote host",
		Action: func(*cli.Context) error {
			log.Println("did something")
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
