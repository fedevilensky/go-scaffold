package inputmodels

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textInput struct {
	textInput    textinput.Model
	header       string
	placeholder  string
	onEnter      func(input string) error
	onEnterEmpty func() error
	next         func() (tea.Model, tea.Cmd)
	nextEmpty    func() (tea.Model, tea.Cmd)
}

type TextInputOptions struct {
	Header       string
	Placeholder  string
	OnEnter      func(input string) error
	OnEnterEmpty func() error
	Next         func() (tea.Model, tea.Cmd)
	NextEmpty    func() (tea.Model, tea.Cmd)
}

func NewTextInput(options TextInputOptions) textInput {
	ti := textinput.New()
	ti.Placeholder = options.Placeholder
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 80
	return textInput{
		textInput:    ti,
		header:       options.Header,
		onEnter:      options.OnEnter,
		onEnterEmpty: options.OnEnterEmpty,
		next:         options.Next,
		nextEmpty:    options.NextEmpty,
	}
}

func (txt textInput) Init() tea.Cmd {
	return textinput.Blink
}

func (txt textInput) View() string {
	return fmt.Sprintf(
		"%s:\n\n%s\n\n%s\n",
		txt.header,
		txt.textInput.View(),
		HelpStyle("(esc to quit)"),
	)
}

func (txt textInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		log.Fatal(msg)
		return txt, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if txt.textInput.Value() == "" {
				if err := txt.onEnterEmpty(); err != nil {
					log.Fatal(err)
				}
				return txt.nextEmpty()
			} else {
				if err := txt.onEnter(txt.textInput.Value()); err != nil {
					log.Fatal(err)
				}
				return txt.next()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return txt, tea.Quit
		}
	}

	var cmd tea.Cmd
	txt.textInput, cmd = txt.textInput.Update(msg)
	return txt, cmd
}
