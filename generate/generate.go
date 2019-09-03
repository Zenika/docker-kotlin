package main

import (
	"io/ioutil"
	"os"
	"path"
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

// FullTag is build's main tag with image name
func (b Build) FullTag() string {
	return "zenika/kotlin:" + b.Tag()
}

// BuildContext is build's build context
func (b Build) BuildContext() (bc string) {
	bc = path.Join(b.Version.Version, "jdk"+b.JDKVersion.JDKVersion)
	if b.JDKVersion.Base.Base != b.Base.Base {
		bc = path.Join(bc, b.Base.Base)
	}
	return
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

func contextWithVersion(version Version) Context {
	return Context{Wd: filepath.Join(wd, version.Version), Version: version.Version, CompilerURL: version.CompilerURL}
}

func (ctxt Context) withJDKVersion(jdkVersion JDKVersion) Context {
	ctxt.Wd = filepath.Join(ctxt.Wd, "jdk"+jdkVersion.JDKVersion)
	ctxt.JDKVersion = jdkVersion.JDKVersion
	return ctxt
}

func (ctxt Context) withBase(base Base, isDefault bool) Context {
	if !isDefault {
		ctxt.Wd = filepath.Join(ctxt.Wd, base.Base)
	}
	ctxt.AdditionalTags = base.AdditionalTags
	ctxt.AdditionalRepositories = base.AdditionalRepositories
	return ctxt
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
	for _, version := range config.Versions {
		if err := generateVersion(version); err != nil {
			return err
		}
	}
	if err := generateReadme(); err != nil {
		return err
	}
	if err := generateCircleci(); err != nil {
		return err
	}
	return nil
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

	if err := generateBase(ctxt, jdkVersion.Base, true); err != nil {
		return err
	}

	for _, variant := range jdkVersion.Variants {
		if err := generateBase(ctxt, variant, false); err != nil {
			return err
		}
	}

	return nil
}

func generateBase(ctxt Context, base Base, isDefault bool) error {
	ctxt = ctxt.withBase(base, isDefault)

	if err := ensureDir(ctxt.Wd); err != nil {
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
	if err := loadBasesTemplates(); err != nil {
		return err
	}
	if err := loadReadmeTemplate(); err != nil {
		return err
	}
	if err := loadCircleciTemplate(); err != nil {
		return err
	}
	return nil
}

func loadBasesTemplates() error {
	for _, base := range config.Bases {
		if err := loadBaseTemplate(base); err != nil {
			return err
		}
	}

	return nil
}

func loadBaseTemplate(base string) error {
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
