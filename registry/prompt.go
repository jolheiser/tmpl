package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go.jolheiser.com/tmpl/config"

	"github.com/AlecAivazis/survey/v2"
)

func prompt(dir string, defaults bool) (templatePrompts, error) {
	templatePath := filepath.Join(dir, "tmpl.yaml")
	fi, err := os.Open(templatePath)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	cfg, err := config.Load(fi)
	if err != nil {
		return nil, err
	}

	prompts := make(templatePrompts, 0, len(cfg.Prompts))
	for _, prompt := range cfg.Prompts {
		tp := templatePrompt{
			Prompt: prompt,
		}
		if tp.Label == "" {
			tp.Label = tp.ID
		}
		if tp.Default == nil {
			tp.Default = ""
		}
		prompts = append(prompts, tp)
	}

	for idx, prompt := range prompts {
		// Check for env variable
		envKey := strings.ToUpper(prompt.ID)
		if e, ok := os.LookupEnv(fmt.Sprintf("TMPL_VAR_%s", envKey)); ok {
			prompts[idx].Value = e
			os.Setenv(fmt.Sprintf("TMPL_PROMPT_%s", envKey), e)
			continue
		}

		// Check if we are using defaults
		if defaults {
			val := prompt.Default
			switch t := prompt.Default.(type) {
			case []string:
				for idy, s := range t {
					t[idy] = os.ExpandEnv(s)
				}
				val = t
			case string:
				val = os.ExpandEnv(t)
			}
			s := fmt.Sprint(val)
			prompts[idx].Value = s
			os.Setenv(fmt.Sprintf("TMPL_PROMPT_%s", envKey), s)
			continue
		}

		var p survey.Prompt
		switch t := prompt.Default.(type) {
		case []string:
			for idy, s := range t {
				t[idy] = os.ExpandEnv(s)
			}
			p = &survey.Select{
				Message: prompt.Label,
				Options: t,
				Help:    prompt.Help,
			}
		case bool:
			p = &survey.Confirm{
				Message: prompt.Label,
				Default: t,
				Help:    prompt.Help,
			}
		case string:
			p = &survey.Input{
				Message: prompt.Label,
				Default: os.ExpandEnv(t),
				Help:    prompt.Help,
			}
		default:
			p = &survey.Input{
				Message: prompt.Label,
				Default: fmt.Sprint(t),
				Help:    prompt.Help,
			}
		}
		var a string
		if err := survey.AskOne(p, &a); err != nil {
			return nil, err
		}
		prompts[idx].Value = a
		os.Setenv(fmt.Sprintf("TMPL_PROMPT_%s", envKey), a)
	}

	return prompts, nil
}

type templatePrompt struct {
	config.Prompt
	Value string
}

type templatePrompts []templatePrompt

// ToMap converts a slice to templatePrompt into a suitable template context
func (t templatePrompts) ToMap() map[string]any {
	m := make(map[string]any)
	for _, p := range t {
		m[p.ID] = p.Value
	}
	return m
}

// ToFuncMap converts a slice of templatePrompt into a suitable template.FuncMap
func (t templatePrompts) ToFuncMap() template.FuncMap {
	m := make(map[string]any)
	for k, v := range t.ToMap() {
		vv := v // Enclosure
		m[k] = func() any {
			return vv
		}
	}
	return m
}
