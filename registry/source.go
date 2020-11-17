package registry

import "fmt"

// Source is a quick way to specify a git source
// e.g. Gitea, GitHub, etc.
type Source struct {
	Name string `toml:"name"`
	URL  string `toml:"url"`
}

// CloneURL constructs a URL suitable for cloning a repository
func (s *Source) CloneURL(namespace string) string {
	return fmt.Sprintf("%s%s", s.URL, namespace)
}
