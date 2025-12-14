package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/noah4ever/totui/internal/app"
)

// Load reads lists from the provided path. If the file does not exist,
// it returns an empty slice without error.
func Load(path string) ([]app.TodoList, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(abspath)
	if errors.Is(err, os.ErrNotExist) {
		return []app.TodoList{}, nil
	}
	if err != nil {
		return nil, err
	}

	var lists []app.TodoList
	if err := json.Unmarshal(b, &lists); err != nil {
		return nil, err
	}
	return lists, nil
}

// Save writes lists to the provided path, creating directories as needed.
func Save(path string, lists []app.TodoList) error {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(abspath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(lists, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(abspath, data, 0o644)
}
