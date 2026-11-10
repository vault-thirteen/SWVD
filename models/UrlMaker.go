package m

import "strconv"

const (
	UrlParam_Page  = `page`
	UrlParam_Query = `query`
	UrlParam_V     = `v`

	UrlQueryPrefix    = "?"
	UrlQuerySeparator = "&"
	UrlEqualSign      = "="
	UrlSlash          = `/`
	UrlShemaPostfix   = `://`

	UrlPath_Watch  = "watch"
	UrlPath_Search = "search"
)

type UrlMaker struct {
	schema string
	domain string
}

func NewUrlMaker(schema string, domain string) (um *UrlMaker) {
	um = &UrlMaker{
		schema: schema,
		domain: domain,
	}

	return um
}

func (um *UrlMaker) MakeAuthorPageNUrl(authorName string, pageNumber int) string {
	if pageNumber < 1 {
		return ""
	}

	if pageNumber == 1 {
		return um.UrlSearch() + UrlQueryPrefix + UrlParam_Query + UrlEqualSign + authorName
	}

	return um.UrlSearch() + UrlQueryPrefix + UrlParam_Query + UrlEqualSign + authorName + UrlQuerySeparator + UrlParam_Page + UrlEqualSign + strconv.Itoa(pageNumber)
}
func (um *UrlMaker) MakeVideoPageUrl(videoId string) string {
	return um.UrlWatch() + UrlQueryPrefix + UrlParam_V + UrlEqualSign + videoId
}

func (um *UrlMaker) UrlWatch() string {
	return um.schema + UrlShemaPostfix + um.domain + UrlSlash + UrlPath_Watch
}
func (um *UrlMaker) UrlSearch() string {
	return um.schema + UrlShemaPostfix + um.domain + UrlSlash + UrlPath_Search
}
