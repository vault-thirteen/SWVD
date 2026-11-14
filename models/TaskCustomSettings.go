package m

type TaskCustomSettings struct {
	// Content.
	AuthorName string
	Quality    QualityInt
	PagesCount int

	// Report Processing.
	DeleteReportAfterSuccess bool
}
