/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/genx/plugs/alert"
	plugCopy "github.com/fitan/genx/plugs/copy"
	"github.com/fitan/genx/plugs/crud"
	"github.com/fitan/genx/plugs/do"
	"github.com/fitan/genx/plugs/enum"
	"github.com/fitan/genx/plugs/fx"
	"github.com/fitan/genx/plugs/gormq"
	"github.com/fitan/genx/plugs/kithttp"
	"github.com/fitan/genx/plugs/kithttpclient"
	"github.com/fitan/genx/plugs/log"
	"github.com/fitan/genx/plugs/otel"
	"github.com/fitan/genx/plugs/temporal"
	"github.com/fitan/genx/plugs/trace"
	"github.com/oklog/run"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc"
	"github.com/urfave/cli/v2"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	app := &cli.App{
		Action: func(ctx *cli.Context) error {
			model := gen.NewModel()

			go func() {
				p := tea.NewProgram(model)
				if _, err := p.Run(); err != nil {
					panic(err)
				}
			}()

			x, err := gen.NewX(staticFiles, "./...", model)
			if err != nil {
				os.Exit(1)
			}

			var moduleWG conc.WaitGroup

			// 创建错误处理器
			errorHandler := common.NewTUIErrorHandler()

			lo.ForEach(x, func(item *gen.X, index int) {
				moduleWG.Go(func() {
					// 使用安全执行包装器来处理插件注册和生成
					err := common.WithRecovery(func() error {
						item.RegImpl(&log.Plug{})
						item.RegImpl(&trace.Plug{})
						item.RegImpl(&otel.Plug{})
						item.RegTypeSpec(&enum.Plug{})
						item.RegStruct(&gormq.Plug{})
						item.RegStruct(&crud.Plug{})
						item.RegImpl(&kithttpclient.Plug{})
						item.RegImpl(&temporal.Plug{})
						item.RegCall(&plugCopy.Plug{})
						item.RegImpl(&kithttp.Plug{})
						item.RegImpl(&kithttp.CEPermissionSqlPlug{})
						item.RegImpl(&kithttp.ObserverPlug{})
						item.RegImpl(&alert.Plug{})
						return item.Gen()
					})

					if err != nil {
						if genxErr, ok := err.(*common.GenxError); ok {
							errorHandler.AddError(genxErr)
							slog.Error("plugin execution failed", "error", errorHandler.FormatError(genxErr))
						} else {
							// 包装非 GenxError
							wrappedErr := common.PluginError("module generation failed").
								WithCause(err).
								WithExtra("module_index", fmt.Sprintf("%d", index)).
								WithDetails("unexpected error during module generation").
								Build()
							errorHandler.AddError(wrappedErr)
							slog.Error("module generation failed", "error", errorHandler.FormatError(wrappedErr))
						}
					}
				})
			})

			moduleWG.Wait()

			x, err = gen.NewX(staticFiles, "./...", model)
			if err != nil {
				// 包装为 GenxError 并添加到错误处理器
				genxErr := common.ConfigError("failed to create global generator").
					WithCause(err).
					WithDetails("unable to initialize global code generator with provided configuration").
					WithExtra("pattern", "./...").
					Build()
				errorHandler.AddError(genxErr)
				slog.Error("global generator creation failed", "error", errorHandler.FormatError(genxErr))

				// 显示错误摘要
				fmt.Println(errorHandler.FormatErrorList())
				os.Exit(1)
			}

			var wg conc.WaitGroup

			wg.Go(func() {
				globalX := gen.NewGlobalX(x, model)
				globalX.RegFunc(do.New())
				globalX.RegFunc(fx.New())
				globalX.Gen()
			})

			g := run.Group{}

			g.Add(func() error {

				wg.Wait()
				time.Sleep(1 * time.Second)
				model.Down()
				// fmt.Println("gen finish !!!")
				return nil
			}, func(err error) {
				if err != nil {
					panic(err)
				}
			})

			err = g.Run()

			if err != nil {
				panic(err)
			}

			return nil
		},
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {

		<-sigChan
		fmt.Println("signal received")
		os.Exit(0)
	}()

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
