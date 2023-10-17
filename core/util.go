package core

import (
	"fmt"
	"github.com/cfanbo/gomod/pkg"
	"github.com/skratchdot/open-golang/open"
	"regexp"
	"strconv"
	"strings"

	"github.com/cfanbo/gomod/core/log"
)

const (
	GITHUB = "github.com"
)

// ParseStar return count of star
func ParseStar(str string) (int, error) {
	exp := regexp.MustCompile(`title="[\d]+(,?)?[\d]*" data-view-component="true" class="Counter js-social-count"`)
	result := exp.FindString(str)
	if result == "" {
		return -1, fmt.Errorf("parse Star Error")
	}

	return parseNumValue(result)
}

// ParseFork return count of fork
func ParseFork(str string) (int, error) {
	exp := regexp.MustCompile(`<span id="repo-network-counter" data-pjax-replace="true" title="[\d]+(,?)?[\d]*" data-view-component="true" class="Counter">`)
	result := exp.FindString(str)
	if result == "" {
		return -1, fmt.Errorf("parse Fork Error")
	}

	return parseNumValue(result)
}

func parseNumValue(result string) (int, error) {
	exp := regexp.MustCompile(`[\d]+(,?)?[\d]*`)
	result = exp.FindString(result)
	result = strings.ReplaceAll(result, ",", "")
	if result == "" {
		return -1, fmt.Errorf("parseNumValue Error")
	}
	return strconv.Atoi(result)
}

// ParseRepoURL return repo's'url on github
func ParseRepoURL(str string) string {
	exp := regexp.MustCompile(`\n(\s*)<a href="(http|https)://[a-zA-Z0-9]+(.)[com|org]{1,}/[a-zA-Z0-9]+(.*)" title="(http|https)://[a-zA-Z0-9]+(.)[com|org]{1,}/[a-zA-Z0-9]+(.*)" target="_blank" rel="noopener">\n(\s*)github.com/([a-zA-Z0-9])(.+)+\n(\s*)</a>`)
	out := exp.FindString(str)
	if out == "" {
		return out
	}

	exp2 := regexp.MustCompile(`https://github.com/[a-zA-Z0-9]+[a-zA-Z0-9_-]*/[a-zA-Z0-9]+[a-zA-Z0-9_-]+`)

	return exp2.FindString(out)
}

// GetRepoStar return
//func GetRepoStar(url string) (int, error) {
//	url, err := GetRepoURL(url)
//	if err != nil {
//		return 0, err
//	}
//
//	homeURL := GetRepoHomeUrl(url)
//	return GetGHStar(homeURL)
//}

// GetRepoURL return request URL, eg: github.com/user/repo
func GetRepoURL(url string) (string, error) {
	if strings.HasPrefix(url, GITHUB) {
		return url, nil
	}

	for prefix, repo := range repoPrefix {
		if strings.HasPrefix(url, prefix) {
			newUrl := strings.Replace(url, prefix, repo, 1)
			if newUrl != url {
				return newUrl, nil
			}
		}
	}
	return getRepoURL(url)
}

func getRepoURL(url string) (string, error) {
	log.Printf("request go.dev/%s", url)
	b, err := fetchBody("https://pkg.go.dev/" + url)
	if err != nil {
		log.Err(err)
		return "", fmt.Errorf("%w", err)
	}

	return ParseRepoURL(string(b)), nil
}

// GetGHStar 获取仓库star, 如果出错，返回-1
//func GetGHStar(url string) (int, error) {
//	body, err := fetchBody(url)
//	if err != nil {
//		return -1, err
//	}
//
//	str := string(body)
//	//exp := regexp.MustCompile(`class="Counter js-social-count">[0-9]+</span>`)
//	exp := regexp.MustCompile(`title="[0-9]+,?[0-9]+" data-view-component="true" class="Counter js-social-count"`)
//	result := exp.FindString(str)
//	if result == "" {
//		return -2, err
//	}
//
//	// 分隔符移除
//	result = strings.ReplaceAll(result, ",", "")
//	exp2 := regexp.MustCompile(`[0-9]+`)
//	result = exp2.FindString(result)
//	if result == "" {
//		return -3, err
//	}
//
//	n, err := strconv.Atoi(result)
//	if err != nil {
//		return -4, err
//	}
//
//	return n, nil
//}

// GetRepoHomeUrl github.com/a/b/c/d => https://github.com/a/b
func GetRepoHomeUrl(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	s := strings.SplitN(url, "/", 4)

	s = s[:3]
	url = strings.Join(s, "/")

	return "https://" + url
}

// ParseGoModURL https://github.com/user/repo => https://raw.githubusercontent.com/user/repo/master/go.mod
func ParseGoModURL(url string) (string, error) {
	if strings.Contains(url, "raw.githubusercontent.com/") {
		if strings.HasPrefix(url, "https://raw.githubusercontent.com/") {
			return url, nil
		}

		return "https://" + url, nil
	}

	if strings.Contains(url, "github.com/") {
		exp := regexp.MustCompile(`github.com/[a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+`)
		ret := exp.FindString(url)
		if ret == "" {
			return url, fmt.Errorf("invalid URL")
		}

		// TODO
		branch := "master"

		newUrl := GetRepoHomeUrl(url)
		newUrl = strings.Replace(newUrl, "github.com/", "raw.githubusercontent.com/", 1)
		newUrl += "/" + branch + "/go.mod"

		return newUrl, nil
	}

	return url, nil
}

func getRepoList(r *Repos) []interface{} {
	var dataList []interface{}
	dataList = make([]interface{}, len(r.Repos))
	for k, v := range r.Repos {
		dataList[k] = v
	}

	return dataList
}

// open window in bs
func openWindow(pager *pkg.Pager, origin string) {
	item := pager.SelectedRecord()
	repo, ok := item.(*Repo)
	if !ok {
		return
	}

	if origin == "github" && repo.RepoUrl != "" {
		open.Run(repo.RepoUrl)
	} else if repo.Mod != "" {
		open.Run("https://pkg.go.dev/" + repo.Mod)
	}
}
