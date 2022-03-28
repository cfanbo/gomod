package core

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"sync"
)

var env map[string]string

var onceEnv sync.Once

func NewGoEnv() map[string]string {
	onceEnv.Do(func() {
		env = make(map[string]string, 30)
		cmd := exec.Command("go", "env")
		out, _ := cmd.StdoutPipe()
		err := cmd.Start()
		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(out)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err == io.EOF {
				break
			}

			l := strings.Split(line, "=")
			v := strings.ReplaceAll(l[1], "\n", "")
			v = strings.ReplaceAll(v, "\"", "")
			env[l[0]] = v
		}
	})

	return env
}
