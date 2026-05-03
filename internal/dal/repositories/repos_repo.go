package repositories

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RepoInfo struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Branches    []string `json:"branches"`
	LastCommit  string   `json:"last_commit"`
	Size        int64    `json:"size"`
}

type ReposRepo struct {
	BasePath string // /home/git/repos или /home/admin/repos
}

// NewReposRepo создает новый репозиторий для работы с git-репозиториями
func NewReposRepo(basePath string) *ReposRepo {
	return &ReposRepo{BasePath: basePath}
}

// GetAllRepos возвращает список всех репозиториев в папке
func (r *ReposRepo) GetAllRepos() ([]RepoInfo, error) {
	var repos []RepoInfo

	entries, err := os.ReadDir(r.BasePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения папки %s: %v", r.BasePath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			repoPath := filepath.Join(r.BasePath, entry.Name())

			// Проверяем, является ли папка git-репозиторием
			if r.isGitRepo(repoPath) {
				repoInfo, err := r.GetRepoInfo(entry.Name())
				if err == nil {
					repos = append(repos, repoInfo)
				}
			}
		}
	}

	return repos, nil
}

// GetRepoInfo возвращает информацию о конкретном репозитории
func (r *ReposRepo) GetRepoInfo(name string) (RepoInfo, error) {
	repoPath := filepath.Join(r.BasePath, name)

	if !r.isGitRepo(repoPath) {
		return RepoInfo{}, fmt.Errorf("репозиторий %s не найден или не является git-репозиторием", name)
	}

	info := RepoInfo{
		Name: name,
		Path: repoPath,
	}

	// Получаем список веток
	branches, _ := r.getBranches(repoPath)
	info.Branches = branches

	// Получаем последний коммит
	lastCommit, _ := r.getLastCommit(repoPath)
	info.LastCommit = lastCommit

	// Получаем размер репозитория
	size, _ := r.getRepoSize(repoPath)
	info.Size = size

	return info, nil
}

// CreateRepo создает новый git-репозиторий
func (r *ReposRepo) CreateRepo(name string) error {
	repoPath := filepath.Join(r.BasePath, name)

	// Проверяем, существует ли уже
	if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
		return fmt.Errorf("репозиторий %s уже существует", name)
	}

	// Создаем папку
	err := os.MkdirAll(repoPath, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания папки: %v", err)
	}

	// Инициализируем git репозиторий
	cmd := exec.Command("git", "init", "--bare", repoPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("ошибка инициализации git: %v", err)
	}

	return nil
}

// DeleteRepo удаляет репозиторий
func (r *ReposRepo) DeleteRepo(name string) error {
	repoPath := filepath.Join(r.BasePath, name)

	if !r.isGitRepo(repoPath) {
		return fmt.Errorf("репозиторий %s не найден", name)
	}

	err := os.RemoveAll(repoPath)
	if err != nil {
		return fmt.Errorf("ошибка удаления репозитория: %v", err)
	}

	return nil
}

// GetLastCommits возвращает последние коммиты
func (r *ReposRepo) GetLastCommits(name string, limit int) ([]string, error) {
	repoPath := filepath.Join(r.BasePath, name)

	if !r.isGitRepo(repoPath) {
		return nil, fmt.Errorf("репозиторий %s не найден", name)
	}

	cmd := exec.Command("git", "--git-dir="+repoPath, "log", "--oneline", fmt.Sprintf("-%d", limit))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения коммитов: %v", err)
	}

	commits := strings.Split(string(output), "\n")
	return commits, nil
}

// isGitRepo проверяет, является ли папка git-репозиторием
func (r *ReposRepo) isGitRepo(path string) bool {
	gitPath := filepath.Join(path, "objects")
	_, err := os.Stat(gitPath)
	return !os.IsNotExist(err)
}

// getBranches возвращает список веток репозитория
func (r *ReposRepo) getBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "--git-dir="+repoPath, "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(string(output), "\n")
	// Удаляем пустые строки
	var result []string
	for _, branch := range branches {
		if branch != "" {
			result = append(result, branch)
		}
	}
	return result, nil
}

// getLastCommit возвращает последний коммит
func (r *ReposRepo) getLastCommit(repoPath string) (string, error) {
	cmd := exec.Command("git", "--git-dir="+repoPath, "log", "--oneline", "-1")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getRepoSize возвращает размер репозитория в байтах
func (r *ReposRepo) getRepoSize(repoPath string) (int64, error) {
	var size int64
	err := filepath.Walk(repoPath, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
