package project

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/muesli/termenv"
)

var (
	term = termenv.ColorProfile()
)

const (
	blueFg = "36"
	redFg  = "31"
)

type Configuration struct {
	Name           string
	DoVendor       bool
	Dependencies   map[string]struct{}
	processedDeps  int
	vendorFinished bool
	currentCmd     string
	Template       func(*Configuration) error
}

func NewConfiguration(strpath string) *Configuration {
	name := path.Base(strpath)
	if name == "." || name == "/" {
		name = "default"
	}
	return &Configuration{
		Name:         name,
		DoVendor:     false,
		Dependencies: map[string]struct{}{},
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
	c.currentCmd = "Creating folders...\n\n"
	err = c.createFolders()
	if err != nil {
		return
	}
	c.currentCmd = "Initializing mod...\n\n"
	err = c.modInit()
	if err != nil {
		return
	}
	c.currentCmd = c.currentCmd + "Installing dependencies...\n\n"
	err = c.installDependencies()
	if err != nil {
		return
	}
	if c.Template != nil {
		c.currentCmd = c.currentCmd + "Creating helper methods and middlewares in pkg/...\n\n"
		err = c.Template(c)
		if err != nil {
			return
		}
	}
	if c.DoVendor {
		c.currentCmd = c.currentCmd + "Vendoring...\n\n"
		err = c.vendor()
		if err != nil {
			return
		}
	}
	c.currentCmd = c.currentCmd + "Finished!"
	return nil
}

func (c *Configuration) createFolders() (err error) {
	dirs := []string{"./internal/models",
		"./internal/services",
		"./internal/controllers",
		"./internal/repositories",
		"./internal/routers",
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
	err = exec.Command("go", "mod", "init", c.Name).Run()
	if err != nil {
		c.currentCmd = fmt.Sprintf(
			"Failed to run command: %s\nError: %s",
			colorFg(fmt.Sprintf("go mod init %s", c.Name), blueFg),
			colorFg(err.Error(), redFg),
		)
	}
	return
}

func (c *Configuration) installDependencies() (err error) {
	for dep := range c.Dependencies {
		startingCmd := c.currentCmd
		c.currentCmd = startingCmd + fmt.Sprintf("Getting dependency: %s\n", colorFg(dep, blueFg))
		err = exec.Command("go", "get", dep).Run()
		if err != nil {
			//log and continue
			c.currentCmd = startingCmd + fmt.Sprintf(
				"Failed to get dependency:%s => %s\n",
				colorFg(dep, blueFg),
				colorFg("Check if it exists", redFg),
			)
		}
		c.processedDeps++
	}
	if len(c.Dependencies) > 0 {
		c.currentCmd = c.currentCmd + "Finished installing dependencies!\n\n"
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

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}
