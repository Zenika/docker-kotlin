package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
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

// Context contains information for the templates
type Context struct {
	Wd                     string
	Version                string
	CompilerURL            string
	JDKVersion             string
	AdditionalTags         []string
	AdditionalRepositories []AdditionalRepository
}

var (
	config         Config
	readmeTemplate *template.Template
	templates      = make(map[string][]*template.Template)
	templatesDir   string
	wd             string
)

func contextWithVersion(version Version) Context {
	return Context{Wd: filepath.Join(wd, version.Version), Version: version.Version, CompilerURL: version.CompilerURL}
}

func (ctxt Context) withJDKVersion(jdkVersion JDKVersion) Context {
	ctxt.Wd = filepath.Join(ctxt.Wd, "jdk"+jdkVersion.JDKVersion)
	ctxt.JDKVersion = jdkVersion.JDKVersion
	return ctxt
}

func (ctxt Context) withBase(base Base) Context {
	if base.Base != "default" {
		ctxt.Wd = filepath.Join(ctxt.Wd, base.Base)
	}
	ctxt.AdditionalTags = base.AdditionalTags
	ctxt.AdditionalRepositories = base.AdditionalRepositories
	return ctxt
}

func init() {
	var err error
	if wd, err = os.Getwd(); err != nil {
		panic(err)
	}
	templatesDir = filepath.Join(wd, "templates")
}

func init() {
	viper.AddConfigPath(filepath.Join(wd))
	viper.SetConfigName("versions")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.Unmarshal(&config)
}

func init() {
	if err := loadTemplates(); err != nil {
		panic(err)
	}
}

func init() {
	if err := loadReadmeTemplate(); err != nil {
		panic(err)
	}
}

func main() {
	for _, version := range config.Versions {
		if err := generateVersion(version); err != nil {
			panic(err)
		}
	}
	if err := generateReadme(); err != nil {
		panic(err)
	}
}

func generateVersion(version Version) error {
	ctxt := contextWithVersion(version)

	if err := ensureDir(ctxt.Wd); err != nil {
		return err
	}

	for _, jdkVersion := range version.JDKVersions {
		if err := generateJDKVersion(ctxt, jdkVersion); err != nil {
			return err
		}
	}

	return nil
}

func generateJDKVersion(ctxt Context, jdkVersion JDKVersion) error {
	ctxt = ctxt.withJDKVersion(jdkVersion)

	if err := ensureDir(ctxt.Wd); err != nil {
		return err
	}

	if err := generateBase(ctxt, jdkVersion.Base); err != nil {
		return err
	}

	for _, variant := range jdkVersion.Variants {
		if err := generateBase(ctxt, variant); err != nil {
			return err
		}
	}

	return nil
}

func generateBase(ctxt Context, base Base) error {
	ctxt = ctxt.withBase(base)

	if err := ensureDir(ctxt.Wd); err != nil {
		return err
	}

	if err := generateTemplates("common", ctxt); err != nil {
		return err
	}
	if err := generateTemplates(base.Base, ctxt); err != nil {
		return err
	}

	return nil
}

func generateTemplates(name string, ctxt Context) error {
	for _, t := range templates[name] {
		if err := generateTemplate(t, ctxt, ctxt.Wd); err != nil {
			return err
		}
	}

	return nil
}

func generateReadme() error {
	return generateTemplate(readmeTemplate, config, wd)
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

func loadTemplates() error {
	for _, base := range config.Bases {
		if err := loadTemplate(base); err != nil {
			return err
		}
	}

	return nil
}

func loadTemplate(base string) error {
	baseDir := filepath.Join(templatesDir, base)

	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		t, err := readTemplateFile(relPath, path)
		if err != nil {
			return err
		}

		templates[base] = append(templates[base], t)

		return nil
	})
}

func loadReadmeTemplate() error {
	var err error
	readmeTemplate, err = readTemplateFile("README.md", filepath.Join(templatesDir, "README.md"))
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
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}

	return err
}
