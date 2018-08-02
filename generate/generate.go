package main

import "fmt"

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
	Bases: []string{"default", "alpine"},
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

func main() {
	fmt.Printf("%+v\n", config)
}
