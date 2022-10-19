package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fedevilensky/go-scaffold/internal/inputmodels"
	"github.com/fedevilensky/go-scaffold/internal/progressloader"
	"github.com/fedevilensky/go-scaffold/internal/project"
	"github.com/fedevilensky/go-scaffold/internal/templates"
)

const (
	projName  = "projName"
	web       = "web"
	db        = "dbLib"
	removeDep = "removeDep"
	addDep    = "addDep"
	vendor    = "vendor"
	build     = "build"
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
		if strpath == "help" || strpath == "h" {
			fmt.Println("go-scaffold:                will create a new project in $PWD")
			fmt.Println("go-scaffold <project-name>: will create a new project in $PWD/<project-name>")
			return
		}
		if err := os.MkdirAll(strpath, 0755); err != nil {
			log.Fatal(err)
		}
		if err := os.Chdir(strpath); err != nil {
			log.Fatal(err)
		}
	} else {
		strpath = filepath.Base(pwd)
	}

	if _, err := os.Stat(strpath); err == nil {
		log.Fatal("go.mod already exists")
	}
	if _, err := os.Stat("go.sum"); err == nil {
		log.Fatal("go.sum already exists")
	}

	proj := project.NewConfiguration(strpath)
	proj.Template = templates.LoadFullTemplates()
	p := tea.NewProgram(projectName(proj))

	if err := p.Start(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func projectName(proj *project.Configuration) tea.Model {
	next := func() tea.Model { return selectWebLibrary(proj) }
	return projectNameWithNext(proj, next)
}

// next func() tea.Model is done in orden to defer build until the last posible moment
func projectNameWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	nextModel := func() (tea.Model, tea.Cmd) {
		nextModel := next()
		cmd := nextModel.Init()
		return nextModel, cmd
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
		Next:      nextModel,
		NextEmpty: nextModel,
	}

	return inputmodels.NewTextInput(opts)
}

func selectWebLibrary(proj *project.Configuration) tea.Model {
	next := func() tea.Model { return selectDBLibrary(proj) }
	return selectWebLibraryWithNext(proj, next)
}
func selectWebLibraryWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	opts := inputmodels.RadioSelectOptions{
		Choices: []string{"Gin", "Fiber",
			"Gorilla/mux",
			"net/http (and other compatible libraries)", "other"},
		Values: []string{project.WebLibraryGin, project.WebLibraryFiber,
			project.WebLibraryGorillamux,
			project.WebLibraryHttp, project.WebLibraryNone},
		Header: "Choose your base library",
		OnEnter: func(selected string, _ int) error {
			switch selected {
			case project.WebLibraryGin:
				proj.WebLibrary = selected
			case project.WebLibraryFiber:
				proj.WebLibrary = selected
			case project.WebLibraryGorillamux:
				proj.WebLibrary = selected
			case project.WebLibraryHttp:
				proj.WebLibrary = selected
			default:
				proj.WebLibrary = project.WebLibraryNone
			}
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			nextModel := next()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}
	return inputmodels.NewRadioSelect(opts)
}

func selectDBLibrary(proj *project.Configuration) tea.Model {
	next := func() tea.Model { return commonPackages(proj) }
	return selectDBLibraryWithNext(proj, next)
}

func selectDBLibraryWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	selectDBProvider := func() tea.Model {
		switch proj.DBLibrary {
		case project.DBLibraryGorm:
			return selectGormDBProviderWithNext(proj, next)
		case project.DBLibrarySql, project.DBLibrarySqlx:
			return selectDBProviderWithNext(proj, next)
		default:
			return next()
		}
	}

	opts := inputmodels.RadioSelectOptions{
		Header:  "Select a db library",
		Choices: []string{"sql", "sqlx", "gorm", "None"},
		Values:  []string{project.DBLibrarySql, project.DBLibrarySqlx, project.DBLibraryGorm, project.DBLibraryNone},
		OnEnter: func(selection string, _ int) error {
			proj.DBLibrary = selection
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			nextModel := selectDBProvider()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}

	return inputmodels.NewRadioSelect(opts)
}

func selectDBProviderWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	opts := inputmodels.RadioSelectOptions{
		Header:  "Select a DBMS",
		Choices: []string{"PostgreSQL", "MySQL/MariaDB"},
		Values:  []string{project.DBProviderPostgres, project.DBProviderMysql},
		OnEnter: func(selection string, _ int) error {
			proj.DBProvider = selection
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			nextModel := next()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}
	return inputmodels.NewRadioSelect(opts)
}

func selectGormDBProviderWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	opts := inputmodels.RadioSelectOptions{
		Header:  "Select a DBMS",
		Choices: []string{"PostgreSQL", "MySQL/MariaDB"},
		Values:  []string{project.DBProviderGormPostgres, project.DBProviderGormMysql},
		OnEnter: func(selection string, _ int) error {
			proj.DBProvider = selection
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			nextModel := next()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}
	return inputmodels.NewRadioSelect(opts)
}

func commonPackages(proj *project.Configuration) tea.Model {
	next := func() tea.Model { return otherPackages(proj) }
	return commonPackagesWithNext(proj, next)
}

func commonPackagesWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	opts := inputmodels.ChoiceModelOptions{
		Choices: []string{"Use spf13/cobra", "Use redis", "Use cleanenv",
			"Use MongoDB", "Use AMQP"},
		Values: []string{"github.com/spf13/cobra", "github.com/go-redis/redis", "github.com/ilyakaznacheev/cleanenv",
			"github.com/mongodb/mongo-go-driver/mongo", "github.com/rabbitmq/amqp091-go"},
		Header: "Choose between common packages",
		OnEnter: func(selected []inputmodels.Selection) error {
			for _, sel := range selected {
				if sel.IsSelected {
					proj.Dependencies[sel.Value()] = struct{}{}
				}
			}
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			nextModel := next()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}
	return inputmodels.NewchoiceModel(opts)
}

func otherPackages(proj *project.Configuration) tea.Model {
	next := func() tea.Model { return selectVendoring(proj) }
	return otherPackagesWithNext(proj, next)
}

func otherPackagesWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	opts := inputmodels.TextInputOptions{
		Header:      "Enter other packages you want to use",
		Placeholder: "Leave empty if you don't want to use any other packages",
		OnEnter: func(input string) error {
			proj.Dependencies[input] = struct{}{}
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			next := otherPackagesWithNext(proj, next)
			cmd := next.Init()
			return next, cmd
		},
		OnEnterEmpty: func() error { return nil },
		NextEmpty: func() (tea.Model, tea.Cmd) {
			nextModel := next()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}
	return inputmodels.NewTextInput(opts)
}

func selectVendoring(proj *project.Configuration) tea.Model {
	next := func() tea.Model { return showSummary(proj, 0) }
	return selectVendorigWithNext(proj, next)
}

func selectVendorigWithNext(proj *project.Configuration, next func() tea.Model) tea.Model {
	choices := []string{"Yes", "No"}
	values := []string{"yes", "no"}
	opts := inputmodels.RadioSelectOptions{
		Header:  "Do you want to use vendoring?",
		Choices: choices,
		Values:  values,
		OnEnter: func(selection string, _ int) error {
			if selection == values[0] {
				proj.DoVendor = true
			} else {
				proj.DoVendor = false
			}

			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			nextModel := next()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}
	return inputmodels.NewRadioSelect(opts)
}

func showSummary(proj *project.Configuration, cursorPosition int) tea.Model {
	var next func() tea.Model
	choices, values := buildChoicesAndValues(proj)
	opts := inputmodels.RadioSelectOptions{
		Header:                 "Summary",
		Choices:                choices,
		CursorStartingPosition: cursorPosition,
		Values:                 values,
		OnEnter: func(input string, cursorPosition int) error {
			switch {
			case input == projName:
				next = func() tea.Model {
					return projectNameWithNext(proj,
						func() tea.Model { return showSummary(proj, cursorPosition) })
				}
			case input == web:
				next = func() tea.Model {
					return selectWebLibraryWithNext(proj,
						func() tea.Model { return showSummary(proj, cursorPosition) })
				}
			case input == db:
				next = func() tea.Model {
					return selectDBLibraryWithNext(proj,
						func() tea.Model { return showSummary(proj, cursorPosition) })
				}
			case strings.HasPrefix(input, removeDep):
				dep := strings.TrimPrefix(input, removeDep)
				delete(proj.Dependencies, dep)
				next = func() tea.Model {
					return showSummary(proj, cursorPosition)
				}
			case input == addDep:
				next = func() tea.Model {
					return otherPackagesWithNext(proj, func() tea.Model { return showSummary(proj, cursorPosition) })
				}
			case input == vendor:
				proj.DoVendor = !proj.DoVendor
				next = func() tea.Model {
					return showSummary(proj, cursorPosition)
				}
			case input == build:
				go func() {
					err := proj.Start()
					if err != nil {
						log.Fatal(err)
					}
				}()
				next = func() tea.Model { return buildingBoilerplate(proj) }
			default:
				return errors.New("unexpected input")
			}
			return nil
		},
		Next: func() (tea.Model, tea.Cmd) {
			nextModel := next()
			cmd := nextModel.Init()
			return nextModel, cmd
		},
	}
	return inputmodels.NewRadioSelect(opts)
}

func buildChoicesAndValues(proj *project.Configuration) ([]string, []string) {
	choices := []string{}
	values := []string{}

	choices = append(choices, "Project name: "+proj.Name)
	values = append(values, projName)

	if proj.WebLibrary != project.WebLibraryNone {
		choices = append(choices, "Web library: "+proj.WebLibrary)
	} else {
		choices = append(choices, "No web library (no helper/util functions will be written)")
	}
	values = append(values, web)

	if proj.DBLibrary != project.DBLibraryNone {
		choices = append(choices, "DB Library: "+proj.DBLibrary+", with driver: "+proj.DBProvider)
	} else {
		choices = append(choices, "No DB library selected (no code example will be generated)")
	}
	values = append(values, db)

	if len(proj.Dependencies) > 0 {
		for dep := range proj.Dependencies {
			choices = append(choices, "Remove dependency: "+dep)
			values = append(values, removeDep+dep)
		}
	}
	choices = append(choices, "Add dependency")
	values = append(values, addDep)

	choices = append(choices, "Vendoring: "+strconv.FormatBool(proj.DoVendor))
	values = append(values, vendor)

	choices = append(choices, "Build")
	values = append(values, build)

	return choices, values
}

func buildingBoilerplate(proj *project.Configuration) tea.Model {
	return progressloader.NewLoader(proj)
}
