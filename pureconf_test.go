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
	"errors"
	"os"
	"strconv"
	"testing"
)

func TestLoad(t *testing.T) {
	type AppConfig struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}

	os.Setenv("MYAPP_HOST", "localhost")
	os.Setenv("MYAPP_PORT", "8080")
	t.Cleanup(func() {
		os.Unsetenv("MYAPP_HOST")
		os.Unsetenv("MYAPP_PORT")
	})

	cfg, err := Load[AppConfig](WithEnvPrefix("MYAPP_"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected Host 'localhost', got '%s'", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected Port 8080, got %d", cfg.Port)
	}
}

func TestLoad_Error(t *testing.T) {
	type ErrorAppConfig struct {
		Port int `env:"PORT"`
	}

	os.Setenv("ERRAPP_PORT", "not-a-number")
	t.Cleanup(func() {
		os.Unsetenv("ERRAPP_PORT")
	})

	cfg, err := Load[ErrorAppConfig](WithEnvPrefix("ERRAPP_"))

	if err == nil {
		t.Fatal("expected errors, but got nil")
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined errors, got %T", err)
	}

	errs := joined.Unwrap()
	if len(errs) != 1 {
		t.Fatalf("expected 1 errors, got %d", len(errs))
	}

	var parseErr *ParseError
	if !errors.As(errs[0], &parseErr) {
		t.Fatalf("expected error to be of type *ParseError, got %T", errs[0])
	}

	if parseErr.Field != "Port" {
		t.Errorf("expected Field 'Port', got '%s'", parseErr.Field)
	}

	if parseErr.Type != "int" {
		t.Errorf("expected Type 'int', got '%s'", parseErr.Type)
	}

	if parseErr.Value != "not-a-number" {
		t.Errorf("expected Value 'not-a-number', got '%s'", parseErr.Value)
	}

	if !errors.Is(parseErr.Err, strconv.ErrSyntax) {
		t.Errorf("expected underlying error to be strconv.ErrSyntax, got %v", parseErr.Err)
	}

	if cfg != nil {
		t.Errorf("expected config to be nil on error, got %v", cfg)
	}
}
