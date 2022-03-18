package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fedevilensky/go-scaffold/inputmodels"
	"github.com/fedevilensky/go-scaffold/progressloader"
	"github.com/fedevilensky/go-scaffold/project"
)

func main() {
	var strpath string
	args := os.Args[1:]
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if len(args) > 0 {
		strpath = args[0]
		if err := os.MkdirAll(strpath, 0755); err != nil {
			log.Fatal(err)
		}
		if err := os.Chdir(strpath); err != nil {
			log.Fatal(err)
		}
	} else {
		strpath = pwd
	}

	proj := project.NewConfiguration(strpath)
	p := tea.NewProgram(projectName(proj))

	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func projectName(proj *project.Configuration) tea.Model {
	next := func() (tea.Model, tea.Cmd) {
		next := commonPackages(proj)
		cmd := next.Init()
		return next, cmd
	}
	opts := inputmodels.TextInputOptions{
		Header:      "Enter your project name",
		Placeholder: proj.Name,
		OnEnter: func(input string) error {
			proj.Name = input
			return nil
		},
		OnEnterEmpty: func() error {
			return nil
		},
		Next:      next,
		NextEmpty: next,
	}
	return inputmodels.NewTextInput(opts)

}

func commonPackages(proj *project.Configuration) tea.Model {
	opts := inputmodels.ChoiceModelOptions{
		Choices: []string{"Use Gin", "Use Gorilla/mux", "Use gorm"},
		Values:  []string{"github.com/gin-gonic/gin", "github.com/gorilla/mux", "github.com/jinzhu/gorm"},
		Header:  "Choose among common packages",
		OnEnter: func(selected []inputmodels.Selection) error {
			for _, sel := range selected {
				if sel.IsSelected {
					proj.Dependencies[sel.Value()] = struct{}{}
				}
			}
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			next := otherPackages(proj)
			cmd := next.Init()
			return next, cmd
		},
	}
	return inputmodels.NewchoiceModel(opts)
}

func otherPackages(proj *project.Configuration) tea.Model {
	opts := inputmodels.TextInputOptions{
		Header:      "Enter other packages you want to use",
		Placeholder: "Leave empty if you don't want to use any other packages",
		OnEnter: func(input string) error {
			proj.Dependencies[input] = struct{}{}
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			next := otherPackages(proj)
			cmd := next.Init()
			return next, cmd
		},
		OnEnterEmpty: func() error { return nil },
		NextEmpty: func() (tea.Model, tea.Cmd) {
			next := selectVendoring(proj)
			cmd := next.Init()
			return next, cmd
		},
	}
	return inputmodels.NewTextInput(opts)
}

func selectVendoring(proj *project.Configuration) tea.Model {
	choices := []string{"Yes", "No"}
	values := []string{"yes", "no"}
	opts := inputmodels.RadioSelectOptions{
		Header:  "Do you want to use vendoring?",
		Choices: choices,
		Values:  values,
		OnEnter: func(selection string) error {
			if selection == values[0] {
				proj.DoVendor = true
			} else {
				proj.DoVendor = false
			}

			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			next := buildingBoilerplate(proj)
			cmd := next.Init()
			go func() {
				err := proj.Start()
				if err != nil {
					log.Fatal(err)
				}
			}()
			return next, cmd
		},
	}
	return inputmodels.NewRadioSelect(opts)
}

func buildingBoilerplate(proj *project.Configuration) tea.Model {
	return progressloader.NewLoader(proj)
}
