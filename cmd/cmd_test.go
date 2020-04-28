package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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
		t.Errorf("wrong default extension: got: '%s'; want '.txt'", ext)
	}
	if err := c.RunE(c, []string{"fR2xOh8CqMM"}); err == nil {
		t.Error("expected error")
	}

	if err := redirectPath(t, func(t *testing.T) {
		c = makeCommand("video", "test videos", ".mp4")
		if err := c.RunE(c, []string{"fR2xOh8CqMM", "O9Ks3_8Nq1s"}); err != nil {
			t.Error("run failed", err)
		}

		c = makeCommand("audio", "test videos", ".mpa")
		if err := c.RunE(c, []string{"fR2xOh8CqMM", "O9Ks3_8Nq1s"}); err != nil {
			t.Error("run failed", err)
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

	baseTestPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("yt_cmd_tests%d", time.Now().UnixNano()),
	)
	testDIR := filepath.Join(baseTestPath, "TESTS")
	path = testDIR
	cwd = testDIR

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	fn(t)

	err = os.RemoveAll(testDIR)
	if err != nil {
		return err
	}
	err = os.Remove(baseTestPath)
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

func TestUtils(t *testing.T) {
	url := "https://www.youtube.com/watch?v=kJQP7kiw5Fk"

	if !isurl(url) {
		t.Error("this is a url")
	}
	if getid(url) != "kJQP7kiw5Fk" {
		t.Error("got wrong video id")
	}
}
