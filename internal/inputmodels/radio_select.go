package inputmodels

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type radioSelect struct {
	choices  []string
	values   []string
	cursor   int
	selected string
	header   string
	onEnter  func(selection string, cursorPosition int) error
	next     func() (tea.Model, tea.Cmd)
}

type RadioSelectOptions struct {
	CursorStartingPosition int
	Choices                []string
	// Values is a map of the choices to their values
	Values  []string
	Header  string
	OnEnter func(selection string, cursorPosition int) error
	Next    func() (tea.Model, tea.Cmd)
}

func NewRadioSelect(options RadioSelectOptions) radioSelect {
	if len(options.Choices) == 0 {
		panic("no choices")
	}

	if len(options.Values) != len(options.Choices) {
		panic("values and choices must be the same length")
	}

	if options.CursorStartingPosition < 0 {
		options.CursorStartingPosition = 0
	} else if options.CursorStartingPosition >= len(options.Choices) {
		options.CursorStartingPosition = len(options.Choices) - 1
	}

	return radioSelect{
		choices:  options.Choices,
		values:   options.Values,
		cursor:   options.CursorStartingPosition,
		selected: options.Choices[0],
		header:   options.Header,
		onEnter:  options.OnEnter,
		next:     options.Next,
	}
}

func (rs radioSelect) Init() tea.Cmd {
	return nil
}

func (rs radioSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), "q":
			return rs, tea.Quit
		case tea.KeyUp.String(), "k":
			if rs.cursor > 0 {
				rs.cursor--
			}
		case tea.KeyHome.String():
			rs.cursor = 0
		case tea.KeyPgUp.String():
			rs.cursor -= 5
			if rs.cursor < 0 {
				rs.cursor = 0
			}
		case tea.KeyDown.String(), "j":
			if rs.cursor < len(rs.choices)-1 {
				rs.cursor++
			}
		case tea.KeyEnd.String():
			rs.cursor = len(rs.choices) - 1
		case tea.KeyDown.String():
			rs.cursor += 5
			if rs.cursor >= len(rs.choices) {
				rs.cursor = len(rs.choices) - 1
			}
		case tea.KeyEnter.String():
			if err := rs.onEnter(rs.values[rs.cursor], rs.cursor); err != nil {
				return rs, tea.Quit
			}
			return rs.next()
		}
	}
	return rs, nil
}

func (rs radioSelect) View() string {

	print := rs.header
	if print == "" {
		print = "Select one option:"
	}
	print += "\n\n"
	for i, choice := range rs.choices {
		cursor := " "
		if i == rs.cursor {
			cursor = ">"
		}

		print += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	print += HelpStyle("\n(press enter to continue, q to quit)\n")

	return print
}
