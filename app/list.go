package app

import (
	"fmt"
	"gomail/list"
	"gomail/mail"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list        list.Model
	vp          viewportModel
	showingList bool
	logger      *slog.Logger
	dump        io.Writer
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				fmt.Printf("vp ready: %t", m.vp.ready)
				return m, nil //, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		if m.showingList {
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)
		}
		_ = m.UpdateAll(msg)
		return m, nil
	}

	var cmd tea.Cmd
	if m.showingList {
		m.list, cmd = m.list.Update(msg)
	} else {
		vpModel, vpCmd := m.vp.Update(msg)
		if vm, ok := vpModel.(viewportModel); ok {
			m.vp = vm
			cmd = vpCmd
		}
	}
	return m, cmd
}

func (m model) View() string {
	if m.showingList {
		return docStyle.Render(m.list.View())
	}
	m.logger.Debug("Running vp view", "ready", m.vp.ready)
	return m.vp.View()
}

type Updateable interface {
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
}

func (m *model) UpdateAll(msg tea.Msg) []tea.Cmd {
	// TODO: Simplification does not work since list does not satisfy tea.Model
	cmds := []tea.Cmd{}
	list, cmd := m.list.Update(msg)
	m.list = list
	cmds = append(cmds, cmd)

	vp, cmd := m.vp.Update(msg)
	m.vp = vp.(viewportModel)
	cmds = append(cmds, cmd)

	return cmds
}

func Run() {
	f, err := os.Create("gomail.log")
	if err != nil {
		log.Fatalf("Could not create log file: %v", err)
	}
	d, err := os.Create("messages.log")
	if err != nil {
		log.Fatalf("Could not create message log file: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	m := model{
		list:        list.New(getTestMails(), list.NewDefaultMailDelegate(), 0, 0),
		showingList: true,
		logger:      logger,
		dump:        d,
		vp: viewportModel{
			dump:   d,
			logger: logger,
		},
	}
	m.list.Title = "Gomail"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getTestMails() []list.Item {
	mails := []mail.Mail{
		{
			ID:        1,
			HistoryId: 1001,
			Received:  time.Date(2023, 10, 25, 14, 0, 0, 0, time.UTC),
			Sender:    "john@example.com",
			Recipient: "jane@example.com",
			Subject:   "Meeting Reminder",
			Msg:       "Hi Jane, just a reminder about our meeting tomorrow at 10 AM.",
		},
		{
			ID:        2,
			HistoryId: 1002,
			Received:  time.Date(2023, 10, 24, 9, 30, 0, 0, time.UTC),
			Sender:    "jane@example.com",
			Recipient: "john@example.com",
			Subject:   "Re: Meeting Reminder",
			Msg:       "Thanks for the reminder, John. I'll be there.",
		},
		{
			ID:        3,
			HistoryId: 1003,
			Received:  time.Date(2023, 10, 23, 16, 15, 0, 0, time.UTC),
			Sender:    "boss.company@example.com",
			Recipient: "john@example.com",
			Subject:   "Project Update",
			Msg:       "Please send me the project update by EOD.",
		},
		{
			ID:        4,
			HistoryId: 1004,
			Received:  time.Date(2023, 10, 22, 11, 0, 0, 0, time.UTC),
			Sender:    "jane@example.com",
			Recipient: "boss@example.com",
			Subject:   "Project Update Response",
			Msg:       "I've attached the latest project updates as requested.",
		},
	}
	// Convert []mail.Mail to []list.Item
	items := make([]list.Item, len(mails))
	for i, m := range mails {
		items[i] = m
	}

	return items
}
