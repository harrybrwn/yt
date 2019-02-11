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
	"errors"
	"os"
	"sync"
	"yt/youtube"

	"github.com/spf13/cobra"
)

var (
	wg sync.WaitGroup

	// flag variabels
	cfgFile string
	path    string
	logging bool
	timed   bool // hidden flag var

	cwd, _     = os.Getwd()
	ytTemplate = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
	
Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.Flags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableFlags}}
	
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`
)

var rootCmd = &cobra.Command{
	Use:          "yt [command]",
	Short:        "A cli tool for youtube videos",
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
	if logging {
		youtube.Logging = true
	}

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", cwd, "download path")
	rootCmd.PersistentFlags().BoolVar(&logging, "log", false, "toggle the internal logger")
	// rootCmd.PersistentFlags().MarkHidden("time")

	rootCmd.SetUsageTemplate(ytTemplate)
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
			log("made vid, about to download...")
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
		log("made only vid, about to download it...")
		return fn(v)
	}
	return nil
}
