package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestVarifyPlaylistPath(t *testing.T) {
	err := pathToTempFile(t, func(t *testing.T) {
		p, err := varifyPlaylistPath("some_id")

		if err != nil {
			t.Error(err)
		}
		if _, err = os.Stat(p); os.IsNotExist(err) {
			t.Error("path should exits")
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestMakeCommand(t *testing.T) {
	c := makeCommand("test", "test command", ".txt")
	if c.Use != "test [ids...]" {
		t.Error("wrong usage message")
	}
	ext, err := c.Flags().GetString("extension")
	if err != nil {
		t.Error(err)
	}
	if ext != ".txt" {
		t.Error("wrong default extention")
	}
	if err := c.RunE(c, []string{"cQ7STILAS0M"}); err == nil {
		t.Error("expected error")
	}

	if err := pathToTempFile(t, func(t *testing.T) {
		c = makeCommand("video", "test videos", ".mp4")
		if err := c.RunE(c, []string{"cQ7STILAS0M"}); err != nil {
			t.Error("run failed")
		}
	}); err != nil {
		t.Error(err)
	}
}

func pathToTempFile(t *testing.T, fn func(t *testing.T)) error {
	tempPath := path
	dir, err := ioutil.TempDir("", "yt_tests")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path = dir
	cwd = dir

	fn(t)

	cwd, err = os.Getwd()
	if err != nil {
		t.Error(err)
	}
	path = tempPath
	return nil
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
