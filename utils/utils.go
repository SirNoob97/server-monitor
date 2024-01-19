package utils

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func GetEnv(key string, def string) string {
	variable, ok := syscall.Getenv(key)
	if ok {
		return variable
	}
	return def
}

func JoinPath(basePath string, segments ...string) string {
	if len(segments) == 0 {
		return basePath
	}
	segList := make([]string, len(segments)+1)
	segList[0] = basePath
	copy(segList[1:], segments)
	return filepath.Join(segList...)
}

func ReadLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	var ret []string

	for sc.Scan() {
		ret = append(ret, sc.Text())
	}

	if sc.Err() != nil {
		return nil, sc.Err()
	}

	return ret, nil
}

func ReadFile(fl FileLocation) ([]string, error) {
	env := GetEnv(fl.Env, fl.EnvDefaultVal)
	filename := JoinPath(env, fl.Segments...)
	return ReadLines(filename)
}

func Sleep(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
