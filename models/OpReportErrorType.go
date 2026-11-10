package m

type OpReportErrorType byte

const (
	OpReportErrorType_None     OpReportErrorType = 0
	OpReportErrorType_Download OpReportErrorType = 1
	OpReportErrorType_Storage  OpReportErrorType = 2
)

func (et OpReportErrorType) ToString() string {
	switch et {
	case OpReportErrorType_Download:
		return "network"
	case OpReportErrorType_Storage:
		return "storage"
	default:
		return "none"
	}
}
