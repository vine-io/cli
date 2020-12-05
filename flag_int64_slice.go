// Copyright 2020 The vine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Int64Slice wraps []int64 to satisfy flag.Value
type Int64Slice struct {
	slice      []int64
	hasBeenSet bool
}

// NewInt64Slice makes a *Int64Slice with default values
func NewInt64Slice(defaults ...int64) *Int64Slice {
	return &Int64Slice{slice: append([]int64{}, defaults...)}
}

// Set parses the value into a int64 and appends it to the list of values
func (f *Int64Slice) Set(value string) error {
	if !f.hasBeenSet {
		f.slice = []int64{}
		f.hasBeenSet = true
	}

	if strings.HasPrefix(value, slPfx) {
		// Deserializing assumes overwrite
		_ = json.Unmarshal([]byte(strings.Replace(value, slPfx, "", 1)), &f.slice)
		f.hasBeenSet = true
		return nil
	}

	tmp, err := strconv.ParseInt(value, 10,64)
	if err != nil {
		return err
	}
	f.slice = append(f.slice, int64(tmp))
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (f *Int64Slice) String() string {
	return fmt.Sprintf("%#v", f.slice)
}

// Serialize allows Int64Slice to fulfill Serializer
func (f *Int64Slice) Serialize() string {
	jsonBytes, _ := json.Marshal(f.slice)
	return fmt.Sprintf("%s%s", slPfx, string(jsonBytes))
}

// Value returns the slice of []int64 set by this flag
func (f *Int64Slice) Value() []int64 {
	return f.slice
}

// Get returns the slice of []int64 set by this flag
func (f *Int64Slice) Get() interface{} {
	return *f
}

// Int64SliceFlag is a flag with type bool
type Int64SliceFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	EnvVars     []string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       *Int64Slice
	DefaultText string
	HasBeenSet  bool
}

// IsSet returns whether or not the flag has been set through env or file
func (f *Int64SliceFlag) IsSet() bool {
	return f.HasBeenSet
}

// String returns a readable representation of this value
// (for usage defaults)
func (f *Int64SliceFlag) String() string {
	return FlagStringer(f)
}

// Names returns the names of the flag
func (f *Int64SliceFlag) Names() []string {
	return flagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *Int64SliceFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise flag
func (f *Int64SliceFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f *Int64SliceFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *Int64SliceFlag) GetValue() string {
	return ""
}

// Apply populates the flag given the flag set and environment
func (f *Int64SliceFlag) Apply(set *flag.FlagSet) error {
	if val, ok := flagFromEnvOrFile(f.EnvVars, f.FilePath); ok {
		if val != "" {
			f.Value = &Int64Slice{}

			for _, s := range strings.Split(val, ",") {
				if err := f.Value.Set(strings.TrimSpace(s)); err != nil {
					return fmt.Errorf("could not parse %q as int64 slice value for flag %s: %v", val, f.Name, err)
				}
			}

			f.HasBeenSet = true
		}
	}

	for _, name := range f.Names() {
		if f.Value != nil {
			f.Value = &Int64Slice{}
		}
		set.Var(f.Value, name, f.Usage)
	}

	return nil
}

// Int64Slice looks up the value of a local Int64SliceFlag, returns
// nil if not found
func (c *Context) Int64Slice(name string) []int64 {
	if fs := lookupFlagSet(name, c); fs != nil {
		return lookupInt64Slice(name, fs)
	}
	return nil
}

func lookupInt64Slice(name string, set *flag.FlagSet) []int64 {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := (f.Value.(*Int64Slice)).Value(), error(nil)
		if err != nil {
			return nil
		}
		return parsed
	}
	return nil
}
