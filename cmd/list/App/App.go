package a

import (
	m "github.com/vault-thirteen/SWVD/models"
	"github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/file"
)

type App struct {
	ts *m.TaskSettings
	um *m.UrlMaker
}

func NewApp(ts *m.TaskSettings) (app *App) {
	app = &App{
		ts: ts,
		um: m.NewUrlMaker(ts.Common.Schema, ts.Common.Domain),
	}

	return app
}

func (a *App) ListAuthorsVideosN() (err error) {
	if a.ts.Custom.PagesCount < 1 {
		return errors.ErrorFromText(m.Err_NoPages)
	}

	var dataFolderPath = a.ts.Common.DataFolder
	err = file.CreateFolderIfNotExists(dataFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	var searchPageUrls []string
	for n := 1; n <= a.ts.Custom.PagesCount; n++ {
		searchPageUrls = append(searchPageUrls, a.um.MakeAuthorPageNUrl(a.ts.Custom.AuthorName, n))
	}

	var authorHtmlFilePaths []string
	var filePath string
	for n := 1; n <= a.ts.Custom.PagesCount; n++ {
		filePath = m.MakeRootFilePath(a.ts.Common.DataFolder, m.MakeAuthorHtmlFileNameN(a.ts.Custom.AuthorName, n))
		authorHtmlFilePaths = append(authorHtmlFilePaths, filePath)
	}

	var allVideos = make([]*m.Video, 0)
	var pageVideos []*m.Video
	for i := 0; i < a.ts.Custom.PagesCount; i++ {
		pageVideos, err = a.ListAuthorsVideos(searchPageUrls[i], authorHtmlFilePaths[i])
		if err != nil {
			return err
		}

		allVideos = append(allVideos, pageVideos...)
	}

	listFilePath := m.MakeRootFilePath(a.ts.Common.DataFolder, m.MakeListFileName(a.ts.Custom.AuthorName))
	err = m.SaveAuthorVideosListOnDisk(a.ts.Custom.AuthorName, searchPageUrls[0], allVideos, listFilePath, a.ts.Common.NewFilePermissions)
	if err != nil {
		return err
	}

	return nil
}
func (a *App) ListAuthorsVideos(searchPageUrl string, authorHtmlFilePath string) (videos []*m.Video, err error) {
	var fileExists bool
	fileExists, err = file.FileExists(authorHtmlFilePath)
	if err != nil {
		return nil, err
	}

	var buf []byte
	if !fileExists {
		buf, err = m.FetchDataFromUrl(searchPageUrl, a.ts)
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
func (a *App) SaveAuthorVideoStreamsFromList() (err error) {
	var authorHtmlFolderPath = m.MakeAuthorHtmlFolder(a.ts.Common.DataFolder, a.ts.Custom.AuthorName)
	err = file.CreateFolderIfNotExists(authorHtmlFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	var authorJsonFolderPath = m.MakeAuthorJsonFolder(a.ts.Common.DataFolder, a.ts.Custom.AuthorName)
	err = file.CreateFolderIfNotExists(authorJsonFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	var listFilePath = m.MakeRootFilePath(a.ts.Common.DataFolder, m.MakeListFileName(a.ts.Custom.AuthorName))
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
		videoPageFilePath = m.MakeAuthorHtmlFilePath(a.ts.Common.DataFolder, a.ts.Custom.AuthorName, m.MakeVideoPageFileName(video.Id))
		fileExists, err = file.FileExists(videoPageFilePath)
		if err != nil {
			return err
		}

		if !fileExists {
			buf, err = m.FetchDataFromUrl(video.PageUrl, a.ts)
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

		var outputFilePath = m.MakeAuthorJsonFilePath(a.ts.Common.DataFolder, a.ts.Custom.AuthorName, m.MakeVideoStreamsFileName(video.Id))
		var authorPage = a.um.MakeAuthorPageNUrl(a.ts.Custom.AuthorName, 1)
		err = m.SaveVideoStreamsOnDisk(a.ts.Custom.AuthorName, authorPage, video.Id, videoStreams, outputFilePath, a.ts.Common.NewFilePermissions)
		if err != nil {
			return err
		}
	}

	return nil
}
