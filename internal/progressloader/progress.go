package progressloader

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/fedevilensky/go-scaffold/internal/inputmodels"
	"github.com/fedevilensky/go-scaffold/internal/project"
)

const (
	maxWidth = 120
	padding  = 4
)

type loader struct {
	conf     *project.Configuration
	progress progress.Model
}
type tickMsg time.Time

func NewLoader(conf *project.Configuration) tea.Model {
	return &loader{
		conf:     conf,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (l loader) Init() tea.Cmd {
	return tickCmd()
}

func (m loader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}
		cmd := m.progress.SetPercent(m.calculatePercent())
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		log.New(os.Stderr, "", 0).Println("unhandled msg:", msg)
		return m, nil
	}
}

func (m loader) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" + m.conf.GetCurrentCmd() + "\n" + pad + m.progress.View() + pad +
		"\n\n" + inputmodels.HelpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (l loader) calculatePercent() float64 {
	return l.conf.CalculateProgress()
}
