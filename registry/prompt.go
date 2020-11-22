package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pelletier/go-toml"
)

type templatePrompt struct {
	Key     string      `toml:"-"`
	Value   interface{} `toml:"-"`
	Message string      `toml:"prompt"`
	Help    string      `toml:"help"`
	Default interface{} `toml:"default"`
}

func prompt(dir string, defaults bool) (templatePrompts, error) {
	templatePath := filepath.Join(dir, "template.toml")
	if _, err := os.Lstat(templatePath); err != nil {
		return nil, err
	}

	tree, err := toml.LoadFile(templatePath)
	if err != nil {
		return nil, err
	}

	prompts := make(templatePrompts, len(tree.Keys()))
	for idx, k := range tree.Keys() {
		v := tree.Get(k)

		obj, ok := v.(*toml.Tree)
		if !ok {
			prompts[idx] = templatePrompt{
				Key:     k,
				Message: k,
				Default: v,
			}
			continue
		}

		var p templatePrompt
		if err := obj.Unmarshal(&p); err != nil {
			return nil, err
		}
		p.Key = k
		if p.Message == "" {
			p.Message = p.Key
		}
		if p.Default == nil {
			p.Default = ""
		}
		prompts[idx] = p
	}

	// Return early if we only want defaults
	if defaults {
		return prompts, nil
	}

	// Sort the prompts so they are consistent
	sort.Sort(prompts)

	for idx, prompt := range prompts {
		var p survey.Prompt
		switch t := prompt.Default.(type) {
		case []string:
			p = &survey.Select{
				Message: prompt.Message,
				Options: t,
				Help:    prompt.Help,
			}
		default:
			p = &survey.Input{
				Message: prompt.Message,
				Default: fmt.Sprintf("%v", t),
				Help:    prompt.Help,
			}
		}
		var a string
		if err := survey.AskOne(p, &a); err != nil {
			return nil, err
		}
		prompts[idx].Value = a
	}

	return prompts, nil
}

type templatePrompts []templatePrompt

func (t templatePrompts) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, p := range t {
		if p.Value != nil {
			m[p.Key] = p.Value
			continue
		}
		m[p.Key] = p.Default
	}
	return m
}

func (t templatePrompts) ToFuncMap() template.FuncMap {
	m := make(map[string]interface{})
	for k, v := range t.ToMap() {
		vv := v // Enclosure
		m[k] = func() string {
			return fmt.Sprintf("%v", vv)
		}
	}
	return m
}

func (t templatePrompts) Len() int {
	return len(t)
}

func (t templatePrompts) Less(i, j int) bool {
	return t[i].Key > t[j].Key
}

func (t templatePrompts) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
