# `genx` Usage Guide

## What is `genx`?

`genx` is a code generator for Go. It works by parsing special annotations (comments) within your Go source code files. The primary purpose of `genx` is to help reduce boilerplate code and maintain consistency across your projects by automating the generation of repetitive code patterns.

## Installation

To install `genx`, you can use the `go install` command. Open your terminal and run the following:

```bash
go install github.com/fitan/genx@latest
```

This command will download and install the `genx` executable to your Go bin directory. Make sure your Go bin directory (usually `$GOPATH/bin` or `$HOME/go/bin`) is in your system's `PATH` environment variable so you can run `genx` from any location.

## Running `genx`

After installation, you can run `genx` from the root directory of your Go project. `genx` will process Go files in the current directory and all its subdirectories by default (equivalent to `./...`).

To execute `genx`, simply run the command:

```bash
genx
```

This will trigger the code generation process based on the `genx` annotations found in your Go files and the configuration in `genx.yaml` (if present).

## Configuration (`genx.yaml`)

`genx` looks for a configuration file named `genx.yaml` in the current directory where it is run, or in any parent directory. This file allows you to define global settings that affect the code generation process.

The following top-level keys are supported in `genx.yaml`:

### `imports`

- **Purpose**: To define global import aliases or add imports that might be necessary for generated code but are not explicitly present in the user's source file where `genx` directives are used. This is useful when generated code relies on packages that the original source file doesn't directly import.
- **Syntax**:
  ```yaml
  imports:
    - alias: "alias_name" # Optional: specifies an alias for the import
      path: "module/path"  # Required: the full module path
    - path: "another/module/path" # Example of an import without an alias
  ```

### `preloads`

- **Purpose**: To parse additional packages that are not directly imported by the files being processed but contain types that `genx` plugins might need to understand. This is crucial for situations where plugins need to access type information from packages that are not part of the immediate dependency graph of the annotated files (e.g., base types, types used in configurations passed via comments).
- **Syntax**:
  ```yaml
  preloads:
    - alias: "preload_alias" # Optional: an alias for the preloaded package, can be used in generated code
      path: "module/path/to/preload" # Required: the full module path to preload
    - path: "another/module/to/preload" # Example of a preload without an alias
  ```

### Example `genx.yaml`

Here is a simple example of a `genx.yaml` file:

```yaml
imports:
  - alias: "customsql"
    path: "github.com/myorg/customsql"
  - path: "github.com/google/uuid"

preloads:
  - alias: "basetypes"
    path: "github.com/myorg/project/internal/basetypes"
  - path: "github.com/anotherorg/externalmodels"
```

## TUI Output

When you run `genx`, it provides a Terminal User Interface (TUI) to display the progress and results of the code generation process. This interface helps you track:

- Which packages are being parsed and processed.
- Which `genx` plugins are being executed for specific types or annotations.
- Which files are being generated or modified.
- The status of each operation (e.g., success, failure, or if the generated file already exists and matches the new content).

This real-time feedback is useful for understanding what `genx` is doing and for quickly identifying any issues during generation.

## Annotation Syntax (Magic Comments)

`genx` is triggered by specially formatted comments in Go source code. These comments, often referred to as "magic comments" or annotations, are parsed by `genx` according to a specific grammar to understand what code to generate and how.

There are three main comment styles that `genx` recognizes:

### 1. `@directive` Style (Primary)

This is the most common and flexible way to invoke `genx` plugins and define code generation rules.

-   **General Forms:**
    -   `@FuncName(arg1, "stringArg", key="value", ...)`
    -   `@FuncName()` (no arguments)
    -   `@FuncName` (no arguments, no parentheses)

-   **`@FuncName` (Directive Name):**
    -   This is an `ATID` token, which means it must start with an `@` symbol.
    -   It can be followed by alphanumeric characters, underscores (`_`), or hyphens (`-`).
    -   Examples: `@gormq`, `@crud.Find`, `@my-custom-plugin.generate`.

-   **Arguments:**
    -   Arguments are enclosed in parentheses `()`.
    -   They can be positional or named.
    -   **Positional Arguments:** These are simple string literals (single or double-quoted).
        -   Example: `@myplugin("value1", "value2")`
    -   **Named Arguments:** Use the format `key = "value"` or `key = 'value'`. The `key` is an identifier (alphanumeric, underscores), and the `value` is a string literal.
        -   Example: `@myplugin(name="user", type="admin")`
    -   Positional and named arguments can be mixed, though by convention, positional arguments usually come first.
    -   All directives are line-based and must end with a newline character within the comment block.

-   **Example:**
    ```go
    // @myDirective("positional_arg", count=10, name="example")
    // @anotherDirective()
    // @simpleDirective
    type MyType struct {
        // ...
    }
    ```

### 2. `@FieldDirective` Style (Legacy/Alternative)

This style is also available and might be encountered in older codebases or used for specific use cases, such as field-level annotations.

-   **General Form:**
    -   `@FieldDirective arg1 "arg2 with spaces" arg3`

-   **`@FieldDirective` (Directive Name):**
    -   This also starts with an `ATID` token (e.g., `@ColumnOptions`, `@Validate`).

-   **Arguments:**
    -   Arguments are space-separated tokens that appear after the directive name.
    -   If an argument needs to contain spaces, it must be enclosed in single or double quotes.
    -   These arguments are typically processed as positional arguments by the respective plugins.
    -   The directive is terminated by a newline character within the comment block.

-   **Example:**
    ```go
    type User struct {
        // @ColumnOptions Name "user_name" PrimaryKey
        // @Validate Required MaxLength="50"
        Name string
    }
    ```

### 3. `INSET` Style (Implicit Content Lines)

If a line within a Go doc comment block *does not* start with an `@` symbol, `genx` treats it as an `INSET` directive.

