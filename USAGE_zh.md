# `genx` 使用指南

## 目录

- [什么是 `genx`?](#什么是-genx)
- [安装](#安装)
- [运行 `genx`](#运行-genx)
- [配置 (`genx.yaml`)](#配置-genxyaml)
  - [`imports`](#imports)
  - [`preloads`](#preloads)
  - [示例 `genx.yaml`](#示例-genxyaml)
- [TUI 输出](#tui-输出)
- [注解语法 (魔术注释)](#注解语法-魔术注释)
  - [1. `@directive` 风格 (主要)](#1-directive-风格-主要)
  - [2. `@FieldDirective` 风格 (旧式/备选)](#2-fielddirective-风格-旧式备选)
  - [3. `INSET` 风格 (隐式内容行)](#3-inset-风格-隐式内容行)
  - [重要说明:](#重要说明)
- [插件](#插件)
  - [`@log`](#log)
  - [`@trace`](#trace)
  - [`@otel`](#otel)
  - [`@enum`](#enum)
  - [`@gq` (GORM 查询作用域)](#gq-gorm-查询作用域)
  - [`@crud`](#crud)
  - [`@kit-http-client`](#kit-http-client)
  - [`@temporal`](#temporal)
  - [`@copy` (结构体复制)](#copy-结构体复制)
  - [`@kit` (Go kit HTTP 服务生成)](#kit-go-kit-http-服务生成)
  - [`@cep` (CE 权限 SQL 生成)](#cep-ce-权限-sql-生成)
  - [`@observer` (日志和追踪中间件生成)](#observer-日志和追踪中间件生成)
  - [`@alert`](#alert)
  - [`@do` (依赖注入设置)](#do-依赖注入设置)

## 什么是 `genx`?

`genx` 是一个用于 Go 的代码生成器。它通过解析 Go 源代码文件中的特殊注解（注释）来工作。`genx` 的主要目的是通过自动化生成重复的代码模式，帮助减少样板代码并在项目中保持一致性。

## 安装

要安装 `genx`，您可以使用 `go install` 命令。打开您的终端并运行以下命令：

```bash
go install github.com/fitan/genx@latest
```

此命令将下载 `genx` 可执行文件并将其安装到您的 Go bin 目录中。请确保您的 Go bin 目录（通常是 `$GOPATH/bin` 或 `$HOME/go/bin`）已添加到系统的 `PATH` 环境变量中，以便您可以从任何位置运行 `genx`。

## 运行 `genx`

安装后，您可以从 Go 项目的根目录运行 `genx`。默认情况下，`genx` 将处理当前目录及其所有子目录中的 Go 文件（相当于 `./...`）。

要执行 `genx`，只需运行以下命令：

```bash
genx
```

这将根据 Go 文件中的 `genx` 注解以及 `genx.yaml`（如果存在）中的配置触发代码生成过程。

## 配置 (`genx.yaml`)

`genx` 会在运行它的当前目录或任何父目录中查找名为 `genx.yaml` 的配置文件。此文件允许您定义影响代码生成过程的全局设置。

`genx.yaml` 支持以下顶级键：

### `imports`

- **目的**：定义全局导入别名或添加生成的代码可能需要但在使用 `genx` 指令的用户源文件中未显式存在的导入。当生成的代码依赖于原始源文件未直接导入的包时，这非常有用。
- **语法**：
  ```yaml
  imports:
    - alias: "alias_name" # 可选：指定导入的别名
      path: "module/path"  # 必需：完整的模块路径
    - path: "another/module/path" # 无别名导入的示例
  ```

### `preloads`

- **目的**：解析那些未被正在处理的文件直接导入，但包含 `genx` 插件可能需要理解的类型的附加包。这对于插件需要访问不属于注解文件直接依赖关系图的包（例如，基本类型、通过注释传递的配置中使用的类型）中的类型信息的情况至关重要。
- **语法**：
  ```yaml
  preloads:
    - alias: "preload_alias" # 可选：预加载包的别名，可在生成的代码中使用
      path: "module/path/to/preload" # 必需：要预加载的完整模块路径
    - path: "another/module/to/preload" # 无别名预加载的示例
  ```

### 示例 `genx.yaml`

这是一个简单的 `genx.yaml` 文件示例：

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

## TUI 输出

当您运行 `genx` 时，它会提供一个终端用户界面 (TUI) 来显示代码生成过程的进度和结果。此界面可帮助您跟踪：

- 哪些包正在被解析和处理。
- 哪些 `genx` 插件正在针对特定类型或注解执行。
- 哪些文件正在被生成或修改。
- 每个操作的状态（例如，成功、失败，或者如果生成的文件已存在且与新内容匹配）。

这种实时反馈对于理解 `genx` 正在做什么以及快速识别生成过程中的任何问题非常有用。

## 注解语法 (魔术注释)

`genx` 是由 Go 源代码中特殊格式的注释触发的。这些注释，通常被称为“魔术注释”或注解，由 `genx` 根据特定语法进行解析，以理解要生成什么代码以及如何生成。

`genx` 可识别三种主要的注释样式：

### 1. `@directive` 风格 (主要)

这是调用 `genx` 插件和定义代码生成规则最常用且最灵活的方式。

-   **通用形式：**
    -   `@FuncName(arg1, "stringArg", key="value", ...)`
    -   `@FuncName()` (无参数)
    -   `@FuncName` (无参数，无括号)

-   **`@FuncName` (指令名称)：**
    -   这是一个 `ATID` 标记，意味着它必须以 `@` 符号开头。
    -   其后可以跟字母数字字符、下划线 (`_`) 或连字符 (`-`)。
    -   示例：`@gormq`, `@crud.Find`, `@my-custom-plugin.generate`。

-   **参数：**
    -   参数括在括号 `()` 中。
    -   它们可以是位置参数或命名参数。
    -   **位置参数：** 这些是简单的字符串字面量（单引号或双引号）。
        -   示例：`@myplugin("value1", "value2")`
    -   **命名参数：** 使用格式 `key = "value"` 或 `key = 'value'`。`key` 是一个标识符（字母数字、下划线），`value` 是一个字符串字面量。
        -   示例：`@myplugin(name="user", type="admin")`
    -   位置参数和命名参数可以混合使用，尽管按照惯例，位置参数通常在前。
    -   所有指令都是基于行的，并且必须在注释块内以换行符结束。

-   **示例：**
    ```go
    // @myDirective("positional_arg", count=10, name="example")
    // @anotherDirective()
    // @simpleDirective
    type MyType struct {
        // ...
    }
    ```

### 2. `@FieldDirective` 风格 (旧式/备选)

这种风格也可用，并且可能在较旧的代码库中遇到或用于特定用例，例如字段级注解。

-   **通用形式：**
    -   `@FieldDirective arg1 "arg2 with spaces" arg3`

-   **`@FieldDirective` (指令名称)：**
    -   这也以 `ATID` 标记开头（例如，`@ColumnOptions`, `@Validate`）。

-   **参数：**
    -   参数是出现在指令名称之后的以空格分隔的标记。
    -   如果参数需要包含空格，则必须用单引号或双引号括起来。
    -   这些参数通常由相应的插件作为位置参数处理。
    -   指令在注释块内以换行符终止。

-   **示例：**
    ```go
    type User struct {
        // @ColumnOptions Name "user_name" PrimaryKey
        // @Validate Required MaxLength="50"
        Name string
    }
    ```

### 3. `INSET` 风格 (隐式内容行)

如果 Go 文档注释块中的某一行*不*以 `@` 符号开头，`genx` 会将其视为 `INSET` 指令。

-   这样一行（在检查初始字符确保它不是 `@` 指令之后）的全部内容将作为单个字符串参数传递给指定处理 `INSET` 指令的插件。
-   这种风格对于直接在注释中向插件传递多行文本块、配置或模板特别有用。
-   处理 `INSET` 指令的具体插件决定了如何解释和使用此文本。

-   **示例：**
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
    *(注意：`INSET` 的具体用法取决于特定的插件。如果发现像 `crud` 或 `gormq` 这样的插件将 `INSET` 用于特定目的，则可以更新此示例以使其更具体。)*

### 重要说明:

-   除非特定插件的文档另有说明，否则所有注解指令名称和参数键都是**区分大小写**的。
-   确保这些注解放置在标准的 Go 注释中（例如，行注释使用 `//`，或在 `/* ... */` 块内）。`genx` 处理在这些注释中找到的内容。

## 插件

### `@log`

#### 目的
为 Go 接口生成日志中间件。此中间件使用 `slog` 库记录方法调用、参数、执行时间和错误。它设计为与 `github.com/samber/do/v2` 进行依赖注入以进行日志记录器配置。

#### 目标
Go 接口。

#### 指令
`@log`

#### 参数
无。

#### 用法

要使用 `@log` 插件，请使用 `// @log` 注解您的服务接口定义。

**示例：**

假设您在 `myservice/service.go` 中定义了一个服务接口：

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

#### 生成的代码 (`myservice/logging.go`)

`genx` 将在 `myservice` 包中生成一个 `logging.go` 文件。此文件将包含：

1.  一个 `logging` 结构体，它包装您的 `Service` 接口并包含一个 `*slog.Logger`。
2.  `logging` 结构体上的方法，对应于您 `Service` 接口中的每个方法。这些方法：
    *   记录输入参数（复杂类型进行 JSON 编码）。
    *   调用包装的 `next` 服务上的实际方法。
    *   记录执行持续时间 (`took`) 和任何返回的错误。成功调用使用 `slog.LevelInfo`，如果返回错误则使用 `slog.LevelError`。
3.  一个 `NewLogging` 构造函数。此函数设计为 `do.Middleware`，并期望一个 `do.Injector` 来解析 `*slog.Logger` 和 `next` 服务实例。

生成的 `logging.go` 可能如下所示的简化片段：

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

#### 设置 (使用 `samber/do/v2`)

通常，在使用 `do` 设置服务时，您会注册此日志中间件：

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
**注意：** 在 `samber/do/v2` 中应用中间件的确切方式可能有所不同（例如，使用 `do.Decorate` 或显式提供和调用中间件链）。生成的 `NewLogging` 对于这些方法足够灵活。`do.Decorate` 方法通常是最简洁的。

### `@trace`

#### 目的
使用 OpenTelemetry 为 Go 接口生成追踪中间件。此中间件为每个方法调用创建一个新的 span，将方法参数记录为 span 属性，捕获错误并相应地设置 span 状态。它设计为与 `github.com/samber/do/v2` 一起使用，用于 `TracerProvider` 的依赖注入。

#### 目标
Go 接口。

#### 指令
`@trace`

#### 参数
无。

#### 用法

使用 `// @trace` 注解您的服务接口定义。

**示例：**

给定 `myservice/service.go` 中的服务接口：

```go
package myservice

import "context"

// @trace
type Service interface {
    CreateItem(ctx context.Context, itemID string, itemName string) (err error)
    GetItem(ctx context.Context, itemID string) (itemName string, err error)
}
```

#### 生成的代码 (`myservice/trace.go`)

`genx` 将在 `myservice` 包中生成一个 `trace.go` 文件。此文件将包括：

1.  一个 `tracing` 结构体，它包装您的 `Service` 接口并持有一个 OpenTelemetry `*sdktrace.TracerProvider`。
2.  `tracing` 结构体上针对 `Service` 中每个方法的方法。这些方法：
    *   启动一个 OpenTelemetry span（例如，以方法名称命名，如 "CreateItem"）。
    *   将方法参数设置为 span 属性（参数被 JSON 编组到 `params` 属性中）。
    *   如果包装的方法返回错误，则会在 span 上记录该错误 (`span.RecordError()`) 并将 span 状态设置为 `codes.Error`。
    *   span 在方法完成时结束。
3.  一个 `NewTracing` 构造函数。此函数是一个 `do.Middleware`，并期望一个 `do.Injector` 来解析 `*sdktrace.TracerProvider` 和 `next` 服务实例。

`trace.go` 的说明性片段：

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

#### 设置 (使用 `samber/do/v2` 和 OpenTelemetry SDK)

您将类似于日志中间件的方式注册此追踪中间件，确保也提供了 OpenTelemetry `TracerProvider`。

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
**注意：** 请记住为您的环境适当配置 OpenTelemetry SDK（例如，导出器、采样器、传播器）。该示例使用一个简单的 stdout 导出器。

### `@otel`

#### 目的
为 Go 接口生成一个全面的 OpenTelemetry 中间件。此中间件结合了追踪、日志（使用 `slog`）和指标（状态计数器、总计数器、持续时间直方图）用于每个方法调用。它依赖于 `github.com/samber/do/v2` 来注入必要的 OpenTelemetry 和日志组件。

#### 目标
Go 接口。

#### 指令
`@otel`

#### 参数
无。

#### 用法

使用 `// @otel` 注解您的服务接口定义。

**示例：**

考虑 `custompkg/service.go` 中的一个接口：

```go
package custompkg

import "context"

// @otel
type StringService interface {
    Uppercase(ctx context.Context, s string) (upper string, err error)
    Count(ctx context.Context, s string) (length int, err error)
}
```

#### 生成的代码 (`custompkg/otel.go`)

`genx` 将在 `custompkg` 包中创建一个 `otel.go` 文件。此文件包含：

1.  一个 `otel` 结构体，它包装您的 `StringService` 并持有 `trace.Tracer`、`*slog.Logger` 和各种 `metric` 工具的实例。它还存储 `pkgName`。
2.  `otel` 结构体上针对接口中每个方法的方法。这些方法：
    *   **追踪：** 启动一个 OpenTelemetry span（命名为 `pkgName + "." + methodName`）。参数作为 JSON 编码的属性添加。错误被记录，span 状态被更新。
    *   **日志：** 使用 `slog` 记录方法调用详细信息（包、方法、参数、持续时间、错误）。
    *   **指标：**
        *   `serviceStatus` (Int64Counter)：使用维度 `service` (pkgName)、`method` 和 `status` ("success" 或 "fail") 递增。
        *   `serviceCounter` (Int64Counter)：使用维度 `service` (pkgName) 和 `method` 递增。
        *   `serviceDuration` (Float64Histogram)：使用维度 `service` (pkgName) 和 `method` 记录方法执行时间（毫秒）。
3.  一个 `NewOtel` 构造函数。这是一个 `do.Middleware`，期望一个 `do.Injector` 来解析追踪器、日志记录器和命名的指标工具 (`serviceStatus`, `serviceCounter`, `serviceDuration`)。

生成的 `otel.go` 的说明性片段：

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

#### 设置 (使用 `samber/do/v2`, OpenTelemetry SDK)

您需要提供所有依赖项：`trace.Tracer`、`*slog.Logger` 和命名的指标。

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
	sloglib "log/slog" // Renamed to avoid conflict
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
	do.Provide(injector, func(i do.Injector) (*sloglib.Logger, error) { // Used sloglib
		return sloglib.New(sloglib.NewJSONHandler(os.Stdout, nil)), nil
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
**重要考虑因素：**
*   确保您在 `do.ProvideNamedValue` 调用中提供的指标名称（`serviceStatus`、`serviceCounter`、`serviceDuration`）与生成的 `NewOtel` 函数期望的名称完全匹配。生成的代码使用 `"serviceStatus"`、`"serviceCounter"` 和 `"serviceDuration"`。
*   OpenTelemetry SDK 的设置可能很复杂。该示例提供了一个基本的 stdout 设置。有关生产就绪配置，请参阅 OpenTelemetry 文档。
*   指标和日志中使用的 `pkgName` 是从被检测接口的包中自动派生的。

### `@enum`

#### 目的
为基于整数的枚举类型生成一组辅助方法和常量，使其更健壮且易于使用。它会创建字符串表示、备注和解析函数。

#### 目标
Go 类型规范（通常是为枚举别名的 `int` 类型）。

#### 指令
`@enum`

#### 参数
*   `@enum` 指令接受一个或多个位置字符串参数。
*   每个参数定义一个枚举成员，格式应为 `"Key:Remark"` 或仅 `"Key"`。
    *   `Key`：枚举成员的标识符（例如，`Active`、`DefaultUser`）。这用于生成常量名称，如 `Status_Active`。
    *   `Remark`：枚举成员的描述性字符串（例如，`User is currently active`）。可通过 `Remark()` 方法访问。如果仅提供 `"Key"`，则备注为空。

#### 用法

为您的枚举定义一个整数类型，并使用 `// @enum(...)` 和成员定义来注解其类型规范。

**示例：**

假设您有 `types/status.go`：

```go
package types

// @enum("Active:User is active", "Pending:Awaiting confirmation", "Disabled:User account is disabled", "Deleted")
type Status int // The base type is typically int
```

#### 生成的代码 (`types/enum.go`)

`genx` 将在 `types` 包（或 `Status` 定义的任何位置）中生成一个 `enum.go` 文件。此文件将包含：

1.  **枚举值常量：** 一个基于 iota 的 `const` 块，用于您的枚举成员。
    ```go
    const (
        _ = iota // Start from 0, actual values will be 1, 2, 3...
        Status_Active
        Status_Pending
        Status_Disabled
        Status_Deleted
    )
    ```
2.  **别名和备注常量：** 每个枚举成员的字符串“键”和“备注”的常量。
    ```go
    const (
        STATUS_ACTIVE_ALIAS   = "Active"
        STATUS_ACTIVE_REMARK  = "User is active"
        STATUS_PENDING_ALIAS  = "Pending"
        STATUS_PENDING_REMARK = "Awaiting confirmation"
        // ... 对于 Disabled、Deleted 也是如此（Deleted 的备注将为空）
    )
    ```
3.  **`String()` 方法：** 返回枚举值的字符串键。
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
4.  **`Remark()` 方法：** 返回枚举值的描述性备注。
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
5.  **`ParseStatus()` 函数：** 将整数转换为您的枚举类型 (`Status`)。
    ```go
    func ParseStatus(id int) (Status, error) {
        if x, ok := _StatusValue[id]; ok { // _StatusValue is a generated map
            return x, nil
        }
        return 0, fmt.Errorf("unknown enum value: %d", id) // Corrected error formatting
    }
    ```
    （还会生成 `var _StatusValue = map[int]Status{...}`，用于将整数值映射到枚举常量。）


#### 使用生成的枚举

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
**注意：** 生成的代码对枚举常量名称使用首字母大写命名法（例如，`Status_Active`），对别名/备注常量使用全大写命名法（例如，`STATUS_ACTIVE_ALIAS`）。整数值从 1 开始（由于 `_ = iota` 后跟成员）。

### `@gq` (GORM 查询作用域)

#### 目的
根据用户定义的查询结构体的字段动态生成 GORM 查询作用域函数。这允许构建可重用且类型安全的查询逻辑。

#### 目标
Go 结构体。

#### 生成的文件
`gorm_scope.go`，与注解的结构体在同一个包中。

#### 主要指令 (在查询结构体本身上)

`// @gq <model.GormModel>`
*   此指令将结构体标记为 GORM 查询定义。
*   `<model.GormModel>`：**必需。** 查询将针对的 GORM 模型结构体的完全限定名称（例如，`myproject/model.User`、`model.Product`）。

#### 字段级指令 (在查询结构体的字段上)

这些指令自定义查询结构体中的每个字段如何转换为 GORM 查询条件。

*   `// @gq-column <db_col_1> [<db_col_2> ...]`
    *   指定字段的数据库列名。
    *   如果提供了多个列名，则会为该字段生成一个 `OR` 条件，使用其值与每个指定的列进行比较。
    *   如果省略，`genx` 会尝试读取 `gorm:"column:..."` 标签。如果该标签也缺失，则假定使用 GORM 对字段名的默认命名策略。
    *   **示例：** `// @gq-column name_alias user_name` 生成 `WHERE name_alias = ? OR user_name = ?`。

*   `// @gq-op <operator>`
    *   定义比较运算符。如果未指定，则默认为 `=`。
    *   **支持的运算符：**
        *   `=`：等于。（大多数类型的默认值；对于切片/数组，默认为 `in`）。
        *   `!=`：不等于。
        *   `>`：大于。
        *   `>=`：大于或等于。
        *   `<`：小于。
        *   `<=`：小于或等于。
        *   `><`：BETWEEN（期望字段是包含两个值的切片或数组，例如 `PriceRange []int`）。生成 `column BETWEEN ? AND ?`。
        *   `!><`：NOT BETWEEN。
        *   `like`：LIKE 运算符（例如 `column LIKE ?`）。生成的代码会自动用 `%` 包装值（例如 `"%value%"`）。
        *   `in`：IN 运算符（例如 `column IN (?)`）。如果 `@gq-op` 为 `=` 或省略，则为切片/数组字段的默认值。
        *   `!in`：NOT IN 运算符。
        *   `null`：IS NULL（查询结构体中该字段的值被忽略）。生成 `column IS NULL`。
        *   `!null`：IS NOT NULL（字段的值被忽略）。生成 `column IS NOT NULL`。

*   `// @gq-clause <clause_type>`
    *   指定 GORM 子句类型。默认为 `Where`。
    *   有效值：`Where`、`Or`、`Not`。（注意：`Or` 和 `Not` 可能与字段值或多个 `@gq-column` 条目有特定的交互模式。`Where` 是最常见的）。

*   `// @gq-sub <foreign_key> <referenced_key>`
    *   用于创建子查询。注解此内容的字段应为指向另一个也具有生成的 `Scope` 方法（即也用 `@gq` 注解）的结构体的指针。
    *   `<foreign_key>`：当前模型（由 `@gq <model.GormModel>` 定义）中将与子查询结果进行比较的列名。
    *   `<referenced_key>`：在子查询中从子查询模型中选择的列名。
    *   **示例：** 如果 `OrderQuery` 有一个字段 `User *UserQuery // @gq-sub user_id id`，它会生成 `WHERE user_id IN (SELECT id FROM users WHERE ...conditions_from_UserQuery...)`。

*   `// @gq-group`
    *   标记一个字段（必须是结构体类型，通常是嵌入式或直接字段），其条件应被分组。该字段的结构体类型本身应该是一个有效的查询结构体（可能带有其自己的 `@gq-*` 注解，但不是主要的 `@gq <model>` 注解）。
    *   生成一个分组条件，如 `db.Where(subScope)`，其中 `subScope` 是嵌套结构体生成的查询逻辑的结果。

*   `// @gq-struct [<value1> <value2> ...]`
    *   将字段（应为结构体或指向结构体的指针）直接传递给 GORM 查询方法，通常是 `Where`。
    *   如果提供了 `<value1>`、`<value2>`，它们将作为附加参数传递给 GORM 方法，通常用于指定要查询结构体的哪些字段（如果它本身不是 GORM 模型）。示例：`db.Where(&q.MyStructField, "field1", "field2")`。
    *   如果没有提供值，则通常为 `db.Where(&q.MyStructField)`。

#### 生成的代码 (`gorm_scope.go`)

*   生成一个方法 `func (q *QueryStructName) Scope(db *gorm.DB) *gorm.DB`。
*   此 `Scope` 方法：
    *   设置 GORM 模型：`db = db.Model(&<model.GormModel>{})`。
    *   迭代 `QueryStructName` 的字段。
    *   根据字段值及其 `@gq-*` 注解构建 GORM 查询条件。
    *   **零值处理：** 通常会跳过具有零值的字段的条件：
        *   `nil` 指针。
        *   空字符串 (`""`)。
        *   整数类型的 `0`。
        *   空切片 (`len() == 0`)。
    *   支持嵌入式结构体：嵌入式结构体的字段被处理，就像它们是父查询结构体的一部分一样。也处理指针嵌入式结构体。
    *   **特殊切片处理：** 如果字段是结构体切片（例如 `RelatedItems []ItemQuery`），并且 `ItemQuery` 具有映射到列的字段，则可以生成类似 `WHERE (item_col_a, item_col_b) IN ((val1_a, val1_b), (val2_a, val2_b), ... )` 的查询。

**示例：**

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

生成的 `queries/gorm_scope.go` (简化概念)：
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

#### 使用作用域

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
**注意：** `plugs/gormq/` 中的 `design.md` 文件为此插件提供了更多示例和上下文。强烈建议查看它以获得更深入的理解。

### `@crud`

#### 目的
为 CRUD (创建、读取、更新、删除) 操作生成样板代码。它支持两种不同的模式：生成基于 GORM 的数据库服务或生成 HTTP 服务结构和类型。

#### 目标
Go 结构体（这些结构体充当 CRUD 生成的配置/标记）。

#### 指令
`@crud`

#### `@crud` 指令的参数 (应用于配置结构体)

*   `type="<mode>"`：**必需。** 指定生成模式。
    *   `"gorm"`：生成基于 GORM 的 CRUD 服务实现。
    *   `"http"`：生成 HTTP 请求/响应类型和基础 HTTP 服务结构。
*   `model="<pkg.ModelName>"`：**必需。** CRUD 操作将针对的 GORM 模型结构体的完全限定名称（例如 `myproject/model.User`）。
*   `idName="<IDFieldName>"`：**必需。** GORM 模型中主键字段的名称（例如 `ID`、`UUID`）。这指的是 Go 结构体字段名称。
*   `idType="<GoType>"`：**必需。** 主键字段的 Go 类型（例如 `uint`、`string`、`uuid.UUID`）。
*   `preload="<field1,field2.subfield>"`：*(仅由 `type="http"` 使用)* 一个逗号分隔的字符串，指定在为 HTTP 响应生成 `GetResponse` 结构体时应深度复制 GORM 模型的哪些字段。支持用于嵌套字段的点表示法（例如 `User,Order.OrderItems`）。

#### 用法

定义一个空结构体，并使用 `@crud` 及其必需参数对其进行注解。此结构体的名称本身对生成的代码的命名影响不像 `model` 参数那么大。

**示例配置结构体：**

```go
package services

// @crud(type="gorm", model="myproject/model.Article", idName="ID", idType="uint")
type ArticleGormCRUDConfig struct{}

// @crud(type="http", model="myproject/model.User", idName="ID", idType="int", preload="Profile,Addresses")
type UserHttpCRUDConfig struct{}
```

---

#### 模式：`type="gorm"`

*   **目的：** 为基于 GORM 的 CRUD 操作生成基础服务层。
*   **生成的文件：**
    *   `crud_base_service.go`：包含 GORM CRUD 服务实现（例如，Create、GetByID、Update、Delete、List 的方法）。
    *   （注意：`crud_gorm_types.go` 的生成在插件源代码中似乎已被注释掉，因此默认情况下可能不会生成。）

*   **生成的代码 (`crud_base_service.go` - 概念性)：**
    *   一个接口（例如，如果 `model` 是 `model.Article`，则为 `ArticleGormCrudBaseImpl`）。
    *   一个实现此接口的结构体，它接受一个 `*gorm.DB` 实例。
    *   **标准 CRUD 方法：**
        *   `Create(m *model.Article) error`
        *   `GetByID(id uint) (*model.Article, error)`
        *   `Update(id uint, m *model.Article) error`
        *   `Delete(id uint) error`
        *   `FindOne(query interface{}, args ...interface{}) (*model.Article, error)`
        *   `Find(query interface{}, args ...interface{}) ([]*model.Article, int64, error)` (返回项目和总计数)
    *   这些方法使用 GORM 进行数据库交互（例如，`db.Create()`、`db.First()`、`db.Save()`、`db.Delete()`、`db.Where().Find()`）。
    *   确切的方法签名和功能由 `crud_gorm.tmpl` 模板定义。

**使用 GORM CRUD 服务：**

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
*(子任务期间的自我修正：“使用”GORM 的部分很难详细说明，因为没有看到 `crud_gorm.tmpl` 中确切的构造函数/结构体名称。文档应说明这一点，并可能在从检查生成的文件或模板中知道这些名称后，展示如何实例化生成的服务。)*

---

#### 模式：`type="http"`

*   **目的：** 为 HTTP 请求/响应负载生成类型，并为 HTTP 服务生成基础结构。
*   **生成的文件：**
    *   `types.go`：包含为 HTTP 请求和响应主体生成的 Go 结构体（例如 `GetResponse`、`CreateRequest`、`UpdateBody`）。
    *   `crud_http_service.go`：包含基础 HTTP 服务接口和可能的存根实现。

*   **生成的代码 (`types.go` - 概念性)：**
    *   `GetResponse` 结构体：根据 `@crud` 注解中指定的 GORM `model` 创建。`preload` 参数控制包含哪些字段（包括通过点表示法如 `User.Profile` 的嵌套字段）。这对于塑造 API 响应很有用。
    *   `CreateRequest` 结构体：用于解组创建请求负载。
    *   `UpdateBody` 结构体：用于解组更新请求负载。
    *   这些结构体通过从 GORM 模型复制字段生成，可能会排除某些字段（例如，标记为 `serializer:"-"` 的字段），并且对于 `GetResponse` 仅包含 `preload` 中指定的字段。

*   **生成的代码 (`crud_http_service.go` - 概念性)：**
    *   一个接口（例如，如果 `model` 是 `model.User`，则为 `UserHttpCrudBaseImpl`）。
    *   **潜在的 HTTP 服务方法：**
        *   `Create(c context.Context, req *CreateRequest) (*GetResponse, error)`
        *   `Get(c context.Context, id int) (*GetResponse, error)`
        *   `Update(c context.Context, id int, req *UpdateBody) (*GetResponse, error)`
        *   `Delete(c context.Context, id int) error`
        *   `List(c context.Context, queryParams ListRequest) ([]*GetResponse, int, error)` (`ListRequest` 也会是为分页/过滤参数生成的类型）
    *   实际的方法签名由 `crud_http.tmpl` 和 `crud_http_types.tmpl` 模板定义。

**使用 HTTP CRUD 服务/类型：**

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
*(子任务期间的自我修正：与 GORM 类似，“使用”HTTP 的部分在没有来自模板的确切生成的接口/结构体名称和方法签名的情况下是概念性的。文档应强调用户需要检查生成的文件或模板以了解这些具体细节。)*

**注意：** `@crud` 插件是模板驱动的。有关生成的方法名称、结构体字段和构造函数的精确详细信息，请参阅 `genx` `static/template` 目录中的 `crud_gorm.tmpl`、`crud_http.tmpl` 和 `crud_http_types.tmpl` 的内容。

### `@kit-http-client`

#### 目的
基于 Go 接口定义生成 Go kit HTTP 客户端实现。此客户端包括服务发现、负载均衡（轮询）和重试等功能。

#### 目标
Go 接口。

#### 生成的文件
`kit_http_client.go`，与接口在同一个包中。

#### 接口级指令

*   `// @basePath <base_path_string>`
    *   **可选。** 放置在接口定义上。
    *   为此接口方法中定义的所有 HTTP URL 指定一个通用的基本路径前缀（例如 `/api/v2`）。

#### 方法级指令 (在注解接口的方法上)

*   `// @kit-http <http_url_path> <http_method>`
    *   **必需。** 定义方法的端点。
    *   `<http_url_path>`：特定的 URL 路径（例如 `/users/{userID}`、`/products`）。路径参数用 `{param_name}` 表示。
    *   `<http_method>`：HTTP 方法（例如 `GET`、`POST`、`PUT`、`DELETE`）。

*   `// @kit-http-request <RequestStructName> [<send_entire_struct_as_body>]`
    *   **必需。** 指定方法的请求对象。
    *   `<RequestStructName>`：封装此方法所有输入参数的 Go 结构体的名称。
    *   `[<send_entire_struct_as_body>]`：**可选。** 如果提供了非空字符串（例如 "true"、"body"），则整个 `<RequestStructName>` 实例将作为 JSON 请求体进行编组。否则，`<RequestStructName>` 中的特定字段必须标记有 `param:"body,..."` 才能构成请求体。

#### 字段级结构体标签 (在 `<RequestStructName>` 内)

请求结构体（由 `@kit-http-request` 指定）中的字段使用 `param` 结构体标签将其映射到 HTTP 请求组件：

*   ``param:"<type>,<name>"``
    *   `<type>`：定义参数类型。
        *   `path`：将字段映射到路径参数。`<name>` 必须与 `@kit-http` URL 路径中的 `{param_name}` 匹配。
        *   `query`：将字段映射到 URL 查询参数。`<name>` 是查询参数键。
        *   `header`：将字段映射到 HTTP 请求头。`<name>` 是头名称。
        *   `body`：指定该字段是 JSON 请求体的一部分。如果在 `@kit-http-request` 中*未*设置 `<send_entire_struct_as_body>`，则生成的客户端通常期望一个字段标记为 `param:"body,..."` 作为请求体的来源。如果*已*设置 `<send_entire_struct_as_body>`，则此标记可能是多余的，或者如果结构体的字段名称不同，则用于 JSON 中的特定字段命名。生成器代码暗示，如果 `RequestBody`（派生自 `[<send_entire_struct_as_body>]`）为 false，它将使用找到的第一个标记为 `param:"body,..."` 的字段作为请求体。
    *   `<name>`：路径变量的名称、查询参数键、头名称或 JSON 字段名称（尽管对于使用整个结构体发送的 `body`，JSON 字段名称通常来自结构体字段名称或 `json` 标签）。

#### 生成的代码 (`kit_http_client.go`)

*   **`HttpClientImpl` 接口：** 一个镜像注解接口的新接口，但其方法签名适用于客户端：
    `MethodName(ctx context.Context, req RequestStructName, option *Option) (res ResponseStructName, err error)`
*   **`HttpClientService` 结构体：** 实现 `HttpClientImpl`。
    *   处理 URL 构建（连接 `@basePath` 和特定于方法的路径）。
    *   管理路径参数替换、查询字符串形成和头注入。
    *   将请求体编码为 JSON（通常使用 `kithttp.EncodeJSONRequest`）。
    *   解码 JSON 响应。
    *   与 Go kit 的服务发现 (`sd.Instancer`)、负载均衡 (`lb.NewRoundRobin`) 和重试机制集成。
*   **`Option` 结构体：** 允许按调用或全局自定义客户端行为：
    *   `PrePath`：覆盖 `@basePath` 或 `HttpClientService` 的全局 `PrePath`。
    *   `Logger`：用于 Go kit 的 `log.Logger`。
    *   `Instancer`：用于服务发现的 `sd.Instancer`。
    *   `RetryMax`、`RetryTimeout`：用于重试配置。
    *   `EndpointOpts`：用于端点创建的 `[]sd.EndpointerOption`。
    *   `ClientOpts`：用于 `http.Client` 自定义的 `[]kithttp.ClientOption`（例如，通过 `kithttp.ClientBefore` 添加头）。
    *   `Encode`：自定义 `kithttp.EncodeRequestFunc`。
    *   `Decode`：用于 `kithttp.DecodeResponseFunc` 的自定义工厂（例如，`func(i interface{}) func(ctx context.Context, res *http.Response) (response interface{}, err error)`，其中 `i` 是指向响应结构体的指针）。
*   **验证：** 使用 `github.com/asaskevich/govalidator` 在发送前验证请求结构体。

**示例：**

接口定义 (`myservice/client.go`)：
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

生成的 `myservice/kit_http_client.go` (概念性)：
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

#### 使用生成的客户端

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
**注意：** 客户端需要为服务发现 (`sd.Instancer`) 进行适当设置，并且在高级场景中可能需要自定义 `ClientOption` 或 `Encode`/`Decode` 函数。“使用”示例对于客户端部分是概念性的，因为确切的构造函数名称需要从生成的代码中获知。

### `@temporal`

#### 目的
通过生成代码将其方法注册为 Temporal activity 并提供可从工作流调用的包装方法来执行这些 activity，从而将 Go 接口（服务）与 Temporal.io 集成。

#### 目标
Go 接口。

#### 生成的文件
`temporal.go`，与注解的接口在同一个包中。

#### 接口级指令

*   `// @temporal`
    *   放置在接口定义上。
    *   作为标记，指示应为此接口处理 Temporal 集成。
    *   不接受任何参数。

#### 方法级指令 (在 `@temporal` 注解接口的方法内)

*   `// @temporal-activity`
    *   在旨在作为 Temporal activity 的方法上**必需**。
    *   标记该方法以注册到 Temporal worker。
    *   在 `Temporal` 结构体上生成相应的包装方法，允许从 Temporal 工作流调用该 activity。

#### 生成的代码 (`temporal.go`)

1.  **`Temporal` 结构体：**
    ```go
    type Temporal struct {
        next Service // Where 'Service' is your annotated interface type
        w    worker.Worker
    }
    ```

2.  **`InitTemporal` 函数 (依赖注入)：**
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
    *   此函数旨在与 `github.com/samber/do/v2` 一起用于依赖注入。
    *   它从 DI 容器中解析 `worker.Worker` 和您的服务实现 (`Service`)。
    *   它将服务中所有使用 `@temporal-activity` 注解的方法注册到 Temporal worker。

3.  **可从工作流调用的方法 (在 `*Temporal` 结构体上)：**
    *   对于接口中每个使用 `// @temporal-activity` 注解的方法，都会在 `*Temporal` 结构体上生成一个相应的包装方法。
    *   **签名更改：**
        *   这些包装方法的第一个参数变为 `ctx workflow.Context` (来自 `go.temporal.io/sdk/workflow`)。
        *   其他参数和返回类型与原始方法匹配。
    *   **实现：** 这些包装器使用 `workflow.ExecuteActivity` 来调用实际注册的 activity。

**示例：**

接口定义 (`example/transfer_service.go`)：
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

生成的 `example/temporal.go` (概念性)：
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

#### 使用 Temporal

1.  **服务实现：** 实现您的 `TransferService` 接口。
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

2.  **Temporal Worker 设置与 DI：**
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

3.  **在您的工作流定义中：**
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
**注意：** 生成的 `(t *Temporal) ProcessTransfer` 方法是一个方便的包装器，它内部调用原始服务方法（例如 `t.next.ProcessTransfer`）上的 `workflow.ExecuteActivity`。在定义工作流时，您可以直接调用此生成的包装方法。关键在于 `InitTemporal` 将*实际的*服务方法（例如 `next.ProcessTransfer`）注册到 Temporal worker。

### `@copy` (结构体复制)

#### 目的
生成类型安全的函数以在两个 Go 结构体（或类型）之间复制数据。它处理嵌套结构、切片、映射、指针，并允许通过注解进行自定义字段映射。

#### 目标
Go 函数调用表达式。注解函数的签名定义了源类型和目标类型。

#### 生成的文件
`copy.go` (在注解函数的包中)。

#### 触发机制

您通过定义一个函数（其主体通常为空或被忽略）并使用 `// @copy` 对其进行注解来触发 `@copy` 插件。函数的签名决定了复制操作：

```go
package myconverters

// @copy
func ConvertUserToUserDTO(dto *UserDTO, user model.User) {
    // This function's body is not executed by genx.
    // Its signature defines the types for the generated copy function.
}

// UserDTO and model.User are your destination and source struct types.
```

*   第一个参数是**目标**（例如 `dto *UserDTO`）。它**必须是指针类型**。
*   第二个参数是**源**（例如 `user model.User`）。

#### 生成的代码 (`copy.go`)

对于上面的 `ConvertUserToUserDTO` 示例，`genx` 将生成：

1.  **原始函数 (已填充)：**
    ```go
    package myconverters

    func ConvertUserToUserDTO(dto *UserDTO, user model.User) {
        ConvertUserToUserDTOCopy{}.Copy(dto, user) // Delegates to the Copy method
        // If ConvertUserToUserDTO had return values, they would be returned here.
    }
    ```

2.  **一个 Copier 结构体：**
    ```go
    type ConvertUserToUserDTOCopy struct{} // Named <YourFunctionName>Copy
    ```

3.  **Copier 结构体上的 `Copy` 方法：**
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
    *   此方法包含将字段从源复制到目标的核心逻辑。
    *   它智能地处理各种类型：
        *   对于兼容的基本类型和命名类型进行直接赋值。
        *   对源指针进行 Nil 检查。
        *   如果目标指针为 nil 且源不为 nil，则为目标指针分配内存（例如 `dest.Nested = new(NestedDTO)`）。
        *   通过为其元素生成复制逻辑来对复杂类型的切片和映射进行深层复制。
        *   递归生成：如果遇到需要复制的嵌套结构体（例如 `User.Address` 到 `UserDTO.AddressDTO`），它将尝试为这些特定类型生成或使用另一个 `Copy` 方法（例如 `AddressToAddressDTOCopy`）。

#### 字段映射自定义 (源或目标结构体字段上的注解)

您可以在源或目标结构体的字段上使用注解来控制它们的映射方式。这些注解通常放在字段上方的注释中。

*   `// @copy-prefix <prefix_to_remove>`：（通常在源字段上）如果源字段名称具有应在与目标字段名称匹配之前删除的前缀（例如 `Source_FieldA` 对应 `FieldA`）。
*   `// @copy-name <alternative_name>`：（通常在源字段上）使用此 `alternative_name` 与目标字段名称进行匹配，而不是实际的 Go 字段名称。
*   `// @copy-must`：如果此注解存在于源字段上，并且在应用名称/前缀规则后在目标中找不到相应的可设置字段，则代码生成将引发 panic。
*   `// @copy-target-path <full.dot.path.to.dest.field>`：（通常在源字段上）将注解的源字段显式映射到目标结构体中的指定完整路径。这提供了最高的匹配优先级。示例：`// @copy-target-path UserDetails.Contact.EmailAddress`。
*   `// @copy-target-name <dest_field_name>`：（通常在源字段上）即使路径不同，也将注解的源字段映射到具有此特定名称的目标字段。
*   `// @copy-target-method <method_name_on_source>`：（通常在源字段上）不是从源直接访问字段，而是在源结构体（或特定的源字段，如果是嵌套结构体）上调用 `<method_name_on_source>()` 来获取值。然后将此值复制到按名称/路径匹配的目标字段。
*   `// @copy-auto-cast`：（布尔值，例如 `@copy-auto-cast`）如果存在，则在直接赋值不可能的情况下，可以启用更宽松的类型转换（例如，在不同的整数类型之间或与字符串之间的转换）。具体的转换行为取决于生成器的实现细节。

#### 匹配策略

插件通过考虑以下因素来尝试匹配源和目标之间的字段：
1.  通过 `@copy-target-path` 或 `@copy-target-name` 进行显式映射。
2.  名称匹配（在对源字段名称应用 `@copy-name` 或 `@copy-prefix` 之后）。
3.  一种“深度查找”算法，即使在不同级别的嵌入式结构体中也尝试匹配字段，从而有效地“扁平化”路径。

**带字段注解的示例：**

源结构体：
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

目标 DTO：
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

生成的 `Copy` 方法 (对于 `ConvertUserToUserDTO` 的概念性示例)：
```go
// ...
dto.ID = user.ID
dto.FullName = user.InternalName // Due to @copy-target-name
dto.Email = user.UserEmail       // Due to @copy-prefix User
dto.Bio = user.Details.Bio       // Matched by path/name
dto.Preferences.Theme = user.Details.ThemePreference // Due to @copy-target-path
// ...
```

此插件在不同结构体类型之间（例如 API DTO、数据库模型和内部领域对象之间）映射时，对于减少样板代码非常有用。

### `@kit` (Go kit HTTP 服务生成)

#### 目的
基于 Go 接口构建 Go kit HTTP 服务。它生成 Go kit 端点、HTTP 传输层（请求/响应编解码函数）和 HTTP 处理程序。

#### 目标
Go 接口。

#### 生成的文件
*   `endpoint.go`：包含注解接口中每个方法的 Go kit 端点定义。
*   `http.go`：包含 HTTP 传输层逻辑，包括请求解码函数、响应编码函数和 Go kit HTTP 服务器处理程序。

#### 接口级指令 (应用于接口定义)

*   `// @basePath <base_path_string>`
    *   **可选。** 指定此接口中方法定义的所有 HTTP URL 的通用基本路径前缀（例如 `/api/v1`、`/catalog`）。
*   `// @tags <comma_separated_tags>`
    *   **可选。** 用于 Swagger/OpenAPI 文档以对端点进行分组（例如 `@tags "Users,Admin"`）。
*   `// @validVersion <version_string>`
    *   **可选。** 指定 API 版本，可能用于生成的文档或路由逻辑中。
*   `// @kit-server-option <OptionName1> [<OptionName2> ...]`
    *   **可选。** 将在创建处理程序时传递给 `kithttp.NewServer()` 的 Go kit HTTP 服务器选项列表（例如 `ErrorEncoder(myErrorHandler)`）。这些值通常是作用域中可用的函数名或完全限定的标识符。
*   `// @swag <true_or_false>`
    *   **可选。** 接口级别的开关，用于为此接口派生的所有端点启用/禁用 Swagger/OpenAPI 文档生成。如果未指定，则默认为 true。单个方法可以覆盖此设置。

#### 方法级指令 (在 `@kit` 注解接口的方法内)

*   `// @kit-http <http_url_path> <http_method>`
    *   对于要公开为 HTTP 端点的方法，此为**必需**。
    *   `<http_url_path>`：此方法特定的 URL 路径。路径参数应使用 Gorilla Mux 语法定义（例如 `/items/{itemID}`、`/documents/{category}/{docId}`）。
    *   `<http_method>`：HTTP 方法（例如 `GET`、`POST`、`PUT`、`DELETE`、`PATCH`）。

*   `// @kit-http-request <RequestStructName> [<is_request_body_bool>]`
    *   如果服务方法具有除 `context.Context` 之外的参数，则此为**必需**。
    *   `<RequestStructName>`：将用于捕获和解码传入 HTTP 请求数据（路径参数、查询参数、标头、请求体）的 Go 结构体的名称。
    *   `[<is_request_body_bool>]`：**可选。** 布尔类型的字符串（"true" 或 "false"，或任何非空字符串表示 true）。
        *   如果为 "true"（或非空），则期望从 HTTP 请求体（通常为 JSON）解码整个 `<RequestStructName>` 实例。
        *   如果为 "false"（或为空/省略），则请求体可能不会直接从 `<RequestStructName>` 整体映射，或者 `<RequestStructName>` 中的特定字段可能使用标签指定来自请求体（尽管此插件主要关注结构体本身作为请求体或来自路径/查询/标头的单个字段）。生成的代码通常会创建单独的 `Decode<MethodName>Request` 函数来处理填充此结构体。

*   `// @swag <true_or_false>`
    *   **可选。** 方法级别的开关，用于为此特定端点启用/禁用 Swagger/OpenAPI 文档生成。覆盖接口级别的 `@swag` 设置。

#### 字段级结构体标签 (在 `<RequestStructName>` 内)

请求结构体（在 `@kit-http-request` 中指定的那些）中的字段使用结构体标签来指定应如何从传入的 HTTP 请求中填充它们：

*   `path:"<name>"`：将字段映射到路径参数。`<name>` 必须与 `@kit-http` URL 路径中定义的占位符名称匹配（例如，如果路径为 `/product/{id}`，则标签为 `path:"id"`）。
*   `query:"<name>"`：将字段映射到 URL 查询参数（例如 `query:"searchTerm"`）。
*   `header:"<name>"`：将字段映射到 HTTP 请求头（例如 `header:"X-API-Key"`）。
*   `json:"<name>"`：用于 JSON 编组/解组的标准 Go 结构体标签。当从 JSON 主体解码请求结构体（或其部分）或将响应结构体编码为 JSON 时使用。
*   `validate:"<rules>"`：指定字段的验证规则，通常使用与 `govalidator` 或 `go-playground/validator` 等库兼容的语法（例如 `validate:"required,min=1,max=100"`）。生成的解码函数通常包含对验证器的调用。
*   `form:"<name>"`：用于 HTML 表单值（例如，来自 `application/x-www-form-urlencoded` 主体）。

#### 生成的代码 (`endpoint.go` 和 `http.go`)

*   **`endpoint.go`：**
    *   为每个使用 `@kit-http` 注解的方法定义 Go kit `endpoint.Endpoint`。
    *   创建包装您的服务方法的 `Make<MethodName>Endpoint` 函数。
    *   通常包含一个对所有生成的端点进行分组的 `Endpoints` 结构体。
*   **`http.go`：**
    *   **请求解码：** 生成 `Decode<MethodName>Request` 函数（例如 `DecodeCreateUserRequest`）。这些函数从 HTTP 请求（路径、查询、标头、主体）中提取数据并填充 `<RequestStructName>`。它们通常包含验证逻辑。
    *   **响应编码：** 生成 `Encode<MethodName>Response` 函数（或通用的 `EncodeResponse`）以将服务方法的结果（和错误）编组为 HTTP 响应，通常为 JSON。
    *   **HTTP 处理程序：** 为每个端点创建 Go kit `kithttp.Handler` 实例，连接解码/编码函数和服务器选项。
    *   一个函数（例如 `NewHTTPHandler` 或 `MakeHandler`），它接受 `Endpoints` 结构体并返回一个 `http.Handler`（通常是来自 Gorilla Mux 的 `*mux.Router`），其中包含所有已配置的路由。

**示例：**

接口 (`productservice/service.go`)：
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

概念性生成的代码：

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

#### 使用生成的服务
您通常会实现 `Service` 接口，然后在您的 `main.go` 或设置代码中：
1. 创建您的服务实现的实例。
2. 调用生成的 `NewEndpoint` 函数来创建 Go kit 端点，可能应用中间件。
3. 调用生成的 `NewHTTPHandler`（或类似名称的函数）来获取 `http.Handler`。
4. 使用此处理程序启动 HTTP 服务器。

此插件自动化了设置 Go kit HTTP 服务所涉及的大部分样板代码。

### `@cep` (CE 权限 SQL 生成)

#### 目的
为名为 `sys_permission` 的权限系统生成 SQL `INSERT` 语句。它为从注解接口派生的父菜单项以及该接口中通过 `@kit-http` 指令公开的每个 HTTP 方法的子权限条目创建 SQL 条目。

#### 目标
Go 接口（与使用 `@kit` 及其子指令注解的接口相同）。

#### 生成的文件
`cep.sql`，与接口在同一个包中。

**注意：** 此插件似乎用于特定的内部用例（“CE”可能指“Cloud Engine”或类似的内部系统），可能与所有 `genx` 用户无关。

#### 先决条件

*   接口应使用 `@kit` 插件进行注解，尤其是在接口上使用 `@basePath` 并在方法上使用 `@kit-http`，因为 `@cep` 使用此信息。

#### 使用的接口级信息

*   接口定义上的主要注释（例如 `// ServiceName is a service for...`）用作 SQL 中父菜单项的 `alias` 和 `description`。
*   `@basePath` 用于构建父菜单项的 `path`（例如 `<basePath>/index`）。

#### 使用的方法级信息

*   对于每个使用 `// @kit-http <path> <http_method>` 注解的方法：
    *   方法的主要注释用作权限条目的 `alias` 和 `description`。
    *   记录 `<http_method>`（例如 GET、POST）。
    *   记录完整路径（basePath + 方法路径）。
    *   为权限条目生成唯一名称（例如，基于基本路径、方法名称和 HTTP 方法）。
    *   通过哈希唯一名称生成唯一 ID（整数）。

#### 生成的 SQL (`cep.sql` - 概念性)

该文件将包含一系列 `INSERT INTO sys_permission (...) VALUES (...);` 语句。

**示例：**

给定一个接口：
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

概念性 `userservice/cep.sql`:
```sql
-- Parent Menu Item for UserService
INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description, created_at, updated_at) VALUES (123456789012345, 0, '', 1, 'GET', 'UserService manages user operations.', 'menu.admin.users.userservice', '/admin/users/index', 'UserService manages user operations.', 'YYYY-MM-DD HH:MM:SS', 'YYYY-MM-DD HH:MM:SS');

-- Permission for ListAllUsers
INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description, created_at, updated_at) VALUES (987654321098765, 123456789012345, '', 0, 'GET', 'ListAllUsers retrieves all users.', '.admin.users.listallusers.get', '/admin/users/list', 'ListAllUsers retrieves all users.', 'YYYY-MM-DD HH:MM:SS', 'YYYY-MM-DD HH:MM:SS');

-- Permission for DeleteUserByID
INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description, created_at, updated_at) VALUES (543210987654321, 123456789012345, '', 0, 'DELETE', 'DeleteUserByID removes a user.', '.admin.users.deleteuserbyid.delete', '/admin/users/{id}', 'DeleteUserByID removes a user.', 'YYYY-MM-DD HH:MM:SS', 'YYYY-MM-DD HH:MM:SS');
```
*(注意：确切的 `id`、`name` 和 `alias` 生成逻辑可能因模板 `ce_permission.tmpl` 内的内部哈希和字符串操作规则而略有不同。)*

如果您要将 `genx` 生成的服务与使用此特定 `sys_permission` 表结构的系统集成，则此插件很有用。

### `@observer` (日志和追踪中间件生成)

#### 目的
生成用于日志记录和追踪的 Go kit 服务中间件。此插件似乎提供了一种替代的或可能较旧的方法来为 `genx` 生成的 Go kit 服务添加可观察性，与更模块化的 `@log` 和 `@trace` 插件不同。

#### 目标
Go 接口（用于为 Go kit 定义的服务）。

#### 生成的文件
*   `logging.go`：包含服务的日志中间件。
*   `tracing.go`：包含追踪中间件（根据其生成器中常见的导入，似乎使用 OpenTracing）。

**注意：** 此插件可能特定于特定设置，或代表 `genx` 中可观察性的早期方法。对于新项目，请考虑 `@log` 和 `@trace`（或全面的 `@otel`）是否提供更大的灵活性。

#### 先决条件

*   理想情况下，该接口也应为其生成 Go kit HTTP 绑定（例如，使用 `@kit`），因为可观察性通常应用于服务端点。

#### 工作原理 (概念性)

`@observer` 插件处理带注解的接口，并使用模板 (`ce_log.tmpl`, `ce_trace.tmpl`) 生成：

1.  **日志中间件 (`logging.go`)：**
    *   一个结构体（例如 `loggingMiddleware`），包装服务接口。
    *   与服务接口中每个方法对应的方法。
    *   这些包装方法在调用实际服务方法之前和/或之后记录有关调用的信息（例如，方法名称、参数、错误、持续时间）。
    *   一个构造函数，用于创建此日志中间件的实例，通常需要一个日志记录器实例（例如 `log.Logger`）和链中的下一个服务。

2.  **追踪中间件 (`tracing.go`)：**
    *   一个结构体（例如 `tracingMiddleware`），包装服务接口。
    *   与服务接口中每个方法对应的方法。
    *   这些包装方法为每个调用创建 span（可能使用 OpenTracing，从典型的导入如 `opentracing-go` 判断），记录标签（如方法名称、错误），并在完成后结束 span。
    *   一个构造函数，用于创建此追踪中间件的实例，通常需要一个追踪器实例（例如 `opentracing.Tracer`）和下一个服务。

#### 用法

使用 `// @observer` 注解您的服务接口。

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

#### 生成的文件 (概念性内容)

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

#### 应用中间件 (在 Go kit 设置中)

如果您将其与 `@kit` 生成的服务一起使用，则通常在构建端点或服务实例时应用这些中间件。

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

此插件提供了一种基于模板的方法，可将日志记录和追踪中间件直接绑定到服务定义。

### `@alert`

#### 目的
生成 Go kit 风格的中间件，当特定服务方法返回错误时启用发送警报。它与外部警报服务集成。

#### 目标
Go 接口。

#### 生成的文件
`alert.go`，与注解的接口在同一个包中。

#### 接口级指令

*   `// @alert`
    *   放置在接口定义上。
    *   标记该接口以供 `@alert` 插件处理。它本身不接受参数。

#### 方法级指令 (在 `@alert` 注解接口的方法内)

*   `// @alert-enable`
    *   对于应激活错误警报的方法，此为**必需**。如果缺少此指令，则该方法的错误不会通过此中间件触发警报。
*   `// @alert-level <level>`
    *   **可选。** 指定警报的严重级别。
    *   有效值：`"info"`、`"warn"`、`"error"`。
    *   如果省略，则使用中间件初始化期间提供的默认警报级别。
    *   这些会映射到假定的 `alarm` 包中的 `alarm.LevelInfo`、`alarm.LevelWarning`、`alarm.LevelError`。
*   `// @alert-metrics <metrics_suffix>`
    *   **可选。** 当推送警报时，附加到默认指标标识符（即 `pkgName.MethodName`）的字符串后缀。例如：如果 `metrics_suffix` 是 `"custom_metric"`，则标识符变为 `pkgName.MethodName.custom_metric`。

#### 生成的代码 (`alert.go`)

1.  **`alert` 结构体：**
    *   生成一个名为 `alert` 的结构体。它包含：
        *   `level alarm.Level`：默认警报级别。
        *   `silencePeriod int`：警报静默期。
        *   `api api.Service`：用于推送警报的外部服务实例。此服务必须公开一个 `Alarm().Push(...)` 方法。
        *   `next Service`：Go kit 中间件链中的下一个服务（您的实际服务实现）。
        *   `logger log.Logger`：用于在推送警报失败时记录内部错误的日志记录器。

2.  **`NewAlert` 构造函数 (中间件工厂)：**
    ```go
    // Conceptual signature from generated code:
    func NewAlert(level alarm.Level, silencePeriod int, api api.Service, log log.Logger) Middleware {
        return func(next Service) Service {
            return &alert{ /* ... initialization ... */ }
        }
    }
    ```
    *   此函数是 Go kit `Middleware` 工厂。
    *   您需要提供默认的 `alarm.Level`、`silencePeriod`、用于发送警报的 `api.Service` 实现以及一个 `log.Logger`。

3.  **包装的方法 (在 `*alert` 结构体上)：**
    *   对于注解接口中的每个方法，都会生成一个相应的包装方法。
    *   如果一个包装方法：
        *   使用 `// @alert-enable` 注解，并且
        *   返回一个错误 (`err != nil`)，
        *   那么，一个 `defer` 块会通过调用以下方法来触发警报：
            `s.api.Alarm().Push(ctx, title, traceId+err.Error(), metricsKey, alertLevel, s.silencePeriod)`
            *   `title`：通常是 `packageName.MethodName`。
            *   `traceId`：从 `ctx.Value("traceId")` 提取。
            *   `metricsKey`：`packageName.MethodName` 或 `packageName.MethodName.metrics_suffix`。
            *   `alertLevel`：由 `@alert-level` 或默认值确定。
    *   然后，该方法调用 `next` 服务的相应方法。

#### 依赖项与假设

*   **`api.Service`：** 您的项目必须定义或导入一个 `api` 包，其中包含一个 `Service` 接口（或具体类型），该接口具有一个 `Alarm()` 方法，该方法又返回一个具有 `Push(...)` 方法的对象。`Push` 的预期签名大致如下：
    `Push(ctx context.Context, title string, content string, metricsKey string, level alarm.Level, silencePeriod int) error`
*   **`alarm` 包：** 假定存在一个 `alarm` 包，提供 `alarm.Level` 和诸如 `alarm.LevelInfo`、`alarm.LevelWarning`、`alarm.LevelError` 之类的常量。
*   **`log.Logger`：** 与 Go kit 日志记录器接口兼容的日志记录器。
*   传递给服务方法的 `context.Context` 可能选择性地包含一个 `traceId` 值，该值将包含在警报内容中。

**示例：**

接口定义 (`payments/service.go`)：
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

概念性用法 (在您的 Go kit 服务设置中)：
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

此插件有助于集中处理服务方法失败的警报逻辑。

### `@do` (依赖注入设置)

#### 目的
使用 `github.com/samber/do/v2` 依赖注入库自动注册服务提供者。它扫描带注解的提供者函数，并生成一个中央 `doInit` 函数以将它们注册到 `do.Injector`。

#### 目标
所有已处理包中的 Go 函数。

#### 生成的文件
`do_init.go`。此文件在任何使用 `@do(type="init")` 注解的函数的包中生成。此生成文件中的 `doInit` 函数将注册来自*所有*扫描包的提供者。

#### 函数级指令

函数应使用 `// @do(...)` 进行注解，以便此插件进行处理。

*   `// @do(type="<type_value>"[, name="<service_name>"])`
    *   `type="<type_value>"`：**必需。**
        *   `provide`：将注解的函数标记为应向 `do.Injector` 注册的构造函数或提供者。函数签名应与 `do.Provide` 兼容（例如 `func(i do.Injector) (MyService, error)` 或 `func() OtherService`）。
        *   `init`：此函数的主要作用是将其包标记为将生成 `do_init.go` 文件的位置。此生成文件中的 `doInit` 函数将包含注册所有发现的 "provide" 函数的逻辑。`genx` 不直接使用标记为 `type="init"` 的函数的主体。
    *   `name="<service_name>"`：**可选。** 仅当 `type="provide"` 时适用。如果指定，则使用 `do.ProvideNamed(injector, "<service_name>", yourProviderFunc)` 注册提供者函数。否则，使用 `do.Provide(injector, yourProviderFunc)`。

#### 生成的代码 (`do_init.go`)

*   在每个至少包含一个使用 `@do(type="init")` 注解的函数的包中，都会生成一个名为 `do_init.go` 的文件。
*   此文件将包含一个函数：`func doInit(i do.Injector)`。
*   `doInit` 函数包括：
    *   需要注册的提供者函数所在的所有包的导入语句。`genx` 通过创建别名来处理潜在的包名冲突。
    *   对所有扫描包中每个使用 `@do(type="provide")` 注解的函数的 `do.Provide(i, ...)` 或 `do.ProvideNamed(i, ...)` 调用。
    *   提供者按确定性顺序注册（按包路径排序，然后按函数名排序）。

#### 工作原理

1.  `genx` 扫描项目中的所有 Go 源文件。
2.  它识别所有使用 `@do(...)` 注解的函数。
3.  它根据 `type="provide"` 或 `type="init"` 对它们进行分类。
4.  对于每个包含一个或多个标记为 `type="init"` 的函数的包，`genx` 在该包中生成一个 `do_init.go` 文件。
5.  任何此类生成的 `do_init.go` 文件中的 `doInit` 函数都将包含对在整个项目扫描中找到的*所有*提供者函数的注册调用。因此，通常您只需要在整个项目中有一个 `// @do(type="init")` 注解（例如，在 `cmd/myapp/main.go` 或 `internal/app/init.go` 文件中）即可生成一个单一、全面的 `doInit` 函数。

**示例：**

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

生成的 `app/do_init.go`:
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
*(注意：确切的别名 `servicesalias` 仅为说明；如果包名 `services` 与其他导入或本地名称冲突，`genx` 将生成一个。)*

#### 使用生成的 `doInit`

在应用程序的主入口点中：
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

此插件通过自动生成提供者注册代码来简化 `samber/do/v2` 的设置，从而减少样板代码和潜在的手动错误。
