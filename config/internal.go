package config

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/prometheus/common/log"
	"os"
	"path/filepath"
	"github.com/fatih/color"
	"reflect"
	"runtime"
	"strings"
	"unicode"
)

func (v *config) isPathShadowedInDeepMap(path []string, m map[string]interface{}) string {
	var parentVal interface{}
	for i := 1; i < len(path); i++ {
		parentVal = v.searchMap(m, path[0:i])
		if parentVal == nil {
			// not found, no need to add more path elements
			return ""
		}
		switch parentVal.(type) {
		case map[interface{}]interface{}:
			continue
		case map[string]interface{}:
			continue
		default:
			// parentVal is a regular value which shadows "path"
			return strings.Join(path[0:i], v.keyDelim)
		}
	}
	return ""
}

func (v *config) isPathShadowedInFlatMap(path []string, mi interface{}) string {
	// unify input map
	var m map[string]interface{}
	switch mi.(type) {
	case map[string]string, map[string]FlagValue:
		m = cast.ToStringMap(mi)
	default:
		return ""
	}

	// scan paths
	var parentKey string
	for i := 1; i < len(path); i++ {
		parentKey = strings.Join(path[0:i], v.keyDelim)
		if _, ok := m[parentKey]; ok {
			return parentKey
		}
	}
	return ""
}

func (v *config) isPathShadowedInAutoEnv(path []string) string {
	var parentKey string
	for i := 1; i < len(path); i++ {
		parentKey = strings.Join(path[0:i], v.keyDelim)
		if _, ok := v.getEnv(v.mergeWithEnvPrefix(parentKey)); ok {
			return parentKey
		}
	}
	return ""
}

func (v *config) searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// Fast path
		if len(path) == 1 {
			return next
		}

		// Nested case
		switch next.(type) {
		case map[interface{}]interface{}:
			return v.searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// Type assertion is safe here since it is only reached
			// if the type of `next` is the same as the type being asserted
			return v.searchMap(next.(map[string]interface{}), path[1:])
		default:
			// got a value but nested key expected, return "nil" for not found
			return nil
		}
	}
	return nil
}

func (v *config) searchMapWithPathPrefixes(source map[string]interface{}, path []string) interface{} {
	log.Debug("event", "searching path with prefix...")

	if len(path) == 0 {
		return source
	}

	// search for path prefixes, starting from the longest one
	for i := len(path); i > 0; i-- {
		prefixKey := strings.ToLower(strings.Join(path[0:i], v.keyDelim))

		next, ok := source[prefixKey]
		if ok {
			// Fast path
			if i == len(path) {
				return next
			}

			// Nested case
			var val interface{}
			switch next.(type) {
			case map[interface{}]interface{}:
				val = v.searchMapWithPathPrefixes(cast.ToStringMap(next), path[i:])
			case map[string]interface{}:
				// Type assertion is safe here since it is only reached
				// if the type of `next` is the same as the type being asserted
				val = v.searchMapWithPathPrefixes(next.(map[string]interface{}), path[i:])
			default:
				// got a value but nested key expected, do nothing and look for next prefix
			}
			if val != nil {
				return val
			}
		}
	}

	// not found
	return nil
}

func (v *config) mergeWithEnvPrefix(in string) string {
	if v.envPrefix != "" {
		return strings.ToUpper(v.envPrefix + "_" + in)
	}

	return strings.ToUpper(in)
}

