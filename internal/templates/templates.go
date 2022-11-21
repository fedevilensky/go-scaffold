package templates

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/fedevilensky/go-scaffold/internal/project"
)

type templateFunc struct {
	build func(proj *project.Configuration) error
}

func (t *templateFunc) Build(proj *project.Configuration) error {
	return t.build(proj)
}

func LoadFullTemplates() *templateFunc {
	return &templateFunc{
		build: func(proj *project.Configuration) error {
			return buildTemplate(proj, selectEmbs(proj)...)
		},
	}
}

func selectEmbs(proj *project.Configuration) []embed.FS {
	var embs []embed.FS
	if isFullProject(proj) {
		embs = []embed.FS{common}
		switch proj.WebLibrary {
		case project.WebLibraryFiber:
			embs = append(embs, fiber)
		case project.WebLibraryGin:
			embs = append(embs, gin)
		case project.WebLibraryGorillamux:
			embs = append(embs, gorillamux)
		case project.WebLibraryHttp:
			embs = append(embs, http)
		}

		switch proj.DBLibrary {
		case project.DBLibraryGorm:
			embs = append(embs, gorm)
		case project.DBLibrarySql:
			embs = append(embs, sql)
		case project.DBLibrarySqlx:
			embs = append(embs, sqlx)
		}

		switch proj.DBProvider {
		case project.DBProviderPostgres:
			embs = append(embs, postgresql)
		case project.DBProviderMysql:
			embs = append(embs, mysql)
		case project.DBProviderGormMysql:
			embs = append(embs, gormMysql)
		case project.DBProviderGormPostgres:
			embs = append(embs, gormPostgresql)
		default:
			embs = append(embs, mysql)
		}

	} else if proj.WebLibrary != project.WebLibraryNone {
		embs = []embed.FS{}
		switch proj.WebLibrary {
		case project.WebLibraryFiber:
			embs = append(embs, fiberLean)
		case project.WebLibraryGin:
			embs = append(embs, ginLean)
		case project.WebLibraryGorillamux:
			embs = append(embs, gorillamuxLean)
		case project.WebLibraryHttp:
			embs = append(embs, httpLean)
		}
	} else if proj.DBLibrary != project.DBLibraryNone {
		embs = []embed.FS{commonLean}
		switch proj.DBProvider {
		case project.DBProviderPostgres:
			embs = append(embs, postgresql)
		case project.DBProviderMysql:
			embs = append(embs, mysql)
		case project.DBProviderGormMysql:
			embs = append(embs, gormMysql)
		case project.DBProviderGormPostgres:
			embs = append(embs, gormPostgresql)
		default:
			embs = append(embs, mysql)
		}
	}

	return embs
}

func isFullProject(proj *project.Configuration) bool {
	return proj.WebLibrary != project.WebLibraryNone && proj.DBLibrary != project.DBLibraryNone
}

func buildTemplate(proj *project.Configuration, embs ...embed.FS) error {
	tmpl, err := createTemplate(embs...)

	if err != nil {
		return err
	}

	err = createFiles(proj, tmpl, embs...)
	return err
}

func createTemplate(embs ...embed.FS) (tmpl *template.Template, err error) {
	tmpl = template.New("")
	for _, emb := range embs {
		patterns := []string{}
		err := fs.WalkDir(emb, ".", func(pathStr string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !strings.HasSuffix(pathStr, ".tmpl") {
				return nil
			}

			patterns = append(patterns, pathStr)
			return nil
		})
		if err != nil {
			return nil, err
		}
		tmpl, err = tmpl.ParseFS(emb, patterns...)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func createFiles(proj *project.Configuration, tmpl *template.Template, embs ...embed.FS) error {
	for _, emb := range embs {
		err := fs.WalkDir(emb, ".", func(pathStr string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !isOneOf(pathStr, ".go.tmpl", "Dockerfile.tmpl") {
				return nil
			}
			pathParts := strings.Split(pathStr, "/")
			destPath := strings.Join(pathParts[2:], "/")
			destPath = strings.TrimRight(destPath, ".tmpl")
			templateName := path.Base(pathStr)
			err = os.MkdirAll(path.Dir(destPath), 0755)
			if err != nil {
				panic(err)
			}
			f, err := os.Create(destPath)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			err = tmpl.Lookup(templateName).Execute(f, proj)
			if err != nil {
				panic(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func isOneOf(pathStr string, strs ...string) bool {
	for _, str := range strs {
		if strings.HasSuffix(pathStr, str) {
			return true
		}
	}
	return false
}
