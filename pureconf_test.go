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
	"strconv"
	"testing"
)

func setEnvs(t *testing.T, envs map[string]string) {
	t.Helper()
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

func TestLoad(t *testing.T) {
	type AppConfig struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	}

	setEnvs(t, map[string]string{
		"MYAPP_HOST": "localhost",
		"MYAPP_PORT": "8080",
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

	setEnvs(t, map[string]string{
		"ERRAPP_PORT": "not-a-number",
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

func TestLoad_Secret(t *testing.T) {
	type SecretConfig struct {
		Token Secret[string] `env:"TOKEN"`
		Pin   Secret[int]    `env:"PIN"`
	}

	setEnvs(t, map[string]string{
		"SEC_TOKEN": "super-secret-token",
		"SEC_PIN":   "1234",
	})

	cfg, err := Load[SecretConfig](WithEnvPrefix("SEC_"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Token.Unmask() != "super-secret-token" {
		t.Errorf("expected Token 'super-secret-token', got '%s'", cfg.Token.Unmask())
	}
	if cfg.Pin.Unmask() != 1234 {
		t.Errorf("expected Pin 1234, got %d", cfg.Pin.Unmask())
	}
}

func TestLoad_Secret_Error(t *testing.T) {
	type SecretErrorConfig struct {
		Pin Secret[int] `env:"PIN"`
	}

	setEnvs(t, map[string]string{
		"SECERR_PIN": "not-a-number",
	})

	cfg, err := Load[SecretErrorConfig](WithEnvPrefix("SECERR_"))

	if err == nil {
		t.Fatal("expected error, but got nil")
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined errors, got %T", err)
	}

	errs := joined.Unwrap()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}

	var parseErr *ParseError
	if !errors.As(errs[0], &parseErr) {
		t.Fatalf("expected error to be of type *ParseError, got %T", errs[0])
	}

	if parseErr.Field != "Pin" {
		t.Errorf("expected Field 'Pin', got '%s'", parseErr.Field)
	}

	if parseErr.Value != "***" {
		t.Errorf("expected Value '***', got '%s'", parseErr.Value)
	}

	var numErr *strconv.NumError
	if !errors.As(parseErr.Err, &numErr) || numErr.Err != strconv.ErrSyntax {
		t.Errorf("expected underlying error to be strconv.ErrSyntax, got %v", parseErr.Err)
	}

	if cfg != nil {
		t.Errorf("expected config to be nil on error, got %v", cfg)
	}
}
