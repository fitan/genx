package common

import (
	"fmt"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrorCode 错误代码类型
type ErrorCode string

const (
	// 解析相关错误
	ErrCodeParseDoc        ErrorCode = "PARSE_DOC"
	ErrCodeParseInterface  ErrorCode = "PARSE_INTERFACE"
	ErrCodeParseStruct     ErrorCode = "PARSE_STRUCT"
	ErrCodeParseAnnotation ErrorCode = "PARSE_ANNOTATION"
	ErrCodeParseType       ErrorCode = "PARSE_TYPE"

	// 验证相关错误
	ErrCodeValidateMethod     ErrorCode = "VALIDATE_METHOD"
	ErrCodeValidateAnnotation ErrorCode = "VALIDATE_ANNOTATION"
	ErrCodeValidateParam      ErrorCode = "VALIDATE_PARAM"
	ErrCodeValidateConfig     ErrorCode = "VALIDATE_CONFIG"

	// 生成相关错误
	ErrCodeGenerate     ErrorCode = "GENERATE"
	ErrCodeGenTemplate  ErrorCode = "GEN_TEMPLATE"
	ErrCodeGenFile      ErrorCode = "GEN_FILE"
	ErrCodeGenCode      ErrorCode = "GEN_CODE"

	// 配置相关错误
	ErrCodeConfig     ErrorCode = "CONFIG"
	ErrCodeConfigFile ErrorCode = "CONFIG_FILE"

	// 插件相关错误
	ErrCodePlugin     ErrorCode = "PLUGIN"
	ErrCodePluginInit ErrorCode = "PLUGIN_INIT"

	// 文件相关错误
	ErrCodeFile      ErrorCode = "FILE"
	ErrCodeFileWrite ErrorCode = "FILE_WRITE"
	ErrCodeFileRead  ErrorCode = "FILE_READ"
)

// GenxError 统一的错误类型
type GenxError struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	Details  string    `json:"details,omitempty"`
	Cause    error     `json:"-"`
	Location *Location `json:"location,omitempty"`
	Context  *Context  `json:"context,omitempty"`
	Stack    []Frame   `json:"stack,omitempty"`
}

// Location 错误位置信息
type Location struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Package  string `json:"package,omitempty"`
	Function string `json:"function,omitempty"`
}

// Context 错误上下文信息
type Context struct {
	Plugin     string            `json:"plugin,omitempty"`
	Interface  string            `json:"interface,omitempty"`
	Method     string            `json:"method,omitempty"`
	Struct     string            `json:"struct,omitempty"`
	Field      string            `json:"field,omitempty"`
	Annotation string            `json:"annotation,omitempty"`
	Extra      map[string]string `json:"extra,omitempty"`
}

// Frame 调用栈帧
type Frame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Error 实现 error 接口
func (e *GenxError) Error() string {
	var parts []string

	// 添加错误代码和消息
	parts = append(parts, fmt.Sprintf("[%s] %s", e.Code, e.Message))

	// 添加位置信息
	if e.Location != nil {
		location := fmt.Sprintf("at %s:%d", filepath.Base(e.Location.File), e.Location.Line)
		if e.Location.Function != "" {
			location += fmt.Sprintf(" in %s", e.Location.Function)
		}
		parts = append(parts, location)
	}

	// 添加上下文信息
	if e.Context != nil {
		var contextParts []string
		if e.Context.Plugin != "" {
			contextParts = append(contextParts, fmt.Sprintf("plugin=%s", e.Context.Plugin))
		}
		if e.Context.Interface != "" {
			contextParts = append(contextParts, fmt.Sprintf("interface=%s", e.Context.Interface))
		}
		if e.Context.Method != "" {
			contextParts = append(contextParts, fmt.Sprintf("method=%s", e.Context.Method))
		}
		if e.Context.Struct != "" {
			contextParts = append(contextParts, fmt.Sprintf("struct=%s", e.Context.Struct))
		}
		if e.Context.Field != "" {
			contextParts = append(contextParts, fmt.Sprintf("field=%s", e.Context.Field))
		}
		if e.Context.Annotation != "" {
			contextParts = append(contextParts, fmt.Sprintf("annotation=%s", e.Context.Annotation))
		}
		if len(contextParts) > 0 {
			parts = append(parts, fmt.Sprintf("context: %s", strings.Join(contextParts, ", ")))
		}
	}

	// 添加详细信息
	if e.Details != "" {
		parts = append(parts, fmt.Sprintf("details: %s", e.Details))
	}

	// 添加原因
	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Cause))
	}

	return strings.Join(parts, " | ")
}

