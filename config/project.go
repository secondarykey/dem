package config

import "github.com/google/uuid"

var currentProject *Project

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
	return currentSetting.Projects, nil
}

func GetCurrentProject() *Project {
	return currentProject
}

func SwitchProject(id string) *Project {
	for _, p := range currentSetting.Projects {
		if p.ID == id {
			currentProject = p
			return p
		}
	}
	return nil
}
