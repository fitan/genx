package main

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

// @fx(type="app")
//
//go:generate genx
func main() {
	runApp()
}

// 实际的 main 函数逻辑
func runApp() {
	app := NewApp(
		// 额外的配置选项
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{}
		}),
		fx.StartTimeout(30*time.Second),
		fx.StopTimeout(30*time.Second),
	)

	// 启动应用
	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	// 等待一段时间让应用运行
	time.Sleep(5 * time.Second)

	// 优雅关闭
	if err := app.Stop(ctx); err != nil {
		panic(err)
	}
}
