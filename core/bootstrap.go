package core

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/cfanbo/gomod/core/log"
	"golang.org/x/mod/modfile"
)

func Enter(modFile string) {
	var b []byte
	var err error

	// #####################
	// Custom Mode
	// #####################
	// args
	if modFile != "" {
		newUrl, err := ParseGoModURL(modFile)
		if err == nil {
			// github
			b, err = fetchBody(newUrl)
			if err != nil {
				fmt.Println(errors.Unwrap(err))
				return
			}
		} else {
			b, err = os.ReadFile(modFile)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		bootstrap(b)
		return
	}

	// #####################
	// Auto Detection Mode
	// #####################

	// current dir
	if len(b) == 0 {
		modFile = "go.mod"
		_, err := os.Stat(modFile)
		if !os.IsNotExist(err) {
			b, err = os.ReadFile(modFile)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	// env
	if len(b) == 0 {
		env := NewGoEnv()
		modFile, ok := env["GOMOD"]
		if !ok {
			fmt.Println("The environment variable GOMOD is Empty")
			return
		}
		b, err = os.ReadFile(modFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if len(b) == 0 {
		fmt.Println("data is empty")
		return
	}

	bootstrap(b)
}

func bootstrap(b []byte) {
	f, err := modfile.ParseLax("", b, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(f.Require) == 0 {
		fmt.Println("No module was found")
		return
	}

	InitCache()
	InitPkgMap()

	repos := NewRepos()

	var wg sync.WaitGroup
	wg.Add(len(f.Require))
	for i, _ := range f.Require {
		go func(wg *sync.WaitGroup, turl string) {
			defer wg.Done()
			repo := NewRepo(turl)
			repo.Do()
			repos.AddRepo(repo)
		}(&wg, f.Require[i].Mod.Path)
	}
	wg.Wait()

	// All failure
	var isPrint bool
	for _, repo := range repos.Repos {
		if repo.Star > -1 {
			isPrint = true
			break
		}
	}

	if !isPrint {
		fmt.Println("Please try again due to server problems")
		return
	}

	// dynamic sync request data for failed repo when print the terminal screen
	var isMonitor bool
	for _, repo := range repos.Repos {
		if !repo.Done {
			isMonitor = true
			break
		}
	}
	if isMonitor {
		log.SetDebugLevel(false)
		log.Debug().Msg("enable monitor")
		monitor := NewMonitor(repos)
		go monitor.Run()
	}

	cache.Remember()
	pkgMap.Remember()

	r := NewRender()
	r.SetRepos(repos)
	r.Run()
}

//func SubString(str string, len int) string {
//	if len < 1 {
//		return str
//	}
//
//	currLen := utf8.RuneCountInString(str)
//	if len >= currLen {
//		return str
//	}
//
//	r := []rune(str)[:len]
//	return string(r) + "..."
//}
