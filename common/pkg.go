package common

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"strings"
)
import "golang.org/x/exp/slog"

const mode packages.LoadMode = packages.NeedName |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedImports |
	//packages.NeedModule |
	//packages.NeedTypesSizes |
	//packages.NeedDeps |
	packages.NeedFiles

func LoadPkg(path string) (*packages.Package, error) {
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

	return pkgs[0], nil
}

func Last2DirName(s string) string {
	pathList := strings.Split(s, "/")
	return strings.Join(pathList[len(pathList)-2:len(pathList)-1], ".")
}
