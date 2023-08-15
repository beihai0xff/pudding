// Package provider implements a koanf.Provider that reads commandline
// parameters as conf maps.
package provider

import (
	"errors"
	"flag"

	"github.com/knadh/koanf/maps"
)

// KoanfIntf is an interface that represents a small subset of methods
// used by this package from Koanf{}. When using this package, a live
// instance of Koanf{} should be passed.
type KoanfIntf interface {
	Exists(string) bool
}

// Flag implements a pflag command line provider.
type Flag struct {
	delim   string
	flagSet *flag.FlagSet
	ko      KoanfIntf
	cb      func(key string) string
	flagCB  func(f *flag.Flag) (string, interface{})
}

// Provider returns a commandline flags provider that returns
// a nested map[string]interface{} of environment variable where the
// nesting hierarchy of keys are defined by delim. For instance, the
// delim "." will convert the key `parent.child.key: 1`
// to `{parent: {child: {key: 1}}}`.
//
// It takes an optional (but recommended) Koanf instance to see if the
// flags defined have been set from other providers, for instance,
// a config file. If they are not, then the default values of the flags
// are merged. If they do exist, the flag values are not merged but only
// the values that have been explicitly set in the command line are merged.
func Provider(f *flag.FlagSet, delim string, ko KoanfIntf) *Flag {
	return &Flag{
		flagSet: f,
		delim:   delim,
		ko:      ko,
	}
}

// ProviderWithKey works exactly the same as Provider except the callback
// takes the variable name allows their modification.
// This is useful for cases where complex types like slices separated by
// custom separators.
func ProviderWithKey(f *flag.FlagSet, delim string, ko KoanfIntf, cb func(key string) string) *Flag {
	return &Flag{
		flagSet: f,
		delim:   delim,
		ko:      ko,
		cb:      cb,
	}
}

// Read reads the flag variables and returns a nested conf map.
func (p *Flag) Read() (map[string]interface{}, error) {
	mp := make(map[string]interface{})

	p.flagSet.VisitAll(func(f *flag.Flag) {
		key := f.Name
		if p.cb != nil {
			// key :=
			key = p.cb(key)
		}

		// if the key is set, and the flag value is the default value, skip it
		if p.ko.Exists(key) && f.Value.(flag.Getter).Get() == f.DefValue {
			return
		}

		mp[key] = f.Value.(flag.Getter).Get()
	})

	return maps.Unflatten(mp, p.delim), nil
}

// ReadBytes is not supported by the flag provider.
func (p *Flag) ReadBytes() ([]byte, error) {
	return nil, errors.New("flag provider does not support this method")
}

// Watch is not supported.
func (p *Flag) Watch(cb func(event interface{}, err error)) error {
	return errors.New("flag provider does not support this method")
}
