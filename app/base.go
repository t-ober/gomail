package app

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type baseModel struct {
	list        listModel
	vp          viewportModel
	showingList bool
	logger      *slog.Logger
}

func (m baseModel) Init() tea.Cmd {
	return nil
}

func (m baseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	LogMessage(m.logger, msg)
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if m.showingList == false {
				m.showingList = true
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.showingList {
				m.showingList = false
				content := "Hello World"
				m.vp.content = content
				m.vp.viewport.SetContent(content)
				return m, nil //, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		_ = m.UpdateAll(msg)
		return m, nil
	}

	var cmd tea.Cmd
	if m.showingList {
		cmd = m.UpdateList(msg)
	} else {
		cmd = m.UpdateVp(msg)
	}
	return m, cmd
}

func (m baseModel) View() string {
	if m.showingList {
		return docStyle.Render(m.list.View())
	}
	m.logger.Debug("Running vp view", "ready", m.vp.ready)
	return m.vp.View()
}

type Updateable interface {
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
}

func (m *baseModel) UpdateAll(msg tea.Msg) []tea.Cmd {
	cmds := []tea.Cmd{}
	lcmd := m.UpdateList(msg)
	vcmd := m.UpdateVp(msg)
	cmds = append(cmds, lcmd, vcmd)
	return cmds
}

func (m *baseModel) UpdateList(msg tea.Msg) tea.Cmd {
	model, cmd := m.list.Update(msg)
	if model, ok := model.(listModel); ok {
		m.list = model
	}
	return cmd
}

func (m *baseModel) UpdateVp(msg tea.Msg) tea.Cmd {
	model, cmd := m.vp.Update(msg)
	if model, ok := model.(viewportModel); ok {
		m.vp = model
	}
	return cmd
}

func Run() {
	f, err := os.Create("gomail.log")
	if err != nil {
		log.Fatalf("Could not create log file: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	m := baseModel{
		list:        NewListModel(),
		showingList: true,
		logger:      logger,
		vp: viewportModel{
			logger: logger,
		},
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
