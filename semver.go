package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	"strconv"
	"regexp"
	"io/ioutil"
	"github.com/urfave/cli"
	"github.com/hashicorp/go-version"
	"github.com/go-yaml/yaml"
)

type Config struct {
    Version string
    Files map[string]string
}

type SemverUpdater struct {
	version *version.Version
	config Config
}

func (s *SemverUpdater) Init(configPath string) {
	if _, err := os.Stat(configPath); err == nil {
		log.Fatal("file " + configPath + " already exists, please remove it before retrying")
	}
	config := Config { Version: "0.0.0" }

	header := "# semver config file -- https://github.com/paullaffitte/semver\n\n"
	configStr, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(configPath, []byte(header + string(configStr)), 0644)
}

func (s *SemverUpdater) ReadConfig(configPath string) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatal(err)
	}

	s.UpdateVersion(config.Version)
	s.config = config
	s.config.Files["semver.yml"] = "version:\\s*.*"
}

func (s *SemverUpdater) UpdateVersion(versionStr string) {
	version, err := version.NewSemver(versionStr)
	if err != nil {
		log.Fatal(err)
	}
	s.version = version
}

func (s *SemverUpdater) UpdatedSegments(incrs []bool) string {
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

	return strings.Join(segmentsStr, ".")
}

func (s *SemverUpdater) Update(major bool, minor bool, patch bool, tag string, metadata string) {
	updated := major || minor || patch

	if len(tag) == 0 && !updated {
		tag = s.version.Prerelease()
	}
	if len(tag) > 0 {
		tag = "-" + tag
	}

	if len(metadata) == 0 && !updated {
		metadata = s.version.Metadata()
	}
	if len(metadata) > 0 {
		metadata = "+" + metadata
	}

	s.UpdateVersion(s.UpdatedSegments([]bool{major, minor, patch}) + tag + metadata)
}

func (s *SemverUpdater) ReplaceVersions(content string, regex string) string {
	re := regexp.MustCompile(regex)
	idx := re.FindStringIndex(content)
	from := idx[0]
	to := idx[1]
	substr := content[from:to]

	newSubstr := strings.Replace(substr, s.config.Version, s.version.String(), -1)
	content = strings.Replace(content, substr, newSubstr, -1)
	return content
}

func (s *SemverUpdater) SyncFiles() {
	for file, regex := range s.config.Files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		newContent := s.ReplaceVersions(string(content), regex)
		ioutil.WriteFile(file, []byte(newContent), 0644)
	}
}

func main() {
	var semver SemverUpdater
	app := cli.NewApp()
	app.Version = "0.3.0"

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
		cli.StringFlag {
			Name: "tag, t",
			Usage: "...",
		},
		cli.StringFlag {
			Name: "metadata, d",
			Usage: "...",
		},
		cli.StringFlag {
			Name: "config, c",
			Value: "./semver.yml",
			Usage: "...",
		},
		cli.BoolFlag {
			Name: "init, i",
			Usage: "...",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("init") {
			semver.Init(c.String("config"))
		}

		semver.ReadConfig(c.String("config"))

		if c.NArg() > 0 {
			newVersion := c.Args().Get(0)
			semver.UpdateVersion(newVersion)
		} else {
			semver.Update(c.Bool("Major"), c.Bool("minor"), c.Bool("patch"), c.String("tag"), c.String("metadata"))
		}

		fmt.Println(semver.version)
		semver.SyncFiles()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
