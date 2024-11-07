package list

import (
	"fmt"
	"gomail/mail"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type DefaultMailItemStyles struct {
	NormalText   lipgloss.Style
	SelectedText lipgloss.Style
	DimmedText   lipgloss.Style
	FilterMatch  lipgloss.Style
}

func NewDefaultMailItemStyles() (s DefaultMailItemStyles) {
	s.NormalText = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)
	s.SelectedText = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)
	s.DimmedText = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)
	s.FilterMatch = lipgloss.NewStyle().Underline(true)
	return s
}

type DefaultMailDelegate struct {
	Styles     DefaultMailItemStyles
	UpdateFunc func(tea.Msg, *Model) tea.Cmd
	height     int
	spacing    int
}

func NewDefaultMailDelegate() DefaultMailDelegate {
	return DefaultMailDelegate{
		Styles: NewDefaultMailItemStyles(),
	}
}

// SetHeight sets delegate's preferred height.
func (d *DefaultMailDelegate) SetHeight(i int) {
	d.height = i
}

// Height returns the delegate's preferred height.
// This has effect only if ShowDescription is true,
// otherwise height is always 1.
func (d DefaultMailDelegate) Height() int {
	return 1
}

// SetSpacing sets the delegate's spacing.
func (d *DefaultMailDelegate) SetSpacing(i int) {
	d.spacing = i
}

// Spacing returns the delegate's spacing.
func (d DefaultMailDelegate) Spacing() int {
	return d.spacing
}

// Update checks whether the delegate's UpdateFunc is set and calls it.
func (d DefaultMailDelegate) Update(msg tea.Msg, m *Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

func (d DefaultMailDelegate) Render(w io.Writer, m Model, index int, item Item) {
	var (
		sender, subject string
		matchedRunes    []int
		s               = &d.Styles
	)

	if i, ok := item.(mail.Mail); ok {
		sender = i.Sender
		subject = i.Subject
	}

	if m.width <= 0 {
		return
	}

	textwidth := m.width - s.NormalText.GetPaddingLeft() - s.NormalText.GetPaddingRight()
	separator := "  "
	separatorWidth := lipgloss.Width(separator)

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == Filtering || m.FilterState() == FilterApplied
	)

	var senderMatches, subjectMatches []int
	if isFiltered && index < len(m.filteredItems) {
		matchedRunes = m.MatchesForItem(index)
		// Split matched runes between sender and subject
		senderLen := lipgloss.Width(sender)
		for _, pos := range matchedRunes {
			if pos < senderLen {
				senderMatches = append(senderMatches, pos)
			} else {
				// Adjust position for subject (accounting for separator)
				subjectMatches = append(subjectMatches, pos-senderLen-separatorWidth+lipgloss.Width(separator)-1)
			}
		}
	}

	senderWidth, renderedSender := renderSender(textwidth, sender, senderMatches, s, isSelected && m.FilterState() != Filtering)
	subjectWidth := textwidth - senderWidth - separatorWidth
	subject = ansi.Truncate(subject, subjectWidth, ellipsis)

	// Apply subject matches if needed
	if len(subjectMatches) > 0 {
		baseStyle := s.NormalText
		if isSelected && m.FilterState() != Filtering {
			baseStyle = s.SelectedText
		}
		unmatched := baseStyle.Inline(true)
		matched := unmatched.Inherit(s.FilterMatch)
		subject = lipgloss.StyleRunes(subject, subjectMatches, matched, unmatched)
	}

	line := fmt.Sprintf("%s%s%s", renderedSender, separator, subject)

	// Apply final styling
	if emptyFilter {
		line = d.Styles.DimmedText.Render(line)
	} else if isSelected && m.FilterState() != Filtering {
		line = s.SelectedText.Render(line)
	} else {
		line = s.NormalText.Render(line)
	}

	fmt.Fprint(w, line)
}
func renderSender(textwidth int, sender string, matchedRunes []int, styles *DefaultMailItemStyles, isSelected bool) (int, string) {
	maxWidth := 30
	proportionalWidth := int(float32(textwidth) * (1.0 / 3.0))
	width := min(proportionalWidth, maxWidth)

	// First truncate if necessary
	if lipgloss.Width(sender) > width {
		sender = ansi.Truncate(sender, width, "...")
	}

	// Apply matched runes styling before padding
	if len(matchedRunes) > 0 {
		baseStyle := styles.NormalText
		if isSelected {
			baseStyle = styles.SelectedText
		}
		unmatched := baseStyle.Inline(true)
		matched := unmatched.Inherit(styles.FilterMatch)
		sender = lipgloss.StyleRunes(sender, matchedRunes, matched, unmatched)
	}

	// Finally apply padding
	return width, lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left).
		Render(sender)
}
