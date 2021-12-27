package tree

import (
	"github.com/spf13/cast"
)

type Tree struct {
	data interface{}
	// config map[string]interface{}

	// typeByDefValue bool
}

func NewTree(v interface{}) {
	// if data, ok := v.(map[string]interface{}); ok {
	// 	return &Tree{
	// 		data: data,
	// 	}
	// }
}

// var v *Tree

// // Get can retrieve any value given the key to use.
// // Get is case-insensitive for a key.
// // Get has the behavior of returning the value associated with the first
// // place from where it is set. Viper will check in the following order:
// // override, flag, env, config file, key/value store, default
// //
// // Get returns an interface. For a specific value use one of the Get____ methods.
// func Get(key string) interface{} { return v.Get(key) }

// func (v *Tree) Get(key string) interface{} {
// 	lcaseKey := strings.ToLower(key)
// 	val := v.find(lcaseKey, true)
// 	if val == nil {
// 		return nil
// 	}

// 	if v.typeByDefValue {
// 		// TODO(bep) this branch isn't covered by a single test.
// 		valType := val
// 		path := strings.Split(lcaseKey, v.keyDelim)
// 		defVal := v.searchMap(v.defaults, path)
// 		if defVal != nil {
// 			valType = defVal
// 		}

// 		switch valType.(type) {
// 		case bool:
// 			return cast.ToBool(val)
// 		case string:
// 			return cast.ToString(val)
// 		case int32, int16, int8, int:
// 			return cast.ToInt(val)
// 		case uint:
// 			return cast.ToUint(val)
// 		case uint32:
// 			return cast.ToUint32(val)
// 		case uint64:
// 			return cast.ToUint64(val)
// 		case int64:
// 			return cast.ToInt64(val)
// 		case float64, float32:
// 			return cast.ToFloat64(val)
// 		case time.Time:
// 			return cast.ToTime(val)
// 		case time.Duration:
// 			return cast.ToDuration(val)
// 		case []string:
// 			return cast.ToStringSlice(val)
// 		case []int:
// 			return cast.ToIntSlice(val)
// 		}
// 	}

// 	return val
// }

// // Given a key, find the value.
// //
// // Viper will check to see if an alias exists first.
// // Viper will then check in the following order:
// // flag, env, config file, key/value store.
// // Lastly, if no value was found and flagDefault is true, and if the key
// // corresponds to a flag, the flag's default value is returned.
// //
// // Note: this assumes a lower-cased key given.
// func (v *Viper) find(lcaseKey string, flagDefault bool) interface{} {
// 	var (
// 		val    interface{}
// 		exists bool
// 		path   = strings.Split(lcaseKey, v.keyDelim)
// 		nested = len(path) > 1
// 	)

// 	// compute the path through the nested maps to the nested value
// 	if nested && v.isPathShadowedInDeepMap(path, castMapStringToMapInterface(v.aliases)) != "" {
// 		return nil
// 	}

// 	// if the requested key is an alias, then return the proper key
// 	lcaseKey = v.realKey(lcaseKey)
// 	path = strings.Split(lcaseKey, v.keyDelim)
// 	nested = len(path) > 1

// 	// Set() override first
// 	val = v.searchMap(v.override, path)
// 	if val != nil {
// 		return val
// 	}
// 	if nested && v.isPathShadowedInDeepMap(path, v.override) != "" {
// 		return nil
// 	}

// 	// PFlag override next
// 	flag, exists := v.pflags[lcaseKey]
// 	if exists && flag.HasChanged() {
// 		switch flag.ValueType() {
// 		case "int", "int8", "int16", "int32", "int64":
// 			return cast.ToInt(flag.ValueString())
// 		case "bool":
// 			return cast.ToBool(flag.ValueString())
// 		case "stringToString":
// 			return stringToStringConv(flag.ValueString())
// 		default:
// 			return flag.ValueString()
// 		}
// 	}
// 	if nested && v.isPathShadowedInFlatMap(path, v.pflags) != "" {
// 		return nil
// 	}

