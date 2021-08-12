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

package altsrc

import (
	"time"

	"github.com/vine-io/cli"
)

// InputSourceContext is an interface used to allow
// other input sources to be implemented as needed.
//
// Source returns an identifier for the input source. In case of file source
// it should return path to the file.
type InputSourceContext interface {
	Source() string

	Int(name string) (int, error)
	Duration(name string) (time.Duration, error)
	Float64(name string) (float64, error)
	String(name string) (string, error)
	StringSlice(name string) ([]string, error)
	IntSlice(name string) ([]int, error)
	Generic(name string) (cli.Generic, error)
	Bool(name string) (bool, error)
}
