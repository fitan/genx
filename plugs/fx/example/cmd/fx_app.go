// Code generated . DO NOT EDIT.
package main

import "go.uber.org/fx"

func FxOptions() []fx.Option {
	return []fx.Option{}
}

func NewApp(opts ...fx.Option) *fx.App {
	allOpts := append(FxOptions(), opts...)
	return fx.New(allOpts...)
}
