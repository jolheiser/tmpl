package env

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Env is tmpl environment variables
type Env map[string]string

// Set sets all environment variables from an Env
func (e Env) Set() error {
	for key, val := range e {
		if err := os.Setenv(key, val); err != nil {
			return err
		}
	}
	return nil
}

// Load loads an env from <path>/env.json
func Load(path string) (Env, error) {
	p := filepath.Join(path, "env.json")
	fi, err := os.Open(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Env{}, nil
		}
		return nil, err
	}
	defer fi.Close()

	var e Env
	if err := json.NewDecoder(fi).Decode(&e); err != nil {
		return nil, err
	}
	return e, nil
}

// Save saves an Env to <path>/env.json
func Save(path string, e Env) error {
	p := filepath.Join(path, "env.json")
	fi, err := os.Create(p)
	if err != nil {
		return err
	}
	return json.NewEncoder(fi).Encode(e)
}
