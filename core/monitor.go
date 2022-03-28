package core

import (
	"time"

	"github.com/cfanbo/gomod/core/log"
)

// var monitor *Monitor

type Monitor struct {
	repos *Repos
}

func NewMonitor(r *Repos) *Monitor {
	return &Monitor{
		repos: r,
	}
}

func (m *Monitor) Run() {
	retryCh := make(chan *Repo, 1)

	go func() {
		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				repo, ok := <-retryCh
				if !ok {
					goto Loop
				}
				log.Debug().Str("Repo", repo.RepoUrl).Msg("retry synchronizing data by monitor server")
				repo.Do()
			}
		}

	Loop:
		cache.Remember()
		pkgMap.Remember()

		//render
		m.repos.sort()
		render.LoadPagerData()
		render.render()
	}()

	go func() {
		retryLimit := 5
		for i := 0; i < retryLimit; i++ {
			foundRetryRepo := false
			for _, repo := range m.repos.Repos {
				if !repo.Done {
					foundRetryRepo = true
					retryCh <- repo
				}
			}

			if !foundRetryRepo {
				break
			}
		}

		close(retryCh)
	}()
}
