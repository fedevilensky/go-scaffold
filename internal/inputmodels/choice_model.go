package inputmodels

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type Selection struct {
	IsSelected bool
	value      string
}

func (s Selection) Value() string {
	return s.value
}

type choiceModel struct {
	choices    []string
	cursor     int
	selections []Selection
	header     string
	next       func() (tea.Model, tea.Cmd)
	onEnter    func(selected []Selection) error
}

type ChoiceModelOptions struct {
	Choices []string
	// the value you wish to receive for every choice, in the same order
	Values  []string
	Header  string
	Next    func() (tea.Model, tea.Cmd)
	OnEnter func(selected []Selection) error
}

func NewchoiceModel(opts ChoiceModelOptions) choiceModel {
	if len(opts.Choices) == 0 {
		panic("no choices")
	}
	if len(opts.Values) != len(opts.Choices) {
		panic("values and choices must be the same length")
	}
	selected := make([]Selection, len(opts.Values))
	for i, value := range opts.Values {
		selected[i] = Selection{false, value}
	}
	return choiceModel{
		choices:    opts.Choices,
		cursor:     0,
		selections: selected,
		header:     opts.Header,
		next:       opts.Next,
		onEnter:    opts.OnEnter,
	}
}
func (m choiceModel) Init() tea.Cmd {
	return nil
}

func (choice choiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return choice, tea.Quit
		case "up", "k":
			if choice.cursor > 0 {
				choice.cursor--
			}
		case "down", "j":
			if choice.cursor < len(choice.choices)-1 {
				choice.cursor++
			}
		case " ":
			choice.selections[choice.cursor].IsSelected = !choice.selections[choice.cursor].IsSelected
		case "enter":
			if err := choice.onEnter(choice.selections); err != nil {
				log.Fatal(err)
			}
			return choice.next()
		}
	}
	return choice, nil
}

func (m choiceModel) View() string {
	print := m.header
	if print == "" {
		print = "Select one or more options:"
	}
	print += "\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		checked := " "
		if m.selections[i].IsSelected {
			checked = "x"
		}

		print += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	print += HelpStyle("\n(press space to select, enter to continue, q to quit)\n")

	return print
}
