package project

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

type Configuration struct {
	Name           string
	DoVendor       bool
	Dependencies   map[string]struct{}
	processedDeps  int
	vendorFinished bool
	log            *log.Logger
	currentCmd     string
}

func NewConfiguration() *Configuration {
	strpath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	name := path.Base(strpath)
	if name == "." || name == "/" {
		name = "default"
	}
	l := log.New(os.Stderr, "", 0)
	return &Configuration{
		Name:         name,
		DoVendor:     false,
		Dependencies: map[string]struct{}{},
		log:          l,
	}
}

func (c *Configuration) GetCurrentCmd() string {
	return c.currentCmd
}

func (c *Configuration) CalculateProgress() float64 {
	// i dont care about data races
	if !c.DoVendor {
		return float64(c.processedDeps) / float64(len(c.Dependencies))
	} else {
		if !c.vendorFinished {
			return float64(c.processedDeps) / (float64(len(c.Dependencies)) + 1)
		}
		return 1
	}
}

func (c *Configuration) Start() (err error) {
	c.currentCmd = "Creating folders..."
	err = c.createFolders()
	if err != nil {
		return
	}
	c.currentCmd = "Initializing mod..."
	err = c.modInit()
	if err != nil {
		return
	}
	err = c.installDependencies()
	if err != nil {
		return
	}
	if c.DoVendor {
		c.currentCmd = "Vendoring..."
		err = c.vendor()
		if err != nil {
			return
		}
	}
	c.currentCmd = "Finished!"
	return nil
}

func (c *Configuration) createFolders() (err error) {
	dirs := []string{"./src/models",
		"./src/services",
		"./src/controllers",
		"./src/repositories",
		"./src/routers",
		"./cmd",
		"./scripts",
		"./docs",
	}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return
		}
	}
	return nil
}

func (c *Configuration) modInit() (err error) {
	return exec.Command("go", "mod", "init", c.Name).Run()
}

func (c *Configuration) installDependencies() (err error) {
	for dep := range c.Dependencies {
		c.currentCmd = fmt.Sprint("Getting dependency: ", dep)
		err = exec.Command("go", "get", dep).Run()
		if err != nil {
			c.log.Printf("Failed to get dependency: %s\nError: %s", dep, err.Error())
		}
		c.processedDeps++
	}
	return nil
}

func (c *Configuration) vendor() (err error) {
	err = exec.Command("go", "mod", "vendor").Run()
	if err != nil {
		return
	}
	c.vendorFinished = true
	return nil
}
