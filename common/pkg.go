package common

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

const mode packages.LoadMode = packages.NeedName |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedImports |
	//packages.NeedModule |
	//packages.NeedTypesSizes |
	//packages.NeedDeps |
	packages.NeedFiles

func LoadPkg(path string) ([]*packages.Package, error) {
	slog.Info("load pkg", slog.String("path", path))
	cfg := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		slog.Error("packages.Load err", slog.String("err", err.Error()))
		return nil, nil
	}

	if len(pkgs) < 1 {
		slog.Error("packages.Load err", slog.String("err", "len(pkgs) < 1"))
		return nil, fmt.Errorf("len(pkgs) < 1")
	}

	return pkgs, nil
}

func Last2DirName(s string) string {
	s = filepath.Clean(s)
	components := strings.Split(s, string(filepath.Separator))

	cLen := len(components)
	if cLen <= 1 {
		return components[0]
	}

	if cLen >= 2 {
		return components[cLen-2] + "." + components[cLen-1]
	}

	return ""
}
