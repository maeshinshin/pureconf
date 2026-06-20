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
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
)

type Secret[T any] struct {
	value T
}

var _ encoding.TextUnmarshaler = (*Secret[string])(nil)

func (s Secret[T]) Unmask() T {
	return s.value
}

func (s *Secret[T]) IsSensitive() bool {
	return true
}

func (s Secret[T]) String() string {
	return "***"
}

func (s Secret[T]) LogValue() slog.Value {
	return slog.StringValue("***")
}

func (s *Secret[T]) UnmarshalText(text []byte) error {
	val := string(text)

	v := reflect.ValueOf(&s.value).Elem()

	switch v.Kind() {
	case reflect.String:
		v.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(val, 10, v.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", v.Kind(), err)
		}
		v.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(val, 10, v.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", v.Kind(), err)
		}
		v.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, v.Type().Bits())
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", v.Kind(), err)
		}
		v.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", v.Kind(), err)
		}
		v.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}

	return nil
}
