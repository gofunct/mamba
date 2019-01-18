package config

import (
	"github.com/spf13/cast"
	"strings"
	"time"
)

func (v *config) GetString(key string) string {
	return cast.ToString(v.Get(key))
}

func (v *config) GetBool(key string) bool {
	return cast.ToBool(v.Get(key))
}

func (v *config) GetInt(key string) int {
	return cast.ToInt(v.Get(key))
}

func (v *config) GetInt32(key string) int32 {
	return cast.ToInt32(v.Get(key))
}

func (v *config) GetInt64(key string) int64 {
	return cast.ToInt64(v.Get(key))
}

func (v *config) GetFloat64(key string) float64 {
	return cast.ToFloat64(v.Get(key))
}

func (v *config) GetTime(key string) time.Time {
	return cast.ToTime(v.Get(key))
}

func (v *config) GetDuration(key string) time.Duration {
	return cast.ToDuration(v.Get(key))
}

func (v *config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(v.Get(key))
}

func (v *config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(v.Get(key))
}

func (v *config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(v.Get(key))
}

func (v *config) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(v.Get(key))
}

func (v *config) Get(key string) interface{} {
	logkv(v, "event", "searching for key value...")

	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey)
	if val == nil {
		return nil
	}

	if v.typeByDefValue {
		// TODO(bep) this branch isn't covered by a single test.
		valType := val
		path := strings.Split(lcaseKey, v.keyDelim)
		defVal := v.searchMap(v.defaults, path)
		if defVal != nil {
			valType = defVal
		}

		switch valType.(type) {
		case bool:
			return cast.ToBool(val)
		case string:
			return cast.ToString(val)
		case int32, int16, int8, int:
			return cast.ToInt(val)
		case int64:
			return cast.ToInt64(val)
		case float64, float32:
			return cast.ToFloat64(val)
		case time.Time:
			return cast.ToTime(val)
		case time.Duration:
			return cast.ToDuration(val)
		case []string:
			return cast.ToStringSlice(val)
		}
	}

	return val
}
