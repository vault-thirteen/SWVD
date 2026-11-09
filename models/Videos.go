package m

type Videos struct {
	AuthorName string   `json:"authorName"`
	AuthorPage string   `json:"authorPage"`
	Videos     []*Video `json:"videos"`
}

func ReadVideosFromFile(filePath string) (vs []*Video, err error) {
	var x *Videos
	x, err = ReadObjectFromDisk[Videos](filePath)
	if err != nil {
		return nil, err
	}

	return x.Videos, nil
}
