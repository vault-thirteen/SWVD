package m

import "os"

type TaskSettings struct {
	UseCustomTransport   bool
	Schema               string
	Domain               string
	Cookies              []Cookie
	Header_UserAgent     string
	PrintUrl             bool
	DataFolder           string
	NewFolderPermissions os.FileMode
	NewFilePermissions   os.FileMode
	AuthorName           string
	IgnoreDownloadErrors bool
	Quality              QualityInt
	PagesCount           int
}
