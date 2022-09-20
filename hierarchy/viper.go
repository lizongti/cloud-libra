package hierarchy

import (
	"io"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Get returns the value associated with the key.
func Get(key string) interface{} {
	return _default.Get(key)
}

// GetString returns the value associated with the key as a string.
func GetString(key string) string {
	return _default.GetString(key)
}

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int {
	return _default.GetInt(key)
}

// GetInt32 returns the value associated with the key as an integer.
func GetInt32(key string) int32 {
	return _default.GetInt32(key)
}

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 {
	return _default.GetInt64(key)
}

// GetUint returns the value associated with the key as an unsigned integer.
func GetUint(key string) uint {
	return _default.GetUint(key)
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func GetUint32(key string) uint32 {
	return _default.GetUint32(key)
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func GetUint64(key string) uint64 {
	return _default.GetUint64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key string) float64 {
	return _default.GetFloat64(key)
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool {
	return _default.GetBool(key)
}

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time {
	return _default.GetTime(key)
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key string) time.Duration {
	return _default.GetDuration(key)
}

// GetIntSlice returns the value associated with the key as a slice of integers.
func GetIntSlice(key string) []int {
	return _default.GetIntSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key string) map[string]interface{} {
	return _default.GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key string) map[string]string {
	return _default.GetStringMapString(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key string) []string {
	return _default.GetStringSlice(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key string) map[string][]string {
	return _default.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size in bytes for the given key.
func GetSizeInBytes(key string) uint {
	return _default.GetSizeInBytes(key)
}

// IsSet checks to see if a key is set in the config.
func IsSet(key string) bool {
	return _default.IsSet(key)
}

// AllSettings returns a map of all settings
func AllSettings() map[string]interface{} {
	return _default.AllSettings()
}

// AllKeys returns all keys
func AllKeys() []string {
	return _default.AllKeys()
}

// SetDefault sets the default value for a key.
func SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

// Set sets the value for a key.
func Set(key string, value interface{}) {
	viper.Set(key, value)
}

// SetConfigFile explicitly sets the config file to use.
func SetConfigFile(in string) {
	viper.SetConfigFile(in)
}

// SetConfigName sets the name of the config file without the extension.
func SetConfigName(in string) {
	viper.SetConfigName(in)
}

// SetConfigType sets the type of the configuration file.
func SetConfigType(in string) {
	viper.SetConfigType(in)
}

// SetEnvPrefix sets the environment variable prefix.
func SetEnvPrefix(in string) {
	viper.SetEnvPrefix(in)
}

// SetEnvKeyReplacer sets the environment variable key replacer.
func SetEnvKeyReplacer(r *strings.Replacer) {
	viper.SetEnvKeyReplacer(r)
}

// AutomaticEnv automatically binds all environment variables with the given prefix.
func AutomaticEnv() {
	viper.AutomaticEnv()
}

// BindPFlag binds a Viper flag to a pflag.Flag
func BindPFlag(key string, flag *pflag.Flag) error {
	return _default.BindPFlag(key, flag)
}

// BindPFlags binds a Viper flag to a pflag.FlagSet
func BindPFlags(flags *pflag.FlagSet) error {
	return _default.BindPFlags(flags)
}

// BindEnv binds a Viper key to an environment variable.
func BindEnv(input ...string) error {
	return _default.BindEnv(input...)
}

// BindEnvWithPrefix binds a Viper key to an environment variable with a prefix.
func BindEnvWithPrefix(prefix, key string) error {
	return _default.BindEnv(prefix + key)
}

// Unmarshal unmarshals the config into a struct.
func Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return _default.Unmarshal(rawVal, opts...)
}

// UnmarshalKey unmarshals a single key from the config into a struct.
func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return _default.UnmarshalKey(key, rawVal, opts...)
}

// AddConfigPath adds a path for Viper to search for the config file in.
func AddConfigPath(in string) {
	viper.AddConfigPath(in)
}

// ReadInConfig reads in a config file.
func ReadInConfig() error {
	return _default.ReadInConfig()
}

// ReadConfig reads in a config file.
func ReadConfig(in io.Reader) error {
	return _default.ReadConfig(in)
}

// MergeConfig merges a config file with the existing config.
func MergeConfig(in io.Reader) error {
	return _default.MergeConfig(in)
}

// MergeInConfig merges a config file with the existing config.
func MergeInConfig() error {
	return _default.MergeInConfig()
}

// WriteConfig writes the current config to a file.
func WriteConfig() error {
	return _default.WriteConfig()
}

// WriteConfigAs writes the current config to a file.
func WriteConfigAs(filename string) error {
	return _default.WriteConfigAs(filename)
}

// SafeWriteConfig writes the current config to a file, but only if it doesn't exist.
func SafeWriteConfig() error {
	return _default.SafeWriteConfig()
}

// SafeWriteConfigAs writes the current config to a file, but only if it doesn't exist.
func SafeWriteConfigAs(filename string) error {
	return _default.SafeWriteConfigAs(filename)
}

// ReadRemoteConfig reads a config from a remote source.
func ReadRemoteConfig() error {
	return _default.ReadRemoteConfig()
}

// WatchConfig watches a config file and calls the callback when the config changes.
func WatchConfig() {
	_default.WatchConfig()
}

// OnConfigChange registers a callback function to be called when the config is changed.
func OnConfigChange(run func(in fsnotify.Event)) {
	_default.OnConfigChange(run)
}

// SetFs sets the filesystem to be used by viper.
func SetFs(fs afero.Fs) {
	_default.SetFs(fs)
}

// AddRemoteProvider adds a remote provider to viper.
func AddRemoteProvider(provider, endpoint, path string) error {
	return _default.AddRemoteProvider(provider, endpoint, path)
}

// AddSecureRemoteProvider adds a secure remote provider to viper.
func AddSecureRemoteProvider(provider, endpoint, path, secret string) error {
	return _default.AddSecureRemoteProvider(provider, endpoint, path, secret)
}

// WatchRemoteCOnfigOnChannel watches a remote config and returns a channel to which all
func WatchRemoteConfigOnChannel() error {
	return _default.WatchRemoteConfigOnChannel()
}