// Unwrap 支持 errors.Unwrap
func (e *GenxError) Unwrap() error {
	return e.Cause
}

// ErrorBuilder 错误构建器
type ErrorBuilder struct {
	err *GenxError
}

// NewError 创建新的错误构建器
func NewError(code ErrorCode, message string) *ErrorBuilder {
	return &ErrorBuilder{
		err: &GenxError{
			Code:    code,
			Message: message,
			Stack:   captureStack(2), // 跳过当前函数和调用者
		},
	}
}

// WithDetails 添加详细信息
func (b *ErrorBuilder) WithDetails(details string) *ErrorBuilder {
	b.err.Details = details
	return b
}

// WithCause 添加原因
func (b *ErrorBuilder) WithCause(cause error) *ErrorBuilder {
	b.err.Cause = cause
	return b
}

// WithLocation 添加位置信息
func (b *ErrorBuilder) WithLocation(file string, line, column int) *ErrorBuilder {
	b.err.Location = &Location{
		File:   file,
		Line:   line,
		Column: column,
	}
	return b
}

// WithTokenPos 从 token.Pos 添加位置信息
func (b *ErrorBuilder) WithTokenPos(fset *token.FileSet, pos token.Pos) *ErrorBuilder {
	if fset != nil && pos.IsValid() {
		position := fset.Position(pos)
		b.err.Location = &Location{
			File:   position.Filename,
			Line:   position.Line,
			Column: position.Column,
		}
	}
	return b
}

// WithContext 添加上下文信息
func (b *ErrorBuilder) WithContext(ctx *Context) *ErrorBuilder {
	b.err.Context = ctx
	return b
}

// WithPlugin 添加插件上下文
func (b *ErrorBuilder) WithPlugin(plugin string) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = &Context{}
	}
	b.err.Context.Plugin = plugin
	return b
}

// WithInterface 添加接口上下文
func (b *ErrorBuilder) WithInterface(interfaceName string) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = &Context{}
	}
	b.err.Context.Interface = interfaceName
	return b
}

// WithMethod 添加方法上下文
func (b *ErrorBuilder) WithMethod(method string) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = &Context{}
	}
	b.err.Context.Method = method
	return b
}

// WithStruct 添加结构体上下文
func (b *ErrorBuilder) WithStruct(structName string) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = &Context{}
	}
	b.err.Context.Struct = structName
	return b
}

// WithField 添加字段上下文
func (b *ErrorBuilder) WithField(field string) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = &Context{}
	}
	b.err.Context.Field = field
	return b
}

// WithAnnotation 添加注解上下文
func (b *ErrorBuilder) WithAnnotation(annotation string) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = &Context{}
	}
	b.err.Context.Annotation = annotation
	return b
}

// WithExtra 添加额外上下文
func (b *ErrorBuilder) WithExtra(key, value string) *ErrorBuilder {
	if b.err.Context == nil {
		b.err.Context = &Context{}
	}
	if b.err.Context.Extra == nil {
		b.err.Context.Extra = make(map[string]string)
	}
	b.err.Context.Extra[key] = value
	return b
}

// Build 构建错误
func (b *ErrorBuilder) Build() *GenxError {
	return b.err
}

// captureStack 捕获调用栈
func captureStack(skip int) []Frame {
	var frames []Frame
	for i := skip; i < skip+10; i++ { // 最多捕获10层
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		frames = append(frames, Frame{
			Function: fn.Name(),
			File:     file,
			Line:     line,
		})
	}
	return frames
}

// 便捷函数

// ParseError 创建解析错误
func ParseError(message string) *ErrorBuilder {
	return NewError(ErrCodeParseDoc, message)
}

// ValidationError 创建验证错误
func ValidationError(message string) *ErrorBuilder {
	return NewError(ErrCodeValidateMethod, message)
}

// GenerateError 创建生成错误
func GenerateError(message string) *ErrorBuilder {
	return NewError(ErrCodeGenerate, message)
}

// ConfigError 创建配置错误
func ConfigError(message string) *ErrorBuilder {
	return NewError(ErrCodeConfig, message)
}

// PluginError 创建插件错误
func PluginError(message string) *ErrorBuilder {
	return NewError(ErrCodePlugin, message)
}

// FileError 创建文件错误
func FileError(message string) *ErrorBuilder {
	return NewError(ErrCodeFile, message)
}
