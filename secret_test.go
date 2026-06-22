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
	"encoding"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
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

func TestSecret_IsSensitive(t *testing.T) {
	s := new(Secret[string])
	if !s.IsSensitive() {
		t.Errorf("expected IsSensitive to return true")
	}
}

func TestSecret_String(t *testing.T) {
	s := Secret[string]{value: "secret-data"}

	if got := s.String(); got != "***" {
		t.Errorf("String() = %v, want ***", got)
	}

	if got := fmt.Sprintf("%v", s); got != "***" {
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

func TestSecret_UnmarshalText_Success(t *testing.T) {
	sStr := new(Secret[string])
	sInt := new(Secret[int])
	sInt8 := new(Secret[int8])
	sInt16 := new(Secret[int16])
	sInt32 := new(Secret[int32])
	sInt64 := new(Secret[int64])
	sUint := new(Secret[uint])
	sUint8 := new(Secret[uint8])
	sUint16 := new(Secret[uint16])
	sUint32 := new(Secret[uint32])
	sUint64 := new(Secret[uint64])
	sFloat32 := new(Secret[float32])
	sFloat64 := new(Secret[float64])
	sBool := new(Secret[bool])

	tests := []struct {
		name        string
		input       string
		unmarshaler encoding.TextUnmarshaler
		check       func() bool
	}{
		{
			name:        "string",
			input:       "my-secret",
			unmarshaler: sStr,
			check:       func() bool { return sStr.Unmask() == "my-secret" },
		},
		{
			name:        "int",
			input:       "42",
			unmarshaler: sInt,
			check:       func() bool { return sInt.Unmask() == 42 },
		},
		{
			name:        "int8",
			input:       "127",
			unmarshaler: sInt8,
			check:       func() bool { return sInt8.Unmask() == 127 },
		},
		{
			name:        "int16",
			input:       "32767",
			unmarshaler: sInt16,
			check:       func() bool { return sInt16.Unmask() == 32767 },
		},
		{
			name:        "int32",
			input:       "2147483647",
			unmarshaler: sInt32,
			check:       func() bool { return sInt32.Unmask() == 2147483647 },
		},
		{
			name:        "int64",
			input:       "9223372036854775807",
			unmarshaler: sInt64,
			check:       func() bool { return sInt64.Unmask() == 9223372036854775807 },
		},
		{
			name:        "uint",
			input:       "42",
			unmarshaler: sUint,
			check:       func() bool { return sUint.Unmask() == 42 },
		},
		{
			name:        "uint8",
			input:       "255",
			unmarshaler: sUint8,
			check:       func() bool { return sUint8.Unmask() == 255 },
		},
		{
			name:        "uint16",
			input:       "65535",
			unmarshaler: sUint16,
			check:       func() bool { return sUint16.Unmask() == 65535 },
		},
		{
			name:        "uint32",
			input:       "4294967295",
			unmarshaler: sUint32,
			check:       func() bool { return sUint32.Unmask() == 4294967295 },
		},
		{
			name:        "uint64",
			input:       "18446744073709551615",
			unmarshaler: sUint64,
			check:       func() bool { return sUint64.Unmask() == 18446744073709551615 },
		},
		{
			name:        "float32",
			input:       "3.141592653589793",
			unmarshaler: sFloat32,
			check:       func() bool { return sFloat32.Unmask() == 3.141592653589793 },
		},
		{
			name:        "float64",
			input:       "2.718281828459045",
			unmarshaler: sFloat64,
			check:       func() bool { return sFloat64.Unmask() == 2.718281828459045 },
		},
		{
			name:        "bool",
			input:       "true",
			unmarshaler: sBool,
			check:       func() bool { return sBool.Unmask() == true },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.unmarshaler.UnmarshalText([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.check() {
				t.Errorf("unmasked value mismatch")
			}
		})
	}
}

func TestSecret_UnmarshalText_Errors(t *testing.T) {
	sInt := new(Secret[int])
	sInt8 := new(Secret[int8])
	sInt16 := new(Secret[int16])
	sInt32 := new(Secret[int32])
	sInt64 := new(Secret[int64])
	sUint := new(Secret[uint])
	sUint8 := new(Secret[uint8])
	sUint16 := new(Secret[uint16])
	sUint32 := new(Secret[uint32])
	sUint64 := new(Secret[uint64])
	sFloat32 := new(Secret[float32])
	sFloat64 := new(Secret[float64])
	sBool := new(Secret[bool])

	type UnsupportedStruct struct{}
	sUnsupport := new(Secret[UnsupportedStruct])

	tests := []struct {
		name        string
		typeName    string
		input       string
		unmarshaler encoding.TextUnmarshaler
	}{
		{
			name:        "invalid int",
			typeName:    "int",
			input:       "not-an-int",
			unmarshaler: sInt,
		},
		{
			name:        "invalid int8",
			typeName:    "int8",
			input:       "128",
			unmarshaler: sInt8,
		},
		{
			name:        "invalid int16",
			typeName:    "int16",
			input:       "32768",
			unmarshaler: sInt16,
		},
		{
			name:        "invalid int32",
			typeName:    "int32",
			input:       "2147483648",
			unmarshaler: sInt32,
		},
		{
			name:        "invalid int64",
			typeName:    "int64",
			input:       "9223372036854775808",
			unmarshaler: sInt64,
		},
		{
			name:        "invalid uint",
			typeName:    "uint",
			input:       "-1",
			unmarshaler: sUint,
		},
		{
			name:        "invalid uint8",
			typeName:    "uint8",
			input:       "256",
			unmarshaler: sUint8,
		},
		{
			name:        "invalid uint16",
			typeName:    "uint16",
			input:       "65536",
			unmarshaler: sUint16,
		},
		{
			name:        "invalid uint32",
			typeName:    "uint32",
			input:       "4294967296",
			unmarshaler: sUint32,
		},
		{
			name:        "invalid uint64",
			typeName:    "uint64",
			input:       "18446744073709551616",
			unmarshaler: sUint64,
		},

		{
			name:        "invalid float32",
			typeName:    "float32",
			input:       "not-a-float",
			unmarshaler: sFloat32,
		},
		{
			name:        "invalid float64",
			typeName:    "float64",
			input:       "not-a-float",
			unmarshaler: sFloat64,
		},
		{
			name:        "invalid bool",
			typeName:    "bool",
			input:       "not-a-bool",
			unmarshaler: sBool,
		},
		{
			name:        "unsupported type",
			typeName:    "UnsupportedStruct",
			input:       "some-value",
			unmarshaler: sUnsupport,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.unmarshaler.UnmarshalText([]byte(tt.input))
			if tt.typeName != "UnsupportedStruct" {
				var numErr *strconv.NumError
				if !errors.As(err, &numErr) {
					t.Errorf("expected error for %s, got: %v", tt.name, err)
				}
			} else {
				if err.Error() != "unsupported type: struct" {
					t.Errorf("expected unsupported type error for %s, got: %v", tt.name, err)
				}
			}
		})
	}
}
