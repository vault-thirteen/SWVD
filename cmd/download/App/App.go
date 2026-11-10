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
func (a *App) DownloadBestStreams(bestStreams []*m.BestStream) (report *m.DownloadReport, err error) {
	var videoFolderPath = m.MakeVideoFolderPath(a.ts.Common.DataFolder)
	err = file.CreateFolderIfNotExists(videoFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return nil, err
	}

	var authorVideoFolderPath = m.MakeAuthorVideoFolderPath(a.ts.Common.DataFolder, a.ts.Custom.AuthorName)
	err = file.CreateFolderIfNotExists(authorVideoFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return nil, err
	}

	// Fetch files.
	report = m.NewDownloadReport(a.ts.Custom.AuthorName)
	var fileExists bool
	var dstFileName string
	var dstFilePath string
	var buf []byte
	for _, bs := range bestStreams {
		dstFileName = m.FormatVideoFileName(bs.VideoId, bs.Stream.Resolution)
		dstFilePath = filepath.Join(authorVideoFolderPath, dstFileName)

		fileExists, err = file.FileExists(dstFilePath)
		if err != nil {
			return nil, err
		}

		if fileExists {
			cio.PrintMessage(cio.FormatMessage(m.Msgf_SkippingExistingFile, dstFileName))
			continue
		}

		buf, err = m.FetchDataFromUrl(bs.Stream.Url, a.ts)
		if err != nil {
			report.AddOpReport(dstFileName, bs.Stream.Url, err, m.OpReportErrorType_Download)
			continue
		}

		err = file.SaveBytesToFile(dstFilePath, buf)
		if err != nil {
			report.AddOpReport(dstFileName, bs.Stream.Url, err, m.OpReportErrorType_Storage)
			continue
		}

		report.AddOpReport(dstFileName, bs.Stream.Url, nil, m.OpReportErrorType_None)
	}

	report.SetTimeEnd()
	return report, nil
}
func (a *App) ProcessReport(report *m.DownloadReport) (err error) {
	report.Print()

	if report.FailedOpsCount == 0 {
		return nil
	}
	if !a.ts.Common.SaveFailedReport {
		return nil
	}

	var reportsFolderPath = a.ts.Common.ReportFolder
	err = file.CreateFolderIfNotExists(reportsFolderPath, a.ts.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	var fReport = report.CopyFailedOpsReport()
	var fileName = m.MakeReportFileName(fReport.AuthorName, fReport.EndTimeForFileName())
	var filePath = m.MakeReportFilePath(a.ts.Common.ReportFolder, fileName)
	err = fReport.SaveFailedOpsToFile(filePath, a.ts.Common.NewFilePermissions)
	if err != nil {
		return err
	}

	return nil
}
