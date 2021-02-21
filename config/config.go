package config

import (
	"fmt"

	"golang.org/x/xerrors"
)

const (
	DefaultNamespace = "[default]"
)

func init() {
	console := ConsoleConfig{}
	console.Port = 8081
	console.Host = "localhost"
	console.ProjectID = "default"
	console.Namespace = DefaultNamespace
	console.Limit = 20

	viewer := ViewerConfig{}
	viewer.Port = 8088
	viewer.ConfigFile = "$HOME/.dem.gob"
	viewer.Limit = 20

	gViewer = &viewer
	gConsole = &console
}

type ConsoleConfig struct {
	ProjectID string
	Host      string
	Port      int
	Namespace string
	Limit     int
}

type ViewerConfig struct {
	Port       int
	ConfigFile string
	Limit      int
}

var (
	gViewer  *ViewerConfig
	gConsole *ConsoleConfig
)

func GetViewer() *ViewerConfig {
	return gViewer
}

func GetConsole() *ConsoleConfig {
	return gConsole
}

func SetViewer(opts []ViewerOption) error {
	for _, opt := range opts {
		err := opt(gViewer)
		if err != nil {
			return xerrors.Errorf("ViewerOption set error: %w", err)
		}
	}
	return nil
}

func SetConsole(opts []ConsoleOption) error {

	for _, opt := range opts {
		err := opt(gConsole)
		if err != nil {
			return xerrors.Errorf("ConsoleOption set error: %w", err)
		}
	}

	currentSetting = NewSetting()
	currentProject = NewProject(fmt.Sprintf("%s:%d", gConsole.Host, gConsole.Port), gConsole.ProjectID)
	currentSetting.limit = gConsole.Limit

	return nil
}

func LoadSetting() error {
	currentSetting = NewSetting()
	conf := GetViewer()

	err := currentSetting.read(conf.ConfigFile)
	if err != nil {
		return xerrors.Errorf("Setting Read() error: %w", err)
	}

	currentSetting.limit = conf.Limit

	return nil
}

func AddProject(p *Project) error {
	currentSetting.AddProject(p)
	conf := GetViewer()

	err := currentSetting.write(conf.ConfigFile)
	if err != nil {
		return xerrors.Errorf("Setting write() error: %w", err)
	}
	return nil
}

func GetDarkMode() bool {
	return currentSetting.DarkMode
}

func SetDarkMode(v bool) error {
	currentSetting.DarkMode = v
	conf := GetViewer()
	err := currentSetting.write(conf.ConfigFile)
	if err != nil {
		return xerrors.Errorf("Setting write() error: %w", err)
	}
	return nil
}

func SetLimit(v int) {
	currentSetting.limit = v
}

func GetLimit() int {
	return currentSetting.limit
}

func SetCursor(v string) {
	currentSetting.cursor = v
}

func GetCursor() string {
	return currentSetting.cursor
}

func SetNamespace(v string) {
	currentSetting.namespace = v
}

func GetNamespace() string {
	return currentSetting.namespace
}
