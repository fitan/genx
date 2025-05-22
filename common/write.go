package common

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/imports"
)

func WriteRaw(name, s string) (err error) {
	err = os.MkdirAll(filepath.Dir(name), os.ModePerm)
	if err != nil {
		slog.Error("os.MkdirAll err", err, "genType", "trace")
		return
	}

	err = os.WriteFile(name, []byte(s), 0664)
	if err != nil {
		slog.Error("ioutil.WriteFile err", err, "gen file name", name)
		return
	}
	return nil
}

func WriteGO(name, s string) (err error) {
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
	return nil
}

type WriteOpt struct {
	Cover bool
	Raw   bool
}

func WriteGoWithOpt(name, s string, opt WriteOpt) (cover bool, err error) {
	if opt.Cover {
		err = lo.Ternary(opt.Raw, WriteRaw, WriteGO)(name, s)
		return true, err
	} else {
		_, err = os.Stat(name)
		if err != nil {
			if os.IsNotExist(err) {
				err = lo.Ternary(opt.Raw, WriteRaw, WriteGO)(name, s)
				return false, err
			} else {
				slog.Error("os.Stat err", err, "gen file name", name)
				return false, err
			}
		}

		return false, nil
	}
}
