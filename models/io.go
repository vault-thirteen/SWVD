package m

import (
	"encoding/json"
	"os"

	"github.com/vault-thirteen/auxie/file"
)

const (
	JsonFileIndent = "\t"
)

func SaveAuthorVideosListOnDisk(authorName string, authorPageUrl string, videos []*Video, filePath string, filePermissions os.FileMode) (err error) {
	var obj = AuthorVideos{
		AuthorName:  authorName,
		AuthorPage:  authorPageUrl,
		VideosCount: len(videos),
		Videos:      videos,
	}

	return SaveObjectOnDisk(obj, filePath, filePermissions)
}
func SaveVideoStreamsOnDisk(authorName string, authorPage string, videoId string, videoStreams []*VideoStream, filePath string, filePermissions os.FileMode) (err error) {
	var obj = VideoStreams{
		AuthorName: authorName,
		AuthorPage: authorPage,
		VideoId:    videoId,
		Streams:    videoStreams,
	}

	return SaveObjectOnDisk(obj, filePath, filePermissions)
}
func SaveObjectOnDisk(object any, filePath string, filePermissions os.FileMode) (err error) {
	var buf []byte
	buf, err = json.MarshalIndent(object, "", JsonFileIndent)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, buf, filePermissions)
	if err != nil {
		return err
	}

	return nil
}
func ReadObjectFromDisk[T any](filePath string) (object *T, err error) {
	var buf []byte
	buf, err = file.GetFileContents(filePath)
	if err != nil {
		return nil, err
	}

	object = new(T)
	err = json.Unmarshal(buf, object)
	if err != nil {
		return nil, err
	}

	return object, nil
}
