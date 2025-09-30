package codegen

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"path/filepath"
)

type GitRepository struct {
	repository *git.Repository

	head *plumbing.Reference
}

func (p *GitRepository) HeadName() string {
	return p.head.Name().Short()
}

func (p *GitRepository) HeadHash() string {
	return p.head.Hash().String()
}

func ParseGit(gitDir string) (*GitRepository, error) {
	gitDir, err := FindGitRoot(gitDir)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	p := &GitRepository{
		repository: repo,
	}

	p.head, err = repo.Head()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return p, nil
}

func FindGitRoot(dir string) (string, error) {
	for filepath.Dir(dir) != dir {
		if osx.IsDir(filepath.Join(dir, ".git")) {
			return dir, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", fmt.Errorf("not found .git directory in %s or its parents", dir)
}
