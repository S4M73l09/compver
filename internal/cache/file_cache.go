package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type FileCache struct {
	directory string
}

func NewFileCache(directory string) (*FileCache, error) {
	if err := os.MkdirAll(directory, 0700); err != nil {
		return nil, err
	}

	return &FileCache{
		directory: directory,
	}, nil
}

func DefaultFileCache() (*FileCache, error) {
	userCacheDirectory, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	return NewFileCache(
		filepath.Join(userCacheDirectory, "compver"),
	)
}

func (c *FileCache) Get(
	key string,
) (Entry, bool, error) {
	path := c.pathFor(key)

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Entry{}, false, nil
	}
	if err != nil {
		return Entry{}, false, err
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, false, err
	}

	return entry, true, nil
}

func (c *FileCache) Set(
	key string,
	entry Entry,
) error {
	entry.Key = key

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		c.pathFor(key),
		data,
		0600,
	)
}

func (c *FileCache) Delete(key string) error {
	err := os.Remove(c.pathFor(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func (c *FileCache) pathFor(key string) string {
	hash := sha256.Sum256([]byte(key))
	filename := hex.EncodeToString(hash[:]) + ".json"

	return filepath.Join(c.directory, filename)
}