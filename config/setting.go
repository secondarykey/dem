package config

import (
	"encoding/gob"
	"os"

	"golang.org/x/xerrors"
)

var currentSetting *Setting

type Setting struct {
	DarkMode bool
	Projects []*Project
}

func getCurrentSetting() *Setting {
	if currentSetting == nil {
		currentSetting = newSetting()
		currentSetting.read(gViewer.ConfigFile)
	}
	return currentSetting
}

func newSetting() *Setting {
	s := Setting{}
	s.Projects = make([]*Project, 0)
	return &s
}

func (s *Setting) addProject(p *Project) {
	s.Projects = append(s.Projects, p)
}

func (s *Setting) deleteProject(id string) error {
	org := s.Projects
	s.Projects = make([]*Project, 0)

	for _, p := range org {
		if p.ID != id {
			s.Projects = append(s.Projects, p)
		}
	}

	if len(org) == len(s.Projects) {
		return xerrors.Errorf("Not Exists ProjectID[%s]", id)
	}
	return nil
}

func (s *Setting) SetDarkMode(f bool) {
	s.DarkMode = f
}

func (s *Setting) GetDarkMode() bool {
	return s.DarkMode
}

func (s *Setting) read(name string) error {

	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return xerrors.Errorf("not exists: %w", err)
	}

	f, err := os.Open(name)
	if err != nil {
		return xerrors.Errorf("Setting file open error: %w", err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(s)
	if err != nil {
		return xerrors.Errorf("Setting decode error: %w", err)
	}

	return nil
}

func (s *Setting) write(name string) error {

	f, err := os.Create(name)
	if err != nil {
		return xerrors.Errorf("Setting file open error: %w", err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(s)
	if err != nil {
		return xerrors.Errorf("Setting encode error: %w", err)
	}

	return nil
}
