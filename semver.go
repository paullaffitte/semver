package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	"strconv"
	"io/ioutil"
	"github.com/urfave/cli"
	"github.com/hashicorp/go-version"
)

type SemverUpdater struct {
	version *version.Version
}

func (s *SemverUpdater) ReadCurrentVersion() {
	version, err := ioutil.ReadFile(".semver")
	if err != nil {
		log.Fatal(err)
	}

	s.UpdateVersion(string(version))
}

func (s *SemverUpdater) UpdateVersion(versionStr string) {
	version, err := version.NewSemver(versionStr)
	if err != nil {
		log.Fatal(err)
	}
	s.version = version
}

func (s *SemverUpdater) UpdateSegments(incrs []bool) {
	reset := false
	segments := s.version.Segments()
	segmentsStr := []string{}
	for i := 0; i < len(segments); i++ {
		if reset {
			segments[i] = 0
		}
		if i < len(incrs) && incrs[i] {
			segments[i]++
			reset = true
		}
		segmentsStr = append(segmentsStr, strconv.Itoa(segments[i]))
	}
	s.UpdateVersion(strings.Join(segmentsStr, "."))
}


func (s *SemverUpdater) Update(major bool, minor bool, patch bool) {
	s.UpdateSegments([]bool{major, minor, patch})
}

func main() {
	semver	:= SemverUpdater{}
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
		semver.ReadCurrentVersion()

		if c.NArg() > 0 {
			newVersion := c.Args().Get(0)
			semver.UpdateVersion(newVersion)
		} else {
			semver.Update(c.Bool("Major"), c.Bool("minor"), c.Bool("patch"))
		}

		fmt.Println(semver.version)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}