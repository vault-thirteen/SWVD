package m

import "time"

type DownloadReportJson struct {
	AuthorName      string         `json:"authorName"`
	TimeStart       time.Time      `json:"timeStart"`
	TimeEnd         time.Time      `json:"timeEnd"`
	AllOpReports    []OpReportJson `json:"allOpReports,omitempty,omitzero"`
	FailedOpReports []OpReportJson `json:"failedOpReports"`
	FailedOpsCount  int            `json:"failedOpsCount"`
}

func (drj DownloadReportJson) ToDownloadReport() (dr DownloadReport) {
	dr = DownloadReport{
		AuthorName:      drj.AuthorName,
		TimeStart:       drj.TimeStart,
		TimeEnd:         drj.TimeEnd,
		AllOpReports:    nil,
		FailedOpReports: nil,
		FailedOpsCount:  drj.FailedOpsCount,
	}

	dr.AllOpReports = make([]OpReport, 0, len(drj.AllOpReports))
	for _, opReportJ := range drj.AllOpReports {
		dr.AllOpReports = append(dr.AllOpReports, opReportJ.ToOpReport())
	}

	dr.FailedOpReports = make([]OpReport, 0, len(drj.FailedOpReports))
	for _, opReportJ := range drj.FailedOpReports {
		dr.FailedOpReports = append(dr.FailedOpReports, opReportJ.ToOpReport())
	}

	return dr
}
