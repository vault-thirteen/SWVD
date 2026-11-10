package main

import (
	a "github.com/vault-thirteen/SWVD/cmd/list/App"
	cla "github.com/vault-thirteen/SWVD/cmd/list/CommandLineArgs"
	m "github.com/vault-thirteen/SWVD/models"
	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
	"github.com/vault-thirteen/auxie/cio"
)

const AppName = "Tool"

func main() {
	showIntro()

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
		PagesCount: args.PagesCount,
	}
	taskSettings.Custom = cts

	var app = a.NewApp(taskSettings)

	err = app.ListAuthorsVideosN()
	m.MustBeNoError(err)

	err = app.SaveAuthorVideoStreamsFromList()
	m.MustBeNoError(err)
}
func showIntro() {
	versioneer, err := ver.New(false)
	m.MustBeNoError(err)
	versioneer.ShowIntroText(AppName)
	versioneer.ShowComponentsInfoText()
	cio.PrintMessage()
}
