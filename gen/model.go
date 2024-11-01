package gen

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

type TreeNode struct {
	Text     string
	Children []TreeNode
	// 0 running
	// 1 success
	// 2 failed
	// 3 exist
	Status    int
	Err       string
	StartTime time.Time
	EndTime   *time.Time
}

func (f TreeNode) RenderTree(s string) string {
	var b strings.Builder

	b.WriteString(f.Text + "\n")
	b.WriteString(f.renderTree(f.Children, "", s))

	return b.String()

}

func (f TreeNode) renderTree(nodes []TreeNode, prefix string, s string) string {
	var b strings.Builder

	for i, node := range nodes {

		connector := "├── "

		if i == len(nodes)-1 { // Last element in this level.
			connector = "└── "
		}

		b.WriteString(fmt.Sprintf("%s%s%s %s\n", prefix, connector, node.Text, lo.Switch[int, string](node.Status).
			CaseF(0, func() string { return s + " " + time.Since(node.StartTime).String() }).
			CaseF(1, func() string { return healthyStyle.Render("✔") + " " + node.EndTime.Sub(node.StartTime).String() }).
			CaseF(2, func() string { return errStyle.Render("✗") + " " + node.EndTime.Sub(node.StartTime).String() }).
			CaseF(3, func() string { return unknownStyle.Render("exist") + " " + node.EndTime.Sub(node.StartTime).String() }).
			Default(unknownStyle.Render("unknown"))))

		if len(node.Children) > 0 { // Recursively render children with increased indentation.

			newPrefix := prefix + "│   "

			if i == len(nodes)-1 { // If it's the last element at this level.
				newPrefix = prefix + "    "

			}

			b.WriteString(f.renderTree(node.Children, newPrefix, s))

		}
	}

	return b.String()
}

type Model struct {
	tree *TreeNode
	// MsgMap   map[string]ModelMsg
	MsgChan <-chan ModelMsg
	Spinner spinner.Model
	m       sync.Mutex
	down    chan struct{}
	endDown chan struct{}
	teaCmds []tea.Cmd
}

func NewModel() *Model {
	return &Model{
		// MsgMap:   map[string]ModelMsg{},
		tree: &TreeNode{
			Text:     "generate",
			Children: []TreeNode{},
			Status:   0,
			Err:      "",
		},
		MsgChan: make(<-chan ModelMsg),
		Spinner: spinner.New(),
		m:       sync.Mutex{},
		down:    make(chan struct{}),
		endDown: make(chan struct{}),
		teaCmds: make([]tea.Cmd, 0),
	}
}

type ModelMsg struct {
	PkgName  string
	PlugName string
	Msg      string
}

type ModelDown struct{}

func (m *Model) Down() struct{} {
	m.down <- struct{}{}

	return <-m.endDown
}

func (m *Model) waitDown() tea.Msg {
	<-m.down
	return ModelDown{}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, m.waitDown)
}

func (m *Model) PkgEnd(req UpdateTreeReq) {
	m.m.Lock()
	defer m.m.Unlock()
	index, ok := m.hasNode(m.tree.Children, req.PkgName)
	if !ok {
		return
	}
	m.tree.Children[index].Status = 1
	m.tree.Children[index].EndTime = lo.ToPtr(time.Now())
}

func (m *Model) hasNode(nodes []TreeNode, name string) (int, bool) {
	_, index, ok := lo.FindIndexOf(nodes, func(item TreeNode) bool {
		return item.Text == name
	})

	return index, ok
}

func (m *Model) PlugStart(req UpdateTreeReq) {
	m.m.Lock()
	defer m.m.Unlock()
	pkgIndex, pkgOk := m.hasNode(m.tree.Children, req.PkgName)
	if !pkgOk {
		m.tree.Children = append(m.tree.Children, TreeNode{Text: req.PkgName, Status: 0, Children: []TreeNode{}, Err: req.Err, StartTime: time.Now()})
		pkgIndex = len(m.tree.Children) - 1
	}

	_, plugOk := m.hasNode(m.tree.Children[pkgIndex].Children, req.PlugName)
	if !plugOk {
		m.tree.Children[pkgIndex].Children = append(m.tree.Children[pkgIndex].Children, TreeNode{
			Text:      req.PlugName,
			Children:  []TreeNode{},
			Status:    0,
			Err:       req.Err,
			StartTime: time.Now(),
		})
	}

}

