package config

import "github.com/google/uuid"

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

func GetProject(id string) *Project {
	for _, p := range currentSetting.Projects {
		if p.ID == id {
			return p
		}
	}
	return nil
}
