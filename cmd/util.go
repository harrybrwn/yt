package cmd

import "regexp"

var youtubeRegex = regexp.MustCompile(`^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)

func isurl(s string) bool {
	return youtubeRegex.MatchString(s)
}

func getid(url string) string {
	groups := youtubeRegex.FindAllStringSubmatch(url, -1)
	return groups[0][5]
}
