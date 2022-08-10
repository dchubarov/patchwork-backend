package extension

import (
	"fmt"
	"strconv"
	"strings"
	"twowls.org/patchwork/commons/logging"
)

// Options contains extension setup options
type Options struct {
	// Debug specifies whether application (and therefore its extension) runs in debug
	Debug bool
	// Log allows extension to send messages to application log
	Log logging.Facade

	// cfg contains extension configuration
	cfg map[string]any
}

// EmptyOptions creates empty Options and returns its address
func EmptyOptions() *Options {
	return &Options{}
}

// PutConfig adds a configuration value
func (o *Options) PutConfig(key string, v any) *Options {
	if o.cfg == nil {
		o.cfg = make(map[string]any)
	}
	o.cfg[normalizeKey(key)] = v
	return o
}

// AnyConfig returns raw config value
func (o *Options) AnyConfig(key string) (any, bool) {
	if o.cfg != nil {
		if v, ok := o.cfg[normalizeKey(key)]; ok {
			return v, true
		}
	}
	return "", false
}

// StrConfig returns configuration value as string
func (o *Options) StrConfig(key string) (string, bool) {
	if v, ok := o.AnyConfig(key); ok {
		return fmt.Sprint(v), true
	}
	return "", false
}

func (o *Options) StrConfigDefault(key string, def string) string {
	if s, ok := o.StrConfig(key); ok {
		return s
	}
	return def
}

// BoolConfig returns a configuration value as bool
func (o *Options) BoolConfig(key string) (bool, bool) {
	if v, ok := o.AnyConfig(key); ok {
		if b, ok := v.(bool); ok {
			return b, true
		} else {
			if b, err := strconv.ParseBool(fmt.Sprint(v)); err == nil {
				return b, true
			}
		}
	}
	return false, false
}

func (o *Options) BoolConfigDefault(key string, def bool) bool {
	if b, ok := o.BoolConfig(key); ok {
		return b
	}
	return def
}

func normalizeKey(key string) string {
	return strings.ToLower(key)
}
