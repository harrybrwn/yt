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
	"path/filepath"
	"yt/youtube"

	"github.com/spf13/cobra"
)

var vExt string

var videoCmd = &cobra.Command{
	Use:   "video [ids...]",
	Short: "A tool for downloading youtube videos.",
	Long:  `To download multiple videos use 'yt video <id> <id>...'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleVideos(args, func(v *youtube.Video) error {
			var err error
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			return v.Download(filepath.Join(path, v.FileName) + vExt)
		})
	},
}

func init() {
	videoCmd.Flags().StringVarP(&vExt, "extension", "e", ".mp4", "file extension used for video download")
	rootCmd.AddCommand(videoCmd)
}
