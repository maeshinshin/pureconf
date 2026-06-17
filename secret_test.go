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
	"fmt"
	"log/slog"
	"testing"
)

func TestSecret_Unmask(t *testing.T) {
	s := Secret[string]{value: "my-super-secret-password"}
	if got := s.Unmask(); got != "my-super-secret-password" {
		t.Errorf("Unmask() = %v, want %v", got, "my-super-secret-password")
	}

	sInt := Secret[int]{value: 12345}
	if got := sInt.Unmask(); got != 12345 {
		t.Errorf("Unmask() = %v, want %v", got, 12345)
	}
}

func TestSecret_String(t *testing.T) {
	s := Secret[string]{value: "secret-data"}

	if got := s.String(); got != "***" {
		t.Errorf("String() = %v, want ***", got)
	}

	if got := fmt.Sprintf("%s", s); got != "***" {
		t.Errorf("fmt.Sprintf() = %v, want ***", got)
	}
}

func TestSecret_LogValue(t *testing.T) {
	s := Secret[string]{value: "secret-data"}
	val := s.LogValue()

	if val.Kind() != slog.KindString {
		t.Errorf("LogValue().Kind() = %v, want %v", val.Kind(), slog.KindString)
	}

	if got := val.String(); got != "***" {
		t.Errorf("LogValue().String() = %v, want ***", got)
	}
}
