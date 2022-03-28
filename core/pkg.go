package core

import (
	"encoding/json"
	"sync"
)

var pkgMap *PkgMap

func InitPkgMap() {
	pkgMap = NewPkgMap(2000)
}

type PkgMap struct {
	mu   sync.RWMutex
	m    map[string]string
	lc   Storage
	size int
}

func NewPkgMap(size int) *PkgMap {
	cs := &FileStorage{}
	cs.SetID(0)

	c := &PkgMap{
		m:    make(map[string]string),
		lc:   cs,
		size: size,
	}
	c.tryLoadLocal()

	return c
}

func (c *PkgMap) Get(pkg string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.m[pkg]
	return v, ok
}

func (c *PkgMap) Set(pkg, github string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[pkg] = github
}

func (c *PkgMap) Remember() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	b, err := json.Marshal(c.m)
	if err != nil {
		return
	}

	c.lc.Write(b)
}

func (c *PkgMap) tryLoadLocal() {
	b := c.lc.Read()
	if b == nil {
		return
	}

	// Unmarshal
	var repos map[string]string
	if err := json.Unmarshal(b, &repos); err != nil {
		return
	}

	// remove objects at random
	if len(repos) > c.size {
		var i int
		deleteCount := len(repos) - c.size
		for k, _ := range repos {
			if i >= deleteCount {
				break
			}
			delete(repos, k)
			i++
		}
	}

	c.m = repos
}
