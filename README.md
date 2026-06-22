# pureconf

[![Go Reference](https://pkg.go.dev/badge/github.com/maeshinshin/pureconf.svg)](https://pkg.go.dev/github.com/maeshinshin/pureconf)
[![Go Report Card](https://goreportcard.com/badge/github.com/maeshinshin/pureconf?style=flat-square)](https://goreportcard.com/report/github.com/maeshinshin/pureconf)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`pureconf` is a modern, generics-first, zero-configuration environment variable mapper for Go. 

Designed specifically for the [12-Factor App](https://12factor.net/config) methodology, it strictly focuses on environment variables, providing a lightweight and secure alternative to multi-format configuration managers.

## Why pureconf? (vs. Viper)

[`spf13/viper`](https://github.com/spf13/viper) is an incredible, battle-tested library and the undisputed king of Go configuration. Its strength lies in its comprehensive versatility—handling JSON, YAML, remote key/value stores, and live-reloading with ease.

However, if your project strictly follows the [12-Factor App](https://12factor.net/config) methodology and relies **exclusively on Environment Variables**, you might not need a multi-format configuration engine. 

`pureconf` is designed to be a focused, generics-native alternative for this specific use case:

| Feature | `spf13/viper` | `pureconf` |
| :--- | :--- | :--- |
| **Philosophy** | Comprehensive & Versatile | **Minimalist & Env-Vars strictly** |
| **Type System** | `map[string]any` + runtime casting | **100% Static typing via Generics (`[T any]`)** |
| **Secret Management** | Standard string handling | **`Secret[T]` auto-masks in logs to prevent leaks** |
| **Mapping Setup** | Requires `mapstructure` tags | **Zero-Config (Auto-infers from struct names)** |
| **Error Handling** | Fails on the first encountered error | **Aggregates ALL errors at once via `errors.Join`** |

By narrowing the focus solely to environment variables and leveraging Go 1.21+, `pureconf` provides a lightweight, type-safe, and highly secure developer experience.

## Features

* **Zero-Configuration:** No `env:"..."` tags required. It automatically maps `AppConfig.DB.Port` to `APP_DB_PORT`.
* **Built-in Secret Masking:** Use the `Secret[T]` type for passwords/tokens. They are automatically masked as `***` when printed or logged, completely preventing accidental credential leaks.
* **Recursive Nested Structs:** Infinitely nest your configuration models natively.
* **Comprehensive Error Aggregation:** Doesn't stop at the first typo. It validates everything and returns an aggregated list of all missing/invalid fields at once.

## Installation

```bash
go get github.com/maeshinshin/pureconf
```
*Requires Go 1.21 or later.*

## Quick Start

Define your configuration purely using Go structs. No tags are needed unless you want to override the default naming convention.

```go
package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/maeshinshin/pureconf"
)

// 1. Define your configuration structure
type DBConfig struct {
	Host string             // Maps to MYAPP_DB_HOST
	Port int                // Maps to MYAPP_DB_PORT
	Pass pureconf.Secret[string] // Maps to MYAPP_DB_PASS (Secured!)
}

type AppConfig struct {
	Debug bool     // Maps to MYAPP_DEBUG
	DB    DBConfig // Automatically creates the "DB_" namespace
}

func main() {
	// 2. Load configuration from environment variables
	// Using WithEnvPrefix("MYAPP_") prevents collision with system env vars.
	cfg, err := pureconf.Load[AppConfig](pureconf.WithEnvPrefix("MYAPP_"))
	if err != nil {
		log.Fatalf("Configuration failed:\n%v", err)
	}

	// 3. Secrets are protected by default!
	slog.Info("Loaded configuration", "db_password", cfg.DB.Pass) 
	// Output: "db_password": "***" (Prevents accidental log leaks)

	// 4. Explicitly unmask when you actually need to use the secret
	connectToDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.Pass.Unmask())
}

func connectToDB(host string, port int, password string) {
	fmt.Printf("Connecting to %s:%d...\n", host, port)
}
```

## Error Aggregation

If multiple environment variables are invalid or missing, `pureconf` tells you exactly what went wrong all at once, saving you from the frustrating "fix one, recompile, find the next error" loop.

```go
cfg, err := pureconf.Load[AppConfig](pureconf.WithEnvPrefix("MYAPP_"))
if err != nil {
    // Returns aggregated errors like:
    // failed to parse 'not-a-number' as int for field Port: invalid syntax
    // failed to parse 'yes' as bool for field Debug: invalid syntax
}
```

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
