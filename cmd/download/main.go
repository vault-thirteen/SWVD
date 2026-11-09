package main

import (
	"fmt"
	"path/filepath"

	m "github.com/vault-thirteen/SWVD/models"
	s "github.com/vault-thirteen/SWVD/settings"
	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/file"
)

const (
	UsageHint = `Usage: download.exe <AuthorName> <CookiesFilePath> <DataFolder> <Quality>`
)

const (
	Err_AuthorFolderIsNotFound = "author folder is not found"
)

const (
	Msgf_SkippingExistingFile = "skipping existing file %s"
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
		IgnoreDownloadErrors: true,
		Quality:              args.quality,
	}

	var tsk = NewTask(taskSettings)

	var bestStreams []*m.BestStream
	bestStreams, err = tsk.collectAuthorsVideos()
	mustBeNoError(err)

	err = tsk.downloadBestStreams(bestStreams)
	mustBeNoError(err)
}
func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Task) collectAuthorsVideos() (bestStreams []*m.BestStream, err error) {
	var authorJsonFolderPath = m.MakeAuthorJsonFolder(t.s.DataFolder, t.s.AuthorName)
	var ok bool
	ok, err = file.FolderExists(authorJsonFolderPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf(Err_AuthorFolderIsNotFound)
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
		bestIdx, err = vss.FindBestStreamByQuality(t.s.Quality)
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
func (t *Task) downloadBestStreams(bestStreams []*m.BestStream) (err error) {
	var videoFolderPath = m.MakeVideoFolderPath(t.s.DataFolder)
	err = file.CreateFolderIfNotExists(videoFolderPath, t.s.NewFolderPermissions)
	if err != nil {
		return err
	}

	var authorVideoFolderPath = m.MakeAuthorVideoFolderPath(t.s.DataFolder, t.s.AuthorName)
	err = file.CreateFolderIfNotExists(authorVideoFolderPath, t.s.NewFolderPermissions)
	if err != nil {
		return err
	}

	// Fetch files.
	var fileExists bool
	var dstFileName string
	var dstFilePath string
	var buf []byte
	for _, bs := range bestStreams {
		dstFileName = m.FormatVideoFileName(bs.VideoId, bs.Stream.Resolution)
		dstFilePath = filepath.Join(authorVideoFolderPath, dstFileName)

		fileExists, err = file.FileExists(dstFilePath)
		if err != nil {
			return err
		}

		if fileExists {
			cio.PrintMessage(cio.FormatMessage(Msgf_SkippingExistingFile, dstFileName))
			continue
		}

		buf, err = m.FetchDataFromUrl(bs.Stream.Url, t.s)
		if err != nil {
			if t.s.IgnoreDownloadErrors {
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

	return nil
}
