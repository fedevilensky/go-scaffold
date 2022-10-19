package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fedevilensky/go-scaffold/internal/project"
	"github.com/fedevilensky/go-scaffold/internal/templates"
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
