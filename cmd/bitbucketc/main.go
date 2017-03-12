package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()
	app.Name = "bitbucketc"
	app.Usage = "bitbucketc provides command line tools for bitbucket pipelines"
	app.Commands = []cli.Command{
		compileCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
