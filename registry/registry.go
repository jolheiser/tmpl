package registry

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mholt/archiver/v3"
	"github.com/pelletier/go-toml"
)

// Registry is a collection of Template
type Registry struct {
	dir       string
	Sources   []*Source   `toml:"sources"`
	Templates []*Template `toml:"templates"`
}

func (r *Registry) save() error {
	fi, err := os.Create(r.MetaFilePath())
	if err != nil {
		return err
	}
	if err := toml.NewEncoder(fi).Encode(r); err != nil {
		return err
	}
	return fi.Close()
}

// MetaFilePath is the path to the Registry meta-file
func (r *Registry) MetaFilePath() string {
	return filepath.Join(r.dir, "registry.toml")
}

// GetTemplate retrieves a Template from the Registry
func (r *Registry) GetTemplate(name string) (*Template, error) {
	for _, t := range r.Templates {
		if strings.EqualFold(name, t.Name) {
			t.reg = r
			return t, nil
		}
	}
	return nil, ErrTemplateNotFound{Name: name}
}

// AddTemplate downloads and adds a new Template to the Registry
func (r *Registry) AddTemplate(name, repo, branch string) (*Template, error) {
	t := &Template{
		reg:        r,
		Name:       name,
		Repository: repo,
		Branch:     branch,
		Created:    time.Now(),
	}
	r.Templates = append(r.Templates, t)

	if err := download(repo, branch, t.ArchivePath()); err != nil {
		return nil, err
	}

	return t, r.save()
}

// RemoveTemplate removes the Template from disk and meta
func (r *Registry) RemoveTemplate(name string) error {
	_, err := r.GetTemplate(name)
	if err != nil {
		return err
	}
	for idx, t := range r.Templates {
		if strings.EqualFold(name, t.Name) {
			r.Templates = append(r.Templates[:idx], r.Templates[idx+1:]...)
			if err := os.Remove(t.ArchivePath()); err != nil {
				return err
			}
		}
	}

	return r.save()
}

// GetSource retrieves a Source from the Registry
func (r *Registry) GetSource(name string) (*Source, error) {
	for _, s := range r.Sources {
		if strings.EqualFold(name, s.Name) {
			return s, nil
		}
	}
	return nil, ErrSourceNotFound{Name: name}
}

// AddSource adds a new Source to the Registry
func (r *Registry) AddSource(url, name string) (*Source, error) {
	url = strings.TrimSuffix(url, "/") + "/"
	s := &Source{
		Name: name,
		URL:  url,
	}
	r.Sources = append(r.Sources, s)

	return s, r.save()
}

// RemoveSource removes the Source from the registry meta
func (r *Registry) RemoveSource(name string) error {
	_, err := r.GetSource(name)
	if err != nil {
		return err
	}
	for idx, s := range r.Sources {
		if strings.EqualFold(name, s.Name) {
			r.Sources = append(r.Sources[:idx], r.Sources[idx+1:]...)
		}
	}

	return r.save()
}

// Open opens a Registry, creating one if none exists at dir
func Open(dir string) (*Registry, error) {
	reg := Registry{
		dir: dir,
	}

	_, err := os.Lstat(reg.MetaFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			if err := create(reg.MetaFilePath()); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	tree, err := toml.LoadFile(reg.MetaFilePath())
	if err != nil {
		return nil, err
	}
	return &reg, tree.Unmarshal(&reg)
}

func create(regFile string) error {
	if err := os.MkdirAll(filepath.Dir(regFile), os.ModePerm); err != nil {
		return err
	}
	fi, err := os.Create(regFile)
	if err != nil {
		return err
	}
	return fi.Close()
}

func download(cloneURL, branch, dest string) error {
	tmp, err := ioutil.TempDir(os.TempDir(), "tmpl")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	// Clone the repo
	if _, err := git.PlainClone(tmp, false, &git.CloneOptions{
		URL:           cloneURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
	}); err != nil {
		return err
	}

	// RemoveTemplate .git
	if err := os.RemoveAll(filepath.Join(tmp, ".git")); err != nil {
		return err
	}

	// Make sure it's a valid template
	if _, err := os.Lstat(filepath.Join(tmp, "template.toml")); err != nil {
		return err
	}
	fi, err := os.Lstat(filepath.Join(tmp, "template"))
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("template found, expected directory")
	}

	// Create archive
	glob, err := filepath.Glob(filepath.Join(tmp, "*"))
	if err != nil {
		return err
	}

	if err := archiver.Archive(glob, dest); err != nil {
		return err
	}

	return nil
}
