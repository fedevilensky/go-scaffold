# go-scaffold
[![Go Reference](https://pkg.go.dev/badge/github.com/fedevilensky/go-scaffold.svg)](https://pkg.go.dev/github.com/fedevilensky/go-scaffold)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/fedevilensky/go-scaffold)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/fedevilensky/go-scaffold)


![alt](go-scaffold.png)

Tool to create project scaffold and provides some boilerplate in case of using gin, fiber or net/http

## Install
To install, use go install
```bash
go install github.com/fedevilensky/go-scaffold@latest
```

## Usage
`go-scaffold h` or `go-scaffold help` to display the help message.

`go-scaffold <folder-name>` to create a new module in a new folder named `<folder-name>`, or just `go-scaffold` to create project in pwd.

Use arrows (or `j` and `k`) to navigate, space to select, `q` to quit and enter to continue. When a text prompt is present you can still quit with `Ctrl+c`

Every file and folder will be created at the end, so feel free to quit and start over

## Known problems
Will panic if go mod already exists


----------
Logo created with [Gopher Konstructor](https://github.com/quasilyte/gopherkon), based on Renee French's design
