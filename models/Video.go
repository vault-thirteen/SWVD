package m

import (
	"net/url"

	"github.com/vault-thirteen/auxie/errors"
)

const (
	Err_VideoIdIsNotFound = "video id is not found"
)

type Video struct {
	Title   string `json:"title"`
	PageUrl string `json:"pageUrl"`
	Id      string `json:"id"`
}

func (v *Video) ParseIdFromUrl() (err error) {
	var u *url.URL
	u, err = url.Parse(v.PageUrl)
	if err != nil {
		return err
	}

	v.Id = u.Query().Get(UrlParam_V)
	if len(v.Id) == 0 {
		return errors.ErrorFromText(Err_VideoIdIsNotFound)
	}

	return nil
}
