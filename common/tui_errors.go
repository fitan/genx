package common

import (
	"fmt"
	"strings"
	"time"
)

// TUIErrorHandler TUI 错误处理器（简化版本）
type TUIErrorHandler struct {
	errors []ErrorItem
}

// ErrorItem 错误项
type ErrorItem struct {
	Error     *GenxError
	Timestamp time.Time
	ID        string
}

// NewTUIErrorHandler 创建 TUI 错误处理器
func NewTUIErrorHandler() *TUIErrorHandler {
	return &TUIErrorHandler{
		errors: make([]ErrorItem, 0),
	}
}

// AddError 添加错误
func (h *TUIErrorHandler) AddError(err *GenxError) {
	item := ErrorItem{
		Error:     err,
		Timestamp: time.Now(),
		ID:        fmt.Sprintf("err_%d", time.Now().UnixNano()),
	}
	h.errors = append(h.errors, item)
}

// GetErrors 获取所有错误
func (h *TUIErrorHandler) GetErrors() []ErrorItem {
	return h.errors
}

// Clear 清空错误
func (h *TUIErrorHandler) Clear() {
	h.errors = make([]ErrorItem, 0)
}

// FormatError 格式化错误为简单文本显示
func (h *TUIErrorHandler) FormatError(err *GenxError) string {
	var parts []string

	// 错误标题
	parts = append(parts, "🚨 Error Occurred")

	// 错误代码和消息
	codeMsg := fmt.Sprintf("[%s] %s", err.Code, err.Message)
	parts = append(parts, codeMsg)

	// 位置信息
	if err.Location != nil {
		location := fmt.Sprintf("📍 %s:%d", err.Location.File, err.Location.Line)
		if err.Location.Function != "" {
			location += fmt.Sprintf(" in %s", err.Location.Function)
		}
		parts = append(parts, location)
	}

	// 上下文信息
	if err.Context != nil {
		var contextParts []string
		if err.Context.Plugin != "" {
			contextParts = append(contextParts, fmt.Sprintf("Plugin: %s", err.Context.Plugin))
		}
		if err.Context.Interface != "" {
			contextParts = append(contextParts, fmt.Sprintf("Interface: %s", err.Context.Interface))
		}
		if err.Context.Method != "" {
			contextParts = append(contextParts, fmt.Sprintf("Method: %s", err.Context.Method))
		}
		if err.Context.Struct != "" {
			contextParts = append(contextParts, fmt.Sprintf("Struct: %s", err.Context.Struct))
		}
		if err.Context.Field != "" {
			contextParts = append(contextParts, fmt.Sprintf("Field: %s", err.Context.Field))
		}
		if err.Context.Annotation != "" {
			contextParts = append(contextParts, fmt.Sprintf("Annotation: %s", err.Context.Annotation))
		}

		if len(contextParts) > 0 {
			context := "🔍 " + strings.Join(contextParts, " | ")
			parts = append(parts, context)
		}
	}

	// 详细信息
	if err.Details != "" {
		details := "💡 " + err.Details
		parts = append(parts, details)
	}

	// 原因
	if err.Cause != nil {
		cause := fmt.Sprintf("🔗 Caused by: %v", err.Cause)
		parts = append(parts, cause)
	}

	// 调用栈（可选，通常只在调试时显示）
	if len(err.Stack) > 0 && len(err.Stack) <= 3 { // 只显示前3层
		parts = append(parts, "📚 Stack:")
		for i, frame := range err.Stack {
			if i >= 3 {
				break
			}
			stackLine := fmt.Sprintf("  %d. %s (%s:%d)",
				i+1, frame.Function, frame.File, frame.Line)
			parts = append(parts, stackLine)
		}
	}

	return strings.Join(parts, "\n")
}

// FormatErrorList 格式化错误列表
func (h *TUIErrorHandler) FormatErrorList() string {
	if len(h.errors) == 0 {
		return "✅ No errors"
	}

	var parts []string
	parts = append(parts, fmt.Sprintf("🚨 %d Error(s) Found", len(h.errors)))
	parts = append(parts, strings.Repeat("=", 50))

	for i, item := range h.errors {
		if i >= 5 { // 最多显示5个错误
			remaining := len(h.errors) - 5
			parts = append(parts, fmt.Sprintf("... and %d more errors", remaining))
			break
		}

		timestamp := item.Timestamp.Format("15:04:05")
		errorSummary := fmt.Sprintf("%s [%s] %s",
			timestamp,
			string(item.Error.Code),
			item.Error.Message)

		if item.Error.Context != nil && item.Error.Context.Method != "" {
			errorSummary += fmt.Sprintf(" (in %s)", item.Error.Context.Method)
		}

		parts = append(parts, errorSummary)
		parts = append(parts, strings.Repeat("-", 30))
	}

	return strings.Join(parts, "\n")
}

// GetSuggestions 获取错误修复建议
func (h *TUIErrorHandler) GetSuggestions(err *GenxError) []string {
	var suggestions []string

	switch err.Code {
	case ErrCodeParseAnnotation:
		suggestions = append(suggestions,
			"• Check annotation syntax: @annotation-name param1 param2",
			"• Ensure annotation is properly formatted",
			"• Verify annotation is supported by the plugin")

	case ErrCodeValidateMethod:
		suggestions = append(suggestions,
			"• Check method signature matches expected format",
			"• Ensure method has correct number of parameters",
			"• Verify return types are correct")

	case ErrCodeValidateAnnotation:
		suggestions = append(suggestions,
			"• Check required annotations are present",
			"• Verify annotation parameters are correct",
			"• Ensure annotation values are valid")

	case ErrCodeGenerate:
		suggestions = append(suggestions,
			"• Check generated code syntax",
			"• Verify all imports are available",
			"• Ensure no naming conflicts")

	case ErrCodeConfig:
		suggestions = append(suggestions,
			"• Check configuration file syntax",
			"• Verify all required fields are present",
			"• Ensure configuration values are valid")

	default:
		suggestions = append(suggestions,
			"• Check the error details above",
			"• Verify your code syntax",
			"• Consult the documentation")
	}

	return suggestions
}

// CreateErrorList 创建错误列表（简化版本）
func (h *TUIErrorHandler) CreateErrorList() string {
	if len(h.errors) == 0 {
		return "No errors found."
	}

	var result strings.Builder
	result.WriteString("GenX Errors:\n")
	result.WriteString(strings.Repeat("=", 50) + "\n")

	for i, err := range h.errors {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, h.FormatError(err.Error)))
		result.WriteString(strings.Repeat("-", 30) + "\n")
	}

	return result.String()
}
