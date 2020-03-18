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
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/harrybrwn/yt/youtube"
	"github.com/spf13/cobra"
)

var (
	path   string // TODO: fix this!!!
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
	rootCmd.AddCommand(newinfoCmd(true))

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", cwd, "download path")
	rootCmd.SetUsageTemplate(ytTemplate)
}

type videoHandler func(v *youtube.Video) error

func makeCommand(name, short, defaultExt string) *cobra.Command {
	c := &cobra.Command{
		Use:     fmt.Sprintf("%s [ids...]", name),
		Short:   fmt.Sprintf("A tool for downloading %s", short),
		Long:    fmt.Sprintf(`To download multiple videos use 'yt %s <id> <id>...'`, name),
		Aliases: []string{name[:1], name[:2]},
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, err := cmd.Flags().GetString("extension")
			if err != nil {
				return err
			}
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			return handleVideos(args, func(v *youtube.Video) (err error) {
				p := filepath.Join(path, v.FileName) + ext
				if name == "audio" {
					return v.DownloadAudio(p)
				} else if name == "video" {
					err = v.Download(p)
					cmd.Printf("Downloaded %s\n", v.Title)
					return err
				}
				return errors.New("bad command name")
			})
		},
	}
	c.Flags().StringP("extension", "e", defaultExt, "file extension used for video download")
	return c
}

func handleVideos(ids []string, fn videoHandler) error {
	if len(ids) == 0 {
		return errors.New("no Arguments\n\nUse \"yt [command] --help\" for more information about a command")
	}
	if len(ids) > 1 {
		return asyncDownload(ids, fn)
	} else if len(ids) == 1 {
		v, err := youtube.NewVideo(ids[0])
		if err != nil {
			return err
		}
		return fn(v)
	}
	return nil
}

func newinfoCmd(hidden bool) *cobra.Command {
	type infocommand struct {
		fflags, playerResp bool
	}
	ic := infocommand{false, false}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get extra information for a youtube video",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := youtube.GetInfo(args[0])
			if err != nil {
				return err
			}
			if ic.fflags {
				return printfflags(info)
			}
			if ic.playerResp {
				fmt.Printf("%s\n", info["player_response"])
				return nil
			}
			for k, v := range info {
				if k == "player_response" || k == "fflags" {
					continue
				}
				fmt.Printf("%s: %s\n", k, v[0])
			}
			return nil
		},
		Hidden: hidden,
	}
	infoCmd.Flags().BoolVar(&ic.fflags, "fflags", ic.fflags, "print out the fflags")
	infoCmd.Flags().BoolVar(&ic.playerResp, "player-response", ic.playerResp, "print out the raw player response data")
	return infoCmd
}

func printfflags(info map[string][][]byte) error {
	f, ok := info["fflags"]
	if !ok || len(f) == 0 {
		return errors.New("could not find fflags")
	}
	data := string(f[0])

	res, err := url.ParseQuery(data)
	if err != nil {
		return err
	}
	for k, v := range res {
		fmt.Println(k, v)
	}
	return nil
}

func asyncDownload(ids []string, fn videoHandler) (err error) {
	var (
		wg sync.WaitGroup
		v  *youtube.Video
		e  error
	)
	wg.Add(len(ids))
	for _, id := range ids {
		v, err = youtube.NewVideo(id)
		if err != nil {
			return err
		}
		go func() {
			e = fn(v)
			if e != nil {
				log.Println(e)
				if err == nil {
					err = e
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return err
}
