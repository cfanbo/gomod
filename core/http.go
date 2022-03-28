package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/cfanbo/gomod/core/log"
)

var (
	client     *http.Client
	repoPrefix map[string]string
)

func init() {
	repoPrefix = map[string]string{
		"golang.org/x/":              "github.com/golang/",
		"k8s.io/":                    "github.com/kubernetes/",
		"sigs.k8s.io/":               "github.com/kubernetes-sigs/",
		"google.golang.org/protobuf": "github.com/protocolbuffers/protobuf-go",
	}

	client = &http.Client{
		Timeout: time.Second * 10,
	}
}
func fetchBody(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		log.Err(err)
		return []byte{}, fmt.Errorf("an error occurred when fetch url %s, error: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, resp.Body)
		return []byte{}, fmt.Errorf("%s error: %w", url, errors.New("bad HTTP Response code:  "+strconv.Itoa(resp.StatusCode)))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
