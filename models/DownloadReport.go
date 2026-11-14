package m

import (
	"os"
	"time"

	"github.com/vault-thirteen/auxie/cio"
)

const (
	OpTimeFormat            = "2006-01-02 15:04:05"
	OpTimeFormatForFileName = "20060102_150405"
)

type DownloadReport struct {
	AuthorName      string     `json:"authorName"`
	TimeStart       time.Time  `json:"timeStart"`
	TimeEnd         time.Time  `json:"timeEnd"`
	AllOpReports    []OpReport `json:"allOpReports,omitempty,omitzero"`
	FailedOpReports []OpReport `json:"failedOpReports"`
	FailedOpsCount  int        `json:"failedOpsCount"`
}

func NewDownloadReport(authorName string) (dr *DownloadReport) {
	dr = &DownloadReport{
		AuthorName:      authorName,
		TimeStart:       time.Now(),
		AllOpReports:    make([]OpReport, 0),
		FailedOpReports: make([]OpReport, 0),
		FailedOpsCount:  0,
	}

	return dr
}
func NewDownloadReportFromFile(filePath string) (dr *DownloadReport, err error) {
	var drj *DownloadReportJson
	drj, err = ReadObjectFromDisk[DownloadReportJson](filePath)
	if err != nil {
		return nil, err
	}

	out := drj.ToDownloadReport()

	return &out, nil
}

func (dr DownloadReport) CopyFailedOpsReport() (out DownloadReport) {
	out = DownloadReport{
		AuthorName:      dr.AuthorName,
		TimeStart:       dr.TimeStart,
		TimeEnd:         dr.TimeEnd,
		AllOpReports:    nil,
		FailedOpReports: make([]OpReport, 0, len(dr.FailedOpReports)),
		FailedOpsCount:  dr.FailedOpsCount,
	}

	var frWTE OpReport
	for _, fr := range dr.FailedOpReports {
		frWTE = fr.CopyWithTextError()
		out.FailedOpReports = append(out.FailedOpReports, frWTE)
	}

	return out
}
func (dr *DownloadReport) AddOpReport(fileName string, url string, err error, errorType OpReportErrorType) {
	if len(fileName) == 0 {
		return
	}
	if len(url) == 0 {
		return
	}

	var or = NewOpReport(fileName, url, err, errorType)
	dr.AllOpReports = append(dr.AllOpReports, or)

	if or.Error != nil {
		dr.FailedOpReports = append(dr.FailedOpReports, or)
		dr.FailedOpsCount += 1
	}
}
func (dr *DownloadReport) SetTimeEnd() {
	dr.TimeEnd = time.Now()
}
func (dr *DownloadReport) Print() {
	if len(dr.AllOpReports) == 0 {
		cio.PrintMessage(Msg_NoDownloadsWereNeeded)
		return
	}

	if dr.FailedOpsCount == 0 {
		cio.PrintMessage(Msg_AllDownloadsSucceeded)
		return
	}

	cio.PrintMessage(cio.FormatMessage(Msgf_FailedUrlsN, dr.FailedOpsCount))
	for _, fop := range dr.FailedOpReports {
		txt := cio.FormatMessage(Msgf_FailedOpText, fop.URL, fop.ErrorType.ToString(), fop.TimeString())
		cio.PrintMessage(txt)
	}
	return
}
func (dr *DownloadReport) SaveFailedOpsToFile(filePath string, filePermissions os.FileMode) (err error) {
	if dr.FailedOpsCount == 0 {
		return nil
	}

	return SaveObjectOnDisk(dr, filePath, filePermissions)
}
func (dr *DownloadReport) EndTimeForFileName() string {
	return dr.TimeEnd.UTC().Format(OpTimeFormatForFileName)
}
