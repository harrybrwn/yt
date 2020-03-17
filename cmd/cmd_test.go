package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestVarifyPlaylistPath(t *testing.T) {
	err := redirectPath(t, func(t *testing.T) {
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
	if err := c.RunE(c, []string{"I9JXjzKVKtA"}); err == nil {
		t.Error("expected error")
	}

	if err := redirectPath(t, func(t *testing.T) {
		c = makeCommand("video", "test videos", ".mp4")
		if err := c.RunE(c, []string{"I9JXjzKVKtA", "O9Ks3_8Nq1s"}); err != nil {
			t.Error("run failed")
		}

		c = makeCommand("audio", "test videos", ".mp4")
		if err := c.RunE(c, []string{"I9JXjzKVKtA", "O9Ks3_8Nq1s"}); err != nil {
			t.Error("run failed")
		}
	}); err != nil {
		t.Error(err)
	}
}

func TestDownloadPlaylist(t *testing.T) {
	if err := redirectPath(t, func(t *testing.T) {
		err := playlistCmd.RunE(playlistCmd, []string{"PLo7FOXNe7Yt9U0Qh1KBDjHQUuQ5BQR9Jt"})
		if err != nil {
			t.Error(err)
		}
	}); err != nil {
		t.Error(err)
	}
}

func redirectPath(t *testing.T, fn func(t *testing.T)) error {
	var err error
	pathCopy := path
	testDIR := filepath.Join(path, "TESTS")
	path = testDIR
	cwd = testDIR
	t.Log(path)

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}

	fn(t)

	err = os.RemoveAll(testDIR)
	if err != nil {
		return err
	}

	cwd, err = os.Getwd()
	if err != nil {
		return err
	}
	path = pathCopy
	return nil
}

func TestRootRun(t *testing.T) {
	err := rootCmd.RunE(rootCmd, []string{})
	if err == nil {
		t.Error("expected error")
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
