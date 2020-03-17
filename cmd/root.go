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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/harrybrwn/yt/youtube"
	"github.com/spf13/cobra"
)

var (
	wg   sync.WaitGroup
	path string // TODO: fix this!!!

	cwd, _ = os.Getwd()

	ytTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{range $i, $alias := .Aliases}}
	{{- if $i}}, {{end -}}{{$alias}}
  {{- end}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{- end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}
{{- end -}}
{{if .HasAvailableFlags}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
{{- end}}
`
)

var rootCmd = &cobra.Command{
	Use:          "yt <command>",
	Short:        "A cli tool for downloading youtube videos.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("no Arguments\n\nUse \"yt help\" for more information")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(0)
	}
}

func init() {
	videoCmd := makeCommand("video", "youtube videos", ".mp4")
	videoCmd.Aliases = append(videoCmd.Aliases, "vid")
	audioCmd := makeCommand("audio", "audio from youtube videos", ".mpa")
	rootCmd.AddCommand(videoCmd)
	rootCmd.AddCommand(audioCmd)

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", cwd, "download path")
	rootCmd.SetUsageTemplate(ytTemplate)
}

func makeCommand(name, short, defaultExt string) *cobra.Command {
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [ids...]", name),
		Short:   fmt.Sprintf("A tool for downloading %s", short),
		Long:    fmt.Sprintf(`To download multiple videos use 'yt %s <id> <id>...'`, name),
		Aliases: []string{name[:1], name[:2]},
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleVideos(args, func(v *youtube.Video) error {
				var err error
				path, err = filepath.Abs(path)
				if err != nil {
					return err
				}

				ext, err := cmd.Flags().GetString("extension")
				if err != nil {
					return err
				}
				p := filepath.Join(path, v.FileName) + ext
				if name == "audio" {
					return v.DownloadAudio(p)
				} else if name == "video" {
					return v.Download(p)
				}
				return errors.New("bad command name")
			})
		},
	}
	c.Flags().StringP("extension", "e", defaultExt, "file extension used for video download")
	return c
}

func handleVideos(ids []string, fn func(*youtube.Video) error) error {
	if len(ids) == 0 {
		return errors.New("no Arguments\n\nUse \"yt [command] --help\" for more information about a command")
	}
	var (
		v   *youtube.Video
		err error
	)
	if len(ids) > 1 {
		wg.Add(len(ids))
		for _, id := range ids {
			v, err = youtube.NewVideo(id)
			if err != nil {
				return err
			}
			go func() {
				err = fn(v)
				if err != nil {
					panic(err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	} else if len(ids) == 1 {
		v, err := youtube.NewVideo(ids[0])
		if err != nil {
			return err
		}
		return fn(v)
	}
	return nil
}
