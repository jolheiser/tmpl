package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
			val := os.ExpandEnv(prompt.Default)
			prompts[idx].Value = val
			os.Setenv(fmt.Sprintf("TMPL_PROMPT_%s", envKey), val)
			continue
		}

		// Otherwise, prompt
		var p survey.Prompt
		switch prompt.Type {
		case config.PromptTypeSelect:
			opts := make([]string, 0, len(prompt.Options))
			for idy, opt := range prompt.Options {
				opts[idy] = os.ExpandEnv(opt)
			}
			p = &survey.Select{
				Message: prompt.Label,
				Options: opts,
				Help:    prompt.Help,
			}
		case config.PromptTypeConfirm:
			def, _ := strconv.ParseBool(os.ExpandEnv(prompt.Default))
			p = &survey.Confirm{
				Message: prompt.Label,
				Help:    prompt.Help,
				Default: def,
			}
		case config.PromptTypeMultiline:
			p = &survey.Multiline{
				Message: prompt.Label,
				Default: os.ExpandEnv(prompt.Default),
				Help:    prompt.Help,
			}
		case config.PromptTypeEditor:
			p = &survey.Editor{
				Message: prompt.Label,
				Default: os.ExpandEnv(prompt.Default),
				Help:    prompt.Help,
			}
		default:
			p = &survey.Input{
				Message: prompt.Label,
				Default: os.ExpandEnv(prompt.Default),
				Help:    prompt.Help,
			}
		}

		m := make(map[string]any)
		if err := survey.AskOne(p, &m); err != nil {
			return nil, fmt.Errorf("could not complete prompt: %w", err)
		}
		a := m[""]
		prompts[idx].Value = a
		os.Setenv(fmt.Sprintf("TMPL_PROMPT_%s", envKey), fmt.Sprint(a))
	}

	return prompts, nil
}

type templatePrompt struct {
	config.Prompt
	Value any
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
