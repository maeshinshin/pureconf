// Copyright 2026 maeshinshin
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

package pureconf

import (
	"testing"
)

func TestWithFile(t *testing.T) {
	opts := &options{}
	optFunc := WithFile("config.json")
	optFunc(opts)

	if opts.filePath != "config.json" {
		t.Errorf("expected filePath to be 'config.json', got '%s'", opts.filePath)
	}
}

func TestWithEnvPrefix(t *testing.T) {
	opts := &options{}
	optFunc := WithEnvPrefix("MYAPP_")
	optFunc(opts)

	if opts.envPrefix != "MYAPP_" {
		t.Errorf("expected envPrefix to be 'MYAPP_', got '%s'", opts.envPrefix)
	}
}
