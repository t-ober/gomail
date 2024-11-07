package app

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

func LogMessage(logger *slog.Logger, msg tea.Msg) {
	logger.Debug("received message",
		"msg_type", fmt.Sprintf("%T", msg),
		"msg", fmt.Sprintf("%+v", msg),
	)
}