func (v *config) find(lcaseKey string) interface{} {

	var (
		val    interface{}
		exists bool
		path   = strings.Split(lcaseKey, v.keyDelim)
		nested = len(path) > 1
	)

	// compute the path through the nested maps to the nested value
	if nested && v.isPathShadowedInDeepMap(path, castMapStringToMapInterface(v.aliases)) != "" {
		return nil
	}

	// if the requested key is an alias, then return the proper key
	lcaseKey = v.realKey(lcaseKey)
	path = strings.Split(lcaseKey, v.keyDelim)
	nested = len(path) > 1

	// Set() override first
	val = v.searchMap(v.override, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.override) != "" {
		return nil
	}

	// PFlag override next
	flag, exists := v.pflags[lcaseKey]
	if exists && flag.HasChanged() {
		switch flag.ValueType() {
		case "int", "int8", "int16", "int32", "int64":
			return cast.ToInt(flag.ValueString())
		case "bool":
			return cast.ToBool(flag.ValueString())
		case "stringSlice":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := v.ReadAsCSV(s)
			return res
		default:
			return flag.ValueString()
		}
	}
	if nested && v.isPathShadowedInFlatMap(path, v.pflags) != "" {
		return nil
	}

	if val, ok := v.getEnv(v.mergeWithEnvPrefix(lcaseKey)); ok {
		return val
	}
	if nested && v.isPathShadowedInAutoEnv(path) != "" {
		return nil
	}
	envkey, exists := v.env[lcaseKey]
	if exists {
		if val, ok := v.getEnv(envkey); ok {
			return val
		}
	}
	if nested && v.isPathShadowedInFlatMap(path, v.env) != "" {
		return nil
	}

	// config file next
	val = v.searchMapWithPathPrefixes(v.config, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.config) != "" {
		return nil
	}

	// K/V store next
	val = v.searchMap(v.kvstore, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.kvstore) != "" {
		return nil
	}

	// Default next
	val = v.searchMap(v.defaults, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.defaults) != "" {
		return nil
	}

	// last chance: if no other value is returned and a flag does exist for the value,
	// get the flag's value even if the flag's value has not changed
	if flag, exists := v.pflags[lcaseKey]; exists {
		switch flag.ValueType() {
		case "int", "int8", "int16", "int32", "int64":
			return cast.ToInt(flag.ValueString())
		case "bool":
			return cast.ToBool(flag.ValueString())
		case "stringSlice":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := v.ReadAsCSV(s)
			return res
		default:
			return flag.ValueString()
		}
	}
	// last item, no need to check shadowing

	return nil
}

func (v *config) registerAlias(alias string, key string) {
	alias = strings.ToLower(alias)
	if alias != key && alias != v.realKey(key) {
		_, exists := v.aliases[alias]

		if !exists {
			// if we alias something that exists in one of the maps to another
			// name, we'll never be able to get that value using the original
			// name, so move the config value to the new realkey.
			if val, ok := v.config[alias]; ok {
				delete(v.config, alias)
				v.config[key] = val
			}
			if val, ok := v.kvstore[alias]; ok {
				delete(v.kvstore, alias)
				v.kvstore[key] = val
			}
			if val, ok := v.defaults[alias]; ok {
				delete(v.defaults, alias)
				v.defaults[key] = val
			}
			if val, ok := v.override[alias]; ok {
				delete(v.override, alias)
				v.override[key] = val
			}
			v.aliases[alias] = key
		}
	} else {
		color.Red("Creating circular reference alias", alias, key, v.realKey(key))
	}
}

func (v *config) realKey(key string) string {
	newkey, exists := v.aliases[key]
	if exists {
		log.Debug("Alias", key, "to", newkey)
		return v.realKey(newkey)
	}
	return key
}

func (v *config) writeconfig(filename string, force bool) error {
	log.Debug("event", "Attempting to write configuration to file.")
	ext := filepath.Ext(filename)
	if len(ext) <= 1 {
		return fmt.Errorf("Filename: %s requires valid extension.", filename)
	}
	configType := ext[1:]
	if !stringInSlice(configType, SupportedExts) {
		return UnsupportedconfigError(configType)
	}
	if v.config == nil {
		v.config = make(map[string]interface{})
	}
	var flags int
	if force == true {
		flags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	} else {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			flags = os.O_WRONLY
		} else {
			return fmt.Errorf(color.RedString("File: %s exists. Use Writeconfig to overwrite."), filename)
		}
	}
	f, err := v.fs.OpenFile(filename, flags, os.FileMode(0644))
	if err != nil {
		return err
	}
	return v.marshalWriter(f, configType)
}

func castMapFlagToMapInterface(src map[string]FlagValue) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}

