package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"

	pathutil "path"

	"github.com/instaunit/instaunit/hunit"
)

type Resource struct {
	Path     string `json:"path"`
	Checksum string `json:"checksum"`
}

func Checksum(path string) (*Resource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, err
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
	Results map[string][]*hunit.Result `json:"results,omitempty"` // checksum -> []results
	suites  map[string]*Resource       `json:"-"`                 // checksum -> Resource
}

func (c *Cache) Suite(checksum string) *Resource {
	return c.suites[checksum]
}

func (c *Cache) AddSuite(suite *Resource, results []*hunit.Result) {
	c.Suites = append(c.Suites, suite)
	if c.suites == nil {
		c.suites = map[string]*Resource{suite.Checksum: suite}
	} else {
		c.suites[suite.Checksum] = suite
	}
	if c.Results == nil {
		c.Results = map[string][]*hunit.Result{suite.Checksum: results}
	} else {
		c.Results[suite.Checksum] = results
	}
}

func (c *Cache) ResultsForSuite(suite *Resource) []*hunit.Result {
	if c.Results != nil {
		return c.Results[suite.Checksum]
	} else {
		return nil
	}
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

	c.suites = make(map[string]*Resource)
	for _, e := range c.Suites {
		c.suites[e.Checksum] = e
	}

	return c, nil
}

func Write(path string, cache *Cache) error {
	err := os.MkdirAll(pathutil.Dir(path), 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cache)
}
