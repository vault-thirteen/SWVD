package main

import (
	"os"
	"path/filepath"

	cla "github.com/vault-thirteen/SWVD/cmd/dwf/CommandLineArgs"
	m "github.com/vault-thirteen/SWVD/models"
	pw "github.com/vault-thirteen/SWVD/models/playwright"
	ver "github.com/vault-thirteen/auxie/Versioneer/classes/Versioneer"
	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/file"
)

const AppName = "Tool"

func main() {
	showIntro()

	args, err := cla.NewCommandLineArgs()
	if err != nil {
		cio.PrintMessage(cla.UsageHint)
	}
	m.MustBeNoError(err)

	var settings *m.TaskSettings
	settings, err = m.NewTaskSettingsFromFile(args.SettingsFilePath)
	m.MustBeNoError(err)
	var cts = m.TaskCustomSettings{
		DeleteReportAfterSuccess: args.DeleteReportAfterSuccess,
	}
	settings.Custom = cts

	err = processReports(settings)
	m.MustBeNoError(err)
}
func processReports(settings *m.TaskSettings) (err error) {
	var reportFileNames []string
	reportFileNames, err = file.ListFileNames(settings.Common.ReportFolder)
	if err != nil {
		return err
	}

	for _, reportFileName := range reportFileNames {
		err = processReport(settings, reportFileName)
		if err != nil {
			return err
		}
	}

	return nil
}
func processReport(settings *m.TaskSettings, reportFileName string) (err error) {
	var dr *m.DownloadReport
	var reportFilePath = filepath.Join(settings.Common.ReportFolder, reportFileName)
	dr, err = m.NewDownloadReportFromFile(reportFilePath)
	if err != nil {
		return err
	}

	err = downloadFailedFiles(settings, dr)
	if err != nil {
		return err
	}

	if settings.Custom.DeleteReportAfterSuccess {
		err = os.Remove(reportFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}
func downloadFailedFiles(settings *m.TaskSettings, report *m.DownloadReport) (err error) {
	var videoFolderPath = m.MakeVideoFolderPath()
	err = file.CreateFolderIfNotExists(videoFolderPath, settings.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	var authorVideoFolderPath = m.MakeAuthorVideoFolderPath(report.AuthorName)
	err = file.CreateFolderIfNotExists(authorVideoFolderPath, settings.Common.NewFolderPermissions)
	if err != nil {
		return err
	}

	var dstFilePath string
	for _, rep := range report.FailedOpReports {
		dstFilePath = filepath.Join(authorVideoFolderPath, rep.FileName)
		cio.PrintMessage(rep.URL)

		err = pw.DownloadFile(rep.URL, dstFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}
func showIntro() {
	versioneer, err := ver.New(false)
	m.MustBeNoError(err)
	versioneer.ShowIntroText(AppName)
	versioneer.ShowComponentsInfoText()
	cio.PrintMessage()
}