// mergeMaps merges two maps. The `itgt` parameter is for handling go-yaml's
// insistence on parsing nested structures as `map[interface{}]interface{}`
// instead of using a `string` as the key for nest structures beyond one level
// deep. Both map types are supported as there is a go-yaml fork that uses
// `map[string]interface{}` instead.
func mergeMaps(
	src, tgt map[string]interface{}, itgt map[interface{}]interface{}) {
	for sk, sv := range src {
		tk := keyExists(sk, tgt)
		if tk == "" {
			log.Debug("tk=\"\", tgt[%s]=%v", sk, sv)
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		tv, ok := tgt[tk]
		if !ok {
			log.Debug("tgt[%s] != ok, tgt[%s]=%v", tk, sk, sv)
			tgt[sk] = sv
			if itgt != nil {
				itgt[sk] = sv
			}
			continue
		}

		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)
		if svType != tvType {
			color.Red(
				"svType != tvType; key=%s, st=%v, tt=%v, sv=%v, tv=%v",
				sk, svType, tvType, sv, tv)
			continue
		}

		color.Red("processing key=%s, st=%v, tt=%v, sv=%v, tv=%v",
			sk, svType, tvType, sv, tv)

		switch ttv := tv.(type) {
		case map[interface{}]interface{}:
			log.Debug("merging maps (must convert)")
			tsv := sv.(map[interface{}]interface{})
			ssv := castToMapStringInterface(tsv)
			stv := castToMapStringInterface(ttv)
			mergeMaps(ssv, stv, ttv)
		case map[string]interface{}:
			log.Debug("merging maps")
			mergeMaps(sv.(map[string]interface{}), ttv, nil)
		default:
			log.Debug("setting value")
			tgt[tk] = sv
			if itgt != nil {
				itgt[tk] = sv
			}
		}
	}
}

// flattenAndMergeMap recursively flattens the given map into a map[string]bool
// of key paths (used as a set, easier to manipulate than a []string):
// - each path is merged into a single key string, delimited with v.keyDelim (= ".")
// - if a path is shadowed by an earlier value in the initial shadow map,
//   it is skipped.
// The resulting set of paths is merged to the given shadow set at the same time.
func (v *config) flattenAndMergeMap(shadow map[string]bool, m map[string]interface{}, prefix string) map[string]bool {
	log.Debug("event", "flattening and merging map...")

	if shadow != nil && prefix != "" && shadow[prefix] {
		// prefix is shadowed => nothing more to flatten
		return shadow
	}
	if shadow == nil {
		shadow = make(map[string]bool)
	}

	var m2 map[string]interface{}
	if prefix != "" {
		prefix += v.keyDelim
	}
	for k, val := range m {
		fullKey := prefix + k
		switch val.(type) {
		case map[string]interface{}:
			m2 = val.(map[string]interface{})
		case map[interface{}]interface{}:
			m2 = cast.ToStringMap(val)
		default:
			// immediate value
			shadow[strings.ToLower(fullKey)] = true
			continue
		}
		// recursively merge to shadow map
		shadow = v.flattenAndMergeMap(shadow, m2, fullKey)
	}
	return shadow
}

// mergeFlatMap merges the given maps, excluding values of the second map
// shadowed by values from the first map.
func (v *config) mergeFlatMap(shadow map[string]bool, m map[string]interface{}) map[string]bool {
	log.Debug("event", "merging flat map...")

	// scan keys
outer:
	for k, _ := range m {
		path := strings.Split(k, v.keyDelim)
		// scan intermediate paths
		var parentKey string
		for i := 1; i < len(path); i++ {
			parentKey = strings.Join(path[0:i], v.keyDelim)
			if shadow[parentKey] {
				// path is shadowed, continue
				continue outer
			}
		}
		// add key
		shadow[strings.ToLower(k)] = true
	}
	return shadow
}

func keyExists(k string, m map[string]interface{}) string {
	lk := strings.ToLower(k)
	for mk := range m {
		lmk := strings.ToLower(mk)
		if lmk == lk {
			return mk
		}
	}
	return ""
}

func castToMapStringInterface(
	src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

func castMapStringToMapInterface(src map[string]string) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[k] = v
	}
	return tgt
}

func (v *config) getconfigType() string {
	log.Debug( "event", "getting config type...")

	if v.configType != "" {
		return v.configType
	}

	cf, err := v.getconfigFile()
	if err != nil {
		return ""
	}

	ext := filepath.Ext(cf)

	if len(ext) > 1 {
		return ext[1:]
	}

	return ""
}

