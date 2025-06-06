# GenX 错误处理系统

## 概述

GenX 的新错误处理系统提供了一个统一、优雅的方式来处理代码生成过程中的错误。它完全替代了之前使用 `panic` 的方式，提供了更好的错误信息、位置定位和 TUI 集成。

## 主要特性

### 🎯 **结构化错误信息**
- 错误代码分类
- 详细的错误消息
- 精确的位置信息
- 丰富的上下文信息

### 🛡️ **Panic 恢复机制**
- 自动捕获和转换 panic
- 安全执行包装器
- 调用栈信息保留

### 🎨 **TUI 集成**
- 美观的错误显示
- 错误列表管理
- 修复建议提示

### 🔗 **错误链支持**
- 错误原因链接
- 上下文传播
- 调试信息保留

## 核心组件

### 1. GenxError

```go
type GenxError struct {
    Code     ErrorCode `json:"code"`
    Message  string    `json:"message"`
    Details  string    `json:"details,omitempty"`
    Cause    error     `json:"-"`
    Location *Location `json:"location,omitempty"`
    Context  *Context  `json:"context,omitempty"`
    Stack    []Frame   `json:"stack,omitempty"`
}
```

### 2. 错误代码

```go
const (
    // 解析相关错误
    ErrCodeParseDoc        ErrorCode = "PARSE_DOC"
    ErrCodeParseInterface  ErrorCode = "PARSE_INTERFACE"
    ErrCodeParseStruct     ErrorCode = "PARSE_STRUCT"
    ErrCodeParseAnnotation ErrorCode = "PARSE_ANNOTATION"
    
    // 验证相关错误
    ErrCodeValidateMethod     ErrorCode = "VALIDATE_METHOD"
    ErrCodeValidateAnnotation ErrorCode = "VALIDATE_ANNOTATION"
    
    // 生成相关错误
    ErrCodeGenerate    ErrorCode = "GENERATE"
    ErrCodeGenTemplate ErrorCode = "GEN_TEMPLATE"
    
    // 插件相关错误
    ErrCodePlugin     ErrorCode = "PLUGIN"
    ErrCodePluginInit ErrorCode = "PLUGIN_INIT"
)
```

## 使用方法

### 1. 创建错误

#### 基本错误创建
```go
err := common.NewError(common.ErrCodeParseAnnotation, "invalid annotation syntax").
    WithPlugin("@kit-http-client").
    WithAnnotation("@kit-http").
    WithDetails("annotation parameters are malformed").
    Build()
```

#### 使用便捷函数
```go
// 解析错误
err := common.ParseError("failed to parse interface").
    WithInterface("UserService").
    WithMethod("GetUser").
    Build()

// 验证错误
err := common.ValidationError("missing required annotation").
    WithPlugin("@crud").
    WithStruct("UserCrud").
    WithAnnotation("@crud").
    Build()

// 生成错误
err := common.GenerateError("template execution failed").
    WithPlugin("@template").
    WithDetails("failed to execute template").
    Build()
```

### 2. 添加上下文信息

```go
err := common.ParseError("method parsing failed").
    WithPlugin("@kit-http-client").           // 插件名称
    WithInterface("UserService").             // 接口名称
    WithMethod("GetUser").                    // 方法名称
    WithStruct("GetUserRequest").             // 结构体名称
    WithField("userID").                      // 字段名称
    WithAnnotation("@kit-http").              // 注解名称
    WithExtra("expected", "GET").             // 额外信息
    WithExtra("actual", "POST").
    WithDetails("HTTP method mismatch").      // 详细描述
    WithCause(originalError).                 // 原始错误
    Build()
```

### 3. 位置信息

```go
// 从 token.Pos 添加位置
err := common.ParseError("syntax error").
    WithTokenPos(fileSet, pos).
    Build()

// 手动添加位置
err := common.ParseError("file error").
    WithLocation("user.go", 42, 10).
    Build()
```

### 4. 安全执行

#### 包装可能 panic 的函数
```go
err := common.WithRecovery(func() error {
    // 可能会 panic 的代码
    riskyOperation()
    return nil
})

if err != nil {
    // 处理错误（包括从 panic 恢复的错误）
}
```

#### 带返回值的安全执行
```go
result, err := common.WithRecoveryResult(func() (string, error) {
    // 可能会 panic 的代码
    return riskyOperationWithResult()
})
```

#### 使用恢复处理器
```go
recoveryHandler := common.NewRecoveryHandler(func(err *common.GenxError) {
    // 自定义错误处理逻辑
    log.Printf("Recovered from panic: %v", err)
})

err := recoveryHandler.SafeExecute(func() error {
    // 可能会 panic 的代码
    return nil
})
```

### 5. TUI 集成

