package config

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/magiconair/properties"
	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
	"io"
	"github.com/prometheus/common/log"
	"strings"
)

// A DecoderconfigOption can be passed to config.Unmarshal to configure
// mapstructure.Decoderconfig options
type DecoderconfigOption func(*mapstructure.DecoderConfig)

// defaultDecoderconfig returns default mapsstructure.Decoderconfig with suppot
// of time.Duration values & string slices
func defaultDecoderconfig(output interface{}, opts ...DecoderconfigOption) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// UnmarshalExact unmarshals the config into a Struct, erroring if a field is nonexistent
// in the destination struct.
func (v *config) UnmarshalExact(rawVal interface{}) error {
	config := defaultDecoderconfig(rawVal)
	config.ErrorUnused = true

	err := decode(v.AllSettings(), config)

	if err != nil {
		return err
	}

	v.insensitiviseMaps()

	return nil
}

func (v *config) UnmarshalKey(key string, rawVal interface{}, opts ...DecoderconfigOption) error {
	err := decode(v.Get(key), defaultDecoderconfig(rawVal, opts...))

	if err != nil {
		return err
	}

	v.insensitiviseMaps()

	return nil
}

func (v *config) Unmarshal(rawVal interface{}, opts ...DecoderconfigOption) error {
	err := decode(v.AllSettings(), defaultDecoderconfig(rawVal, opts...))

	if err != nil {
		return err
	}

	v.insensitiviseMaps()

	return nil
}

func (v *config) unmarshalReader(in io.Reader, c map[string]interface{}) error {
	log.Debug("event", "unmarshaling reader...")

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(in)

	if err != nil {
		return err
	}

	switch strings.ToLower(v.getconfigType()) {
	case "yaml", "yml":
		if err := yaml.Unmarshal(buf.Bytes(), &c); err != nil {
			return configParseError{err}
		}

	case "json":
		if err := json.Unmarshal(buf.Bytes(), &c); err != nil {
			return configParseError{err}
		}

	case "hcl":
		obj, err := hcl.Parse(string(buf.Bytes()))
		if err != nil {
			return configParseError{err}
		}
		if err = hcl.DecodeObject(&c, obj); err != nil {
			return configParseError{err}
		}

	case "toml":
		tree, err := toml.LoadReader(buf)
		if err != nil {
			return configParseError{err}
		}
		tmap := tree.ToMap()
		for k, v := range tmap {
			c[k] = v
		}

	case "properties", "props", "prop":
		v.properties = properties.NewProperties()
		var err error
		if v.properties, err = properties.Load(buf.Bytes(), properties.UTF8); err != nil {
			return configParseError{err}
		}
		for _, key := range v.properties.Keys() {
			value, _ := v.properties.Get(key)
			// recursively build nested maps
			path := strings.Split(key, ".")
			lastKey := strings.ToLower(path[len(path)-1])
			deepestMap := deepSearch(c, path[0:len(path)-1])
			// set innermost value
			deepestMap[lastKey] = value
		}
	}

	insensitiviseMap(c)
	return nil
}

func (v *config) marshalWriter(f afero.File, configType string) error {
	log.Debug("event", "marshaling writer...")

	c := v.AllSettings()
	switch configType {
	case "json":
		b, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			return configMarshalError{err}
		}
		_, err = f.WriteString(string(b))
		if err != nil {
			return configMarshalError{err}
		}

	case "hcl":
		b, err := json.Marshal(c)
		ast, err := hcl.Parse(string(b))
		if err != nil {
			return configMarshalError{err}
		}
		err = printer.Fprint(f, ast.Node)
		if err != nil {
			return configMarshalError{err}
		}

	case "prop", "props", "properties":
		if v.properties == nil {
			v.properties = properties.NewProperties()
		}
		p := v.properties
		for _, key := range v.AllKeys() {
			_, _, err := p.Set(key, v.GetString(key))
			if err != nil {
				return configMarshalError{err}
			}
		}
		_, err := p.WriteComment(f, "#", properties.UTF8)
		if err != nil {
			return configMarshalError{err}
		}

	case "toml":
		t, err := toml.TreeFromMap(c)
		if err != nil {
			return configMarshalError{err}
		}
		s := t.String()
		if _, err := f.WriteString(s); err != nil {
			return configMarshalError{err}
		}

	case "yaml", "yml":
		b, err := yaml.Marshal(c)
		if err != nil {
			return configMarshalError{err}
		}
		if _, err = f.WriteString(string(b)); err != nil {
			return configMarshalError{err}
		}
	}
	return nil
}
