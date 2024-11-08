package app

import (
	"gomail/list"
	"gomail/mail"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type listModel struct {
	list list.Model
}

func (m listModel) Init() tea.Cmd {
	return nil
}
func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		model, cmd := m.list.Update(msg)
		m.list = model
		return m, cmd
	}
	model, cmd := m.list.Update(msg)
	m.list = model
	return m, cmd
}

func (m listModel) View() string {
	return m.list.View()
}

func NewListModel() listModel {
	list := list.New(getTestMails(), list.NewDefaultMailDelegate(), 0, 0)
	list.Title = "Gomail"
	return listModel{
		list,
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
	items := make([]list.Item, len(mails))
	for i, m := range mails {
		items[i] = m
	}

	return items
}
