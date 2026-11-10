package main

import (
	a "github.com/vault-thirteen/SWVD/cmd/download/App"
	cla "github.com/vault-thirteen/SWVD/cmd/download/CommandLineArgs"
	m "github.com/vault-thirteen/SWVD/models"
	"github.com/vault-thirteen/auxie/cio"
)

func main() {
	args, err := cla.NewCommandLineArgs()
	if err != nil {
		cio.PrintMessage(cla.UsageHint)
	}
	m.MustBeNoError(err)

	var taskSettings *m.TaskSettings
	taskSettings, err = m.NewTaskSettingsFromFile(args.SettingsFilePath)
	m.MustBeNoError(err)
	var cts = m.TaskCustomSettings{
		AuthorName: args.AuthorName,
		Quality:    args.Quality,
	}
	taskSettings.Custom = cts

	var app = a.NewApp(taskSettings)

	var bestStreams []*m.BestStream
	bestStreams, err = app.CollectAuthorsVideos()
	m.MustBeNoError(err)

	err = app.DownloadBestStreams(bestStreams)
	m.MustBeNoError(err)
}
