package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/fedevilensky/go-scaffold/internal/project"
)

//go:embed "embedded/gin"
var gin embed.FS

//go:embed "embedded/fiber"
var fiber embed.FS

//go:embed "embedded/http"
var http embed.FS

func LoadGinTemplates(proj *project.Configuration) error {
	return loadTemplates(gin, proj)
}

func LoadFiberTemplates(proj *project.Configuration) error {
	return loadTemplates(fiber, proj)
}

func LoadHttpTemplates(proj *project.Configuration) error {
	return loadTemplates(http, proj)
}

func loadTemplates(embFS embed.FS, proj *project.Configuration) error {
	fs.WalkDir(embFS, ".", func(pathStr string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		fullDir, templateFileName := path.Split(pathStr)
		dir := path.Base(fullDir)
		fileName := strings.TrimSuffix(templateFileName, ".tmpl")

		newPath := fmt.Sprintf("pkg/%s/%s", dir, fileName)
		err = os.MkdirAll(path.Dir(newPath), 0755)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(newPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		tmpl, err := template.ParseFS(embFS, pathStr)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(f, proj)
		if err != nil {
			panic(err)
		}

		fmt.Println(pathStr)
		return nil
	})
	return nil
}
