package core

import (
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/cfanbo/gomod/core/log"
	"golang.org/x/sync/singleflight"
)

var (
	singleRequest singleflight.Group
)

type Repo struct {
	Star     int    `json:"star"`
	Fork     int    `json:"fork"`
	Watch    int    `json:"watch"`
	ImportBy int    `json:"import_by"`
	RepoUrl  string `json:"url"`
	Mod      string `json:"mod"`

	// 多个依赖是否共用一个仓库
	Shared bool `json:"shared"`
	Done   bool `json:"done"`

	UpdateAt time.Time `json:"update_at"`
}

func NewRepo(mod string) *Repo {
	return &Repo{
		Mod:      mod,
		Star:     -1,
		Fork:     -1,
		Watch:    -1,
		ImportBy: -1,
	}
}

func (r *Repo) Do() error {
	// check cache
	if url, ok := pkgMap.Get(r.Mod); ok {
		log.Printf("hit mapcache %s => %s", r.Mod, url)
		if url == "" {
			r.Star = 0
			r.Fork = 0
			r.Watch = 0
			r.ImportBy = 0
			r.Done = true
			r.UpdateAt = time.Now()
			return nil
		}

		r.RepoUrl = GetRepoHomeUrl(url)
	} else {
		url, err := GetRepoURL(r.Mod)
		if err != nil {
			log.Err(errors.Unwrap(err)).Send()
			return err
		}

		log.Printf("miss mapcache %s => %s", r.Mod, url)
		// set map cache
		pkgMap.Set(r.Mod, url)

		// eg: gitlab/user/repo
		if url == "" {
			log.Info().Msg(r.Mod)
			r.Star = 0
			r.Fork = 0
			r.Watch = 0
			r.ImportBy = 0
			r.Done = true
			r.UpdateAt = time.Now()
			return nil
		}

		r.RepoUrl = GetRepoHomeUrl(url)
	}

	// check cache
	if repo, ok := cache.Get(r.RepoUrl); ok {
		log.Printf("hit mcache %s", r.RepoUrl)
		r.Star = repo.Star
		r.Fork = repo.Fork
		r.Watch = repo.Watch
		r.Shared = repo.Shared
		r.Done = repo.Done
		r.UpdateAt = repo.UpdateAt
		return nil
	}

	log.Printf("miss mcache url= %s", r.RepoUrl)
	return r.do()
}

func (r *Repo) do() error {
	result, err, shared := singleRequest.Do(r.RepoUrl, func() (interface{}, error) {
		b, err := fetchBody(r.RepoUrl)
		if err != nil {
			return nil, err
		}
		body := string(b)

		result := make(map[string]int, 2)
		if n, err := ParseStar(body); err == nil {
			result["Star"] = n
		}

		if n, err := ParseFork(body); err == nil {
			result["Fork"] = n
		}

		return result, nil
	})
	singleRequest.Forget(r.RepoUrl)

	if err != nil {
		log.Printf("%s; %s", r.RepoUrl, err)
		return nil
	}

	ret := result.(map[string]int)
	if n, ok := ret["Star"]; ok {
		r.Star = n
	}
	if n, ok := ret["Fork"]; ok {
		r.Fork = n
	}
	r.Shared = shared
	r.Done = true
	r.UpdateAt = time.Now()

	// set cache
	cache.Set(r.RepoUrl, r)

	return nil
}

func (r *Repo) GetStar() string {
	if !r.Done && r.Star == -1 {
		return "?"
	}

	return strconv.Itoa(r.Star)
}
func (r *Repo) GetFork() string {
	if !r.Done && r.Fork == -1 {
		return "?"
	}

	return strconv.Itoa(r.Fork)
}

func (r *Repo) GetShared() string {
	if r.Shared {
		return "Y"
	}
	return ""
}

type Repos struct {
	Ch          <-chan *Repo
	mu          sync.RWMutex
	Repos       []*Repo
	lastKey     rune   // s、w、f、m、g
	lastAscDesc string // ASC or DESC
}

func NewRepos() *Repos {
	r := &Repos{
		Ch:          make(<-chan *Repo, 1),
		Repos:       make([]*Repo, 0),
		lastAscDesc: "DESC",
		lastKey:     's',
	}

	go func() {
		for {
			select {
			case repo, ok := <-r.Ch:
				if !ok {
					return
				}
				r.AddRepo(repo)
			}
		}
	}()

	return r
}

func (r *Repos) AddRepo(repo *Repo) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Repos = append(r.Repos, repo)
}

func (r *Repos) detectKey(key rune) {
	r.checkDuplicateKey(key)
	r.lastKey = key

	r.sort()
}

func (r *Repos) checkDuplicateKey(key rune) {
	if r.lastKey == key {
		if r.lastAscDesc == "DESC" {
			r.lastAscDesc = "ASC"
		} else {
			r.lastAscDesc = "DESC"
		}
	} else {
		r.lastAscDesc = "DESC"
	}
}

func (r *Repos) sort() {
	sort.Slice(r.Repos, func(i, j int) bool {
		var ret bool

		switch r.lastKey {
		case 'i', 'I':
			if r.lastAscDesc == "DESC" {
				ret = r.Repos[i].ImportBy > r.Repos[j].ImportBy
			} else {
				ret = r.Repos[i].ImportBy < r.Repos[j].ImportBy
			}

		case 's', 'S':
			if r.lastAscDesc == "DESC" {
				ret = r.Repos[i].Star > r.Repos[j].Star
			} else {
				ret = r.Repos[i].Star < r.Repos[j].Star
			}

		case 'f', 'F':
			if r.lastAscDesc == "DESC" {
				ret = r.Repos[i].Fork > r.Repos[j].Fork
			} else {
				ret = r.Repos[i].Fork < r.Repos[j].Fork
			}

		case 'w', 'W':
			if r.lastAscDesc == "DESC" {
				ret = r.Repos[i].Watch > r.Repos[j].Watch
			} else {
				ret = r.Repos[i].Watch < r.Repos[j].Watch
			}
		case 'm', 'M':
			if r.lastAscDesc == "DESC" {
				ret = r.Repos[i].Mod > r.Repos[j].Mod
			} else {
				ret = r.Repos[i].Mod < r.Repos[j].Mod
			}
		case 'g', 'G':
			if r.lastAscDesc == "DESC" {
				ret = r.Repos[i].RepoUrl > r.Repos[j].RepoUrl
			} else {
				ret = r.Repos[i].RepoUrl < r.Repos[j].RepoUrl
			}
		}

		return ret
	})
}
