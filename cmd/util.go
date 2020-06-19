package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/harrybrwn/yt/pkg/terminal"
)

var youtubeRegex = regexp.MustCompile(`^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)

func isurl(s string) bool {
	return youtubeRegex.MatchString(s)
}

func getid(url string) string {
	groups := youtubeRegex.FindAllStringSubmatch(url, -1)
	return groups[0][5]
}

func setCursorOnHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rStopped.          ")
		terminal.CursorOn()
		os.Exit(0)
	}()

}
