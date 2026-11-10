package m

import (
	"path/filepath"
	"strconv"

	"github.com/vault-thirteen/auxie/cio"
)

const (
	Folder_Video = "video"

	Ext_Html  = "html"
	Ext_Json  = "json"
	Ext_Mpeg4 = "mp4"
	Ext_Video = Ext_Mpeg4

	Dot                 = "."
	Underscore          = "_"
	PageMarker          = "page"
	VideoFileNameFormat = "%s-%sp" + Dot + Ext_Video
)

func MakeAuthorHtmlFileNameN(authorName string, pageNumber int) string {
	return authorName + Underscore + PageMarker + Underscore + strconv.Itoa(pageNumber) + Dot + Ext_Html
}
func MakeListFileName(authorName string) string {
	return authorName + Dot + Ext_Json
}
func MakeVideoPageFileName(videoId string) string {
	return videoId + Dot + Ext_Html
}
func MakeVideoStreamsFileName(videoId string) string {
	return videoId + Dot + Ext_Json
}
func FormatVideoFileName(videoId string, resolution string) string {
	return cio.FormatMessage(VideoFileNameFormat, videoId, resolution)
}
func MakeReportFileName(authorName string, timeString string) string {
	return authorName + Underscore + timeString + Dot + Ext_Json
}
func MakeReportFilePath(reportFolder string, fileName string) string {
	return filepath.Join(reportFolder, fileName)
}

func MakeRootFilePath(dataFolder string, filename string) string {
	return filepath.Join(dataFolder, filename)
}
func MakeVideoFolderPath(dataFolder string) string {
	return filepath.Join(dataFolder, Folder_Video)
}
func MakeAuthorVideoFolderPath(dataFolder string, authorName string) string {
	return filepath.Join(dataFolder, Folder_Video, authorName)
}
func MakeAuthorJsonFilePath(dataFolder string, authorName string, filename string) string {
	return filepath.Join(MakeAuthorJsonFolder(dataFolder, authorName), filename)
}
func MakeAuthorHtmlFilePath(dataFolder string, authorName string, filename string) string {
	return filepath.Join(MakeAuthorHtmlFolder(dataFolder, authorName), filename)
}
func MakeAuthorJsonFolder(dataFolder string, authorName string) string {
	return filepath.Join(dataFolder, authorName+Underscore+Ext_Json)
}
func MakeAuthorHtmlFolder(dataFolder string, authorName string) string {
	return filepath.Join(dataFolder, authorName+Underscore+Ext_Html)
}
