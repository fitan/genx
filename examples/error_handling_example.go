package main

import (
	"fmt"
	"log"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/genx/plugs/kithttpclient"
)

// 这个示例展示了如何使用新的错误处理系统

func main() {
	// 创建 TUI 错误处理器
	errorHandler := common.NewTUIErrorHandler()

	// 创建恢复处理器，集成 TUI 错误显示
	recoveryHandler := common.NewRecoveryHandler(func(err *common.GenxError) {
		// 添加错误到 TUI 处理器
		errorHandler.AddError(err)
		
		// 格式化并显示错误
		fmt.Println(errorHandler.FormatError(err))
		
		// 显示修复建议
		suggestions := errorHandler.GetSuggestions(err)
		if len(suggestions) > 0 {
			fmt.Println("\n💡 Suggestions:")
			for _, suggestion := range suggestions {
				fmt.Println(suggestion)
			}
		}
	})

	// 示例1: 安全执行可能 panic 的代码
	fmt.Println("=== Example 1: Safe execution ===")
	err := recoveryHandler.SafeExecute(func() error {
		// 模拟一个会产生错误的操作
		return simulateParsingError()
	})
	
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	// 示例2: 使用 WithRecovery 包装函数
	fmt.Println("\n=== Example 2: WithRecovery wrapper ===")
	err = common.WithRecovery(func() error {
		return simulateValidationError()
	})
	
	if err != nil {
		if genxErr, ok := err.(*common.GenxError); ok {
			errorHandler.AddError(genxErr)
			fmt.Println(errorHandler.FormatError(genxErr))
		}
	}

	// 示例3: 插件错误处理
	fmt.Println("\n=== Example 3: Plugin error handling ===")
	demonstratePluginErrorHandling(errorHandler)

	// 显示所有错误的摘要
	fmt.Println("\n=== Error Summary ===")
	fmt.Println(errorHandler.FormatErrorList())
}

// simulateParsingError 模拟解析错误
func simulateParsingError() error {
	return common.ParseError("failed to parse interface").
		WithPlugin("@kit-http-client").
		WithInterface("UserService").
		WithMethod("GetUser").
		WithAnnotation("@kit-http").
		WithDetails("@kit-http annotation is malformed. Expected format: @kit-http <url> <method>").
		WithExtra("expected", "@kit-http /users/{id} GET").
		WithExtra("actual", "@kit-http /users").
		Build()
}

// simulateValidationError 模拟验证错误
func simulateValidationError() error {
	return common.ValidationError("invalid method signature").
		WithPlugin("@crud").
		WithStruct("UserCrud").
		WithMethod("CreateUser").
		WithDetails("method must return exactly 2 values (response, error), got 1").
		WithExtra("signature", "CreateUser(ctx context.Context, req CreateUserRequest) UserResponse").
		WithExtra("expected", "CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)").
		Build()
}

// demonstratePluginErrorHandling 演示插件错误处理
func demonstratePluginErrorHandling(errorHandler *common.TUIErrorHandler) {
	// 模拟插件执行
	plug := &kithttpclient.Plug{}
	
	// 创建一个模拟的选项和元数据
	option := gen.Option{
		// 这里应该有真实的配置，但为了示例我们使用空值
	}
	
	var metas []gen.InterfaceGoTypeMeta
	
	// 使用安全执行来调用插件
	_, err := common.SafeExecuteWithResult(common.DefaultRecoveryHandler, func() ([]gen.GenResult, error) {
		return plug.Gen(option, metas)
	})
	
	if err != nil {
		if genxErr, ok := err.(*common.GenxError); ok {
			errorHandler.AddError(genxErr)
			fmt.Println("Plugin error caught:")
			fmt.Println(errorHandler.FormatError(genxErr))
		} else {
			// 如果不是 GenxError，包装它
			wrappedErr := common.PluginError("plugin execution failed").
				WithCause(err).
				WithPlugin(plug.Name()).
				WithDetails("unexpected error during plugin execution").
				Build()
			errorHandler.AddError(wrappedErr)
			fmt.Println("Wrapped plugin error:")
			fmt.Println(errorHandler.FormatError(wrappedErr))
		}
	}
}

// ExampleTUIIntegration 展示如何在 TUI 应用中集成错误处理
func ExampleTUIIntegration() {
	// 这个函数展示了如何在实际的 TUI 应用中使用错误处理
	
	errorHandler := common.NewTUIErrorHandler()
	
	// 在 TUI 应用的更新循环中
	handleError := func(err error) {
		if genxErr, ok := err.(*common.GenxError); ok {
			errorHandler.AddError(genxErr)
			
			// 在 TUI 中显示错误
			// 这里可以触发 TUI 的错误显示组件
			log.Printf("Error added to TUI: %s", genxErr.Error())
		}
	}
	
	// 示例：处理插件错误
	err := common.WithRecovery(func() error {
		// 模拟插件操作
		return common.GenerateError("code generation failed").
			WithPlugin("@example").
			WithDetails("failed to generate output file").
			Build()
	})
	
	if err != nil {
		handleError(err)
	}
	
	// 创建错误列表用于 TUI 显示
	errorList := errorHandler.CreateErrorList()
	_ = errorList // 在实际的 TUI 应用中，这会被添加到界面中
}

// BestPracticesExample 展示最佳实践
func BestPracticesExample() {
	fmt.Println("=== Best Practices ===")
	
	// 1. 总是使用具体的错误代码
	err1 := common.NewError(common.ErrCodeParseAnnotation, "invalid annotation syntax").
		WithPlugin("@my-plugin").
		WithAnnotation("@my-annotation").
		WithDetails("annotation parameters are malformed").
		Build()
	
	// 2. 提供足够的上下文信息
	err2 := common.ValidationError("missing required field").
		WithPlugin("@validation").
		WithStruct("UserRequest").
		WithField("email").
		WithDetails("email field is required for user creation").
		WithExtra("validation_rule", "required").
		Build()
	
	// 3. 使用恢复机制包装可能 panic 的代码
	result, err3 := common.WithRecoveryResult(func() (string, error) {
		// 可能会 panic 的代码
		return "success", nil
	})
	
	fmt.Printf("Results: %v, %v, %v, %s\n", err1, err2, err3, result)
	
	// 4. 在插件中使用错误链
	err4 := common.GenerateError("template execution failed").
		WithCause(err1). // 链接原始错误
		WithPlugin("@template").
		WithDetails("failed to execute template due to validation error").
		Build()
	
	fmt.Printf("Chained error: %v\n", err4)
}

// 运行示例
func init() {
	// 这些示例函数可以在需要时调用
	_ = ExampleTUIIntegration
	_ = BestPracticesExample
}
