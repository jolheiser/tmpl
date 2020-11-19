package registry

import (
	"strings"
	"testing"
)

func TestSource(t *testing.T) {
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
			if !strings.EqualFold(tc.CloneURL, cloneURL) {
				t.Logf("incorrect clone URL:\n\tExpected: %s\n\tGot: %s\n", tc.CloneURL, cloneURL)
				t.Fail()
			}
		})
	}
}
