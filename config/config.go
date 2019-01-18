package config

import (
	"bytes"
	"encoding/csv"
	"github.com/Masterminds/sprig"
	"github.com/fsnotify/fsnotify"
	"github.com/huandu/xstrings"
	"github.com/magiconair/properties"
	"github.com/spf13/afero"
	"io"
	"os"
	"strings"
)

// SupportedExts are universally supported extensions.
var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl"}

// SupportedRemoteProviders are universally supported remote providers.
var SupportedRemoteProviders = []string{"etcd", "consul"}

var c *config

type config struct {
	// Name of file to look for inside the path
	Name string
	File string
	Path string

	fs afero.Fs

	// A set of remote providers to search for the configuration
	remoteProviders []*defaultRemoteProvider

	configType     string
	envPrefix      string
	envKeyReplacer *strings.Replacer
	config         map[string]interface{}
	override       map[string]interface{}
	defaults       map[string]interface{}
	kvstore        map[string]interface{}
	pflags         map[string]FlagValue
	env            map[string]string
	aliases        map[string]string
	typeByDefValue bool
	keyDelim       string
	// Store read properties on the object so that we can write back in order with comments.
	// This will only be used if the configuration read is a properties file.
	properties     *properties.Properties
	onconfigChange func(fsnotify.Event)
}

// New returns an initialized config instance.
func New(name, file, pth string) *config {
	c.keyDelim = "."
	c.Name = name
	c.File = file
	c.Path = pth
	c.fs = afero.NewOsFs()
	c.config = make(map[string]interface{})
	c.override = make(map[string]interface{})
	c.defaults = make(map[string]interface{})
	c.kvstore = make(map[string]interface{})
	c.pflags = make(map[string]FlagValue)
	c.aliases = make(map[string]string)
	c.typeByDefValue = true
	c.env = make(map[string]string)
	for _, val := range os.Environ() {
		slice := strings.Split(val, "=")
		c.env[strings.ToLower(slice[0])] = strings.ToLower(slice[1])
	}

	for _, val := range c.AllKeys() {
		c.regAlias(strings.ToUpper(val), val)
		c.regAlias(strings.ToTitle(val), val)
		c.regAlias(xstrings.ToCamelCase(val), val)
		c.regAlias(xstrings.ToSnakeCase(val), val)
		c.regAlias(xstrings.ToKebabCase(val), val)
	}
	return c
}

func (v *config) setEnvPrefix(in string) {
	if in != "" {
		v.envPrefix = in
	}
}

func (v *config) getEnv(key string) (string, bool) {
	if v.envKeyReplacer != nil {
		key = v.envKeyReplacer.Replace(key)
	}

	return os.LookupEnv(key)
}

func (v *config) configFileUsed() string { return v.File }

func (v *config) SetTypeByDynameicValue(enable bool) {
	v.typeByDefValue = false
}

func (v *config) GenerateFuncMap() map[string]interface{} {
	m := v.AllSettings()
	for k, val := range sprig.GenericFuncMap() {
		m[k] = val
	}
	return m
}

func (v *config) BindEnv(e map[string]string) {
	for k, val := range e {
		v.env[strings.ToLower(k)] = val
	}
}

func (v *config) ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func (v *config) IsSet(key string) bool {
	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey)
	return val != nil
}

func (v *config) setEnvKeyReplacer(r *strings.Replacer) {
	v.envKeyReplacer = r
}

func (v *config) regAlias(alias string, key string) {
	v.registerAlias(alias, strings.ToLower(key))
}

func (v *config) Inconfig(key string) bool {
	// if the requested key is an alias, then return the proper key
	key = v.realKey(key)

	_, exists := v.config[key]
	return exists
}

func (v *config) SetDefault(key string, value interface{}) {
	switch x := value.(type) {
	case string:
		_ = os.Setenv(key, x)
	case int:
		_ = os.Setenv(key, string(x))
	}

	// If alias passed in, then set the proper default
	key = v.realKey(strings.ToLower(key))
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(v.defaults, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value
}

func (v *config) Set(key string, value interface{}) {
	// If alias passed in, then set the proper override
	key = v.realKey(strings.ToLower(key))
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(v.override, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value
}

func (v *config) ReadInconfig() error {
	green("%s\n", "Attempting to read in config file...")
	filename, err := v.getconfigFile()
	if err != nil {
		return err
	}

	if !stringInSlice(v.getconfigType(), SupportedExts) {
		return UnsupportedconfigError(v.getconfigType())
	}

	green("%s\n", "Reading file:"+filename+"...")
	file, err := afero.ReadFile(v.fs, filename)
	if err != nil {
		return err
	}

	config := make(map[string]interface{})

	err = v.unmarshalReader(bytes.NewReader(file), config)
	if err != nil {
		return err
	}

	v.config = config
	return nil
}

func (v *config) MergeInconfig() error {
	green("%s\n", "Attempting to merge in config file...")
	filename, err := v.getconfigFile()
	if err != nil {
		return err
	}

	if !stringInSlice(v.getconfigType(), SupportedExts) {
		return UnsupportedconfigError(v.getconfigType())
	}

	file, err := afero.ReadFile(v.fs, filename)
	if err != nil {
		return err
	}

	return v.Mergeconfig(bytes.NewReader(file))
}

func (v *config) Readconfig(in io.Reader) error {
	v.config = make(map[string]interface{})
	return v.unmarshalReader(in, v.config)
}

func (v *config) Mergeconfig(in io.Reader) error {
	cfg := make(map[string]interface{})
	if err := v.unmarshalReader(in, cfg); err != nil {
		return err
	}
	return v.MergeconfigMap(cfg)
}

func (v *config) MergeconfigMap(cfg map[string]interface{}) error {
	if v.config == nil {
		v.config = make(map[string]interface{})
	}
	insensitiviseMap(cfg)
	mergeMaps(cfg, v.config, nil)
	return nil
}

func (v *config) Writeconfig() error {
	filename, err := v.getconfigFile()
	if err != nil {
		return err
	}
	return v.writeconfig(filename, true)
}

func (v *config) SafeWriteconfig() error {
	filename, err := v.getconfigFile()
	if err != nil {
		return err
	}
	return v.writeconfig(filename, false)
}

func (v *config) WriteconfigAs(filename string) error {
	return v.writeconfig(filename, true)
}

func (v *config) SafeWriteconfigAs(filename string) error {
	return v.writeconfig(filename, false)
}

func (v *config) AllKeys() []string {
	m := map[string]bool{}
	// add all paths, by order of descending priority to ensure correct shadowing
	m = v.flattenAndMergeMap(m, castMapStringToMapInterface(v.aliases), "")
	m = v.flattenAndMergeMap(m, v.override, "")
	m = v.mergeFlatMap(m, castMapFlagToMapInterface(v.pflags))
	m = v.mergeFlatMap(m, castMapStringToMapInterface(v.env))
	m = v.flattenAndMergeMap(m, v.config, "")
	m = v.flattenAndMergeMap(m, v.kvstore, "")
	m = v.flattenAndMergeMap(m, v.defaults, "")

	// convert set of paths to list
	a := []string{}
	for x := range m {
		a = append(a, x)
	}
	return a
}

func (v *config) SetFs(fs afero.Fs) {
	v.fs = fs
}

func (v *config) SetconfigType(in string) {
	if in != "" {
		v.configType = in
	}
}
