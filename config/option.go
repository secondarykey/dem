package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ViewerOption func(*ViewerConfig) error

func SetMyPort(port int) ViewerOption {
	return func(c *ViewerConfig) error {
		c.Port = port
		return nil
	}
}

func SetConfigFile(path string) ViewerOption {
	return func(c *ViewerConfig) error {

		p := filepath.Clean(path)
		slice := strings.Split(p, string(filepath.Separator))

		fp := ""

		for _, v := range slice {
			d := v
			if strings.Index(v, "$") == 0 {
				if len(v) <= 1 {
					return fmt.Errorf("environment variable charactor error[%s]", v)
				}
				env := v[1:]
				if runtime.GOOS == "windows" && env == "HOME" {
					env = "USERPROFILE"
				}
				d = os.Getenv(env)
				if d == "" {
					return fmt.Errorf("environment variable variable error[%s]", env)
				}
			}

			if fp != "" {
				fp += string(filepath.Separator)
			}
			fp += d
		}

		c.ConfigFile = fp
		return nil
	}
}

type ConsoleOption func(*ConsoleConfig) error

func SetEmulatorPort(p int) ConsoleOption {
	return func(c *ConsoleConfig) error {
		c.Port = p
		return nil
	}
}

func SetHost(h string) ConsoleOption {
	return func(c *ConsoleConfig) error {
		c.Host = h
		return nil
	}
}

func SetProjectID(p string) ConsoleOption {
	return func(c *ConsoleConfig) error {
		c.ProjectID = p
		return nil
	}
}

func SetNamespace(n string) ConsoleOption {
	return func(c *ConsoleConfig) error {
		c.Namespace = n
		return nil
	}
}
