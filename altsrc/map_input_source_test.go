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
	"testing"
	"time"
)

func TestMapDuration(t *testing.T) {
	inputSource := &MapInputSource{
		file: "test",
		valueMap: map[interface{}]interface{}{
			"duration_of_duration_type": time.Minute,
			"duration_of_string_type":   "1m",
			"duration_of_int_type":      1000,
		},
	}
	d, err := inputSource.Duration("duration_of_duration_type")
	expect(t, time.Minute, d)
	expect(t, nil, err)
	d, err = inputSource.Duration("duration_of_string_type")
	expect(t, time.Minute, d)
	expect(t, nil, err)
	_, err = inputSource.Duration("duration_of_int_type")
	refute(t, nil, err)
}
