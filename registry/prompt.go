package registry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	templateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	// Expand the template with environment variables
	templateContents := os.ExpandEnv(string(templateBytes))

	tree, err := toml.Load(templateContents)
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

	// Sort the prompts so they are consistent
	sort.Sort(prompts)

	for idx, prompt := range prompts {
		// Check for env variable
		if e, ok := os.LookupEnv(fmt.Sprintf("TMPL_VAR_%s", strings.ToUpper(prompt.Key))); ok {
			prompts[idx].Value = e
			continue
		}

		// Check if we are using defaults
		if defaults {
			prompts[idx].Value = prompt.Default
			continue
		}

		var p survey.Prompt
		switch t := prompt.Default.(type) {
		case []string:
			p = &survey.Select{
				Message: prompt.Message,
				Options: t,
				Help:    prompt.Help,
			}
		case bool:
			p = &survey.Confirm{
				Message: prompt.Message,
				Default: t,
				Help:    prompt.Help,
			}
		case string:
			p = &survey.Input{
				Message: prompt.Message,
				Default: t,
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

// ToMap converts a slice to templatePrompt into a suitable template context
func (t templatePrompts) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, p := range t {
		m[p.Key] = p.Value
	}
	return m
}

// ToFuncMap converts a slice of templatePrompt into a suitable template.FuncMap
func (t templatePrompts) ToFuncMap() template.FuncMap {
	m := make(map[string]interface{})
	for k, v := range t.ToMap() {
		vv := v // Enclosure
		m[k] = func() interface{} {
			return vv
		}
	}
	return m
}

// Len is for sort.Sort
func (t templatePrompts) Len() int {
	return len(t)
}

// Less is for sort.Sort
func (t templatePrompts) Less(i, j int) bool {
	return t[i].Key > t[j].Key
}

// Swap is for sort.Sort
func (t templatePrompts) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