func (m *Model) PlugEnd(req UpdateTreeReq) {
	m.m.Lock()
	defer m.m.Unlock()

	pkgIndex, _ := m.hasNode(m.tree.Children, req.PkgName)

	plugIndex, _ := m.hasNode(m.tree.Children[pkgIndex].Children, req.PlugName)

	m.tree.Children[pkgIndex].Children[plugIndex].Status = req.Status
	m.tree.Children[pkgIndex].Children[plugIndex].Err = req.Err
	m.tree.Children[pkgIndex].Children[plugIndex].EndTime = lo.ToPtr(time.Now())
}

func (m *Model) FileStart(req UpdateTreeReq) {
	m.m.Lock()
	defer m.m.Unlock()

	pkgIndex, _ := m.hasNode(m.tree.Children, req.PkgName)
	plugIndex, _ := m.hasNode(m.tree.Children[pkgIndex].Children, req.PlugName)

	m.tree.Children[pkgIndex].Children[plugIndex].Children = append(m.tree.Children[pkgIndex].Children[plugIndex].Children, TreeNode{
		Text:      req.FileName,
		Children:  []TreeNode{},
		Status:    0,
		Err:       req.Err,
		StartTime: time.Now(),
	})
}

func (m *Model) FileEnd(req UpdateTreeReq) {
	m.m.Lock()
	defer m.m.Unlock()
	pkgIndex, _ := m.hasNode(m.tree.Children, req.PkgName)
	plugIndex, _ := m.hasNode(m.tree.Children[pkgIndex].Children, req.PlugName)
	fileIndex, _ := m.hasNode(m.tree.Children[pkgIndex].Children[plugIndex].Children, req.FileName)

	fileNode := m.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex]

	switch req.Status {
	case 3:
		m.teaCmds = append(m.teaCmds, tea.Printf("%s %s %s",
			existStyle.Render("exist"),
			req.FileName,
			time.Since(fileNode.StartTime).String()))
	case 1:
		m.teaCmds = append(m.teaCmds, tea.Printf("%s %s %s",
			healthyStyle.Render("✔"),
			req.FileName,
			time.Since(fileNode.StartTime).String()))
	case 2:
		m.teaCmds = append(m.teaCmds, tea.Printf("%s %s %s",
			errStyle.Render("✗"),
			req.FileName,
			time.Since(fileNode.StartTime).String()))
	default:
		m.teaCmds = append(m.teaCmds, tea.Printf("%s %s %s",
			unknownStyle.Render("unknown"),
			req.FileName,
			time.Since(fileNode.StartTime).String()))
	}

	m.tree.Children[pkgIndex].Children[plugIndex].Children = lo.DropByIndex(m.tree.Children[pkgIndex].Children[plugIndex].Children, fileIndex)

	/* m.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex].Status = req.Status
	m.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex].Err = req.Err
	m.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex].EndTime = lo.ToPtr(time.Now()) */
}

func (m *Model) View() string {
	return baseStyle.Render(m.tree.RenderTree(m.Spinner.View()))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case spinner.TickMsg:
		m.m.Lock()
		defer m.m.Unlock()
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		// m.TickUpdateMsgTime()
		cmds := append(m.teaCmds, cmd)
		m.teaCmds = make([]tea.Cmd, 0)
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch v.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		return m, nil
	case tea.QuitMsg:
		return m, tea.Quit
	case ModelDown:
		// m.TickUpdateMsgTime()
		defer func() {
			m.endDown <- struct{}{}
		}()
		return m, tea.Quit
	default:
		return m, nil
	}
}

var baseStyle = lipgloss.NewStyle()

/* 	BorderStyle(lipgloss.NormalBorder()).
BorderForeground(lipgloss.Color("240") )*/

var healthyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
var errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
var unknownStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
var existStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

type UpdateTreeReq struct {
	PkgName  string
	PlugName string
	FileName string
	Status   int
	Err      string
}
