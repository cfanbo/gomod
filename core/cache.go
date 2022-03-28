package core

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/cfanbo/gomod/core/log"
)

var (
	cache *Cache
)

func InitCache() {
	cache = NewCache(time.Hour * 24 * 7)
}

// Cache repo cache layer
type Cache struct {
	lifetime time.Duration
	mu       sync.RWMutex
	github   map[string]*Repo
	lc       Storage
}

func NewCache(d time.Duration) *Cache {
	cs := &FileStorage{}
	cs.SetID(1)

	c := &Cache{
		github: make(map[string]*Repo),
		lc:     cs,
	}
	if c != nil {
		c.lifetime = d
	}

	// try load data from cache
	c.tryLoadLocal()
	return c
}

func (c *Cache) tryLoadLocal() {
	c.mu.Lock()
	defer c.mu.Unlock()

	b := c.lc.Read()
	if b == nil {
		return
	}

	// Unmarshal
	var repos map[string]*Repo
	if err := json.Unmarshal(b, &repos); err != nil {
		return
	}

	// remove expire cache item
	for k, repo := range repos {
		if c.lifetime > 0 && time.Since(repo.UpdateAt) > c.lifetime {
			delete(repos, k)
		}
	}
	c.github = repos
}

func (c *Cache) Get(repoUrl string) (*Repo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if repo, ok := c.github[repoUrl]; ok {
		return repo, true
	}

	return nil, false
}

func (c *Cache) Set(repoUrl string, repo *Repo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.github[repoUrl] = repo
}

func (c *Cache) Remember() {
	c.mu.Lock()
	defer c.mu.Unlock()

	tmap := make(map[string]*Repo, 0)
	for k, repo := range c.github {
		if !repo.Done {
			continue
		}
		tmap[k] = repo
	}

	if len(tmap) == 0 {
		return
	}

	b, err := json.Marshal(tmap)
	if err != nil {
		return
	}

	c.lc.Write(b)
}

type Storage interface {
	SetID(Id int)
	Write(b []byte)
	Read() []byte
}

type FileStorage struct {
	uniqID int
}

func (c *FileStorage) SetID(Id int) {
	c.uniqID = Id
}

func (c *FileStorage) Write(b []byte) {
	log.Info().Msgf("write %s", c.cachefile())
	if err := os.WriteFile(c.cachefile(), b, 0644); err != nil {
		return
	}
}

func (c *FileStorage) Read() []byte {
	file := c.cachefile()
	log.Info().Msgf("read from file:%s", file)

	_, err := os.Stat(file)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
	}

	// read cache
	b, err := os.ReadFile(file)
	if err != nil {
		if err != io.EOF {
			return nil
		}
	}

	return b
}

func (c *FileStorage) cachefile() string {
	return filepath.Join(os.TempDir(), "_gomod.cache"+strconv.Itoa(c.uniqID))
}