#### 创建 TUI 错误处理器
```go
errorHandler := common.NewTUIErrorHandler()

// 添加错误
errorHandler.AddError(genxError)

// 格式化单个错误
formatted := errorHandler.FormatError(genxError)
fmt.Println(formatted)

// 格式化错误列表
list := errorHandler.FormatErrorList()
fmt.Println(list)

// 获取修复建议
suggestions := errorHandler.GetSuggestions(genxError)
for _, suggestion := range suggestions {
    fmt.Println(suggestion)
}
```

#### 创建错误列表组件
```go
errorList := errorHandler.CreateErrorList()
// 在 TUI 应用中使用 errorList
```

## 插件开发指南

### 1. 替换 panic

**之前的代码：**
```go
func (p *Plugin) Parse() {
    if err != nil {
        panic(err)  // ❌ 不好的做法
    }
}
```

**改进后的代码：**
```go
func (p *Plugin) Parse() error {
    if err != nil {
        return common.ParseError("failed to parse").
            WithCause(err).
            WithPlugin("@my-plugin").
            WithDetails("detailed error description").
            Build()  // ✅ 好的做法
    }
    return nil
}
```

### 2. 插件接口更新

**更新插件接口：**
```go
// 之前
func (p *Plugin) Gen(option gen.Option, metas []gen.InterfaceGoTypeMeta) ([]gen.GenResult, error) {
    plugin := &MyPlugin{option: option, metas: metas}
    plugin.Parse()  // 可能 panic
    // ...
}

// 改进后
func (p *Plugin) Gen(option gen.Option, metas []gen.InterfaceGoTypeMeta) ([]gen.GenResult, error) {
    plugin := &MyPlugin{option: option, metas: metas}
    
    if err := plugin.Parse(); err != nil {
        if genxErr, ok := err.(*common.GenxError); ok {
            return nil, genxErr
        }
        return nil, common.PluginError("plugin execution failed").
            WithCause(err).
            WithPlugin(p.Name()).
            Build()
    }
    // ...
}
```

### 3. 错误处理最佳实践

#### 提供具体的错误信息
```go
// ❌ 模糊的错误
return common.ParseError("parsing failed").Build()

// ✅ 具体的错误
return common.ParseError("missing required annotation").
    WithPlugin("@kit-http-client").
    WithInterface("UserService").
    WithMethod("GetUser").
    WithAnnotation("@kit-http").
    WithDetails("@kit-http annotation is required for HTTP client generation. Format: @kit-http <url> <method>").
    Build()
```

#### 使用适当的错误代码
```go
// 解析错误
common.ParseError("...")
common.NewError(common.ErrCodeParseAnnotation, "...")

// 验证错误
common.ValidationError("...")
common.NewError(common.ErrCodeValidateMethod, "...")

// 生成错误
common.GenerateError("...")
common.NewError(common.ErrCodeGenTemplate, "...")
```

#### 保留错误链
```go
if err := someOperation(); err != nil {
    return common.GenerateError("operation failed").
        WithCause(err).  // 保留原始错误
        WithPlugin("@my-plugin").
        Build()
}
```

## 错误显示示例

### 控制台输出
```
🚨 Error Occurred
[PARSE_ANNOTATION] missing required annotation
📍 user.go:42 in parseMethod
🔍 Plugin: @kit-http-client | Interface: UserService | Method: GetUser | Annotation: @kit-http
💡 @kit-http annotation is required for HTTP client generation. Format: @kit-http <url> <method>
🔗 Caused by: annotation not found
```

### TUI 显示
```
┌─ 🚨 2 Error(s) Found ─────────────────────────────────────┐
│ 15:04:05 [PARSE_ANNOTATION] missing required annotation  │
│          (in GetUser)                                     │
│ 15:04:06 [VALIDATE_METHOD] invalid method signature      │
│          (in CreateUser)                                  │
└───────────────────────────────────────────────────────────┘
```

## 迁移指南

### 1. 识别 panic 使用
搜索代码中的 `panic(` 调用，特别是在插件代码中。

### 2. 更新方法签名
将不返回错误的方法更新为返回 `error`。

### 3. 替换 panic 调用
使用适当的错误构建器替换 `panic` 调用。

### 4. 添加错误处理
在调用可能返回错误的方法时添加错误处理。

### 5. 测试错误场景
确保所有错误路径都被正确处理和测试。

## 总结

新的错误处理系统提供了：
- 🎯 **更好的错误信息** - 结构化、详细、可定位
- 🛡️ **更安全的执行** - 自动 panic 恢复
- 🎨 **更好的用户体验** - TUI 集成和美观显示
- 🔧 **更容易调试** - 完整的上下文和调用栈信息

通过使用这个系统，GenX 的错误处理变得更加专业和用户友好。
