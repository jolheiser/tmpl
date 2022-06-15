package registry

import "fmt"

type ErrTemplateExists struct {
	Name string
}

func (e ErrTemplateExists) Error() string {
	return fmt.Sprintf("template %s already exists", e.Name)
}

type ErrTemplateNotFound struct {
	Name string
}

func (e ErrTemplateNotFound) Error() string {
	return fmt.Sprintf("template not found for %s", e.Name)
}

type ErrSourceNotFound struct {
	Name string
}

func (e ErrSourceNotFound) Error() string {
	return fmt.Sprintf("Source not found for %s", e.Name)
}
