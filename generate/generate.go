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
	Base               string
	AdditionalVersions []string
}

// Context contains information for the templates
type Context struct {
	Wd                 string
	Version            string
	CompilerURL        string
	JDKVersion         string
	AdditionalVersions []string
}

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
	ctxt.AdditionalVersions = base.AdditionalVersions
	return ctxt
}

var (
	config    Config
	templates = make(map[string][]*template.Template)
	wd        string
)

func init() {
	var err error
	if wd, err = os.Getwd(); err != nil {
		panic(err)
	}
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

func main() {
	for _, version := range config.Versions {
		if err := generateVersion(version); err != nil {
			panic(err)
		}
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

	if err := generateTemplate(ctxt, "common"); err != nil {
		return err
	}
	if err := generateTemplate(ctxt, base.Base); err != nil {
		return err
	}

	return nil
}

func generateTemplate(ctxt Context, name string) error {
	for _, template := range templates[name] {
		fName := filepath.Join(ctxt.Wd, template.Name())

		if err := ensureDir(filepath.Dir(fName)); err != nil {
			return err
		}

		f, err := os.Create(fName)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := template.Execute(f, ctxt); err != nil {
			return err
		}
	}

	return nil
}

func loadTemplates() error {
	templatesDir := filepath.Join(wd, "templates")

	for _, base := range config.Bases {
		if err := loadTemplate(templatesDir, base); err != nil {
			return err
		}
	}

	return nil
}

func loadTemplate(templatesDir, base string) error {
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

		s, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		t := template.New(relPath)
		t = t.Delims("#{", "}")
		t = t.Funcs(template.FuncMap{
			"join": func(sep string, args ...string) string { return strings.Join(args, sep) },
		})
		t = template.Must(t.Parse(string(s)))

		templates[base] = append(templates[base], t)

		return nil
	})
}

func ensureDir(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}

	return err
}
