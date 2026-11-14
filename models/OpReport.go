package m

import (
	"time"
)

type OpReport struct {
	FileName  string
	URL       string
	Error     error
	ErrorType OpReportErrorType
	Time      time.Time
}

func NewOpReport(fileName string, url string, err error, errorType OpReportErrorType) OpReport {
	return OpReport{
		FileName:  fileName,
		URL:       url,
		Error:     err,
		ErrorType: errorType,
		Time:      time.Now(),
	}
}

func (op OpReport) TimeString() string {
	return op.Time.UTC().Format(OpTimeFormat)
}
func (op OpReport) CopyWithTextError() (out OpReport) {
	if op.Error == nil {
		return
	}

	out = OpReport{
		FileName:  op.FileName,
		URL:       op.URL,
		Error:     NewTextError(op.Error.Error()),
		ErrorType: op.ErrorType,
		Time:      op.Time,
	}

	return out
}
