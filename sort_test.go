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

import "testing"

var lexicographicLessTests = []struct {
	i        string
	j        string
	expected bool
}{
	{"", "a", true},
	{"a", "", false},
	{"a", "a", false},
	{"a", "A", false},
	{"A", "a", true},
	{"aa", "a", false},
	{"a", "aa", true},
	{"a", "b", true},
	{"a", "B", true},
	{"A", "b", true},
	{"A", "B", true},
}

func TestLexicographicLess(t *testing.T) {
	for _, test := range lexicographicLessTests {
		actual := lexicographicLess(test.i, test.j)
		if test.expected != actual {
			t.Errorf(`expected string "%s" to come before "%s"`, test.i, test.j)
		}
	}
}
