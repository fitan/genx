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
	index, ok := m.hasPkg(req)
	if !ok {
		return
	}
	m.tree.Children[index].Status = 1
	m.tree.Children[index].EndTime = lo.ToPtr(time.Now())
}

func (m *Model) hasPkg(req UpdateTreeReq) (int, bool) {
	_, index, ok := lo.FindIndexOf(m.tree.Children, func(item TreeNode) bool {
		return item.Text == req.PkgName
	})

	return index, ok
}

func (m *Model) PlugStart(req UpdateTreeReq) int {
	m.m.Lock()
	defer m.m.Unlock()
	index, ok := m.hasPkg(req)
	if !ok {
		m.tree.Children = append(m.tree.Children, TreeNode{Text: req.PkgName, Status: 0, Children: []TreeNode{}, Err: req.Err, StartTime: time.Now()})
		index = len(m.tree.Children) - 1
	}

	m.tree.Children[index].Children = append(m.tree.Children[index].Children, TreeNode{
		Text:      req.PlugName,
		Children:  []TreeNode{},
		Status:    0,
		Err:       req.Err,
		StartTime: time.Now(),
	})

	return len(m.tree.Children[index].Children) - 1
}

func (m *Model) PlugEnd(plugIndex int, req UpdateTreeReq) {
	m.m.Lock()
	defer m.m.Unlock()

	pkgIndex, _ := m.hasPkg(req)

	m.tree.Children[pkgIndex].Children[plugIndex].Status = req.Status
	m.tree.Children[pkgIndex].Children[plugIndex].Err = req.Err
	m.tree.Children[pkgIndex].Children[plugIndex].EndTime = lo.ToPtr(time.Now())
}

func (m *Model) FileStart(plugIndex int, req UpdateTreeReq) int {
	m.m.Lock()
	defer m.m.Unlock()

	pkgIndex, _ := m.hasPkg(req)

	m.tree.Children[pkgIndex].Children[plugIndex].Children = append(m.tree.Children[pkgIndex].Children[plugIndex].Children, TreeNode{
		Text:      req.FileName,
		Children:  []TreeNode{},
		Status:    0,
		Err:       req.Err,
		StartTime: time.Now(),
	})

	return len(m.tree.Children[pkgIndex].Children[plugIndex].Children) - 1
}

func (m *Model) FileEnd(plugIndex, fileIndex int, req UpdateTreeReq) {
	m.m.Lock()
	defer m.m.Unlock()
	pkgIndex, _ := m.hasPkg(req)

	m.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex].Status = req.Status
	m.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex].Err = req.Err
	m.tree.Children[pkgIndex].Children[plugIndex].Children[fileIndex].EndTime = lo.ToPtr(time.Now())
}

func (m *Model) View() string {
	return baseStyle.Render(m.tree.RenderTree(m.Spinner.View()))
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		// m.TickUpdateMsgTime()
		return m, cmd
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

var healthyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
var errStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
var unknownStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("8"))

type UpdateTreeReq struct {
	PkgName  string
	PlugName string
	FileName string
	Status   int
	Err      string
}
