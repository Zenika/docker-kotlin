package main

import (
	"html/template"
	"os"
	"path/filepath"
)

// Config contains information about the different versions
type Config struct {
	Bases    []string
	Versions []Version
}

// Version contains information about a version
type Version struct {
	Version     string
	CompilerURL string
	JDKVersions []JDKVersion
}

// JDKVersion contains information about a JDK version
type JDKVersion struct {
	Base
	JDKVersion string
	Variants   []Base
}

// Base contains information about a base
type Base struct {
	Base               string
	AdditionalVersions []string
}

var config = Config{
	Bases: []string{"common", "default", "alpine"},
	Versions: []Version{
		{
			Version: "1.1",
			JDKVersions: []JDKVersion{
				{
					JDKVersion: "8",
					Base: Base{
						Base:               "default",
						AdditionalVersions: []string{"1.1.61-jdk8"},
					},
					Variants: []Base{
						{
							Base:               "alpine",
							AdditionalVersions: []string{"1.1.61-jdk8-alpine"},
						},
					},
				},
			},
		},
	},
}

var (
	templates = template.New("root")
)

func main() {
	if err := loadTemplates(); err != nil {
		panic(err)
	}
}

func loadTemplates() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	templatesDir := filepath.Join(wd, "templates")

	for _, base := range config.Bases {
		if err := loadTemplate(templatesDir, base); err != nil {
			return err
		}
	}

	return nil
}

func loadTemplate(templatesDir, base string) error {
	t := templates.New(base)
	baseDir := filepath.Join(templatesDir, base)

	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		_, err = t.New(path).Parse(path)

		return err
	})
}
