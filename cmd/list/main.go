package main

import (
	"errors"

	m "github.com/vault-thirteen/SWVD/models"
	s "github.com/vault-thirteen/SWVD/settings"
	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/file"
)

const (
	Err_NoPages = "no pages"
)

const (
	UsageHint = `Usage: list.exe <AuthorName> <PagesCount> <CookiesFilePath> <DataFolder>`
)

func main() {
	args, err := NewCommandLineArgs()
	if err != nil {
		cio.PrintMessage(UsageHint)
	}
	mustBeNoError(err)

	var cookies []m.Cookie
	cookies, err = m.ReadCookiesFromFile(args.cookiesFilePath)
	mustBeNoError(err)

	taskSettings := &m.TaskSettings{
		UseCustomTransport:   false,
		Schema:               s.SchemaCommon + s.SchemaPostfix,
		Domain:               s.Domain,
		Cookies:              cookies,
		Header_UserAgent:     s.BrowserUserAgent,
		PrintUrl:             true,
		DataFolder:           args.dataFolder,
		NewFolderPermissions: s.NewFolderPermissions,
		NewFilePermissions:   s.NewFilePermissions,
		AuthorName:           args.authorName,
		PagesCount:           args.pagesCount,
	}

	var tsk = NewTask(taskSettings)

	err = tsk.listAuthorsVideosN()
	mustBeNoError(err)

	err = tsk.saveAuthorVideoStreamsFromList()
	mustBeNoError(err)
}
func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Task) listAuthorsVideosN() (err error) {
	if t.s.PagesCount < 1 {
		return errors.New(Err_NoPages)
	}

	var searchPageUrls []string
	for n := 1; n <= t.s.PagesCount; n++ {
		searchPageUrls = append(searchPageUrls, t.um.MakeAuthorPageNUrl(t.s.AuthorName, n))
	}

	var authorHtmlFilePaths []string
	var filePath string
	for n := 1; n <= t.s.PagesCount; n++ {
		filePath = m.MakeRootFilePath(t.s.DataFolder, m.MakeAuthorHtmlFileNameN(t.s.AuthorName, n))
		authorHtmlFilePaths = append(authorHtmlFilePaths, filePath)
	}

	var allVideos = make([]*m.Video, 0)
	var pageVideos []*m.Video
	for i := 0; i < t.s.PagesCount; i++ {
		pageVideos, err = t.listAuthorsVideos(searchPageUrls[i], authorHtmlFilePaths[i])
		if err != nil {
			return err
		}

		allVideos = append(allVideos, pageVideos...)
	}

	listFilePath := m.MakeRootFilePath(t.s.DataFolder, m.MakeListFileName(t.s.AuthorName))
	err = m.SaveAuthorVideosListOnDisk(t.s.AuthorName, searchPageUrls[0], allVideos, listFilePath, t.s.NewFilePermissions)
	if err != nil {
		return err
	}

	return nil
}
func (t *Task) listAuthorsVideos(searchPageUrl string, authorHtmlFilePath string) (videos []*m.Video, err error) {
	var fileExists bool
	fileExists, err = file.FileExists(authorHtmlFilePath)
	if err != nil {
		return nil, err
	}

	var buf []byte
	if !fileExists {
		buf, err = m.FetchDataFromUrl(searchPageUrl, t.s)
		if err != nil {
			return nil, err
		}

		err = file.SaveBytesToFile(authorHtmlFilePath, buf)
		if err != nil {
			return nil, err
		}
	} else {
		buf, err = file.GetFileContents(authorHtmlFilePath)
		if err != nil {
			return nil, err
		}
	}

	videos, err = m.FindVideos(buf)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
func (t *Task) saveAuthorVideoStreamsFromList() (err error) {
	var authorHtmlFolderPath = m.MakeAuthorHtmlFolder(t.s.DataFolder, t.s.AuthorName)
	err = file.CreateFolderIfNotExists(authorHtmlFolderPath, t.s.NewFolderPermissions)
	if err != nil {
		return err
	}

	var authorJsonFolderPath = m.MakeAuthorJsonFolder(t.s.DataFolder, t.s.AuthorName)
	err = file.CreateFolderIfNotExists(authorJsonFolderPath, t.s.NewFolderPermissions)
	if err != nil {
		return err
	}

	var listFilePath = m.MakeRootFilePath(t.s.DataFolder, m.MakeListFileName(t.s.AuthorName))
	var videos []*m.Video
	videos, err = m.ReadVideosFromFile(listFilePath)
	if err != nil {
		return err
	}

	var videoPageFilePath string
	var fileExists bool
	var buf []byte
	for _, video := range videos {
		// Get HTML file.
		videoPageFilePath = m.MakeAuthorHtmlFilePath(t.s.DataFolder, t.s.AuthorName, m.MakeVideoPageFileName(video.Id))
		fileExists, err = file.FileExists(videoPageFilePath)
		if err != nil {
			return err
		}

		if !fileExists {
			buf, err = m.FetchDataFromUrl(video.PageUrl, t.s)
			if err != nil {
				return err
			}

			err = file.SaveBytesToFile(videoPageFilePath, buf)
			if err != nil {
				return err
			}
		} else {
			buf, err = file.GetFileContents(videoPageFilePath)
			if err != nil {
				return err
			}
		}

		var videoStreams []*m.VideoStream
		videoStreams, err = m.FindVideoStreams(buf)
		if err != nil {
			return err
		}

		var outputFilePath = m.MakeAuthorJsonFilePath(t.s.DataFolder, t.s.AuthorName, m.MakeVideoStreamsFileName(video.Id))
		var authorPage = t.um.MakeAuthorPageNUrl(t.s.AuthorName, 1)
		err = m.SaveVideoStreamsOnDisk(t.s.AuthorName, authorPage, video.Id, videoStreams, outputFilePath, t.s.NewFilePermissions)
		if err != nil {
			return err
		}
	}

	return nil
}
