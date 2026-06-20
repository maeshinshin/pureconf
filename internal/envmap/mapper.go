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
)

func Apply[T any](target *T, prefix string) error {
	v := reflect.ValueOf(target).Elem()

	var errs []error

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}

		envVar := prefix + tag
		envVal, exists := os.LookupEnv(envVar)
		if !exists {
			continue
		}

		fieldVal := v.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		if fieldVal.CanAddr() {
			addr := fieldVal.Addr().Interface()
			if unmarshaler, ok := addr.(encoding.TextUnmarshaler); ok {
				if err := unmarshaler.UnmarshalText([]byte(envVal)); err != nil {
					displayVal := envVal
					if sensitiveObj, isSensitive := addr.(interface{ IsSensitive() bool }); isSensitive && sensitiveObj.IsSensitive() {
						displayVal = "***"
					}
					errs = append(errs, &ParseError{
						Field: field.Name,
						Type:  fieldVal.Type().String(),
						Value: displayVal,
						Err:   err,
					})
				}
				continue
			}
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(envVal)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if n, err := strconv.ParseInt(envVal, 10, fieldVal.Type().Bits()); err == nil {
				fieldVal.SetInt(n)
			} else {
				errs = append(errs, &ParseError{
					Field: field.Name,
					Type:  fieldVal.Type().String(),
					Value: envVal,
					Err:   err,
				})
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if n, err := strconv.ParseUint(envVal, 10, fieldVal.Type().Bits()); err == nil {
				fieldVal.SetUint(n)
			} else {
				errs = append(errs, &ParseError{
					Field: field.Name,
					Type:  fieldVal.Type().String(),
					Value: envVal,
					Err:   err,
				})
			}
		case reflect.Float32, reflect.Float64:
			if n, err := strconv.ParseFloat(envVal, fieldVal.Type().Bits()); err == nil {
				fieldVal.SetFloat(n)
			} else {
				errs = append(errs, &ParseError{
					Field: field.Name,
					Type:  fieldVal.Type().String(),
					Value: envVal,
					Err:   err,
				})
			}
		case reflect.Bool:
			if b, err := strconv.ParseBool(envVal); err == nil {
				fieldVal.SetBool(b)
			} else {
				errs = append(errs, &ParseError{
					Field: field.Name,
					Type:  fieldVal.Type().String(),
					Value: envVal,
					Err:   err,
				})
			}
		default:
			errs = append(errs, &UnsupportedTypeError{
				Field: field.Name,
				Type:  fieldVal.Type().String(),
			})
		}
	}
	return errors.Join(errs...)
}