-   The entire content of such a line (after the initial character, which is checked to ensure it's not an `@` directive) is passed as a single string argument to the plugin designated to handle `INSET` directives.
-   This style is particularly useful for passing multi-line text blocks, configurations, or templates directly within comments to a plugin.
-   The specific plugin that handles `INSET` directives determines how this text is interpreted and used.

-   **Example:**
    ```go
    // @MyPluginWithInsetSupport
    // This is the first line of inset content.
    // And this is the second line.
    // It can be json: {"key": "value"}
    // Or any other text the plugin expects.
    func MyFunction() {
        // ...
    }
    ```
    *(Note: The concrete usage of `INSET` depends on the specific plugin. If a plugin like `crud` or `gormq` is found to use `INSET` for a specific purpose, this example could be updated to be more specific.)*

### Important Notes:

-   All annotation directive names and argument keys are **case-sensitive** unless a specific plugin's documentation states otherwise.
-   Ensure that these annotations are placed within standard Go comments (e.g., `//` for line comments, or inside `/* ... */` blocks). `genx` processes the content found within these comments.

## Plugins

### `@log`

*   **Purpose:** Generates a logging middleware for a Go interface. This middleware logs method calls, parameters, execution time, and errors using the `slog` library. It's designed to work with dependency injection using `github.com/samber/do/v2` for logger provisioning.
*   **Target:** Go interfaces.
*   **Directive:** `@log`
*   **Arguments:** None.

**Usage:**

To use the `@log` plugin, annotate your service interface definition with `// @log`.

**Example:**

Suppose you have a service interface defined in `myservice/service.go`:

```go
package myservice

import "context"

// @log
type Service interface {
    CreateUser(ctx context.Context, username string, email string) (id string, err error)
    GetUser(ctx context.Context, id string) (username string, email string, err error)
    NoReturn(ctx context.Context, data string)
}
```

**Generated Code (`myservice/logging.go`):**

`genx` will generate a `logging.go` file in the `myservice` package. This file will contain:

1.  A `logging` struct that wraps your `Service` interface and includes a `*slog.Logger`.
2.  Methods on the `logging` struct corresponding to each method in your `Service` interface. These methods:
    *   Log input parameters (complex types are JSON encoded).
    *   Call the actual method on the wrapped `next` service.
    *   Log the execution duration (`took`) and any returned error. `slog.LevelInfo` is used for successful calls, and `slog.LevelError` if an error is returned.
3.  A `NewLogging` constructor function. This function is designed to be a `do.Middleware` and expects a `do.Injector` to resolve the `*slog.Logger` and the `next` service instance.

An illustrative snippet of what the generated `logging.go` might look like (simplified):

```go
package myservice

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/samber/do/v2"
)

type logging struct {
	logger *slog.Logger
	next   Service
}

func (s *logging) CreateUser(ctx context.Context, username string, email string) (id string, err error) {
	// JSON marshal non-basic params (example)
	usernameJSON, _ := json.Marshal(username) // Simplified for illustration
	emailJSON, _ := json.Marshal(email)       // Simplified for illustration

	defer func(begin time.Time) {
		level := slog.LevelInfo
		if err != nil {
			level = slog.LevelError
		}
		s.logger.Log(ctx, level, "",
			"method", "CreateUser",
			"username", string(usernameJSON),
			"email", string(emailJSON),
			"took", time.Since(begin).String(),
			"err", err,
		)
	}(time.Now())

	return s.next.CreateUser(ctx, username, email)
}

// ... other wrapped methods (GetUser, NoReturn) ...

func NewLogging(i do.Injector) func(next Service) Service {
	return func(next Service) Service {
		return &logging{
			logger: do.MustInvoke[*slog.Logger](i), // Logger obtained via DI
			next:   next,
		}
	}
}
```

**Setup (using `samber/do/v2`):**

You would typically register this logging middleware when setting up your services with `do`:

```go
package main

import (
	"log/slog"
	"os"
	"context" // Added missing import for context.Background()

	"github.com/fitan/genx/plugs/log/example/myservice" // Assuming this is your service package
	"github.com/samber/do/v2"
)

func main() {
	injector := do.New()

	// Provide slog.Logger
	do.Provide(injector, func(i do.Injector) (*slog.Logger, error) {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil)), nil
	})

	// Provide your actual service implementation
	do.Provide(injector, func(i do.Injector) (myservice.Service, error) {
		return &myServiceImpl{}, nil // myServiceImpl is your concrete implementation
	})

	// Provide the logging middleware
	// The generated NewLogging function acts as a middleware provider
	do.ProvideNamedValue(injector, "LoggingMiddleware", myservice.NewLogging(injector))


    // To get an instance of the service with logging:
    // 1. Invoke the middleware
    // loggingMiddleware := do.MustInvokeNamed[func(myservice.Service) myservice.Service](injector, "LoggingMiddleware")
    // 2. Invoke the original service
    // originalService := do.MustInvoke[myservice.Service](injector)
    // 3. Apply the middleware
    // loggedService := loggingMiddleware(originalService)

    // More typically, if using a framework that supports layered services:
    // You might register the middleware to be applied automatically.
    // Or, if you resolve services by type and it supports decorators:
    do.Decorate(injector, myservice.NewLogging(injector))
    
    // Now when you invoke myservice.Service, it will be wrapped with logging
    finalService := do.MustInvoke[myservice.Service](injector)
    finalService.CreateUser(context.Background(), "test", "test@example.com")


	// ... your application logic ...
}

type myServiceImpl struct{} // Your actual service implementation
// ... implement myservice.Service for myServiceImpl ...
func (s *myServiceImpl) CreateUser(ctx context.Context, username string, email string) (id string, err error) {
    // ... actual logic ...
    return "123", nil
}
func (s *myServiceImpl) GetUser(ctx context.Context, id string) (username string, email string, err error) {
    // ... actual logic ...
    return "test", "test@example.com", nil
}
func (s *myServiceImpl) NoReturn(ctx context.Context, data string) {
    // ... actual logic ...
}

```
**Note:** The exact way to apply the middleware in `samber/do/v2` can vary (e.g. using `do.Decorate` or by explicitly providing and invoking middleware chains). The generated `NewLogging` is flexible enough for these approaches. The `do.Decorate` approach is often the cleanest.

### `@trace`

*   **Purpose:** Generates a tracing middleware for a Go interface using OpenTelemetry. This middleware creates a new span for each method call, records method parameters as span attributes, captures errors, and sets the span status accordingly. It's designed for use with `github.com/samber/do/v2` for dependency injection of the `TracerProvider`.
*   **Target:** Go interfaces.
*   **Directive:** `@trace`
*   **Arguments:** None.

**Usage:**

Annotate your service interface definition with `// @trace`.

**Example:**

Given a service interface in `myservice/service.go`:

```go
package myservice

import "context"

// @trace
type Service interface {
    CreateItem(ctx context.Context, itemID string, itemName string) (err error)
    GetItem(ctx context.Context, itemID string) (itemName string, err error)
}
```

**Generated Code (`myservice/trace.go`):**

`genx` will generate a `trace.go` file in the `myservice` package. This file will include:

1.  A `tracing` struct that wraps your `Service` interface and holds an OpenTelemetry `*sdktrace.TracerProvider`.
2.  Methods on the `tracing` struct for each method in your `Service`. These methods:
    *   Start an OpenTelemetry span (e.g., named after the method like "CreateItem").
    *   Set method parameters as span attributes (parameters are JSON marshaled into a `params` attribute).
    *   If the wrapped method returns an error, it's recorded on the span (`span.RecordError()`) and the span status is set to `codes.Error`.
    *   The span is ended upon method completion.
3.  A `NewTracing` constructor function. This function is a `do.Middleware` and expects a `do.Injector` to resolve the `*sdktrace.TracerProvider` and the `next` service instance.

Illustrative snippet of `trace.go`:

```go
package myservice

import (
	"context"
	"encoding/json"

	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace" // Assuming this import
)

type tracing struct {
	next   Service
	tracer *sdktrace.TracerProvider
}

func (s *tracing) CreateItem(ctx context.Context, itemID string, itemName string) (err error) {
	ctx, span := s.tracer.Tracer("CreateItem").Start(ctx, "CreateItem") // Simplified tracer name for example
	defer func() {
		params := map[string]interface{}{
			"itemID":   itemID,
			"itemName": itemName,
		}
		paramsB, _ := json.Marshal(params)
		span.SetAttributes(attribute.String("params", string(paramsB)))

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()
	return s.next.CreateItem(ctx, itemID, itemName)
}

// ... other wrapped methods (GetItem) ...

func NewTracing(i do.Injector) func(next Service) Service {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: do.MustInvoke[*sdktrace.TracerProvider](i),
		}
	}
}
```

**Setup (using `samber/do/v2` and OpenTelemetry SDK):**

You would register this tracing middleware similarly to the logging middleware, ensuring an OpenTelemetry `TracerProvider` is also provided.

```go
package main

import (
	"context"
	"log"

	"github.com/fitan/genx/plugs/trace/example/myservice" // Your service package
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Setup OTel Tracer Provider
func newTracerProvider() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Adjust sampler in production
	)
	otel.SetTracerProvider(tp) // Global tracer provider
	// otel.SetTextMapPropagator(...) // Configure propagator if needed
	return tp, nil
}

func main() {
	injector := do.New()

	// Provide OpenTelemetry TracerProvider
	do.Provide(injector, func(i do.Injector) (*sdktrace.TracerProvider, error) {
		return newTracerProvider()
	})

	// Provide your actual service implementation
	do.Provide(injector, func(i do.Injector) (myservice.Service, error) {
		return &myServiceImpl{}, nil
	})

	// Decorate the service with the tracing middleware
	do.Decorate(injector, myservice.NewTracing(injector))
    
    // Example usage
	finalService := do.MustInvoke[myservice.Service](injector)
	_, err := finalService.GetItem(context.Background(), "item123")
    if err != nil {
        log.Printf("Error calling GetItem: %v", err)
    }

    // Shutdown TracerProvider when application exits
    tp := do.MustInvoke[*sdktrace.TracerProvider](injector)
    if err := tp.Shutdown(context.Background()); err != nil {
        log.Printf("Error shutting down tracer provider: %v", err)
    }
}

type myServiceImpl struct{} // Your actual service implementation
func (s *myServiceImpl) CreateItem(ctx context.Context, itemID string, itemName string) (err error) {
	// ... actual logic ...
	return nil
}
func (s *myServiceImpl) GetItem(ctx context.Context, itemID string) (itemName string, err error) {
	// ... actual logic ...
	return "Test Item", nil
}
```
**Note:** Remember to configure your OpenTelemetry SDK appropriately (e.g., exporters, samplers, propagators) for your environment. The example uses a simple stdout exporter.

### `@otel`

*   **Purpose:** Generates a comprehensive OpenTelemetry middleware for a Go interface. This middleware combines tracing, logging (using `slog`), and metrics (status counter, total counter, duration histogram) for each method call. It relies on `github.com/samber/do/v2` for injecting necessary OpenTelemetry and logging components.
*   **Target:** Go interfaces.
*   **Directive:** `@otel`
*   **Arguments:** None.

**Usage:**

Annotate your service interface definition with `// @otel`.

**Example:**

Consider an interface in `custompkg/service.go`:

```go
package custompkg

import "context"

// @otel
type StringService interface {
    Uppercase(ctx context.Context, s string) (upper string, err error)
    Count(ctx context.Context, s string) (length int, err error)
}
```

**Generated Code (`custompkg/otel.go`):**

`genx` will create an `otel.go` file in the `custompkg` package. This file contains:

1.  An `otel` struct that wraps your `StringService` and holds instances of `trace.Tracer`, `*slog.Logger`, and various `metric` instruments. It also stores the `pkgName`.
2.  Methods on the `otel` struct for each method in your interface. These methods:
    *   **Tracing:** Start an OpenTelemetry span (named `pkgName + "." + methodName`). Parameters are added as JSON encoded attributes. Errors are recorded, and span status is updated.
    *   **Logging:** Log method invocation details (pkg, method, params, duration, error) using `slog`.
    *   **Metrics:**
        *   `serviceStatus` (Int64Counter): Incremented with dimensions `service` (pkgName), `method`, and `status` ("success" or "fail").
        *   `serviceCounter` (Int64Counter): Incremented with dimensions `service` (pkgName) and `method`.
        *   `serviceDuration` (Float64Histogram): Records method execution time in milliseconds with dimensions `service` (pkgName) and `method`.
3.  A `NewOtel` constructor function. This is a `do.Middleware` expecting a `do.Injector` to resolve the tracer, logger, and the named metric instruments (`serviceStatus`, `serviceCounter`, `serviceDuration`).

Illustrative snippet of the generated `otel.go`:

```go
package custompkg

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type otel struct {
	pkgName              string
	next                 StringService // Your interface name
	tracer               trace.Tracer
	logger               *slog.Logger
	meterServiceStatus   metric.Int64Counter
	meterServiceCounter  metric.Int64Counter
	meterServiceDuration metric.Float64Histogram
}

func (s *otel) Uppercase(ctx context.Context, S string) (upper string, err error) { // Renamed 's' to 'S' to match user example
	_method := "Uppercase"
	ctx, span := s.tracer.Start(ctx, s.pkgName+"."+_method)

	defer func(begin time.Time) {
		_endTime := time.Since(begin)
		_level := slog.LevelInfo
		// _statusAttr := "success" // Variable not used, removed

		if err != nil {
			_level = slog.LevelError
			// _statusAttr = "fail" // Variable not used, removed
			s.meterServiceStatus.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("service", s.pkgName),
					attribute.String("method", _method),
					attribute.String("status", "fail"),
				))
		} else {
			s.meterServiceStatus.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("service", s.pkgName),
					attribute.String("method", _method),
					attribute.String("status", "success"),
				))
		}

		s.meterServiceDuration.Record(ctx, float64(_endTime.Milliseconds()),
			metric.WithAttributes(
				attribute.String("service", s.pkgName),
				attribute.String("method", _method),
			))
		s.meterServiceCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("service", s.pkgName),
				attribute.String("method", _method),
			))

		_params := map[string]interface{}{"s": S} // Param name 's' from interface
		_paramsB, _ := json.Marshal(_params)

		s.logger.Log(ctx, _level, "",
			"pkg", s.pkgName, "method", _method, "params", string(_paramsB),
			"took", _endTime.String(), "err", err,
		)

		span.SetAttributes(attribute.String("params", string(_paramsB)), attribute.String("took", _endTime.String()))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}(time.Now())

	return s.next.Uppercase(ctx, S)
}

// ... other wrapped methods (Count) ...

func NewOtel(i do.Injector) func(next StringService) StringService { // Your interface name
	return func(next StringService) StringService { // Your interface name
		return &otel{
			pkgName:              "custompkg", // This is derived by genx from the package
			next:                 next,
			tracer:               do.MustInvoke[trace.Tracer](i),
			logger:               do.MustInvoke[*slog.Logger](i),
			meterServiceStatus:   do.MustInvokeNamed[metric.Int64Counter](i, "serviceStatus"),
			meterServiceCounter:  do.MustInvokeNamed[metric.Int64Counter](i, "serviceCounter"),
			meterServiceDuration: do.MustInvokeNamed[metric.Float64Histogram](i, "serviceDuration"),
		}
	}
}
```

**Setup (using `samber/do/v2`, OpenTelemetry SDK):**

You'll need to provide all the dependencies: `trace.Tracer`, `*slog.Logger`, and the named metrics.

```go
package main

import (
	"context"
	"log"
	"os"
    "strings" // Added for S_ToUpper

	"github.com/fitan/genx/plugs/otel/example/custompkg" // Your service package
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/global" 
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace" 
	"log/slog"
)

// Setup OTel Tracer Provider
func newTracerProvider() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil { return nil, err }
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	return tp, nil
}

// Setup OTel Meter Provider
func newMeterProvider() (*sdkmetric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil { return nil, err }
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)))
	global.SetMeterProvider(mp) 
	return mp, nil
}

func main() {
	injector := do.New()

	// Provide slog.Logger
	do.Provide(injector, func(i do.Injector) (*slog.Logger, error) {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil)), nil
	})

	// Provide OTel TracerProvider and Tracer
	tp, err := newTracerProvider()
	if err != nil { log.Fatal(err) }
	do.ProvideValue(injector, tp) 
	do.Provide(injector, func(i do.Injector) (trace.Tracer, error) { 
		return tp.Tracer("my-app-tracer"), nil 
	})


	// Provide OTel MeterProvider and Metrics
	mp, err := newMeterProvider()
	if err != nil { log.Fatal(err) }
	do.ProvideValue(injector, mp)
	meter := mp.Meter("my-app-meter")

	serviceStatusCounter, _ := meter.Int64Counter("serviceStatus") 
	serviceCounter, _ := meter.Int64Counter("serviceCounter")      
	serviceDurationHistogram, _ := meter.Float64Histogram("serviceDuration")

	do.ProvideNamedValue(injector, "serviceStatus", serviceStatusCounter)
	do.ProvideNamedValue(injector, "serviceCounter", serviceCounter) 
	do.ProvideNamedValue(injector, "serviceDuration", serviceDurationHistogram)


	// Provide your actual service implementation
	do.Provide(injector, func(i do.Injector) (custompkg.StringService, error) {
		return &myStringServiceImpl{}, nil
	})

	// Decorate with Otel middleware
	do.Decorate(injector, custompkg.NewOtel(injector))

	finalService := do.MustInvoke[custompkg.StringService](injector)
	finalService.Uppercase(context.Background(), "hello")

	// Shutdown providers
	if err := tp.Shutdown(context.Background()); err != nil { log.Printf("Error shutting down tracer provider: %v", err) }
	if err := mp.Shutdown(context.Background()); err != nil { log.Printf("Error shutting down meter provider: %v", err) }
}

type myStringServiceImpl struct{}
func (s *myStringServiceImpl) Uppercase(ctx context.Context, S string) (upper string, err error) { return S_ToUpper(S), nil } 
func (s *myStringServiceImpl) Count(ctx context.Context, S string) (length int, err error) { return len(S), nil } 

// Helper to avoid conflict with param name 's'
func S_ToUpper(s string) string {
    return strings.ToUpper(s) // Corrected to use strings.ToUpper
}
```
**Important Considerations:**
*   Ensure the names of the provided metrics (`serviceStatus`, `serviceCounter`, `serviceDuration`) in your `do.ProvideNamedValue` calls exactly match the names expected by the generated `NewOtel` function. The generated code uses `"serviceStatus"`, `"serviceCounter"`, and `"serviceDuration"`.
*   The OpenTelemetry SDK setup can be complex. The example provides a basic stdout setup. Refer to OpenTelemetry documentation for production-ready configurations.
*   The `pkgName` used in metrics and logging is automatically derived from the package of the instrumented interface.

### `@enum`

*   **Purpose:** Generates a set of helper methods and constants for an integer-based enum type, making it more robust and easier to use. It creates string representations, remarks, and parsing functions.
*   **Target:** Go type specifications (typically an `int` type aliased for the enum).
*   **Directive:** `@enum`
*   **Arguments:**
    *   The `@enum` directive takes one or more positional string arguments.
    *   Each argument defines an enum member and should be in the format `"Key:Remark"` or simply `"Key"`.
        *   `Key`: The identifier for the enum member (e.g., `Active`, `DefaultUser`). This is used to generate constant names like `Status_Active`.
        *   `Remark`: A descriptive string for the enum member (e.g., `User is currently active`). This is accessible via the `Remark()` method. If only `"Key"` is provided, the remark is empty.

**Usage:**

Define an integer type for your enum and annotate its type specification with `// @enum(...)` and the member definitions.

**Example:**

Suppose you have `types/status.go`:

```go
package types

// @enum("Active:User is active", "Pending:Awaiting confirmation", "Disabled:User account is disabled", "Deleted")
type Status int // The base type is typically int
```

**Generated Code (`types/enum.go`):**

`genx` will generate an `enum.go` file in the `types` package (or wherever `Status` is defined). This file will contain:

1.  **Enum Value Constants:** An iota-based `const` block for your enum members.
    ```go
    const (
        _ = iota // Start from 0, actual values will be 1, 2, 3...
        Status_Active
        Status_Pending
        Status_Disabled
        Status_Deleted
    )
    ```
2.  **Alias and Remark Constants:** Constants for the string "key" and "remark" of each enum member.
    ```go
    const (
        STATUS_ACTIVE_ALIAS   = "Active"
        STATUS_ACTIVE_REMARK  = "User is active"
        STATUS_PENDING_ALIAS  = "Pending"
        STATUS_PENDING_REMARK = "Awaiting confirmation"
        // ... and so on for Disabled, Deleted (remark for Deleted will be empty)
    )
    ```
3.  **`String()` Method:** Returns the string key of the enum value.
    ```go
    func (e Status) String() string {
        switch e {
        case Status_Active:
            return STATUS_ACTIVE_ALIAS
        case Status_Pending:
            return STATUS_PENDING_ALIAS
        // ...
        default:
            return fmt.Sprintf("unknown %d", e)
        }
    }
    ```
4.  **`Remark()` Method:** Returns the descriptive remark of the enum value.
    ```go
    func (e Status) Remark() string {
        switch e {
        case Status_Active:
            return STATUS_ACTIVE_REMARK
        case Status_Pending:
            return STATUS_PENDING_REMARK
        // ...
        default:
            return fmt.Sprintf("unknown %d", e)
        }
    }
    ```
5.  **`ParseStatus()` Function:** Converts an integer to your enum type (`Status`).
    ```go
    func ParseStatus(id int) (Status, error) {
        if x, ok := _StatusValue[id]; ok { // _StatusValue is a generated map
            return x, nil
        }
        return 0, fmt.Errorf("unknown enum value: %d", id) // Corrected error formatting
    }
    ```
    (Also generates `var _StatusValue = map[int]Status{...}` mapping int values to enum constants.)


**Using the Generated Enum:**

```go
package main

import (
	"fmt"
	"your_project_module/types" // Adjust import path
)

func main() {
	activeStatus := types.Status_Active
	fmt.Println(activeStatus.String()) // Output: Active
	fmt.Println(activeStatus.Remark()) // Output: User is active
	fmt.Println(int(activeStatus))     // Output: 1 (or its iota value)

	parsedStatus, err := types.ParseStatus(2)
	if err == nil {
		fmt.Println(parsedStatus.String()) // Output: Pending
	}

    // Accessing alias/remark constants directly
    fmt.Println(types.STATUS_DISABLED_ALIAS)  // Output: Disabled
    fmt.Println(types.STATUS_DISABLED_REMARK) // Output: User account is disabled
}
```
**Note:** The generated code uses title casing for enum constant names (e.g., `Status_Active`) and uppercase for alias/remark constants (e.g., `STATUS_ACTIVE_ALIAS`). The integer values start from 1 (due to `_ = iota` followed by members).

### `@gq` (GORM Query Scopes)

*   **Purpose:** Generates dynamic GORM query scope functions based on the fields of a user-defined query struct. This allows for building reusable and type-safe query logic.
*   **Target:** Go structs.
*   **Generated File:** `gorm_scope.go` in the same package as the annotated struct.

**Primary Directive (on the query struct itself):**

*   `// @gq <model.GormModel>`
    *   This directive marks a struct as a GORM query definition.
    *   `<model.GormModel>`: **Required.** The fully qualified name of the GORM model struct that the query will target (e.g., `myproject/model.User`, `model.Product`).

**Field-Level Directives (on fields within the query struct):**

These directives customize how each field in the query struct translates to a GORM query condition.

*   `// @gq-column <db_col_1> [<db_col_2> ...]`
    *   Specifies the database column name(s) for the field.
    *   If multiple column names are provided, an `OR` condition is generated for that field using its value against each specified column.
    *   If omitted, `genx` attempts to read the `gorm:"column:..."` tag. If that's also missing, GORM's default naming strategy for the field name is assumed.
    *   **Example:** `// @gq-column name_alias user_name` generates `WHERE name_alias = ? OR user_name = ?`.

*   `// @gq-op <operator>`
    *   Defines the comparison operator. Defaults to `=` if not specified.
    *   **Supported Operators:**
        *   `=` : Equal. (Default for most types; for slices/arrays, defaults to `in`).
        *   `!=`: Not equal.
        *   `>` : Greater than.
        *   `>=`: Greater than or equal to.
        *   `<` : Less than.
        *   `<=`: Less than or equal to.
        *   `><`: BETWEEN (expects the field to be a slice or array of two values, e.g., `PriceRange []int`). Generates `column BETWEEN ? AND ?`.
        *   `!><`: NOT BETWEEN.
        *   `like`: LIKE operator (e.g., `column LIKE ?`). The generated code automatically wraps the value with `%` (e.g., `"%value%"`).
        *   `in`: IN operator (e.g., `column IN (?)`). Default for slice/array fields if `@gq-op` is `=` or omitted.
        *   `!in`: NOT IN operator.
        *   `null`: IS NULL (the field's value in the query struct is ignored). Generates `column IS NULL`.
        *   `!null`: IS NOT NULL (field's value is ignored). Generates `column IS NOT NULL`.

*   `// @gq-clause <clause_type>`
    *   Specifies the GORM clause type. Defaults to `Where`.
    *   Valid values: `Where`, `Or`, `Not`. (Note: `Or` and `Not` might have specific interaction patterns with field values or multiple `@gq-column` entries. `Where` is the most common).

*   `// @gq-sub <foreign_key> <referenced_key>`
    *   Used for creating subqueries. The field this annotates should be a pointer to another struct that also has a generated `Scope` method (i.e., also annotated with `@gq`).
    *   `<foreign_key>`: The column name in the current model (defined by `@gq <model.GormModel>`) that will be compared against the subquery results.
    *   `<referenced_key>`: The column name selected in the subquery from the subquery's model.
    *   **Example:** If `OrderQuery` has a field `User *UserQuery // @gq-sub user_id id`, it generates `WHERE user_id IN (SELECT id FROM users WHERE ...conditions_from_UserQuery...)`.

*   `// @gq-group`
    *   Marks a field (which must be a struct type, usually embedded or a direct field) whose conditions should be grouped. The field's struct type should itself be a valid query struct (potentially with its own `@gq-*` annotations but not the main `@gq <model>` one).
    *   Generates a grouped condition like `db.Where(subScope)` where `subScope` is the result of the nested struct's generated query logic.

*   `// @gq-struct [<value1> <value2> ...]`
    *   Passes the field (which should be a struct or a pointer to a struct) directly to a GORM query method, typically `Where`.
    *   If `<value1>`, `<value2>` are provided, they are passed as additional arguments to the GORM method, often used to specify which fields of the struct to query against if it's not a GORM model itself. Example: `db.Where(&q.MyStructField, "field1", "field2")`.
    *   If no values are provided, it's typically `db.Where(&q.MyStructField)`.

**Generated Code (`gorm_scope.go`):**

*   A method `func (q *QueryStructName) Scope(db *gorm.DB) *gorm.DB` is generated.
*   This `Scope` method:
    *   Sets the GORM model: `db = db.Model(&<model.GormModel>{})`.
    *   Iterates through the fields of `QueryStructName`.
    *   Builds GORM query conditions based on field values and their `@gq-*` annotations.
    *   **Zero Value Handling:** Conditions for fields with zero values are typically skipped:
        *   `nil` pointers.
        *   Empty strings (`""`).
        *   `0` for integer types.
        *   Empty slices (`len() == 0`).
    *   Supports embedded structs: Fields from embedded structs are processed as if they were part of the parent query struct. Pointer embedded structs are also handled.
    *   **Special Slice Handling:** If a field is a slice of structs (e.g., `RelatedItems []ItemQuery`), and `ItemQuery` has fields that map to columns, this can generate queries like `WHERE (item_col_a, item_col_b) IN ((val1_a, val1_b), (val2_a, val2_b), ... )`.

**Example:**

`queries/user_query.go`:
```go
package queries

import "myproject/model" // Assuming model.User is your GORM model
import "time" // Added import for time.Time

// @gq model.User
type UserQuery struct {
    // @gq-column id
    // @gq-op =
    ID *uint `json:"id"` // Pointer to allow nil (and skip query if nil)

    // @gq-op like
    // @gq-column name
    NameSearch *string `json:"nameSearch"`

    // @gq-op in
    // @gq-column email
    Emails []string `json:"emails"`

    // @gq-column created_at
    // @gq-op ><
    CreatedAtRange []*time.Time `json:"createdAtRange"` // Expects two time.Time pointers

    // @gq-column status
    // @gq-op =
    Status *int `json:"status"`

    // @gq-sub user_id id
    Profile *ProfileQuery `json:"profile"` // ProfileQuery is another @gq struct for model.Profile
}

// @gq model.Profile
type ProfileQuery struct {
    // @gq-column bio
    // @gq-op like
    BioSearch *string `json:"bioSearch"`
}

```

Generated `queries/gorm_scope.go` (simplified concept):
```go
package queries

import (
	"gorm.io/gorm"
	"myproject/model"
	"time" // Assuming time is used
)

func (q *UserQuery) Scope(db *gorm.DB) *gorm.DB {
	db = db.Model(&model.User{})

	if q.ID != nil {
		db = db.Where("id = ?", *q.ID)
	}
	if q.NameSearch != nil {
		db = db.Where("name LIKE ?", "%"+*q.NameSearch+"%")
	}
	if len(q.Emails) > 0 {
		db = db.Where("email IN (?)", q.Emails)
	}
	if q.CreatedAtRange != nil && len(q.CreatedAtRange) == 2 && q.CreatedAtRange[0] != nil && q.CreatedAtRange[1] != nil {
		db = db.Where("created_at BETWEEN ? AND ?", *q.CreatedAtRange[0], *q.CreatedAtRange[1])
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}
	if q.Profile != nil {
		// For subqueries, a new GORM session is typically used to avoid interference.
		// The actual generated code might use db.Session(&gorm.Session{NewDB: true}) or similar.
		profileSubQuery := q.Profile.Scope(db.Session(&gorm.Session{NewDB: true})).Select("id")
		db = db.Where("user_id IN (?)", profileSubQuery) // 'user_id' from @gq-sub foreign_key, 'id' from referenced_key
	}
	return db
}

func (q *ProfileQuery) Scope(db *gorm.DB) *gorm.DB {
    db = db.Model(&model.Profile{}) // Sets model for ProfileQuery
    if q.BioSearch != nil {
        db = db.Where("bio LIKE ?", "%"+*q.BioSearch+"%")
    }
    return db
}
```

**Using the Scope:**

```go
package main // Or your relevant package

import (
    "gorm.io/gorm"
    // "your_project/model" // Assuming your GORM models are here
    // "your_project/queries" // Assuming your query structs are here
    // "time" // If using time
)

// Helper functions for example clarity (not part of genx output)
func stringPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }

func main() {
    // var db *gorm.DB // Your initialized GORM DB instance
    // Assume db is initialized appropriately

    // query := queries.UserQuery{
    //     NameSearch: stringPtr("admin"),
    //     Status:     intPtr(1),
    // }

    // var users []model.User
    // resultDB := query.Scope(db).Find(&users)
    // if resultDB.Error != nil {
    //     // handle error
    // }
    // users now contains filtered users
    // fmt.Println("Filtered users:", users)
}
```
**Note:** The `design.md` file in `plugs/gormq/` provides additional examples and context for this plugin. It's highly recommended to review it for a deeper understanding.

### `@crud`

*   **Purpose:** Generates boilerplate code for CRUD (Create, Read, Update, Delete) operations. It supports two distinct modes: generating GORM-based database services or generating HTTP service structures and types.
*   **Target:** Go structs (these structs act as configuration/markers for the CRUD generation).
*   **Directive:** `@crud`

**Arguments for the `@crud` directive (applied to the configuration struct):**

*   `type="<mode>"`: **Required.** Specifies the generation mode.
    *   `"gorm"`: Generates a GORM-based CRUD service implementation.
    *   `"http"`: Generates HTTP request/response types and a base HTTP service structure.
*   `model="<pkg.ModelName>"`: **Required.** The fully qualified name of the GORM model struct that the CRUD operations will target (e.g., `myproject/model.User`).
*   `idName="<IDFieldName>"`: **Required.** The name of the primary key field within the GORM model (e.g., `ID`, `UUID`). This refers to the Go struct field name.
*   `idType="<GoType>"`: **Required.** The Go type of the primary key field (e.g., `uint`, `string`, `uuid.UUID`).
*   `preload="<field1,field2.subfield>"`: *(Used only by `type="http"`)* A comma-separated string specifying which fields of the GORM model should be deeply copied when generating the `GetResponse` struct for HTTP responses. Supports dot notation for nested fields (e.g., `User,Order.OrderItems`).

**Usage:**

Define an empty struct and annotate it with `@crud` and its required arguments. The name of this struct itself doesn't directly influence the generated code's naming as much as the `model` argument does.

**Example Configuration Struct:**

```go
package services

// @crud(type="gorm", model="myproject/model.Article", idName="ID", idType="uint")
type ArticleGormCRUDConfig struct{}

// @crud(type="http", model="myproject/model.User", idName="ID", idType="int", preload="Profile,Addresses")
type UserHttpCRUDConfig struct{}
```

---

**Mode: `type="gorm"`**

*   **Purpose:** Generates a base service layer for GORM-based CRUD operations.
*   **Generated Files:**
    *   `crud_base_service.go`: Contains the GORM CRUD service implementation (e.g., methods for Create, GetByID, Update, Delete, List).
    *   (Note: `crud_gorm_types.go` generation seems to be commented out in the plugin's source, so it's likely not produced by default.)

*   **Generated Code (`crud_base_service.go` - conceptual):**
    *   An interface (e.g., `ArticleGormCrudBaseImpl` if `model` was `model.Article`).
    *   A struct implementing this interface, which takes a `*gorm.DB` instance.
    *   **Standard CRUD Methods:**
        *   `Create(m *model.Article) error`
        *   `GetByID(id uint) (*model.Article, error)`
        *   `Update(id uint, m *model.Article) error`
        *   `Delete(id uint) error`
        *   `FindOne(query interface{}, args ...interface{}) (*model.Article, error)`
        *   `Find(query interface{}, args ...interface{}) ([]*model.Article, int64, error)` (returns items and total count)
    *   These methods use GORM for database interactions (e.g., `db.Create()`, `db.First()`, `db.Save()`, `db.Delete()`, `db.Where().Find()`).
    *   The exact method signatures and capabilities are defined by the `crud_gorm.tmpl` template.

**Using the GORM CRUD Service:**

```go
package main

import (
    "log"
    "myproject/model" // Your GORM models
    "myproject/services" // Where crud_base_service.go is generated
    "gorm.io/gorm"
    // ... other imports ...
)

func main() {
    var db *gorm.DB // Your initialized GORM DB instance

    // Assuming ArticleGormCrudBaseImpl struct and NewArticleGormCrudBaseImpl constructor are generated
    // The actual names will depend on the 'model' parameter given to @crud
    // For model="myproject/model.Article", it might be NewArticleGormCrudBaseImpl
    // Let's assume a generic constructor for illustration if the template provides one,
    // or direct struct instantiation.
    // articleService := services.NewArticleGormCrudBaseImpl(db) // Or direct instantiation if no constructor

    // Example: If generated struct is ArticleService (derived from model.Article)
    // articleService := &services.ArticleService{DB: db}


    // The exact name of the generated service struct and its constructor
    // needs to be inferred from the template `crud_gorm.tmpl`.
    // For now, let's assume a conceptual `NewGeneratedCrudService` for `model.Article`.
    // crudService := services.NewArticleCrudService(db) // This name is an assumption

    // To use it, you'd need to know the actual generated struct and constructor name.
    // The documentation should ideally list the exact names from the template.
    // For now, we'll focus on the concept.
}
```
*(Self-correction during subtask: The "Using" part for GORM is hard to detail without seeing the exact constructor/struct names from `crud_gorm.tmpl`. The documentation should state this and perhaps show how to instantiate the generated service once those names are known from inspecting generated files or the template.)*

---

**Mode: `type="http"`**

*   **Purpose:** Generates types for HTTP request/response payloads and a base structure for an HTTP service.
*   **Generated Files:**
    *   `types.go`: Contains generated Go structs for HTTP request and response bodies (e.g., `GetResponse`, `CreateRequest`, `UpdateBody`).
    *   `crud_http_service.go`: Contains a base HTTP service interface and potentially a stub implementation.

*   **Generated Code (`types.go` - conceptual):**
    *   `GetResponse` struct: Created based on the GORM `model` specified in the `@crud` annotation. The `preload` argument controls which fields (including nested ones via dot notation like `User.Profile`) are included. This is useful for shaping API responses.
    *   `CreateRequest` struct: For unmarshaling create request payloads.
    *   `UpdateBody` struct: For unmarshaling update request payloads.
    *   These structs are generated by copying fields from the GORM model, potentially excluding certain fields (e.g., those tagged `serializer:"-"`) and including only those specified in `preload` for `GetResponse`.

*   **Generated Code (`crud_http_service.go` - conceptual):**
    *   An interface (e.g., `UserHttpCrudBaseImpl` if `model` was `model.User`).
    *   **Potential HTTP Service Methods:**
        *   `Create(c context.Context, req *CreateRequest) (*GetResponse, error)`
        *   `Get(c context.Context, id int) (*GetResponse, error)`
        *   `Update(c context.Context, id int, req *UpdateBody) (*GetResponse, error)`
        *   `Delete(c context.Context, id int) error`
        *   `List(c context.Context, queryParams ListRequest) ([]*GetResponse, int, error)` (ListRequest would also be a generated type for pagination/filtering parameters)
    *   The actual method signatures are defined by the `crud_http.tmpl` and `crud_http_types.tmpl` templates.

**Using the HTTP CRUD Service/Types:**

```go
package main

import (
    // "myproject/services" // Where crud_http_service.go and types.go are generated
    // "myproject/model"
    // "github.com/gin-gonic/gin" // Or your preferred HTTP framework
    // ...
)

// func setupRouter(userService *services.UserHttpService) *gin.Engine { // Assuming UserHttpService
//     r := gin.Default()
//     // Example endpoint using generated types
//     r.POST("/users", func(c *gin.Context) {
//         var req services.CreateRequest // Generated type
//         if err := c.ShouldBindJSON(&req); err != nil {
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }
//         // resp, err := userService.Create(c.Request.Context(), &req)
//         // ... handle response ...
//     })
//     return r
// }
```
*(Self-correction during subtask: Similar to GORM, the "Using" part for HTTP is conceptual without the exact generated interface/struct names and method signatures from the templates. The documentation should highlight that users need to inspect the generated files or templates for these specifics.)*

**Note:** The `@crud` plugin is template-driven. For precise details on generated method names, struct fields, and constructor functions, refer to the content of `crud_gorm.tmpl`, `crud_http.tmpl`, and `crud_http_types.tmpl` within the `genx` `static/template` directory.

### `@kit-http-client`

*   **Purpose:** Generates a Go kit HTTP client implementation based on a Go interface definition. This client includes features like service discovery, load balancing (round-robin), and retries.
*   **Target:** Go interfaces.
*   **Generated File:** `kit_http_client.go` in the same package as the interface.

**Interface-Level Directives:**

*   `// @basePath <base_path_string>`
    *   **Optional.** Placed on the interface definition.
    *   Specifies a common base path prefix for all HTTP URLs defined in the methods of this interface (e.g., `/api/v2`).

**Method-Level Directives (on methods within the annotated interface):**

*   `// @kit-http <http_url_path> <http_method>`
    *   **Required.** Defines the endpoint for the method.
    *   `<http_url_path>`: The specific URL path (e.g., `/users/{userID}`, `/products`). Path parameters are denoted by `{param_name}`.
    *   `<http_method>`: The HTTP method (e.g., `GET`, `POST`, `PUT`, `DELETE`).

*   `// @kit-http-request <RequestStructName> [<send_entire_struct_as_body>]`
    *   **Required.** Specifies the request object for the method.
    *   `<RequestStructName>`: The name of the Go struct that encapsulates all input parameters for this method.
    *   `[<send_entire_struct_as_body>]`: **Optional.** If a non-empty string (e.g., "true", "body") is provided, the entire `<RequestStructName>` instance is marshaled as the JSON request body. Otherwise, specific fields within `<RequestStructName>` must be tagged with `param:"body,..."` to contribute to the request body.

**Field-Level Struct Tags (within the `<RequestStructName>`):**

Fields in the request struct (specified by `@kit-http-request`) use the `param` struct tag to map them to HTTP request components:

*   ``param:"<type>,<name>"``
    *   `<type>`: Defines the parameter type.
        *   `path`: Maps the field to a path parameter. `<name>` must match a `{param_name}` in the `@kit-http` URL path.
        *   `query`: Maps the field to a URL query parameter. `<name>` is the query parameter key.
        *   `header`: Maps the field to an HTTP request header. `<name>` is the header name.
        *   `body`: Designates the field as part of the JSON request body. If `<send_entire_struct_as_body>` was *not* set in `@kit-http-request`, the generated client typically expects one field to be tagged `param:"body,..."` to be the source of the request body. If `<send_entire_struct_as_body>` *was* set, this tag might be redundant or used for specific field naming within the JSON if the struct's field name differs. The generator code implies that if `RequestBody` (derived from `[<send_entire_struct_as_body>]`) is false, it will use the *first* field it finds tagged with `param:"body,..."` as the request body.
    *   `<name>`: The name of the path variable, query parameter key, header name, or JSON field name (though for `body` with whole-struct sending, JSON field names usually come from struct field names or `json` tags).

**Generated Code (`kit_http_client.go`):**

*   **`HttpClientImpl` Interface:** A new interface mirroring the annotated one, but with method signatures suitable for a client:
    `MethodName(ctx context.Context, req RequestStructName, option *Option) (res ResponseStructName, err error)`
*   **`HttpClientService` Struct:** Implements `HttpClientImpl`.
    *   Handles URL construction (joining `@basePath` and method-specific paths).
    *   Manages path parameter substitution, query string formation, and header injection.
    *   Encodes the request body to JSON (typically using `kithttp.EncodeJSONRequest`).
    *   Decodes JSON responses.
    *   Integrates with Go kit's service discovery (`sd.Instancer`), load balancing (`lb.NewRoundRobin`), and retry mechanisms.
*   **`Option` Struct:** Allows customization of client behavior per call or globally:
    *   `PrePath`: Overrides `@basePath` or the `HttpClientService`'s global `PrePath`.
    *   `Logger`: `log.Logger` for Go kit.
    *   `Instancer`: `sd.Instancer` for service discovery.
    *   `RetryMax`, `RetryTimeout`: For retry configuration.
    *   `EndpointOpts`: `[]sd.EndpointerOption` for endpoint creation.
    *   `ClientOpts`: `[]kithttp.ClientOption` for `http.Client` customization (e.g., adding headers via `kithttp.ClientBefore`).
    *   `Encode`: Custom `kithttp.EncodeRequestFunc`.
    *   `Decode`: Custom factory for `kithttp.DecodeResponseFunc` (e.g., `func(i interface{}) func(ctx context.Context, res *http.Response) (response interface{}, err error)` where `i` is a pointer to the response struct).
*   **Validation:** Uses `github.com/asaskevich/govalidator` to validate the request struct before sending.

**Example:**

Interface definition (`myservice/client.go`):
```go
package myservice

import "context"

// @basePath /api/v1
type MyServiceClient interface {
    // @kit-http /users/{userID} GET
    // @kit-http-request GetUserRequest
    GetUser(ctx context.Context, req GetUserRequest) (UserResponse, error)

    // @kit-http /users POST
    // @kit-http-request CreateUserRequest true // Send entire CreateUserRequest as body
    CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)
}

type GetUserRequest struct {
    UserID string `param:"path,userID"` // Mapped to {userID} in path
    Token  string `param:"header,X-Auth-Token"`
}

type CreateUserRequest struct {
    Name  string `json:"name"` // Standard json tags for body
    Email string `json:"email"`
    Token string `param:"header,X-Auth-Token"` // Headers can still be separate
}

type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

Generated `myservice/kit_http_client.go` (conceptual):
```go
package myservice

// ... imports including context, net/http, go-kit components ...
import (
    "context"
    "encoding/json" // Added for conceptual decode example
    "net/http"      // Added for conceptual decode example
    // kithttp "github.com/go-kit/kit/transport/http" // For Option struct fields
    // gokitlog "github.com/go-kit/log" // For Option struct field
    // sd "github.com/go-kit/kit/sd" // For Option struct field
)


type Option struct { /* ... fields as described above ... */ }
type HttpClientService struct{ Option } // May hold global options
func NewHttpClientService(option Option) HttpClientImpl { /* ... */ return &HttpClientService{Option: option} } // Constructor

type HttpClientImpl interface {
    GetUser(ctx context.Context, req GetUserRequest, option *Option) (res UserResponse, err error)
    CreateUser(ctx context.Context, req CreateUserRequest, option *Option) (res UserResponse, err error)
}

func (s *HttpClientService) GetUser(ctx context.Context, req GetUserRequest, option *Option) (res UserResponse, err error) {
    // ... generated Go kit client logic for GET /api/v1/users/{userID} ...
    // ... uses req.UserID for path, req.Token for header ...
    // ... decodes response into res ...
    return
}

func (s *HttpClientService) CreateUser(ctx context.Context, req CreateUserRequest, option *Option) (res UserResponse, err error) {
    // ... generated Go kit client logic for POST /api/v1/users ...
    // ... sends req (CreateUserRequest) as JSON body, uses req.Token for header ...
    // ... decodes response into res ...
    return
}

// ... other helper types and functions ...
```

**Using the Generated Client:**

```go
package main

import (
    "context"
    "encoding/json" // Added for example decode function
    "fmt"
    "log"
    "myproject/myservice" // Where kit_http_client.go is generated
    "net/http" // Added for example decode function

    gokitlog "github.com/go-kit/log"
    // For service discovery, e.g., Consul or etcd
    // "github.com/go-kit/kit/sd/consul"
    // "github.com/hashicorp/consul/api"
    // kithttp "github.com/go-kit/kit/transport/http" // For ClientOpts
)

func main() {
    logger := gokitlog.NewLogfmtLogger(gokitlog.StdlibWriter{})

    // Setup Service Discovery (example, replace with your actual SD)
    // var instancer sd.Instancer
    // {
    //     consulClient, err := api.NewClient(api.DefaultConfig())
    //     if err != nil {
    //         log.Fatalf("consul client: %v", err)
    //     }
    //     instancer = consul.NewInstancer(consul.NewClient(consulClient), logger, "my-remote-service-name", nil, true)
    // }


    clientOptions := myservice.Option{
        Logger: logger,
        // Instancer: instancer, // From service discovery
        // ClientOpts: []kithttp.ClientOption{ /* ... custom client options ... */ },
        Decode: func(i interface{}) func(ctx context.Context, resp *http.Response) (response interface{}, err error) {
            // Default JSON decoder, i is a pointer to the success response struct
            return func(ctx context.Context, resp *http.Response) (response interface{}, err error) {
                if resp.StatusCode != http.StatusOK { // Example error handling
                    return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
                }
                if i == nil { // For methods that don't expect a body in success
                    return nil, nil
                }
                // Ensure i is a pointer, and then decode into it
                // This is a common pattern for generic decoders
                return i, json.NewDecoder(resp.Body).Decode(i)
            }
        },
    }
    
    // The generated code might provide a NewHttpClientService or similar constructor
    // The exact name depends on the generator's template.
    // Assuming the constructor is NewHttpClientService as per the conceptual generated code:
    client := myservice.NewHttpClientService(clientOptions)


    // Example call (if client is initialized)
    userResp, err := client.GetUser(context.Background(), myservice.GetUserRequest{
        UserID: "123",
        Token:  "mysecrettoken",
    }, nil) // Pass nil for call-specific options or an &myservice.Option{}

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("User: %+v\n", userResp)
    }
}
```
**Note:** The client requires proper setup for service discovery (`sd.Instancer`) and potentially custom `ClientOption` or `Encode`/`Decode` functions for advanced scenarios. The "Using" example is conceptual for the client part as the exact constructor name needs to be known from the generated code.

### `@temporal`

*   **Purpose:** Integrates a Go interface (service) with Temporal.io by generating code to register its methods as Temporal activities and providing workflow-callable wrapper methods to execute these activities.
*   **Target:** Go interfaces.
*   **Generated File:** `temporal.go` in the same package as the annotated interface.

**Interface-Level Directive:**

*   `// @temporal`
    *   Placed on the interface definition.
    *   Acts as a marker to indicate that this interface should be processed for Temporal integration.
    *   Does not take any arguments.

**Method-Level Directives (on methods within the `@temporal` annotated interface):**

*   `// @temporal-activity`
    *   **Required** on methods intended to be Temporal activities.
    *   Marks the method for registration with a Temporal worker.
    *   Causes a corresponding wrapper method to be generated on a `Temporal` struct, allowing the activity to be called from a Temporal workflow.

**Generated Code (`temporal.go`):**

1.  **`Temporal` Struct:**
    ```go
    type Temporal struct {
        next Service // Where 'Service' is your annotated interface type
        w    worker.Worker
    }
    ```

2.  **`InitTemporal` Function (Dependency Injection):**
    ```go
    // @do provide
    func InitTemporal(i do.Injector) (t *Temporal, err error) {
        w := do.MustInvoke[worker.Worker](i)
        next := do.MustInvoke[Service](i) // Your actual service implementation

        t = &Temporal{next: next, w: w}

        // Registers methods marked with @temporal-activity
        // Example: w.RegisterActivity(next.YourActivityMethod)
        // ...
        return t, nil
    }
    ```
    *   This function is intended to be used with `github.com/samber/do/v2` for dependency injection.
    *   It resolves a `worker.Worker` and your service implementation (`Service`) from the DI container.
    *   It registers all methods annotated with `@temporal-activity` from your service with the Temporal worker.

3.  **Workflow-Callable Methods (on `*Temporal` struct):**
    *   For each method in your interface that was annotated with `// @temporal-activity`, a corresponding wrapper method is generated on the `*Temporal` struct.
    *   **Signature Change:**
        *   The first parameter of these wrapper methods becomes `ctx workflow.Context` (from `go.temporal.io/sdk/workflow`).
        *   Other parameters and return types match the original method.
    *   **Implementation:** These wrappers use `workflow.ExecuteActivity` to call the actual registered activity.

**Example:**

Interface definition (`example/transfer_service.go`):
```go
package example

import "context"

// @temporal
type TransferService interface {
    // @temporal-activity
    ProcessTransfer(ctx context.Context, transactionID string, amount float64) (confirmationID string, err error)

    // This method will not be registered as an activity
    UtilityMethod(ctx context.Context, data string) error
}
```

Generated `example/temporal.go` (conceptual):
```go
package example

import (
	"go.temporal.io/sdk/workflow"
	"go.temporal.io/sdk/worker"
	// "github.com/samber/do/v2" // Assuming Injector is from here
	// ... other necessary imports from your original service
)

type Temporal struct {
	next TransferService // Your interface
	w    worker.Worker
}

// @do provide
func InitTemporal(i do.Injector) (t *Temporal, err error) { // Assuming do.Injector
	w := do.MustInvoke[worker.Worker](i)
	next := do.MustInvoke[TransferService](i)
	t = &Temporal{next: next, w: w}

	// Generated registration:
	w.RegisterActivity(next.ProcessTransfer)
	return t, nil
}

// Generated workflow-callable wrapper for ProcessTransfer
func (t *Temporal) ProcessTransfer(ctx workflow.Context, transactionID string, amount float64) (confirmationID string, err error) {
	err = workflow.ExecuteActivity(ctx, t.next.ProcessTransfer, transactionID, amount).Get(ctx, &confirmationID)
	return confirmationID, err
}
```

**Using with Temporal:**

1.  **Service Implementation:** Implement your `TransferService` interface.
    ```go
    package example

    import "context"
    import "fmt"

    type transferServiceImpl struct{}

    func NewTransferServiceImpl() TransferService { return &transferServiceImpl{} }

    func (s *transferServiceImpl) ProcessTransfer(ctx context.Context, transactionID string, amount float64) (string, error) {
        fmt.Printf("Processing transfer %s for amount %f\n", transactionID, amount)
        // Your actual business logic
        return "confirm-" + transactionID, nil
    }

    func (s *transferServiceImpl) UtilityMethod(ctx context.Context, data string) error {
        // ...
        return nil
    }
    ```

2.  **Temporal Worker Setup & DI:**
    ```go
    package main

    import (
        "log"
        "myproject/example" // Your service and generated temporal code package

        "github.com/samber/do/v2"
        "go.temporal.io/sdk/client"
        "go.temporal.io/sdk/worker"
    )

    func main() {
        injector := do.New()

        // 1. Provide Temporal Client
        do.Provide(injector, func(i do.Injector) (client.Client, error) {
            // return client.Dial(client.Options{ /* ... host, namespace ... */ })
            // For local testing without a full Temporal server:
            return client.NewLazyClient(nil, client.Options{}, nil, nil)
        })

        // 2. Provide Temporal Worker
        do.Provide(injector, func(i do.Injector) (worker.Worker, error) {
            c := do.MustInvoke[client.Client](i)
            // Define your task queue
            return worker.New(c, "MY_TASK_QUEUE", worker.Options{}), nil
        })

        // 3. Provide your actual service implementation
        do.Provide(injector, func(i do.Injector) (example.TransferService, error) {
            return example.NewTransferServiceImpl(), nil
        })

        // 4. Provide the generated Temporal struct (which also registers activities)
        // The generated InitTemporal is already a provider.
        do.Provide(injector, example.InitTemporal)


        // Ensure Temporal struct is initialized (which registers activities)
        // The _ syntax indicates we're calling it for its side effect (activity registration)
        _ = do.MustInvoke[*example.Temporal](injector) 
        // if err != nil { // MustInvoke will panic on error, so direct err check isn't typical here
        //     log.Fatalf("Failed to initialize Temporal integration: %v", err)
        // }
        
        // Start the worker
        temporalWorker := do.MustInvoke[worker.Worker](injector) // Corrected 'i' to 'injector'
        err := temporalWorker.Run(worker.InterruptCh())
        if err != nil {
            log.Fatalf("Worker failed: %v", err)
        }
    }
    ```

3.  **In your Workflow Definition:**
    ```go
    package main // Or your workflow package

    import (
    	"time" // Added for ActivityOptions
    	"go.temporal.io/sdk/workflow"
    	"myproject/example" // Assuming Temporal struct is here
    )

    func MyWorkflow(ctx workflow.Context, transactionID string, amount float64) (string, error) {
        ao := workflow.ActivityOptions{
            TaskQueue:           "MY_TASK_QUEUE",
            StartToCloseTimeout: time.Minute * 10,
        }
        ctx = workflow.WithActivityOptions(ctx, ao)

        var temporalActivities *example.Temporal // Zero value is fine for activity execution
        var confirmationID string
        // The following line directly calls the generated wrapper.
        // Alternatively, you could pass the function reference:
        // err := workflow.ExecuteActivity(ctx, temporalActivities.ProcessTransfer, transactionID, amount).Get(ctx, &confirmationID)
        // However, the generated wrapper itself uses ExecuteActivity.
        // So, you call the wrapper which calls ExecuteActivity on the *actual* service method.
        confirmationID, err := temporalActivities.ProcessTransfer(ctx, transactionID, amount)
        if err != nil {
            return "", err
        }
        return confirmationID, nil
    }
    ```
**Note:** The generated `(t *Temporal) ProcessTransfer` method is a convenience wrapper that internally calls `workflow.ExecuteActivity` on the original service method (e.g., `t.next.ProcessTransfer`). When defining workflows, you can call this generated wrapper method directly. The key is that `InitTemporal` registers the *actual* service methods (e.g., `next.ProcessTransfer`) with the Temporal worker.

### `@copy` (Struct Copying)

*   **Purpose:** Generates type-safe functions to copy data between two Go structs (or types). It handles nested structures, slices, maps, pointers, and allows for custom field mapping through annotations.
*   **Target:** Go function call expressions. The signature of the annotated function defines the source and destination types.
*   **Generated File:** `copy.go` (in the package of the annotated function).

**Triggering Mechanism:**

You trigger the `@copy` plugin by defining a function (whose body is usually empty or ignored) and annotating it with `// @copy`. The function's signature dictates the copy operation:

```go
package myconverters

// @copy
func ConvertUserToUserDTO(dto *UserDTO, user model.User) {
    // This function's body is not executed by genx.
    // Its signature defines the types for the generated copy function.
}

// UserDTO and model.User are your destination and source struct types.
```

*   The first parameter is the **destination** (e.g., `dto *UserDTO`). It **must be a pointer type**.
*   The second parameter is the **source** (e.g., `user model.User`).

**Generated Code (`copy.go`):**

For the example `ConvertUserToUserDTO` above, `genx` would generate:

1.  **The Original Function (Filled In):**
    ```go
    package myconverters

    func ConvertUserToUserDTO(dto *UserDTO, user model.User) {
        ConvertUserToUserDTOCopy{}.Copy(dto, user) // Delegates to the Copy method
        // If ConvertUserToUserDTO had return values, they would be returned here.
    }
    ```

2.  **A Copier Struct:**
    ```go
    type ConvertUserToUserDTOCopy struct{} // Named <YourFunctionName>Copy
    ```

3.  **A `Copy` Method on the Copier Struct:**
    ```go
    func (d ConvertUserToUserDTOCopy) Copy(dto *UserDTO, user model.User) {
        // Auto-generated field copying logic:
        dto.ID = user.ID
        dto.Name = user.Name
        // ... handles pointers, slices, maps, nested structs recursively ...

        // Example for a nested struct (if User and UserDTO had Address fields):
        // if user.Address != nil && dto.Address == nil { // Assuming Address is a pointer
        //     dto.Address = new(AddressDTO)
        // }
        // if user.Address != nil { // Check if source Address is not nil
        //     ConvertUserToUserDTOCopy{}.AddressToAddressDTOCopy(dto.Address, *user.Address) // Recursive call for nested types
        // }
        return
    }
    ```
    *   This method contains the core logic for copying fields from the source to the destination.
    *   It intelligently handles various types:
        *   Direct assignment for compatible basic types and named types.
        *   Nil checks for source pointers.
        *   Allocation for destination pointers if they are nil and the source is not (e.g., `dest.Nested = new(NestedDTO)`).
        *   Deep copying for slices and maps of complex types by generating copy logic for their elements.
        *   Recursive generation: If it encounters nested structs that need copying (e.g., `User.Address` to `UserDTO.AddressDTO`), it will try to generate or use another `Copy` method for those specific types (e.g., `AddressToAddressDTOCopy`).

**Field Mapping Customization (Annotations on Source or Destination Struct Fields):**

You can annotate fields within your source or destination structs to control how they are mapped. These annotations are typically placed as comments above the field.

*   `// @copy-prefix <prefix_to_remove>`: (Usually on source field) If the source field name has a prefix (e.g., `Source_FieldA`) that should be removed before matching with a destination field name (e.g., `FieldA`).
*   `// @copy-name <alternative_name>`: (Usually on source field) Use this `alternative_name` for matching against destination field names, instead of the actual Go field name.
*   `// @copy-must`: If this annotation is present on a source field and no corresponding settable field can be found in the destination (after applying name/prefix rules), the code generation will panic.
*   `// @copy-target-path <full.dot.path.to.dest.field>`: (Usually on source field) Explicitly maps the annotated source field to the specified full path in the destination struct. This offers the highest precedence for matching. Example: `// @copy-target-path UserDetails.Contact.EmailAddress`.
*   `// @copy-target-name <dest_field_name>`: (Usually on source field) Maps the annotated source field to a destination field with this specific name, even if paths differ.
*   `// @copy-target-method <method_name_on_source>`: (Usually on source field) Instead of direct field access from the source, call `<method_name_on_source>()` on the source struct (or the specific source field if it's a nested struct) to get the value. This value is then copied to the destination field matched by name/path.
*   `// @copy-auto-cast`: (Boolean, e.g., `@copy-auto-cast`) If present, may enable more lenient type conversions (e.g., between different integer types or to/from string) if direct assignment is not possible. Specific casting behaviors depend on the generator's implementation details.

**Matching Strategy:**

The plugin attempts to match fields between source and destination by considering:
1.  Explicit mapping via `@copy-target-path` or `@copy-target-name`.
2.  Name matching (after applying `@copy-name` or `@copy-prefix` to the source field name).
3.  A "depth find" algorithm that tries to match fields even within different levels of embedded structs, effectively "flattening" paths where appropriate.

**Example with Field Annotations:**

Source struct:
```go
package model

type User struct {
    ID           int
    InternalName string // @copy-target-name FullName (map to DTO's FullName)
    UserEmail    string // @copy-prefix User (map to DTO's Email)
    Details      UserDetails
}

type UserDetails struct {
    Bio string
    // @copy-target-path Preferences.Theme
    ThemePreference string
}
```

Destination DTO:
```go
package myconverters

type UserDTO struct {
    ID         int
    FullName   string
    Email      string
    Bio        string
    Preferences UserPreferencesDTO
}

type UserPreferencesDTO struct {
    Theme string
}
```

Generated `Copy` method (conceptual for `ConvertUserToUserDTO`):
```go
// ...
dto.ID = user.ID
dto.FullName = user.InternalName // Due to @copy-target-name
dto.Email = user.UserEmail       // Due to @copy-prefix User
dto.Bio = user.Details.Bio       // Matched by path/name
dto.Preferences.Theme = user.Details.ThemePreference // Due to @copy-target-path
// ...
```

This plugin is highly useful for reducing boilerplate when mapping between different struct types, such as between API DTOs, database models, and internal domain objects.

### `@kit` (Go kit HTTP Service Generation)

*   **Purpose:** Scaffolds a Go kit HTTP service based on a Go interface. It generates Go kit endpoints, HTTP transport layers (request/response encode/decode functions), and HTTP handlers.
*   **Target:** Go interfaces.
*   **Generated Files:**
    *   `endpoint.go`: Contains Go kit endpoint definitions for each method in the annotated interface.
    *   `http.go`: Contains HTTP transport layer logic, including request decoding functions, response encoding functions, and Go kit HTTP server handlers.

**Interface-Level Directives (applied to the interface definition):**

*   `// @basePath <base_path_string>`
    *   **Optional.** Specifies a common base path prefix for all HTTP URLs defined for the methods in this interface (e.g., `/api/v1`, `/catalog`).
*   `// @tags <comma_separated_tags>`
    *   **Optional.** Used for Swagger/OpenAPI documentation to group endpoints (e.g., `@tags "Users,Admin"`).
*   `// @validVersion <version_string>`
    *   **Optional.** Specifies an API version, potentially used in generated documentation or routing logic.
*   `// @kit-server-option <OptionName1> [<OptionName2> ...]`
    *   **Optional.** A list of Go kit HTTP server options (e.g., `ErrorEncoder(myErrorHandler)`) that will be passed to `kithttp.NewServer()` when creating the handlers. The values are typically function names or fully qualified identifiers available in the scope.
*   `// @swag <true_or_false>`
    *   **Optional.** An interface-level switch to enable/disable Swagger/OpenAPI documentation generation for all endpoints derived from this interface. Defaults to true if not specified. Individual methods can override this.

**Method-Level Directives (on methods within the `@kit` annotated interface):**

*   `// @kit-http <http_url_path> <http_method>`
    *   **Required** for a method to be exposed as an HTTP endpoint.
    *   `<http_url_path>`: The URL path specific to this method. Path parameters should be defined using Gorilla Mux syntax (e.g., `/items/{itemID}`, `/documents/{category}/{docId}`).
    *   `<http_method>`: The HTTP method (e.g., `GET`, `POST`, `PUT`, `DELETE`, `PATCH`).

*   `// @kit-http-request <RequestStructName> [<is_request_body_bool>]`
    *   **Required** if the service method has parameters other than `context.Context`.
    *   `<RequestStructName>`: The name of the Go struct that will be used to capture and decode incoming HTTP request data (path parameters, query parameters, headers, request body).
    *   `[<is_request_body_bool>]`: **Optional.** A boolean-like string ("true" or "false", or any non-empty string for true).
        *   If "true" (or non-empty), the entire `<RequestStructName>` instance is expected to be decoded from the HTTP request body (typically JSON).
        *   If "false" (or empty/omitted), the request body might not be directly mapped from the `<RequestStructName>` as a whole, or specific fields within `<RequestStructName>` might be designated to come from the body using tags (though this plugin primarily focuses on the struct itself as the body or individual fields from path/query/header). The generated code will often create separate `Decode<MethodName>Request` functions that handle populating this struct.

*   `// @swag <true_or_false>`
    *   **Optional.** Method-level switch to enable/disable Swagger/OpenAPI documentation generation for this specific endpoint. Overrides the interface-level `@swag` setting.

**Field-Level Struct Tags (within the `<RequestStructName>`):**

Fields within your request structs (those named in `@kit-http-request`) use struct tags to specify how they should be populated from the incoming HTTP request:

*   `path:"<name>"`: Maps the field to a path parameter. `<name>` must match a placeholder name defined in the `@kit-http` URL path (e.g., if path is `/product/{id}`, tag would be `path:"id"`).
*   `query:"<name>"`: Maps the field to a URL query parameter (e.g., `query:"searchTerm"`).
*   `header:"<name>"`: Maps the field to an HTTP request header (e.g., `header:"X-API-Key"`).
*   `json:"<name>"`: Standard Go struct tag for JSON marshaling/unmarshaling. Used when the request struct (or parts of it) is decoded from a JSON body or when the response struct is encoded to JSON.
*   `validate:"<rules>"`: Specifies validation rules for the field, often using syntax compatible with libraries like `govalidator` or `go-playground/validator` (e.g., `validate:"required,min=1,max=100"`). The generated decode functions typically incorporate calls to a validator.
*   `form:"<name>"`: For HTML form values (e.g., from `application/x-www-form-urlencoded` bodies).

**Generated Code (`endpoint.go` and `http.go`):**

*   **`endpoint.go`:**
    *   Defines Go kit `endpoint.Endpoint` for each method annotated with `@kit-http`.
    *   Creates `Make<MethodName>Endpoint` functions that wrap your service methods.
    *   Often includes an `Endpoints` struct that groups all generated endpoints.
*   **`http.go`:**
    *   **Request Decoding:** Generates `Decode<MethodName>Request` functions (e.g., `DecodeCreateUserRequest`). These functions extract data from the HTTP request (path, query, headers, body) and populate the `<RequestStructName>`. They often include validation logic.
    *   **Response Encoding:** Generates `Encode<MethodName>Response` functions (or a generic `EncodeResponse`) to marshal the service method's results (and errors) into an HTTP response, typically JSON.
    *   **HTTP Handlers:** Creates Go kit `kithttp.Handler` instances for each endpoint, wiring up the decode/encode functions and server options.
    *   A function (e.g., `NewHTTPHandler` or `MakeHandler`) that takes the `Endpoints` struct and returns an `http.Handler` (often a `*mux.Router` from Gorilla Mux) with all routes configured.

**Example:**

Interface (`productservice/service.go`):
```go
package productservice

import "context"

// @basePath /v1/products
// @tags "Products"
type Service interface {
    // @kit-http /{id} GET
    // @kit-http-request GetProductRequest
    GetProduct(ctx context.Context, req GetProductRequest) (ProductResponse, error)

    // @kit-http / POST
    // @kit-http-request CreateProductRequest true // Request struct is the body
    CreateProduct(ctx context.Context, req CreateProductRequest) (ProductResponse, error)
}

type GetProductRequest struct {
    ID string `path:"id" validate:"required,uuid"`
    // @swag false // Example: disable swag for this field if supported by template
}

type CreateProductRequest struct {
    Name  string  `json:"name" validate:"required"`
    Price float64 `json:"price" validate:"required,gt=0"`
}

type ProductResponse struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

Conceptual Generated Code:

`productservice/endpoint.go`:
```go
// package productservice
// ... imports ...
// type Endpoints struct { GetProductEndpoint endpoint.Endpoint; CreateProductEndpoint endpoint.Endpoint; }
// func MakeGetProductEndpoint(s Service) endpoint.Endpoint { ... }
// func MakeCreateProductEndpoint(s Service) endpoint.Endpoint { ... }
// func NewEndpoint(s Service, dmw ...middleware) Endpoints { ... } // Simplified
```

`productservice/http.go`:
```go
// package productservice
// ... imports ...
// func DecodeGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) { ... }
// func DecodeCreateProductRequest(_ context.Context, r *http.Request) (interface{}, error) { ... }
// func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error { ... }
// func NewHTTPHandler(ep Endpoints, serverOptions ...kithttp.ServerOption) http.Handler {
//    router := mux.NewRouter()
//    router.Methods("GET").Path("/v1/products/{id}").Handler(kithttp.NewServer(ep.GetProductEndpoint, DecodeGetProductRequest, EncodeResponse, serverOptions...))
//    router.Methods("POST").Path("/v1/products").Handler(kithttp.NewServer(ep.CreateProductEndpoint, DecodeCreateProductRequest, EncodeResponse, serverOptions...))
//    return router
// }
```

**Using the Generated Service:**
You would typically implement the `Service` interface, then in your `main.go` or setup code:
1. Create an instance of your service implementation.
2. Call the generated `NewEndpoint` function to create the Go kit endpoints, possibly applying middleware.
3. Call the generated `NewHTTPHandler` (or similarly named function) to get an `http.Handler`.
4. Start an HTTP server with this handler.

This plugin automates much of the boilerplate involved in setting up Go kit HTTP services.

### `@cep` (CE Permission SQL Generation)

*   **Purpose:** Generates SQL `INSERT` statements for a permissions system, apparently named `sys_permission`. It creates SQL entries for a parent menu item (derived from the annotated interface) and child permission entries for each HTTP method exposed via the `@kit-http` directive within that interface.
*   **Target:** Go interfaces (the same ones annotated with `@kit` and its sub-directives).
*   **Generated File:** `cep.sql` in the same package as the interface.
*   **Note:** This plugin seems to be for a specific internal use case ("CE" likely refers to "Cloud Engine" or a similar internal system) and might not be relevant for all users of `genx`.

**Prerequisites:**

*   The interface should be annotated for the `@kit` plugin, especially with `@basePath` on the interface and `@kit-http` on methods, as `@cep` uses this information.

**Interface-Level Information Used:**

*   The primary comment on the interface definition (e.g., `// ServiceName is a service for...`) is used as the `alias` and `description` for the parent menu item in the SQL.
*   `@basePath` is used to construct the `path` for the parent menu item (e.g., `<basePath>/index`).

**Method-Level Information Used:**

*   For each method annotated with `// @kit-http <path> <http_method>`:
    *   The method's primary comment is used as the `alias` and `description` for the permission entry.
    *   The `<http_method>` (e.g., GET, POST) is recorded.
    *   The full path (basePath + method path) is recorded.
    *   A unique name for the permission entry is generated (e.g., based on base path, method name, and HTTP method).
    *   A unique ID (integer) is generated by hashing the unique name.

**Generated SQL (`cep.sql` - conceptual):**

The file will contain a series of `INSERT INTO sys_permission (...) VALUES (...);` statements.

**Example:**

Given an interface:
```go
package userservice

import "context"

// @basePath /admin/users
// @tags "User Management"
// UserService manages user operations.
type UserService interface {
    // @kit-http /list GET
    // @kit-http-request ListUsersRequest
    // ListAllUsers retrieves all users.
    ListAllUsers(ctx context.Context, req ListUsersRequest) (UserListResponse, error)

    // @kit-http /{id} DELETE
    // @kit-http-request DeleteUserRequest
    // DeleteUserByID removes a user.
    DeleteUserByID(ctx context.Context, req DeleteUserRequest) (error)
}
// ... request/response structs ...
```

Conceptual `userservice/cep.sql`:
```sql
-- Parent Menu Item for UserService
INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description, created_at, updated_at) VALUES (123456789012345, 0, '', 1, 'GET', 'UserService manages user operations.', 'menu.admin.users.userservice', '/admin/users/index', 'UserService manages user operations.', 'YYYY-MM-DD HH:MM:SS', 'YYYY-MM-DD HH:MM:SS');

-- Permission for ListAllUsers
INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description, created_at, updated_at) VALUES (987654321098765, 123456789012345, '', 0, 'GET', 'ListAllUsers retrieves all users.', '.admin.users.listallusers.get', '/admin/users/list', 'ListAllUsers retrieves all users.', 'YYYY-MM-DD HH:MM:SS', 'YYYY-MM-DD HH:MM:SS');

-- Permission for DeleteUserByID
INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description, created_at, updated_at) VALUES (543210987654321, 123456789012345, '', 0, 'DELETE', 'DeleteUserByID removes a user.', '.admin.users.deleteuserbyid.delete', '/admin/users/{id}', 'DeleteUserByID removes a user.', 'YYYY-MM-DD HH:MM:SS', 'YYYY-MM-DD HH:MM:SS');
```
*(Note: The exact `id`, `name`, and `alias` generation logic might vary slightly based on internal hashing and string manipulation rules within the template `ce_permission.tmpl`.)*

This plugin is useful if you are integrating your `genx`-generated services with a system that uses this specific `sys_permission` table structure.

### `@observer` (Logging and Tracing Middleware Generation)

*   **Purpose:** Generates Go kit service middleware for logging and tracing. This plugin seems to provide an alternative or possibly older way to add observability to Go kit services generated by `genx`, distinct from the more modular `@log` and `@trace` plugins.
*   **Target:** Go interfaces (intended for services defined for Go kit).
*   **Generated Files:**
    *   `logging.go`: Contains logging middleware for the service.
    *   `tracing.go`: Contains tracing middleware (appears to use OpenTracing, based on common imports seen in its generator).
*   **Note:** This plugin might be specific to a particular setup or represent an earlier approach to observability within `genx`. For new projects, consider if `@log` and `@trace` (or the comprehensive `@otel`) offer more flexibility.

**Prerequisites:**

*   The interface should ideally be one for which Go kit HTTP bindings are also being generated (e.g., using `@kit`), as observability is typically applied to service endpoints.

**How it Works (Conceptual):**

The `@observer` plugin processes the annotated interface and uses templates (`ce_log.tmpl`, `ce_trace.tmpl`) to generate:

1.  **Logging Middleware (`logging.go`):**
    *   A struct (e.g., `loggingMiddleware`) that wraps the service interface.
    *   Methods corresponding to each method in the service interface.
    *   These wrapper methods log information about the call (e.g., method name, parameters, errors, duration) before and/or after calling the actual service method.
    *   A constructor function to create an instance of this logging middleware, typically expecting a logger instance (e.g., `log.Logger`) and the next service in the chain.

2.  **Tracing Middleware (`tracing.go`):**
    *   A struct (e.g., `tracingMiddleware`) that wraps the service interface.
    *   Methods corresponding to each method in the service interface.
    *   These wrapper methods create spans (likely using OpenTracing, judging by typical imports like `opentracing-go`) for each call, record tags (like method name, errors), and finish the span upon completion.
    *   A constructor function to create an instance of this tracing middleware, typically expecting a tracer instance (e.g., `opentracing.Tracer`) and the next service.

**Usage:**

Annotate your service interface with `// @observer`.

```go
package myservice

import "context"

// @observer
// @kit (and other @kit-* annotations if also generating HTTP server)
type MyObservedService interface {
    ProcessData(ctx context.Context, data string) (result string, err error)
    FetchItem(ctx context.Context, itemID int) (itemData string, err error)
}
```

**Generated Files (Conceptual Content):**

`myservice/logging.go`:
```go
// package myservice
// ... imports including a logger ...
// type loggingMiddleware struct {
//     next   MyObservedService
//     logger log.Logger // Or other logger type
// }
//
// func NewLoggingMiddleware(logger log.Logger, next MyObservedService) MyObservedService {
//     return &loggingMiddleware{logger, next}
// }
//
// func (mw *loggingMiddleware) ProcessData(ctx context.Context, data string) (result string, err error) {
//     defer func(begin time.Time) {
//         mw.logger.Log("method", "ProcessData", "input", data, "result", result, "err", err, "took", time.Since(begin))
//     }(time.Now())
//     result, err = mw.next.ProcessData(ctx, data)
//     return
// }
// ... other methods ...
```

`myservice/tracing.go`:
```go
// package myservice
// ... imports including opentracing ...
// type tracingMiddleware struct {
//     next   MyObservedService
//     tracer opentracing.Tracer
// }
//
// func NewTracingMiddleware(tracer opentracing.Tracer, next MyObservedService) MyObservedService {
//     return &tracingMiddleware{tracer, next}
// }
//
// func (mw *tracingMiddleware) ProcessData(ctx context.Context, data string) (result string, err error) {
//     span := mw.tracer.StartSpan("ProcessData", opentracing.ChildOf(opentracing.SpanFromContext(ctx)))
//     defer span.Finish()
//     span.SetTag("input.data", data)
//     result, err = mw.next.ProcessData(ctx, data)
//     if err != nil {
//         ext.Error.Set(span, true)
//         span.LogKV("error", err.Error())
//     }
//     span.SetTag("output.result", result)
//     return
// }
// ... other methods ...
```

**Applying the Middleware (in Go kit setup):**

If you are using this with services generated by `@kit`, you would typically apply these middlewares when constructing your endpoints or service instances.

```go
// import (
//    "myproject/myservice"
//    "github.com/go-kit/kit/endpoint"
//    "github.com/go-kit/log"
//    "github.com/opentracing/opentracing-go"
// )

// var serviceImpl myservice.MyObservedService // Your concrete implementation
// var logger log.Logger
// var tracer opentracing.Tracer

// serviceWithLogging := myservice.NewLoggingMiddleware(logger, serviceImpl)
// serviceWithAllObs := myservice.NewTracingMiddleware(tracer, serviceWithLogging)

// Pass 'serviceWithAllObs' when creating Go kit Endpoints:
// endpoints := myservice.NewEndpoint(serviceWithAllObs, ...) // Assuming NewEndpoint from @kit generation
```

This plugin provides a template-based approach to add logging and tracing middleware directly tied to the service definition.

### `@alert`

*   **Purpose:** Generates a Go kit style middleware that enables sending alerts when specific service methods return errors. It integrates with an external alerting service.
*   **Target:** Go interfaces.
*   **Generated File:** `alert.go` in the same package as the annotated interface.

**Interface-Level Directive:**

*   `// @alert`
    *   Placed on the interface definition.
    *   Marks the interface for processing by the `@alert` plugin. It does not take arguments itself.

**Method-Level Directives (on methods within the `@alert` annotated interface):**

*   `// @alert-enable`
    *   **Required** on methods for which error alerting should be active. If this directive is absent, errors from that method will not trigger alerts via this middleware.
*   `// @alert-level <level>`
    *   **Optional.** Specifies the severity level for the alert.
    *   Valid values: `"info"`, `"warn"`, `"error"`.
    *   If omitted, the default alert level provided during the middleware's initialization is used.
    *   These are mapped to an assumed `alarm.LevelInfo`, `alarm.LevelWarning`, `alarm.LevelError` from an `alarm` package.
*   `// @alert-metrics <metrics_suffix>`
    *   **Optional.** A string suffix that gets appended to the default metrics identifier (which is `pkgName.MethodName`) when an alert is pushed. Example: if `metrics_suffix` is `"custom_metric"`, the identifier becomes `pkgName.MethodName.custom_metric`.

**Generated Code (`alert.go`):**

1.  **`alert` Struct:**
    *   A struct named `alert` is generated. It holds:
        *   `level alarm.Level`: The default alert level.
        *   `silencePeriod int`: A period for silencing alerts.
        *   `api api.Service`: An instance of an external service used to push alerts. This service must expose an `Alarm().Push(...)` method.
        *   `next Service`: The next service in the Go kit middleware chain (your actual service implementation).
        *   `logger log.Logger`: A logger for internal errors if pushing an alert fails.

2.  **`NewAlert` Constructor Function (Middleware Factory):**
    ```go
    // Conceptual signature from generated code:
    func NewAlert(level alarm.Level, silencePeriod int, api api.Service, log log.Logger) Middleware {
        return func(next Service) Service {
            return &alert{ /* ... initialization ... */ }
        }
    }
    ```
    *   This function is a Go kit `Middleware` factory.
    *   You provide the default `alarm.Level`, `silencePeriod`, the `api.Service` implementation for sending alerts, and a `log.Logger`.

3.  **Wrapped Methods (on `*alert` struct):**
    *   For every method in your annotated interface, a corresponding wrapper method is generated.
    *   If a wrapped method:
        *   Is annotated with `// @alert-enable`, AND
        *   Returns an error (`err != nil`),
        *   Then, a `defer` block triggers an alert by calling:
            `s.api.Alarm().Push(ctx, title, traceId+err.Error(), metricsKey, alertLevel, s.silencePeriod)`
            *   `title`: Typically `packageName.MethodName`.
            *   `traceId`: Extracted from `ctx.Value("traceId")`.
            *   `metricsKey`: `packageName.MethodName` or `packageName.MethodName.metrics_suffix`.
            *   `alertLevel`: Determined by `@alert-level` or the default.
    *   The method then calls the `next` service's corresponding method.

**Dependencies & Assumptions:**

*   **`api.Service`:** Your project must define or import an `api` package with a `Service` interface (or concrete type) that has an `Alarm()` method, which in turn returns an object with a `Push(...)` method. The expected signature for `Push` is approximately:
    `Push(ctx context.Context, title string, content string, metricsKey string, level alarm.Level, silencePeriod int) error`
*   **`alarm` package:** An `alarm` package is assumed to exist, providing `alarm.Level` and constants like `alarm.LevelInfo`, `alarm.LevelWarning`, `alarm.LevelError`.
*   **`log.Logger`:** A logger compatible with the Go kit logger interface.
*   The `context.Context` passed to service methods may optionally contain a `traceId` value, which will be included in the alert content.

**Example:**

Interface definition (`payments/service.go`):
```go
package payments

import "context"

// @alert
type PaymentService interface {
    // @alert-enable
    // @alert-level error
    // @alert-metrics credit_card_failures
    ProcessPayment(ctx context.Context, amount float64, cardToken string) (transactionID string, err error)

    // @alert-enable
    // @alert-level warn
    RefundPayment(ctx context.Context, transactionID string, reason string) (err error)

    // No @alert-enable, so no alerts for this method's errors
    GetPaymentStatus(ctx context.Context, transactionID string) (status string, err error)
}
```

Conceptual usage (in your Go kit service setup):
```go
// import (
//     "myproject/payments"
//     "myproject/pkg/alarm" // Your alarm package
//     "myproject/pkg/alertapi" // Your alert sending API service
//     "github.com/go-kit/log"
// )

// var serviceImpl payments.PaymentService // Your concrete implementation
// var logger log.Logger
// var alertApiService alertapi.Service // Implements the expected api.Service for alerts

// alertMiddleware := payments.NewAlert(
//     alarm.LevelError, // Default level if not specified in annotation
//     300,              // Example silence period (e.g., 5 minutes)
//     alertApiService,
//     logger,
// )

// serviceWithAlerts := alertMiddleware(serviceImpl)

// Use 'serviceWithAlerts' when creating Go kit Endpoints
// endpoints := payments.NewEndpoint(serviceWithAlerts, ...)
```

This plugin helps in centralizing alerting logic for service method failures.

### `@do` (Dependency Injection Setup)

*   **Purpose:** Automates the registration of service providers with the `github.com/samber/do/v2` dependency injection library. It scans for annotated provider functions and generates a central `doInit` function to register them with a `do.Injector`.
*   **Target:** Go functions across all processed packages.
*   **Generated File:** `do_init.go`. This file is generated in the package of any function that is annotated with `@do(type="init")`. The `doInit` function within this file will register providers from *all* scanned packages.

**Function-Level Directives:**

Functions should be annotated with `// @do(...)` to be processed by this plugin.

*   `// @do(type="<type_value>"[, name="<service_name>"])`
    *   `type="<type_value>"`: **Required.**
        *   `provide`: Marks the annotated function as a constructor or provider that should be registered with the `do.Injector`. The function signature should be compatible with `do.Provide` (e.g., `func(i do.Injector) (MyService, error)` or `func() OtherService`).
        *   `init`: This function's primary role is to mark its package as a location where a `do_init.go` file will be generated. The `doInit` function inside this generated file will contain the logic to register all discovered "provide" functions. The body of the function marked `type="init"` is not directly used by `genx`.
    *   `name="<service_name>"`: **Optional.** Only applicable if `type="provide"`. If specified, the provider function is registered using `do.ProvideNamed(injector, "<service_name>", yourProviderFunc)`. Otherwise, `do.Provide(injector, yourProviderFunc)` is used.

**Generated Code (`do_init.go`):**

*   A file named `do_init.go` is generated in each package that contains at least one function annotated with `@do(type="init")`.
*   This file will contain a function: `func doInit(i do.Injector)`.
*   The `doInit` function includes:
    *   Import statements for all packages containing the provider functions it needs to register. `genx` handles potential package name collisions by creating aliases.
    *   Calls to `do.Provide(i, ...)` or `do.ProvideNamed(i, ...)` for every function that was annotated with `@do(type="provide")` across all scanned packages.
    *   Providers are registered in a deterministic order (sorted by package path, then by function name).

**How It Works:**

1.  `genx` scans all Go source files in the project.
2.  It identifies all functions annotated with `@do(...)`.
3.  It categorizes them based on `type="provide"` or `type="init"`.
4.  For every package that contains one or more functions marked `type="init"`, `genx` generates a `do_init.go` file in that package.
5.  The `doInit` function within any such generated `do_init.go` file will contain the registration calls for *all* provider functions found in the entire project scan. Therefore, you typically only need one `// @do(type="init")` annotation in your entire project (e.g., in a `cmd/myapp/main.go` or an `internal/app/init.go` file) to generate a single, comprehensive `doInit` function.

**Example:**

`services/user_service.go`:
```go
package services

// @do(type="provide")
func NewUserService() string { // Example: simple string provider
    return "Hello from UserService"
}
```

`services/another_service.go`:
```go
package services

// @do(type="provide", name="namedService")
func NewNamedService() int { // Example: int provider, named
    return 123
}
```

`app/setup.go`:
```go
package app

// @do(type="init")
// This function's body is not used by genx; its annotation triggers do_init.go generation in this 'app' package.
func InitializeAppContainer() {}
```

Generated `app/do_init.go`:
```go
package app

import (
	"github.com/samber/do/v2"
	servicesalias "myproject/services" // Alias may or maynot be needed
)

// doInit registers all @do(type="provide") functions.
func doInit(i do.Injector) {
	do.Provide(i, servicesalias.NewUserService)
	do.ProvideNamed(i, "namedService", servicesalias.NewNamedService)
}
```
*(Note: The exact alias `servicesalias` is illustrative; `genx` will generate one if the package name `services` conflicts with other imports or local names.)*

**Using the Generated `doInit`:**

In your application's main entry point:
```go
package main

import (
	"fmt"
	"myproject/app" // Package where do_init.go was generated
	"github.com/samber/do/v2"
)

func main() {
	injector := do.New()

	// Call the generated initializer
	app.doInit(injector) // Or app.Invoke(injector) if that's the pattern

	// Now you can resolve your services
	userServiceValue := do.MustInvoke[string](injector)
	fmt.Println(userServiceValue) // Output: Hello from UserService

	namedServiceValue := do.MustInvokeNamed[int](injector, "namedService")
	fmt.Println(namedServiceValue) // Output: 123
}
```

This plugin streamlines the setup of `samber/do/v2` by automatically generating the provider registration code, reducing boilerplate and potential for manual errors.
