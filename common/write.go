package common

import (
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
	"golang.org/x/tools/imports"
)

func WriteGO(name, s string) {
	processedSource, err := imports.Process(name, []byte(s), nil)
	if err != nil {
		slog.Error("imports.Process err", err, "genType", "trace")
		return
	}
	err = os.MkdirAll(filepath.Dir(name), os.ModePerm)
	if err != nil {
		slog.Error("os.MkdirAll err", err, "genType", "trace")
		return
	}
	header := `// Code generated . DO NOT EDIT.
`

	fileBody := []byte(header + string(processedSource))

	err = os.WriteFile(name, []byte(fileBody), 0664)
	if err != nil {
		slog.Error("ioutil.WriteFile err", err, "gen file name", name)
		return
	}
}

type WriteOpt struct {
	Cover bool
}

func WriteGoWithOpt(name, s string, opt WriteOpt) (cover bool, err error) {
	if opt.Cover {
		WriteGO(name, s)
		return true, nil
	} else {
		_, err = os.Stat(name)
		if err != nil {
			if os.IsNotExist(err) {
				WriteGO(name, s)
				return false, nil
			} else {
				slog.Error("os.Stat err", err, "gen file name", name)
				return false, err
			}
		}

		return false, nil
	}
}
