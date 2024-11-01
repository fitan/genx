/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"
	"log/slog"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fitan/genx/gen"
	plugCopy "github.com/fitan/genx/plugs/copy"
	"github.com/fitan/genx/plugs/crud"
	"github.com/fitan/genx/plugs/do"
	"github.com/fitan/genx/plugs/enum"
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

			x, err := gen.NewX(staticFiles, "./...", model)
			if err != nil {
				slog.Error("new x error", "err", err)
				os.Exit(1)
			}

			var wg conc.WaitGroup

			lo.ForEach(x, func(item *gen.X, index int) {
				wg.Go(func() {
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
					item.RegImpl(&kithttp.ObserverPlug{})
					item.Gen()
				})
			})

			wg.Go(func() {
				globalX := gen.NewGlobalX(x, model)
				globalX.RegFunc(do.New())
				globalX.Gen()
			})

			g := run.Group{}

			signal, cancel := run.SignalHandler(ctx.Context, os.Interrupt)

			g.Add(signal, func(err error) {
				if err != nil {
					panic(err)
				}
			})
			g.Add(func() error {
				p := tea.NewProgram(model)

				if _, err := p.Run(); err != nil {
					return err
				}

				return nil
			}, func(err error) {
				if err != nil {
					panic(err)
				}
			})

			g.Add(func() error {
				wg.Wait()
				time.Sleep(time.Second * 3)
				model.Down()
				cancel(nil)
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

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
