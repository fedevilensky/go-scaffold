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
	onEnter  func(selection string) error
	next     func() (tea.Model, tea.Cmd)
}

type RadioSelectOptions struct {
	Choices []string
	// Values is a map of the choices to their values
	Values  []string
	Header  string
	OnEnter func(selection string) error
	Next    func() (tea.Model, tea.Cmd)
}

func NewRadioSelect(options RadioSelectOptions) radioSelect {
	if len(options.Choices) == 0 {
		panic("no choices")
	}

	return radioSelect{
		choices:  options.Choices,
		values:   options.Values,
		cursor:   0,
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
		case "ctrl+c", "q":
			return rs, tea.Quit
		case "up", "k":
			if rs.cursor > 0 {
				rs.cursor--
			}
		case "down", "j":
			if rs.cursor < len(rs.choices)-1 {
				rs.cursor++
			}
		case " ":
			rs.selected = rs.choices[rs.cursor]
		case "enter":
			if err := rs.onEnter(rs.selected); err != nil {
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
		print = "Select one or more options:"
	}
	print += "\n\n"
	for i, choice := range rs.choices {
		cursor := " "
		if i == rs.cursor {
			cursor = ">"
		}

		checked := " "
		if choice == rs.selected {
			checked = "x"
		}

		print += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	print += HelpStyle("\n(press space to select, enter to continue, q to quit)\n")

	return print
}
