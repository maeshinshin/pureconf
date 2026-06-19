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

package envmap

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseError(t *testing.T) {
	baseErr := strconv.ErrSyntax

	err := &ParseError{
		Field: "Port",
		Type:  "int",
		Value: "not-a-number",
		Err:   baseErr,
	}

	expectedMsg := "failed to parse 'not-a-number' as int for field Port: invalid syntax"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}

	if !errors.Is(err, baseErr) {
		t.Errorf("expected errors.Is to match the base error")
	}

	var target *ParseError
	if !errors.As(err, &target) {
		t.Errorf("expected errors.As to match ParseError type")
	}
}

func TestUnsupportedTypeError(t *testing.T) {
	err := &UnsupportedTypeError{
		Field: "Items",
		Type:  "slice",
	}

	expectedMsg := "unsupported field type slice for field Items"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}
