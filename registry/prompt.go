package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"go.jolheiser.com/tmpl/config"

	"github.com/charmbracelet/huh"
)

func prompt(dir string, defaults, accessible bool) (templatePrompts, error) {
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
		var f huh.Field
		switch prompt.Type {
		case config.PromptTypeSelect:
			opts := make([]huh.Option[string], 0, len(prompt.Options))
			for idy, opt := range prompt.Options {
				o := os.ExpandEnv(opt)
				opts[idy] = huh.NewOption(o, o)
			}
			def := prompt.Default
			f = huh.NewSelect[string]().
				Title(prompt.Label).
				Description(prompt.Help).
				Options(opts...).
				Key(prompt.ID).
				Value(&def)
		case config.PromptTypeConfirm:
			def, _ := strconv.ParseBool(os.ExpandEnv(prompt.Default))
			f = huh.NewConfirm().
				Title(prompt.Label).
				Description(prompt.Help).
				Key(prompt.ID).
				Value(&def)
		case config.PromptTypeMultiline, config.PromptTypeEditor:
			def := os.ExpandEnv(prompt.Default)
			f = huh.NewText().
				Title(prompt.Label).
				Description(prompt.Help).
				Key(prompt.ID).
				Value(&def)
		default:
			def := os.ExpandEnv(prompt.Default)
			f = huh.NewInput().
				Title(prompt.Label).
				Description(prompt.Help).
				Key(prompt.ID).
				Value(&def)
		}
		if err := huh.NewForm(huh.NewGroup(f)).WithAccessible(accessible).WithTheme(huh.ThemeCatppuccin()).Run(); err != nil {
			return nil, fmt.Errorf("could not run field: %w", err)
		}
		prompts[idx].Value = f.GetValue()
		os.Setenv(fmt.Sprintf("TMPL_PROMPT_%s", envKey), fmt.Sprint(f.GetValue()))
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
