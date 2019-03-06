package main

import (
	"fmt"
	"os"
	"log"
	"io/ioutil"
	"github.com/urfave/cli"
)

type Semver struct {
	major uint16
	minor uint16
	patch uint16
}

func (s Semver) readCurrentVersion() {
	version, err := ioutil.ReadFile(".semver")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File contents: %s", version)
}

func (s Semver) updateVersion(version string) {
	fmt.Println("New version: %s", version)
}

func (s Semver) update(major bool, minor bool, patch bool) {
	if major {
		fmt.Println("Major")
		s.major += 1
		s.minor = 0
		s.patch = 0
	}
	if minor {
		fmt.Println("minor")
		s.minor += 1
		s.patch = 0
	}
	if patch {
		fmt.Println("patch")
		s.patch += 1
	}
}

func main() {
	semver	:= Semver{}
	app		:= cli.NewApp()
	app.Version = "0.0.0"

	app.Flags = []cli.Flag {
		cli.BoolFlag{Name: "compgen", Hidden: true},
		cli.BoolFlag {
			Name: "Major, M",
			Usage: "...",
		},
		cli.BoolFlag {
			Name: "minor, m",
			Usage: "...",
		},
		cli.BoolFlag {
			Name: "patch, p",
			Usage: "...",
		},
	}

	app.Action = func(c *cli.Context) error {
		semver.readCurrentVersion()
		if c.NArg() > 0 {
			newVersion := c.Args().Get(0)
			semver.updateVersion(newVersion)
		}
		semver.update(c.Bool("Major"), c.Bool("minor"), c.Bool("patch"))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}