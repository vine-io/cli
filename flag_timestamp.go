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
	"flag"
	"fmt"
	"strings"
	"time"
)

// Timestamp wraps to satisfy golang's flag interface
type Timestamp struct {
	timestamp *time.Time
	hasBeenSet bool
	layout string
}

func NewTimestamp(timestamp time.Time) *Timestamp {
	return &Timestamp{timestamp: &timestamp}
}

// Set the timestamp value directly
func (t *Timestamp) SetTimestamp(value time.Time) {
	if !t.hasBeenSet {
		t.timestamp = &value
		t.hasBeenSet = true
	}
}

// Set the timestamp string layout for future parsing
func (t *Timestamp) SetLayout(layout string) {
	t.layout = layout
}

// Set parses the value into a int and appends it to the list of values
func (t *Timestamp) Set(value string) error {
	timpstamp, err := time.Parse(t.layout, value)
	if err != nil {
		return err
	}

	t.timestamp = &timpstamp
	t.hasBeenSet = true
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (t *Timestamp) String() string {
	return fmt.Sprintf("%#v", t.timestamp)
}

// Value returns the slice of timestamp set by this flag
func (t *Timestamp) Value() *time.Time {
	return t.timestamp
}

// Get returns the slice of timestamp set by this flag
func (t *Timestamp) Get() interface{} {
	return *t
}

// TimestampFlag is a flag with type bool
type TimestampFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	EnvVars     []string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       *Timestamp
	DefaultText string
	HasBeenSet  bool
}

// IsSet returns whether or not the flag has been set through env or file
func (f *TimestampFlag) IsSet() bool {
	return f.HasBeenSet
}

// String returns a readable representation of this value
// (for usage defaults)
func (f *TimestampFlag) String() string {
	return FlagStringer(f)
}

// Names returns the names of the flag
func (f *TimestampFlag) Names() []string {
	return flagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *TimestampFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise flag
func (f *TimestampFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f *TimestampFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *TimestampFlag) GetValue() string {
	return ""
}

// Apply populates the flag given the flag set and environment
func (f *TimestampFlag) Apply(set *flag.FlagSet) error {
	if val, ok := flagFromEnvOrFile(f.EnvVars, f.FilePath); ok {
		if val != "" {
			f.Value = &Timestamp{}

			for _, s := range strings.Split(val, ",") {
				if err := f.Value.Set(strings.TrimSpace(s)); err != nil {
					return fmt.Errorf("could not parse %q as timestamp slice value for flag %s: %v", val, f.Name, err)
				}
			}

			f.HasBeenSet = true
		}
	}

	for _, name := range f.Names() {
		if f.Value != nil {
			f.Value = &Timestamp{}
		}
		set.Var(f.Value, name, f.Usage)
	}

	return nil
}

// Timestamp looks up the value of a local TimestampFlag, returns
// nil if not found
func (c *Context) Timestamp(name string) *time.Time {
	if fs := lookupFlagSet(name, c); fs != nil {
		return lookupTimestamp(name, fs)
	}
	return nil
}

func lookupTimestamp(name string, set *flag.FlagSet) *time.Time {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := (f.Value.(*Timestamp)).Value(), error(nil)
		if err != nil {
			return nil
		}
		return parsed
	}
	return nil
}
