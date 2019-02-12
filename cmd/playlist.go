// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path/filepath"
	"yt/youtube"

	"github.com/spf13/cobra"
)

var pExt string

var playlistCmd = &cobra.Command{
	Use:   "playlist [ids...]",
	Short: "A tool for downloading youtube playlists.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		audio, err := cmd.Flags().GetBool("audio")
		if err != nil {
			return err
		}

		for _, id := range args {
			err = downloadPlaylist(id, audio)
			if err != nil {
				return err
			}
		}
		wg.Wait()
		return nil
	},
}

func varifyPlaylistPath(id string) (string, error) {
	var err error
	p := path
	if p == cwd {
		p = filepath.Join(path, id)
	}
	p, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir)
	}
	return path, err
}

func downloadPlaylist(id string, getAudio bool) error {
	var err error
	var v *youtube.Video

	if err != nil {
		return err
	}
	path, err = varifyPlaylistPath(id)
	if err != nil {
		return err
	}

	plst, err := youtube.NewPlaylist(id)
	if err != nil {
		return err
	}
	for vID := range plst.VideoIds() {
		wg.Add(1)
		v, err = youtube.NewVideo(vID)
		if err != nil {
			return err
		}
		go func() {
			if getAudio {
				err = v.DownloadAudio(filepath.Join(path, v.FileName) + aExt)
			} else {
				err = v.Download(filepath.Join(path, v.FileName) + pExt)
			}

			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}
	return nil
}

func init() {
	playlistCmd.Flags().BoolP("audio", "a", false, "download the audio from all the videos in the specifies playlist")
	playlistCmd.Flags().StringVarP(&pExt, "extension", "e", ".mp4", "file extension used for video download")
	rootCmd.AddCommand(playlistCmd)
}
