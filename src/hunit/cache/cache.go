package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"

	"github.com/instaunit/instaunit/hunit"
)

type Resource struct {
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
}

func Checksum(path string) (*Resource, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return &Resource{
		Path:     path,
		Checksum: fmt.Sprintf("%x", h.Sum(nil)),
	}, nil
}

type Cache struct {
	Version string                     `json:"version"`
	Binary  *Resource                  `json:"binary,omitempty"`
	Suites  []*Resource                `json:"suites,omitempty"`
	Results map[string][]*hunit.Result `json:"results,omitempty"`
}

func Read(path string) (*Cache, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := &Cache{}
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Write(path string, cache *Cache) error {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cache)
}
