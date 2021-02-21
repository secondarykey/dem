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
	return getCurrentSetting().Projects, nil
}

func GetCurrentProject() *Project {
	return getProject(current.ID)
}

func SwitchProject(id string) *Project {
	p := getProject(id)
	GetCurrent().ID = id
	return p
}

func getProject(id string) *Project {
	for _, p := range getCurrentSetting().Projects {
		if p.ID == id {
			return p
		}
	}
	return nil
}
