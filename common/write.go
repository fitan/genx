package common

import (
	"golang.org/x/exp/slog"
	"golang.org/x/tools/imports"
	"io/ioutil"
)

func WriteGO(name, s string) {
	processedSource, err := imports.Process(name, []byte(s), nil)
	if err != nil {
		slog.Error("imports.Process err", err, "genType", "trace")
		return
	}
	err = ioutil.WriteFile(name, processedSource, 0664)
	if err != nil {
		slog.Error("ioutil.WriteFile err", err, "gen file name", name)
		return
	}
}
