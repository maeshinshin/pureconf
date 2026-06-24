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
	"encoding"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Apply[T any](target *T, prefix string) error {
	if target == nil {
		return errors.New("target cannot be nil")
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}

	return applyRecursive(v.Elem(), prefix)
}

func applyRecursive(v reflect.Value, prefix string) error {
	var errs []error
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		if err := processField(v.Field(i), field, prefix); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func processField(fieldVal reflect.Value, field reflect.StructField, prefix string) error {
	tag := field.Tag.Get("env")
	if tag == "" {
		tag = strings.ToUpper(field.Name)
	}
	envVarName := prefix + tag

	if fieldVal.CanAddr() {
		addr := fieldVal.Addr().Interface()
		if unmarshaler, ok := addr.(encoding.TextUnmarshaler); ok {
			envVal, exists := resolveEnvValue(envVarName, field, fieldVal)
			if !exists {
				return nil
			}

			if err := unmarshaler.UnmarshalText([]byte(envVal)); err != nil {
				displayVal := envVal
				if sensitiveObj, isSensitive := addr.(interface{ IsSensitive() bool }); isSensitive && sensitiveObj.IsSensitive() {
					displayVal = "***"
				}
				return &ParseError{
					Field: field.Name,
					Type:  fieldVal.Type().String(),
					Value: displayVal,
					Err:   err,
				}
			}
			return nil
		}
	}

	if fieldVal.Kind() == reflect.Struct {
		if fieldVal.CanAddr() {
			return applyRecursive(fieldVal, envVarName+"_")
		}
	}

	envVal, exists := resolveEnvValue(envVarName, field, fieldVal)
	if !exists {
		return nil
	}

	return setPrimitive(fieldVal, field, envVal)
}

func resolveEnvValue(envVar string, field reflect.StructField, fieldVal reflect.Value) (string, bool) {
	if envVal, exists := os.LookupEnv(envVar); exists {
		return envVal, true
	}

	defaultVal, hasDefault := field.Tag.Lookup("default")
	if hasDefault && fieldVal.IsZero() {
		return defaultVal, true
	}

	return "", false
}

func setPrimitive(fieldVal reflect.Value, field reflect.StructField, envVal string) error {
	var err error

	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(envVal)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		if n, err = strconv.ParseInt(envVal, 10, fieldVal.Type().Bits()); err == nil {
			fieldVal.SetInt(n)
			return nil
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var n uint64
		if n, err = strconv.ParseUint(envVal, 10, fieldVal.Type().Bits()); err == nil {
			fieldVal.SetUint(n)
			return nil
		}

	case reflect.Float32, reflect.Float64:
		var n float64
		if n, err = strconv.ParseFloat(envVal, fieldVal.Type().Bits()); err == nil {
			fieldVal.SetFloat(n)
			return nil
		}

	case reflect.Bool:
		var b bool
		if b, err = strconv.ParseBool(envVal); err == nil {
			fieldVal.SetBool(b)
			return nil
		}

	default:
		return &UnsupportedTypeError{
			Field: field.Name,
			Type:  fieldVal.Type().String(),
		}
	}

	return &ParseError{
		Field: field.Name,
		Type:  fieldVal.Type().String(),
		Value: envVal,
		Err:   err,
	}
}
