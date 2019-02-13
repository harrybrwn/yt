package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestVarifyPlaylistPath(t *testing.T) {
	tempPath := path

	dir, err := ioutil.TempDir("yt_tests", "yt")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	path = dir
	cwd = dir

	p, err := varifyPlaylistPath("some_id")
	if err != nil {
		t.Error(err)
	}
	if _, err = os.Stat(p); os.IsNotExist(err) {
		t.Error("path should exits")
	}
	cwd, err = os.Getwd()
	if err != nil {
		t.Error(err)
	}
	path = tempPath
	t.Error(p)
	t.Error(path)
}

func tempfile() string {
	f, err := ioutil.TempFile("", "yt")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	dir, err := ioutil.TempDir("", "yt")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	return filepath.Join(dir, filepath.Base(f.Name()))
}
