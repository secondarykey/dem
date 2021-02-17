package config

import "golang.org/x/xerrors"

func init() {
	console := ConsoleConfig{}
	console.Port = 8081
	console.Host = "localhost"
	console.ProjectID = "default"
	console.Namespace = "[default]"

	viewer := ViewerConfig{}
	viewer.Port = 8088
	viewer.ConfigFile = "$HOME/.dem.gob"

	gViewer = &viewer
	gConsole = &console
}

type ConsoleConfig struct {
	ProjectID string
	Host      string
	Port      int
	Namespace string
}

type ViewerConfig struct {
	Port       int
	ConfigFile string
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
	return nil
}

func LoadSetting() error {
	currentSetting = NewSetting()
	conf := GetViewer()

	err := currentSetting.read(conf.ConfigFile)
	if err != nil {
		return xerrors.Errorf("Setting Read() error: %w", err)
	}

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
