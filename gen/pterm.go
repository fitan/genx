package gen

import (
	"time"

	"github.com/pterm/pterm"
	"github.com/samber/lo"
)

type TUI struct {
	tree pterm.TreeNode

	down bool
}

func NewTUI() *TUI {
	return &TUI{
		tree: pterm.TreeNode{
			Children: []pterm.TreeNode{},
			Text:     "gen file work",
		},
	}
}
func (t *TUI) HasFileName(req TUIAddReq) (fileNode pterm.TreeNode, fileIndex int, fileOk bool) {
	plugNode, _, plugOk := t.HasPlugName(req)
	if plugOk {
		return lo.FindIndexOf(plugNode.Children, func(
			node pterm.TreeNode,
		) bool {
			return node.Text == req.FileName
		})
	}

	return
}

func (t *TUI) HasPlugName(req TUIAddReq) (plugNode pterm.TreeNode, plugIndex int, plugOk bool) {
	pkgNode, _, pkgOk := t.HasPkgName(req)

	if pkgOk {
		return lo.FindIndexOf(pkgNode.Children, func(
			node pterm.TreeNode,
		) bool {
			return node.Text == req.PlugName
		})
	}

	return
}

func (t *TUI) HasPkgName(req TUIAddReq) (pkgNode pterm.TreeNode, pkgIndex int, pkgOk bool) {
	return lo.FindIndexOf(t.tree.Children, func(
		node pterm.TreeNode,
	) bool {
		return node.Text == req.PkgName
	})
}

func (t *TUI) Update(req TUIAddReq) error {
	pkgNode, pkgIndex, pkgOk := t.HasPkgName(req)

	if !pkgOk {
		pkgNode = pterm.TreeNode{
			Text: req.PkgName,
		}
		t.tree.Children = append(t.tree.Children, pkgNode)
	}

	plugNode, plugIndex, plugOk := t.HasPlugName(req)

	if !plugOk {
		plugNode = pterm.TreeNode{
			Text: req.PlugName,
		}
		t.tree.Children[pkgIndex].Children = append(pkgNode.Children, plugNode)
	}

	fileNode, fileIndex, fileOk := t.HasFileName(req)

	if !fileOk {
		fileNode = pterm.TreeNode{
			Text: req.FileName,
		}
		t.tree.Children[pkgIndex].Children[plugIndex].Children = append(plugNode.Children, fileNode)
	} else {
		t.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex].Text = req.FileName
	}

	return nil
}

func (t *TUI) Start() error {
	area, err := pterm.DefaultArea.Start()
	if err != nil {
		return err
	}

	for !t.down {
		area.Update(pterm.DefaultTree.WithRoot(t.tree).Srender())
		time.Sleep(250 * time.Millisecond)
	}

	return area.Stop()
}

type TUIAddReq struct {
	PkgName       string
	PlugName      string
	FileName      string
	GenFileStatus string
}
