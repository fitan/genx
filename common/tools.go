package common

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"golang.org/x/tools/go/packages"
)

func init() {

}

func DownFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToLower(v)) + s[len(string(v)):]
	}
	return ""
}

func UpFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToUpper(v)) + s[len(string(v)):]
	}
	return ""
}

// 获取项目的根目录
func GetProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// 检查当前目录是否包含 go.mod 文件
		goModPath := filepath.Join(currentDir, "go.mod")
		_, err := os.Stat(goModPath)
		if err == nil {
			return currentDir, nil
		}

		// 到达文件系统的根目录时停止查找
		if currentDir == filepath.Dir(currentDir) {
			break
		}

		// 向上一级目录继续查找
		currentDir = filepath.Dir(currentDir)
	}

	return "", fmt.Errorf("未找到项目根目录")
}

func DepPkg(pkg *packages.Package, record map[string]*packages.Package) {
	for k, v := range pkg.Imports {
		if _, ok := record[k]; !ok {
			record[k] = v

			for _, v1 := range v.Imports {
				DepPkg(v1, record)
			}
		}
	}
	//var pkgs []*packages.Package
	//pkgs = append(pkgs, pkg)
	//for {
	//	var p *packages.Package
	//	if len(pkgs) == 0 {
	//		return
	//	}
	//	p = pkgs[0]
	//	pkgs = pkgs[1:]
	//
	//	slog.Info("load pkg dep pkg", "name", pkg.Name)
	//
	//	for _, v := range p.Imports {
	//		pkgs = append(pkgs, v)
	//	}
	//}
}

func GenFileByTemplate(efs embed.FS, fileName string, input any) (string, error) {
	f, err := fs.ReadFile(efs, fmt.Sprintf("static/template/%s.tmpl", fileName))
	if err != nil {
		err = fmt.Errorf("read template file failed: %w", err)
		return "", err
	}
	tt, err := template.New("template").Funcs(helperFuncs).Parse(string(f))
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tt.Execute(&buffer, input)

	return buffer.String(), err
}

var helperFuncs = template.FuncMap{
	"up":         strings.ToUpper,
	"down":       strings.ToLower,
	"upFirst":    upFirst,
	"downFirst":  downFirst,
	"replace":    strings.ReplaceAll,
	"snake":      toSnakeCase,
	"plural":     inflection.Plural,
	"camel":      strcase.ToCamel,
	"lowerCamel": strcase.ToLowerCamel,
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	result := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	result = matchAllCap.ReplaceAllString(result, "${1}_${2}")
	return strings.ToLower(result)
}

func downFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToLower(v)) + s[len(string(v)):]
	}
	return ""
}

func upFirst(s string) string {
	for _, v := range s {
		return string(unicode.ToUpper(v)) + s[len(string(v)):]
	}
	return ""
}
