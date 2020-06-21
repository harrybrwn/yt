package youtube

import "fmt"

// Represents video meta-data for a youtube video
type baseVideo struct {
	Author        string `json:"author"`
	ChannelID     string `json:"channelId"`
	LengthSeconds string `json:"lengthSeconds"`
	Title         string `json:"title"`
	ID            string `json:"videoId"`
	ViewCount     string `json:"viewCount"`
}

// The VideoData struct is an intermediate struct between the video's raw json string
// and the Video object.
type VideoData struct {
	StreamingData struct {
		AdaptiveFormats []Stream `json:"adaptiveFormats"`
		Formats         []Stream `json:"formats"`
	} `json:"streamingData"`
	VideoDetails struct {
		baseVideo
		Keywords  []string `json:"keywords"`
		Thumbnail struct {
			Thumbnails []Thumbnail
		}
	} `json:"videoDetails"`
	PlayabilityStatus *playabilityStatus `json:"playabilityStatus"`
}

type playabilityStatus struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (ps *playabilityStatus) Error() string {
	return fmt.Sprintf("%s: %s", ps.Status, ps.Reason)
}

// PlaylistInitData is meant to be an intermediate struct for going from raw
// data to a Playlist object
type PlaylistInitData struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Selected       bool   `json:"selected"`
					TrackingParams string `json:"trackingParams"`
					Content        struct {
						SectionListRenderer struct {
							TrackingParams string `json:"trackingParams"`
							Contents       []struct {
								ItemSectionRenderer struct {
									TrackingParams string `json:"trackingParams"`
									Contents       []struct {
										PlaylistVideoListRenderer struct {
											Contents []struct {
												PlaylistVideoRenderer struct {
													VideoID string `json:"videoId"`
												} `json:"playlistVideoRenderer"`
											} `json:"contents"`
										} `json:"playlistVideoListRenderer"`
									} `json:"contents"`
								} `json:"itemSectionRenderer"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}
