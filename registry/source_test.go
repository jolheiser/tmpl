package registry

import (
	"testing"

	"github.com/matryer/is"
)

func TestSource(t *testing.T) {
	assert := is.New(t)

	tt := []struct {
		Name     string
		Source   *Source
		CloneURL string
	}{
		{
			Name: "Gitea",
			Source: &Source{
				URL: "https://gitea.com/",
			},
			CloneURL: "https://gitea.com/user/repo.git",
		},
		{
			Name: "GitHub",
			Source: &Source{
				URL: "https://github.com/",
			},
			CloneURL: "https://github.com/user/repo.git",
		},
		{
			Name: "GitLab",
			Source: &Source{
				URL: "https://gitlab.com/",
			},
			CloneURL: "https://gitlab.com/user/repo.git",
		},
	}

	namespace := "user/repo"
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			cloneURL := tc.Source.CloneURL(namespace)
			assert.Equal(tc.CloneURL, cloneURL) // Clone URLs should match
		})
	}
}
