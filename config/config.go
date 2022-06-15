package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is a tmpl config
type Config struct {
	Prompts []Prompt `yaml:"prompts"`
}

// Prompt is a tmpl prompt
type Prompt struct {
	ID      string   `yaml:"id"`
	Label   string   `yaml:"label"`
	Help    string   `yaml:"help"`
	Default string   `yaml:"default"`
	Options []string `yaml:"options"`
}

// Load loads a tmpl config
func Load(r io.Reader) (*Config, error) {
	configBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	configContents := os.Expand(string(configBytes), func(s string) string {
		if strings.HasPrefix(s, "TMPL_PROMPT") {
			return fmt.Sprintf("${%s}", s)
		}
		return os.Getenv(s)
	})

	var c Config
	return &c, yaml.Unmarshal([]byte(configContents), &c)
}
