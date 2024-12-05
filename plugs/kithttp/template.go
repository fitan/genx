package kithttp

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"go/ast"
	"math"
	"path"
	"strings"
	"time"

	"github.com/fitan/genx/common"
	"github.com/fitan/genx/gen"
	"github.com/fitan/genx/parser"
	"github.com/samber/lo"
)

type TemplateInputInterface struct {
	Name    string
	Methods map[string]Method
	Doc     common.Doc
	RawDoc  ast.CommentGroup

	Opt gen.Option
}

func (t *TemplateInputInterface) Instance() string {
	return common.Last2DirName(t.Opt.Dir)
}

func (t *TemplateInputInterface) KitServerOption() (res string) {
	line := t.Doc.ByFuncName("@kit-server-option")

	if line == nil {
		return ""
	}

	return strings.Join(lo.Map(line.Args, func(item parser.FuncArg, index int) string {
		return item.Value
	}), ",")
}

func (t *TemplateInputInterface) Tags() string {
	var tag string
	t.Doc.ByFuncNameAndArgs("@tags", &tag)

	return tag
}

func (t *TemplateInputInterface) ValidVersion() string {
	var validVersion string
	t.Doc.ByFuncNameAndArgs("@validVersion", &validVersion)

	return validVersion
}

func (t *TemplateInputInterface) BasePath() string {
	var basePath string
	t.Doc.ByFuncNameAndArgs("@basePath", &basePath)

	return basePath
}

func (t *TemplateInputInterface) EnableSwag(name string) bool {
	var swag string
	t.Doc.ByFuncNameAndArgs("@swag", &swag)
	if swag == "false" {
		return false
	}

	return t.Methods[name].EnableSwag()
}

func (t *TemplateInputInterface) HasMethodPath(name string) bool {
	return t.Methods[name].RawKit.Conf.Url != ""
}

func (t *TemplateInputInterface) MethodPath(name string) string {
	return strings.TrimSuffix(path.Join(t.BasePath(), t.Methods[name].RawKit.Conf.Url), "/")
}

func (t TemplateInputInterface) hashToID(s string) int64 {
	hash := sha256.Sum256([]byte(s))
	return int64(binary.BigEndian.Uint64(hash[:8]) % uint64(int64(math.Pow(10, 15))))
}

func (t TemplateInputInterface) Annotation() string {
	if t.Doc == nil {
		return ""
	}
	for _, c := range t.RawDoc.List {
		docFormat := DocFormat(c.Text)
		if strings.HasPrefix(docFormat, "// "+t.Name) {
			return strings.TrimPrefix(docFormat, "// "+t.Name)
		}
	}
	return strings.TrimPrefix(DocFormat(t.RawDoc.List[0].Text), "// ")
}

func (t TemplateInputInterface) CEPermissionSql() string {
	var (
		parentID          int64
		parentIcon        string
		parentMenu        int
		parentMethod      string
		parentAlias       string
		parentName        string
		parentPath        string
		parentDescription string
		sqls              []string
	)

	//INSERT INTO spider_dev.sys_permission (id, parent_id, icon, menu, method, alias, name, path, description, created_at, updated_at, deleted_at) VALUES (878, 877, '', 1, 'GET', 'Redis实例', 'menu.cdb.redis', '/cdb/index', '', '2022-12-29 14:21:22', '2022-12-29 14:21:22', null);

	parentMenu = 1
	parentMethod = "GET"
	parentAlias = t.Annotation()
	parentName = strings.ToLower(strings.Join([]string{"menu", strings.Trim(strings.Replace(t.BasePath(), "/", ".", -1), "."), t.Name}, "."))
	parentPath = t.BasePath() + "/index"
	parentDescription = t.Annotation()
	parentID = t.hashToID(parentName)
	parentSql := fmt.Sprintf(`INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description,created_at, updated_at) VALUES (%v, %v, '%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v');`,
		parentID, 0, parentIcon, parentMenu, parentMethod, strings.TrimSpace(parentAlias), parentName, parentPath, strings.TrimSpace(parentDescription), time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	sqls = append(sqls, parentSql)

	for _, v := range t.Methods {
		var (
			id          int64
			icon        string
			menu        int
			method      string
			alias       string
			name        string
			mPath       string
			description string
		)

		mPath = t.MethodPath(v.Name)
		menu = 0
		method = strings.ToUpper(v.RawKit.Conf.UrlMethod)
		if method == "" {
			continue
		}
		alias = v.Annotation()
		name = strings.Trim(strings.ToLower(strings.Join([]string{strings.Replace(t.BasePath(), "/", ".", -1), v.Name, method}, ".")), ".")
		description = v.Annotation()
		id = t.hashToID(name)
		sql := fmt.Sprintf(`INSERT INTO sys_permission (id, parent_id, icon, menu, method, alias, name, path, description,created_at, updated_at) VALUES (%v, %v, '%v', %v, '%v', '%v', '%v', '%v', '%v', '%v', '%v');`,
			id, parentID, icon, menu, method, strings.TrimSpace(alias), name, mPath, strings.TrimSpace(description), time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

		sqls = append(sqls, sql)
	}

	return strings.Join(sqls, "\n")
}
