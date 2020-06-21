// Copyright © 2019 Harrison Brown harrybrown98@gmail.com
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
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/harrybrwn/yt/pkg/terminal"
	"github.com/harrybrwn/yt/youtube"
	"github.com/spf13/cobra"
)

var pExt string

var playlistCmd = &cobra.Command{
	Use:     "playlist [ids...]",
	Short:   "A tool for downloading youtube playlists.",
	Aliases: []string{"p", "plst"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err error
			wg  sync.WaitGroup
		)
		audio, err := cmd.Flags().GetBool("audio")
		if err != nil {
			return err
		}

		setCursorOnHandler()
		terminal.CursorOff()
		defer terminal.CursorOn()
		go func() {
			for i := 0; ; i++ {
				fmt.Printf("\r%s...  ", terminal.Red("Downloading"))
				printLoadingChar(i)
				time.Sleep(loadingInterval)
			}
		}()

		for _, plstID := range args {
			err = downloadPlaylist(plstID, audio, &wg)
			if err != nil {
				panic(err)
			}
		}
		wg.Wait()
		return nil
	},
}

func verifyPlaylistPath(id string) (string, error) {
	var err error
	p := path
	if p == cwd {
		p = filepath.Join(p, id)
	}
	p, err = filepath.Abs(p)
	if err != nil {
		return "", err
	}

	if _, err = os.Stat(p); os.IsNotExist(err) {
		err = os.Mkdir(p, 0755)
		if err != nil {
			return "", err
		}
	}
	return p, nil
}

func downloadPlaylist(id string, getAudio bool, wg *sync.WaitGroup) error {
	var err error
	var v *youtube.Video

	plst, err := youtube.NewPlaylist(id)
	if err != nil {
		return err
	}
	path = filepath.Join(path, plst.Title)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			return err
		}
	}

	for vID := range plst.VideoIds() {
		wg.Add(1)
		v, err = youtube.NewVideo(vID)
		if err != nil {
			return err
		}
		go func() {
			name := filepath.Join(path, v.FileName)
			if getAudio {
				name += ".mpa"
				err = v.DownloadAudio(name)
			} else {
				name += pExt
				err = v.Download(name)
			}

			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Printf("\r%s %s\n", terminal.Green("Downloaded"), name)
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