func (v *config) getconfigFile() (string, error) {
	log.Debug( "event", "getting config file...")

	if v.File == "" {
		cf, err := v.findconfigFile()
		if err != nil {
			return "", err
		}
		v.File = cf
	}
	return v.File, nil
}

func (v *config) searchInPath(in string) (filename string) {
	log.Debug( "event", "searchingfor config file...")
	for _, ext := range SupportedExts {
		log.Debug("check",(filepath.Join(in, v.Name+"."+ext)))
		if b, _ := exists(v.fs, filepath.Join(in, v.Name+"."+ext)); b {
			log.Debug("check",(filepath.Join(in, v.Name+"."+ext)))
			return filepath.Join(in, v.Name+"."+ext)
		}
	}

	return ""
}

func (v *config) insensitiviseMaps() {
	insensitiviseMap(v.config)
	insensitiviseMap(v.defaults)
	insensitiviseMap(v.override)
	insensitiviseMap(v.kvstore)
}

// configParseError denotes failing to parse configuration file.
type configParseError struct {
	err error
}

// Error returns the formatted configuration error.
func (pe configParseError) Error() string {
	return color.RedString("While parsing config: %s", pe.err.Error())
}

// toCaseInsensitiveValue checks if the value is a  map;
// if so, create a copy and lower-case the keys recursively.
func toCaseInsensitiveValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		value = copyAndInsensitiviseMap(cast.ToStringMap(v))
	case map[string]interface{}:
		value = copyAndInsensitiviseMap(v)
	}

	return value
}

// copyAndInsensitiviseMap behaves like insensitiviseMap, but creates a copy of
// any map it makes case insensitive.
func copyAndInsensitiviseMap(m map[string]interface{}) map[string]interface{} {
	nm := make(map[string]interface{})

	for key, val := range m {
		lkey := strings.ToLower(key)
		switch v := val.(type) {
		case map[interface{}]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(cast.ToStringMap(v))
		case map[string]interface{}:
			nm[lkey] = copyAndInsensitiviseMap(v)
		default:
			nm[lkey] = v
		}
	}

	return nm
}

func insensitiviseMap(m map[string]interface{}) {
	for key, val := range m {
		switch val.(type) {
		case map[interface{}]interface{}:
			// nested map: cast and recursively insensitivise
			val = cast.ToStringMap(val)
			insensitiviseMap(val.(map[string]interface{}))
		case map[string]interface{}:
			// nested map: recursively insensitivise
			insensitiviseMap(val.(map[string]interface{}))
		}

		lower := strings.ToLower(key)
		if key != lower {
			// remove old key (not lower-cased)
			delete(m, key)
		}
		// update map
		m[lower] = val
	}
}

func absPathify(inPath string) string {
	log.Debug("resolve", inPath)

	if strings.HasPrefix(inPath, "$HOME") {
		inPath = userHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	color.Red("Couldn't discover absolute path")
	color.Red(err.Error())
	return ""
}

// Check if File / Directory Exists
func exists(fs afero.Fs, path string) (bool, error) {
	_, err := fs.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func safeMul(a, b uint) uint {
	c := a * b
	if a > 1 && b > 1 && c/b != a {
		return 0
	}
	return c
}

// parseSizeInBytes converts strings like 1GB or 12 mb into an unsigned integer number of bytes
func parseSizeInBytes(sizeStr string) uint {
	sizeStr = strings.TrimSpace(sizeStr)
	lastChar := len(sizeStr) - 1
	multiplier := uint(1)

	if lastChar > 0 {
		if sizeStr[lastChar] == 'b' || sizeStr[lastChar] == 'B' {
			if lastChar > 1 {
				switch unicode.ToLower(rune(sizeStr[lastChar-1])) {
				case 'k':
					multiplier = 1 << 10
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'm':
					multiplier = 1 << 20
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				case 'g':
					multiplier = 1 << 30
					sizeStr = strings.TrimSpace(sizeStr[:lastChar-1])
				default:
					multiplier = 1
					sizeStr = strings.TrimSpace(sizeStr[:lastChar])
				}
			}
		}
	}

	size := cast.ToInt(sizeStr)
	if size < 0 {
		size = 0
	}

	return safeMul(uint(size), multiplier)
}

func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}