// 	// Env override next
// 	if v.automaticEnvApplied {
// 		// even if it hasn't been registered, if automaticEnv is used,
// 		// check any Get request
// 		if val, ok := v.getEnv(v.mergeWithEnvPrefix(lcaseKey)); ok {
// 			return val
// 		}
// 		if nested && v.isPathShadowedInAutoEnv(path) != "" {
// 			return nil
// 		}
// 	}
// 	envkeys, exists := v.env[lcaseKey]
// 	if exists {
// 		for _, envkey := range envkeys {
// 			if val, ok := v.getEnv(envkey); ok {
// 				return val
// 			}
// 		}
// 	}
// 	if nested && v.isPathShadowedInFlatMap(path, v.env) != "" {
// 		return nil
// 	}

// 	// Config file next
// 	val = v.searchIndexableWithPathPrefixes(v.config, path)
// 	if val != nil {
// 		return val
// 	}
// 	if nested && v.isPathShadowedInDeepMap(path, v.config) != "" {
// 		return nil
// 	}

// 	// K/V store next
// 	val = v.searchMap(v.kvstore, path)
// 	if val != nil {
// 		return val
// 	}
// 	if nested && v.isPathShadowedInDeepMap(path, v.kvstore) != "" {
// 		return nil
// 	}

// 	// Default next
// 	val = v.searchMap(v.defaults, path)
// 	if val != nil {
// 		return val
// 	}
// 	if nested && v.isPathShadowedInDeepMap(path, v.defaults) != "" {
// 		return nil
// 	}

// 	if flagDefault {
// 		// last chance: if no value is found and a flag does exist for the key,
// 		// get the flag's default value even if the flag's value has not been set.
// 		if flag, exists := v.pflags[lcaseKey]; exists {
// 			switch flag.ValueType() {
// 			case "int", "int8", "int16", "int32", "int64":
// 				return cast.ToInt(flag.ValueString())
// 			case "bool":
// 				return cast.ToBool(flag.ValueString())
// 			case "stringToString":
// 				return stringToStringConv(flag.ValueString())
// 			default:
// 				return flag.ValueString()
// 			}
// 		}
// 		// last item, no need to check shadowing
// 	}

// 	return nil
// }

// // mostly copied from pflag's implementation of this operation here https://github.com/spf13/pflag/blob/master/string_to_string.go#L79
// // alterations are: errors are swallowed, map[string]interface{} is returned in order to enable cast.ToStringMap
// func stringToStringConv(val string) interface{} {
// 	val = strings.Trim(val, "[]")
// 	// An empty string would cause an empty map
// 	if len(val) == 0 {
// 		return map[string]interface{}{}
// 	}
// 	r := csv.NewReader(strings.NewReader(val))
// 	ss, err := r.Read()
// 	if err != nil {
// 		return nil
// 	}
// 	out := make(map[string]interface{}, len(ss))
// 	for _, pair := range ss {
// 		kv := strings.SplitN(pair, "=", 2)
// 		if len(kv) != 2 {
// 			return nil
// 		}
// 		out[kv[0]] = kv[1]
// 	}
// 	return out
// }

// func castToMapStringInterface(
// 	src map[interface{}]interface{}) map[string]interface{} {
// 	tgt := map[string]interface{}{}
// 	for k, v := range src {
// 		tgt[fmt.Sprintf("%v", k)] = v
// 	}
// 	return tgt
// }

// func castMapStringSliceToMapInterface(src map[string][]string) map[string]interface{} {
// 	tgt := map[string]interface{}{}
// 	for k, v := range src {
// 		tgt[k] = v
// 	}
// 	return tgt
// }

// func castMapStringToMapInterface(src map[string]string) map[string]interface{} {
// 	tgt := map[string]interface{}{}
// 	for k, v := range src {
// 		tgt[k] = v
// 	}
// 	return tgt
// }

// searchMap recursively searches for a value for path in source map.
// Returns nil if not found.
// Note: This assumes that the path entries and map keys are lower cased.
func (v *Tree) searchMap(source map[string]interface{}, path []string) interface{} {
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
