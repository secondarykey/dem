package config_test

import (
	"os"
	"testing"

	"github.com/secondarykey/dem/config"
)

func TestSettingFile(t *testing.T) {

	name := "test.gob"
	if _, err := os.Stat(name); err == nil {
		t.Errorf("%s file exist.", name)
	}

	defer os.Remove(name)
	s := config.NewSetting()

	p := config.Project{}
	p.Endpoint = "endpoint"
	p.ID = "111"
	p.ProjectID = "projectid"

	s.AddProject(&p)

	err := s.Write(name)
	if err != nil {
		t.Errorf("write error: %+v", err)
	}

	if _, err := os.Stat(name); err != nil {
		t.Errorf("%s file not exist.", name)
	}

	l := config.NewSetting()
	err = l.Read(name)
	if err != nil {
		t.Errorf("read error")
	}

	if len(l.Projects) != 1 {
		t.Errorf("length not 1")
	}

	if s.Projects[0].ProjectID != l.Projects[0].ProjectID ||
		s.Projects[0].Endpoint != l.Projects[0].Endpoint ||
		s.Projects[0].ID != l.Projects[0].ID {
		t.Errorf("read write not equals")
	}

}
