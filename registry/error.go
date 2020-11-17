package registry

import "fmt"

type ErrTemplateNotFound struct {
	Name string
}

func (e ErrTemplateNotFound) Error() string {
	return fmt.Sprintf("template not found for %s", e.Name)
}

func IsErrTemplateNotFound(err error) bool {
	_, ok := err.(ErrTemplateNotFound)
	return ok
}

type ErrSourceNotFound struct {
	Name string
}

func (e ErrSourceNotFound) Error() string {
	return fmt.Sprintf("Source not found for %s", e.Name)
}

func IsErrSourceNotFound(err error) bool {
	_, ok := err.(ErrSourceNotFound)
	return ok
}
