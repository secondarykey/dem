package config

import (
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

type Project struct {
	ID        string
	ProjectID string
	Endpoint  string
}

func NewProject(endpoint string, projectid string) *Project {
	p := Project{}
	p.Endpoint = endpoint
	p.ProjectID = projectid
	p.ID = uuid.New().String()
	return &p
}

func GetProjects() ([]*Project, error) {
	return getCurrentSetting().Projects, nil
}

func GetCurrentProject() *Project {
	return getProject(currentEmbed.ID)
}

func SwitchProject(id string) *Project {
	p := getProject(id)
	GetCurrentEmbed().ID = id
	return p
}

func AddProject(p *Project) error {

	s := getCurrentSetting()
	s.addProject(p)

	conf := GetViewer()
	err := s.write(conf.ConfigFile)
	if err != nil {
		return xerrors.Errorf("Setting write() error: %w", err)
	}
	return nil
}

func DeleteProject(id string) error {
	s := getCurrentSetting()
	err := s.deleteProject(id)
	if err != nil {
		return xerrors.Errorf("deleteProject() error: %w", err)
	}
	conf := GetViewer()
	err = s.write(conf.ConfigFile)
	if err != nil {
		return xerrors.Errorf("Setting write() error: %w", err)
	}
	return nil
}

func getProject(id string) *Project {
	for _, p := range getCurrentSetting().Projects {
		if p.ID == id {
			return p
		}
	}
	return nil
}
