package m

import "os"

type TaskCommonSettings struct {
	// Network.
	UseCustomTransport bool         `json:"useCustomTransport"`
	Schema             string       `json:"schema"`
	Domain             string       `json:"domain"`
	Cookies            []Cookie     `json:"cookies"`
	HttpHeaders        []HttpHeader `json:"httpHeaders"`

	// Storage.
	DataFolder           string      `json:"dataFolder"`
	NewFolderPermissions os.FileMode `json:"newFolderPermissions"`
	NewFilePermissions   os.FileMode `json:"newFilePermissions"`

	// Behaviour.
	PrintUrl             bool `json:"printUrl"`
	IgnoreDownloadErrors bool `json:"ignoreDownloadErrors"`
}
