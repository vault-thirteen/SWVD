package a

import (
	"path/filepath"

	m "github.com/vault-thirteen/SWVD/models"
	"github.com/vault-thirteen/auxie/cio"
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

func (a *App) CollectAuthorsVideos() (bestStreams []*m.BestStream, err error) {
	var authorJsonFolderPath = m.MakeAuthorJsonFolder(a.ts.Common.DataFolder, a.ts.Custom.AuthorName)
	var ok bool
	ok, err = file.FolderExists(authorJsonFolderPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.ErrorFromText(m.Err_AuthorFolderIsNotFound)
	}

	var fileNames []string
	fileNames, err = file.ListFileNames(authorJsonFolderPath)
	if err != nil {
		return nil, err
	}

	bestStreams = make([]*m.BestStream, 0, len(fileNames))

	var filePath string
	var vss *m.VideoStreams
	var bs *m.BestStream
	for _, fileName := range fileNames {
		filePath = filepath.Join(authorJsonFolderPath, fileName)

		vss, err = m.ReadVideoStreamsFromFile(filePath)
		if err != nil {
			return nil, err
		}

		var bestIdx int
		bestIdx, err = vss.FindBestStreamByQuality(a.ts.Custom.Quality)
		if err != nil {
			return nil, err
		}

		bs = &m.BestStream{
			VideoId: vss.VideoId,
			Stream:  vss.Streams[bestIdx],
		}
		bestStreams = append(bestStreams, bs)
	}

	return bestStreams, nil
}
func (a *App) DownloadBestStreams(bestStreams []*m.BestStream) (err error) {
	var videoFolderPath = m.MakeVideoFolderPath(a.ts.Common.DataFolder)
	err = file.CreateFolderIfNotExists(videoFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	var authorVideoFolderPath = m.MakeAuthorVideoFolderPath(a.ts.Common.DataFolder, a.ts.Custom.AuthorName)
	err = file.CreateFolderIfNotExists(authorVideoFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	// Fetch files.
	var fileExists bool
	var dstFileName string
	var dstFilePath string
	var buf []byte
	var failedUrls = make([]*m.FailedUrl, 0)
	var fu *m.FailedUrl
	for _, bs := range bestStreams {
		dstFileName = m.FormatVideoFileName(bs.VideoId, bs.Stream.Resolution)
		dstFilePath = filepath.Join(authorVideoFolderPath, dstFileName)

		fileExists, err = file.FileExists(dstFilePath)
		if err != nil {
			return err
		}

		if fileExists {
			cio.PrintMessage(cio.FormatMessage(m.Msgf_SkippingExistingFile, dstFileName))
			continue
		}

		buf, err = m.FetchDataFromUrl(bs.Stream.Url, a.ts)
		if err != nil {
			if a.ts.Common.IgnoreDownloadErrors {
				fu = m.NewFailedUrl(bs.Stream.Url, err)
				failedUrls = append(failedUrls, fu)
				errors.LogError(err)
				continue
			} else {
				return err
			}
		}

		err = file.SaveBytesToFile(dstFilePath, buf)
		if err != nil {
			return err
		}
	}

	if a.ts.Common.IgnoreDownloadErrors {
		// Visual separator.
		cio.PrintMessage()

		if len(failedUrls) > 0 {
			cio.PrintMessage(cio.FormatMessage(m.Msgf_FailedUrlsN, len(failedUrls)))
			for _, fu = range failedUrls {
				cio.PrintMessage(fu.Url)
			}
		} else {
			cio.PrintMessage(m.Msg_AllDownloadsSucceeded)
		}
	}

	return nil
}
