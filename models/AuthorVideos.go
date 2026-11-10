package m

type AuthorVideos struct {
	AuthorName  string   `json:"authorName"`
	AuthorPage  string   `json:"authorPage"`
	VideosCount int      `json:"videosCount"`
	Videos      []*Video `json:"videos"`
}

func ReadVideosFromFile(filePath string) (vs []*Video, err error) {
	var x *AuthorVideos
	x, err = ReadObjectFromDisk[AuthorVideos](filePath)
	if err != nil {
		return nil, err
	}

	return x.Videos, nil
}
