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
	"os"
	"testing"
)

func TestApply_AllPrimitives(t *testing.T) {
	type AllPrimitivesConfig struct {
		StringVal  string  `env:"STRING_VAL"`
		BoolVal    bool    `env:"BOOL_VAL"`
		IntVal     int     `env:"INT_VAL"`
		Int8Val    int8    `env:"INT8_VAL"`
		Int16Val   int16   `env:"INT16_VAL"`
		Int32Val   int32   `env:"INT32_VAL"`
		Int64Val   int64   `env:"INT64_VAL"`
		UintVal    uint    `env:"UINT_VAL"`
		Uint8Val   uint8   `env:"UINT8_VAL"`
		Uint16Val  uint16  `env:"UINT16_VAL"`
		Uint32Val  uint32  `env:"UINT32_VAL"`
		Uint64Val  uint64  `env:"UINT64_VAL"`
		Float32Val float32 `env:"FLOAT32_VAL"`
		Float64Val float64 `env:"FLOAT64_VAL"`
	}

	envs := map[string]string{
		"MYAPP_STRING_VAL":  "pureconf-test",
		"MYAPP_BOOL_VAL":    "true",
		"MYAPP_INT_VAL":     "-42",
		"MYAPP_INT8_VAL":    "-128",
		"MYAPP_INT16_VAL":   "-32768",
		"MYAPP_INT32_VAL":   "-2147483648",
		"MYAPP_INT64_VAL":   "-9223372036854775808",
		"MYAPP_UINT_VAL":    "42",
		"MYAPP_UINT8_VAL":   "255",
		"MYAPP_UINT16_VAL":  "65535",
		"MYAPP_UINT32_VAL":  "4294967295",
		"MYAPP_UINT64_VAL":  "18446744073709551615",
		"MYAPP_FLOAT32_VAL": "3.1415",
		"MYAPP_FLOAT64_VAL": "2.718281828459",
	}

	for k, v := range envs {
		os.Setenv(k, v)
	}
	t.Cleanup(func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	})

	target := &AllPrimitivesConfig{}

	err := Apply(target, "MYAPP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target.StringVal != "pureconf-test" {
		t.Errorf("StringVal: got %v", target.StringVal)
	}
	if target.BoolVal != true {
		t.Errorf("BoolVal: got %v", target.BoolVal)
	}
	if target.IntVal != -42 {
		t.Errorf("IntVal: got %v", target.IntVal)
	}
	if target.Int8Val != -128 {
		t.Errorf("Int8Val: got %v", target.Int8Val)
	}
	if target.Int16Val != -32768 {
		t.Errorf("Int16Val: got %v", target.Int16Val)
	}
	if target.Int32Val != -2147483648 {
		t.Errorf("Int32Val: got %v", target.Int32Val)
	}
	if target.Int64Val != -9223372036854775808 {
		t.Errorf("Int64Val: got %v", target.Int64Val)
	}
	if target.UintVal != 42 {
		t.Errorf("UintVal: got %v", target.UintVal)
	}
	if target.Uint8Val != 255 {
		t.Errorf("Uint8Val: got %v", target.Uint8Val)
	}
	if target.Uint16Val != 65535 {
		t.Errorf("Uint16Val: got %v", target.Uint16Val)
	}
	if target.Uint32Val != 4294967295 {
		t.Errorf("Uint32Val: got %v", target.Uint32Val)
	}
	if target.Uint64Val != 18446744073709551615 {
		t.Errorf("Uint64Val: got %v", target.Uint64Val)
	}
	if target.Float32Val != 3.1415 {
		t.Errorf("Float32Val: got %v", target.Float32Val)
	}
	if target.Float64Val != 2.718281828459 {
		t.Errorf("Float64Val: got %v", target.Float64Val)
	}
}

func TestApply_Errors(t *testing.T) {
	type ErrorConfig struct {
		NoTag       string
		unexported  string   `env:"UNEXPORTED"`
		MissingEnv  string   `env:"MISSING_ENV"`
		BadInt      int      `env:"BAD_INT"`
		BadUint     uint     `env:"BAD_UINT"`
		BadFloat    float64  `env:"BAD_FLOAT"`
		BadBool     bool     `env:"BAD_BOOL"`
		Unsupported []string `env:"UNSUPPORTED"`
	}

	envs := map[string]string{
		"ERRTEST_UNEXPORTED":  "secret",
		"ERRTEST_BAD_INT":     "not-an-int",
		"ERRTEST_BAD_UINT":    "-1",
		"ERRTEST_BAD_FLOAT":   "not-a-float",
		"ERRTEST_BAD_BOOL":    "not-a-bool",
		"ERRTEST_UNSUPPORTED": "some-value",
	}

	for k, v := range envs {
		os.Setenv(k, v)
	}
	t.Cleanup(func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	})

	target := &ErrorConfig{}

	err := Apply(target, "ERRTEST_")

	if err == nil {
		t.Fatal("expected errors, but got nil")
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined errors, got %T", err)
	}

	errs := joined.Unwrap()
	if len(errs) != 5 {
		t.Fatalf("expected 5 errors, got %d", len(errs))
	}

	var parseErr *ParseError
	if !errors.As(errs[0], &parseErr) || parseErr.Field != "BadInt" {
		t.Errorf("error 0: expected ParseError for BadInt, got %s", errs[0].(*ParseError).Type)
	}

	if !errors.As(errs[1], &parseErr) || parseErr.Field != "BadUint" {
		t.Errorf("error 1: expected ParseError for BadUint, got %s", errs[1].(*ParseError).Type)
	}

	if !errors.As(errs[2], &parseErr) || parseErr.Field != "BadFloat" {
		t.Errorf("error 2: expected ParseError for BadFloat, got %s", errs[2].(*ParseError).Type)
	}

	if !errors.As(errs[3], &parseErr) || parseErr.Field != "BadBool" {
		t.Errorf("error 3: expected ParseError for BadBool, got %s", errs[3].(*ParseError).Type)
	}

	var unsupportErr *UnsupportedTypeError
	if !errors.As(errs[4], &unsupportErr) || unsupportErr.Field != "Unsupported" {
		t.Errorf("error 4: expected UnsupportedTypeError for Unsupported, got %T", errs[4].(*UnsupportedTypeError).Type)
	}
}
