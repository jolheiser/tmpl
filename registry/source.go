package registry

import "fmt"

// Source is a quick way to specify a git source
// e.g. Gitea, GitHub, etc.
type Source struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// CloneURL constructs a URL suitable for cloning a repository
func (s *Source) CloneURL(namespace string) string {
	return fmt.Sprintf("%s%s.git", s.URL, namespace)
}
