package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// Config contains information about the different versions
type Config struct {
	Bases    []string
	Versions []Version
}

func (c Config) Builds() (builds []Build) {
	for _, v := range c.Versions {
		builds = append(builds, v.Builds()...)
	}
	return
}

// Version contains information about a version
type Version struct {
	Version     string
	CompilerURL string
	JDKVersions []JDKVersion
}

// VersionSnakeCased returns v.Version snake-cased
func (v Version) VersionSnakeCased() string {
	return string(regexp.MustCompile("\\W").ReplaceAll(([]byte)(v.Version), ([]byte)("_")))
}

func (v Version) Builds() (builds []Build) {
	for _, jdkVersion := range v.JDKVersions {
		builds = append(builds, Build{v, jdkVersion, jdkVersion.Base})
		for _, variant := range jdkVersion.Variants {
			builds = append(builds, Build{v, jdkVersion, variant})
		}
	}
	return
}

// JDKVersion contains information about a JDK version
type JDKVersion struct {
	JDKVersion string
	Base
	Variants []Base
}

// Base contains information about a base
type Base struct {
	Base                   string
	AdditionalTags         []string
	AdditionalRepositories []AdditionalRepository
}

// AdditionalRepository contains information about an additional repository
type AdditionalRepository struct {
	Repository string
	Tags       []string
}

// Build is a specific build of the image
type Build struct {
	Version
	JDKVersion
	Base
}

// Name is build's name in CI
func (b Build) Name() (n string) {
	n = "build_" + b.Version.VersionSnakeCased() + "_jdk" + b.JDKVersion.JDKVersion
	if b.JDKVersion.Base.Base != b.Base.Base {
		n += "_" + b.Base.Base
	}
	return
}

// Tag is build's main tag
func (b Build) Tag() (t string) {
	t = b.Version.Version + "-jdk" + b.JDKVersion.JDKVersion
	if b.JDKVersion.Base.Base != b.Base.Base {
		t += "-" + b.Base.Base
	}
	return
}

func (b Build) Source() (s string) {
	s = "openjdk:" + b.JDKVersion.JDKVersion + "-jdk"
	if b.JDKVersion.Base.Base != b.Base.Base {
		s += "-" + b.Base.Base
	}
	return
}

// FullTag is build's main tag with image name
func (b Build) FullTag() string {
	return "zenika/kotlin:" + b.Tag()
}

// AdditionalTags is build's additional tags
func (b Build) AdditionalTags() (tags []string) {
	for _, t := range b.Base.AdditionalTags {
		tags = append(tags, "zenika/kotlin:"+t)
	}
	for _, r := range b.Base.AdditionalRepositories {
		for _, t := range r.Tags {
			tags = append(tags, r.Repository+":"+t)
		}
	}
	return
}

var (
	config           Config
	readmeTemplate   *template.Template
	circleciTemplate *template.Template
	templates        = make(map[string][]*template.Template)
	templatesDir     string
	wd               string
)

func main() {
	if err := initTemplatesDir(); err != nil {
		panic(err)
	}
	if err := loadConfig(); err != nil {
		panic(err)
	}
	if err := loadAllTemplates(); err != nil {
		panic(err)
	}
	if err := generateAll(); err != nil {
		panic(err)
	}
}

func initTemplatesDir() error {
	var err error
	if wd, err = os.Getwd(); err != nil {
		return err
	}
	templatesDir = filepath.Join(wd, "templates")
	return nil
}

func loadConfig() error {
	viper.AddConfigPath(filepath.Join(wd))
	viper.SetConfigName("versions")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&config)
}

func generateAll() error {
	if err := generateReadme(); err != nil {
		return err
	}
	if err := generateCircleci(); err != nil {
		return err
	}
	return nil
}

func generateReadme() error {
	return generateTemplate(readmeTemplate, config, wd)
}

func generateCircleci() error {
	return generateTemplate(circleciTemplate, config, wd)
}

func generateTemplate(t *template.Template, ctxt interface{}, outDir string) error {
	fName := filepath.Join(outDir, t.Name())

	if err := ensureDir(filepath.Dir(fName)); err != nil {
		return err
	}

	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, ctxt); err != nil {
		return err
	}

	return nil
}

func loadAllTemplates() error {
	if err := loadReadmeTemplate(); err != nil {
		return err
	}
	if err := loadCircleciTemplate(); err != nil {
		return err
	}
	return nil
}

func loadReadmeTemplate() error {
	var err error
	readmeTemplate, err = readTemplateFile("README.md", filepath.Join(templatesDir, "README.md"))
	return err
}

func loadCircleciTemplate() error {
	var err error
	circleciTemplate, err = readTemplateFile(".circleci/config.yml", filepath.Join(templatesDir, ".circleci/config.yml"))
	return err
}

func readTemplateFile(name, path string) (*template.Template, error) {
	s, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t := template.New(name)
	t = t.Delims("#{", "}")
	t = t.Funcs(template.FuncMap{
		"join": strings.Join,
	})
	t = template.Must(t.Parse(string(s)))

	return t, nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
