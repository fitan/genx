package gen

import (
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// MsgMap   map[string]ModelMsg
	MsgIndex map[string]MsgIndexMsg
	MsgChan  <-chan ModelMsg
	Spinner  spinner.Model
	m        sync.Mutex
	down     chan struct{}
	endDown  chan struct{}
	table    table.Model
}

type MsgIndexMsg struct {
	index     int
	startTime time.Time
	status    string
}

func NewModel() *Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)
	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle()
	s.Cell = lipgloss.NewStyle().Padding(0, 1)
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
		/* 	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false */
	t.SetStyles(s)
	return &Model{
		// MsgMap:   map[string]ModelMsg{},
		MsgIndex: make(map[string]MsgIndexMsg),
		MsgChan:  make(<-chan ModelMsg),
		Spinner:  spinner.New(),
		m:        sync.Mutex{},
		down:     make(chan struct{}),
		table:    t,
		endDown:  make(chan struct{}),
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

func (m *Model) AddMsg(key string, msg ModelMsg) {
	m.m.Lock()
	defer m.m.Unlock()
	rows := m.table.Rows()
	rows = append(rows, []string{"", msg.PkgName, msg.PlugName, "0s"})
	m.MsgIndex[fmt.Sprintf("%s/%s", msg.PkgName, msg.PlugName)] = MsgIndexMsg{
		index:     len(rows) - 1,
		startTime: time.Now(),
		status:    "running",
	}
	m.table.SetRows(rows)
}

func (m *Model) TickUpdateMsgTime() {
	m.m.Lock()
	defer m.m.Unlock()
	rows := m.table.Rows()
	var res []table.Row
	for _, item := range rows {
		indexMsg := m.MsgIndex[fmt.Sprintf("%s/%s", item[1], item[2])]
		if indexMsg.status != "success" {
			item[0] = m.Spinner.View()
			item[3] = time.Since(indexMsg.startTime).String()
		} else {
			item[0] = healthyStyle.Render("✔")
		}

		res = append(res, item)
	}

	m.table.SetRows(res)
}

func (m *Model) UpdateMsg(key string, msg ModelMsg) {
	m.m.Lock()
	defer m.m.Unlock()
	index := fmt.Sprintf("%s/%s", msg.PkgName, msg.PlugName)
	indexMsg := m.MsgIndex[index]
	indexMsg.status = msg.Msg
	m.MsgIndex[index] = indexMsg
}

func (m *Model) View() string {
	m.m.Lock()
	defer m.m.Unlock()
	return baseStyle.Render(m.table.View()) + "\n  "
	/* var s []string
	for _, key := range m.MsgIndex {

		if m.MsgMap[key].Msg != "success" {
			msg := m.MsgMap[key]
			msg.endTime = time.Since(msg.startTime)
			m.MsgMap[key] = msg
		}

		s = append(s, fmt.Sprintf("%s %s/%s: %s %v", m.Spinner.View(), m.MsgMap[key].PkgName, m.MsgMap[key].PlugName, m.MsgMap[key].Msg, m.MsgMap[key].endTime))
	}

	return strings.Join(s, "\n") + "\n" */

}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		m.TickUpdateMsgTime()
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
		m.TickUpdateMsgTime()
		defer func() {
			m.endDown <- struct{}{}
		}()
		return m, tea.Quit
	default:
		return m, nil
	}
}

var columns = []table.Column{
	table.Column{
		Title: "Status",
		Width: 10,
	},
	{
		Title: "PkgPath",
		Width: 50,
	},
	{
		Title: "TagName",
		Width: 10,
	},
	table.Column{
		Title: "Time",
		Width: 20,
	},
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var healthyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
