package m

import (
	"time"

	"github.com/vault-thirteen/auxie/errors"
)

type OpReportJson struct {
	FileName  string
	URL       string
	Error     TextError
	ErrorType OpReportErrorType
	Time      time.Time
}

func (orj OpReportJson) ToOpReport() (or OpReport) {
	or = OpReport{
		FileName:  orj.FileName,
		URL:       orj.URL,
		Error:     nil,
		ErrorType: orj.ErrorType,
		Time:      orj.Time,
	}

	if len(orj.Error.Text) != 0 {
		or.Error = errors.ErrorFromText(orj.Error.Text)
	}

	return or
}
