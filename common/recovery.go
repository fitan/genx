package common

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// RecoveryHandler 恢复处理器
type RecoveryHandler struct {
	OnPanic func(err *GenxError)
}

// NewRecoveryHandler 创建恢复处理器
func NewRecoveryHandler(onPanic func(err *GenxError)) *RecoveryHandler {
	return &RecoveryHandler{
		OnPanic: onPanic,
	}
}

// Recover 恢复 panic 并转换为 GenxError
func (r *RecoveryHandler) Recover() {
	if rec := recover(); rec != nil {
		var err *GenxError

		switch v := rec.(type) {
		case *GenxError:
			// 已经是 GenxError，直接使用
			err = v
		case error:
			// 普通错误，包装为 GenxError
			err = NewError(ErrCodePlugin, "panic recovered").
				WithCause(v).
				WithDetails(fmt.Sprintf("panic: %v", v)).
				Build()
		case string:
			// 字符串错误
			err = NewError(ErrCodePlugin, "panic recovered").
				WithDetails(fmt.Sprintf("panic: %s", v)).
				Build()
		default:
			// 其他类型
			err = NewError(ErrCodePlugin, "panic recovered").
				WithDetails(fmt.Sprintf("panic: %v", v)).
				Build()
		}

		// 添加调用栈信息
		if len(err.Stack) == 0 {
			err.Stack = captureStackFromPanic()
		}

		// 调用处理函数
		if r.OnPanic != nil {
			r.OnPanic(err)
		}
	}
}

// SafeExecute 安全执行函数，捕获 panic
func (r *RecoveryHandler) SafeExecute(fn func() error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			switch v := rec.(type) {
			case *GenxError:
				err = v
			case error:
				err = NewError(ErrCodePlugin, "panic in safe execution").
					WithCause(v).
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			default:
				err = NewError(ErrCodePlugin, "panic in safe execution").
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			}

			// 调用处理函数
			if r.OnPanic != nil {
				r.OnPanic(err.(*GenxError))
			}
		}
	}()

	return fn()
}

// SafeExecuteWithResult 安全执行带返回值的函数
func SafeExecuteWithResult[T any](r *RecoveryHandler, fn func() (T, error)) (result T, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			var zero T
			result = zero

			switch v := rec.(type) {
			case *GenxError:
				err = v
			case error:
				err = NewError(ErrCodePlugin, "panic in safe execution").
					WithCause(v).
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			default:
				err = NewError(ErrCodePlugin, "panic in safe execution").
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			}

			// 调用处理函数
			if r.OnPanic != nil {
				r.OnPanic(err.(*GenxError))
			}
		}
	}()

	return fn()
}

// captureStackFromPanic 从 panic 中捕获调用栈
func captureStackFromPanic() []Frame {
	var frames []Frame

	// 获取调用栈（用于调试，当前使用 runtime.Caller）
	_ = debug.Stack()

	// 解析调用栈（简化版本）
	for i := 3; i < 13; i++ { // 跳过前几层，最多10层
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

	// 如果没有获取到足够信息，使用 debug.Stack 的信息
	if len(frames) == 0 {
		frames = append(frames, Frame{
			Function: "panic",
			File:     "unknown",
			Line:     0,
		})
	}

	return frames
}

// MustNotPanic 确保函数不会 panic，如果 panic 则转换为错误
func MustNotPanic(fn func()) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			switch v := rec.(type) {
			case *GenxError:
				err = v
			case error:
				err = NewError(ErrCodePlugin, "unexpected panic").
					WithCause(v).
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			default:
				err = NewError(ErrCodePlugin, "unexpected panic").
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			}
		}
	}()

	fn()
	return nil
}

// MustNotPanicWithResult 确保带返回值的函数不会 panic
func MustNotPanicWithResult[T any](fn func() T) (result T, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			var zero T
			result = zero

			switch v := rec.(type) {
			case *GenxError:
				err = v
			case error:
				err = NewError(ErrCodePlugin, "unexpected panic").
					WithCause(v).
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			default:
				err = NewError(ErrCodePlugin, "unexpected panic").
					WithDetails(fmt.Sprintf("panic: %v", v)).
					Build()
			}
		}
	}()

	result = fn()
	return result, nil
}

// TryRecover 尝试恢复并返回错误
func TryRecover() error {
	if rec := recover(); rec != nil {
		switch v := rec.(type) {
		case *GenxError:
			return v
		case error:
			return NewError(ErrCodePlugin, "recovered from panic").
				WithCause(v).
				WithDetails(fmt.Sprintf("panic: %v", v)).
				Build()
		default:
			return NewError(ErrCodePlugin, "recovered from panic").
				WithDetails(fmt.Sprintf("panic: %v", v)).
				Build()
		}
	}
	return nil
}

// DefaultRecoveryHandler 默认的恢复处理器
var DefaultRecoveryHandler = NewRecoveryHandler(func(err *GenxError) {
	// 默认情况下，只记录错误，不做其他处理
	// 实际的错误处理由调用方决定
})

// WithRecovery 为函数添加恢复机制
func WithRecovery(fn func() error) error {
	return DefaultRecoveryHandler.SafeExecute(fn)
}

// WithRecoveryResult 为带返回值的函数添加恢复机制
func WithRecoveryResult[T any](fn func() (T, error)) (T, error) {
	return SafeExecuteWithResult(DefaultRecoveryHandler, fn)
}
