package m

import (
	"strings"

	htmldom "github.com/vault-thirteen/auxie/HTML-DOM"
	"github.com/vault-thirteen/auxie/errors"
	"golang.org/x/net/html"
)

const (
	Err_DomNodeIsNotFound = "DOM node is not found"

	Errf_AttributeIsNotFound = "attribute is not found: %s"
)

func FindVideos(pageData []byte) (videos []*Video, err error) {
	videos = make([]*Video, 0)

	var domNode *html.Node
	domNode, err = html.Parse(strings.NewReader(string(pageData)))
	if err != nil {
		return nil, err
	}
	if domNode == nil {
		return nil, errors.ErrorFromText(htmldom.Err_DomNodeIsNotFound)
	}

	var node = htmldom.GetChildNodeByTag(domNode, htmldom.TagHtml) // -> <doctype>
	node = htmldom.GetSiblingNodeByTag(node, htmldom.TagHtml)      // -> <html>
	node = htmldom.GetChildNodeByTag(node, htmldom.TagHead)        // -> <head>
	node = htmldom.GetSiblingNodeByTag(node, htmldom.TagBody)      // -> <body>
	node = htmldom.GetChildNodeByTagAndClass(node, htmldom.TagDiv, "app-main-content-wrapper")
	node = htmldom.GetChildNodeByTag(node, htmldom.TagDiv)
	node = htmldom.GetSiblingNodeByTag(node, htmldom.TagDiv)
	node = htmldom.GetSiblingNodeByTag(node, htmldom.TagDiv)
	node = htmldom.GetChildNodeByTagAndClass(node, htmldom.TagDiv, "search-rows-wrapper search-rows-section-margin-top")
	node = htmldom.GetChildNodeByTagAndClass(node, htmldom.TagDiv, "content-padding-new")
	node = htmldom.GetChildNodeByTag(node, htmldom.TagDiv)
	node = htmldom.GetChildNodeByTagAndClass(node, htmldom.TagDiv, "horizontal-row")
	node = htmldom.GetChildNodeByTagAndClass(node, htmldom.TagDiv, "video-item-container")

	var video *Video
	video, err = FindSearchVideoDetails(node)
	if err != nil {
		return nil, err
	}
	videos = append(videos, video)

	// Next video.
	node = htmldom.GetSiblingNodeByTagAndClass(node, htmldom.TagDiv, "video-item-container")
	for node != nil {
		video, err = FindSearchVideoDetails(node)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)

		// Next.
		node = htmldom.GetSiblingNodeByTagAndClass(node, htmldom.TagDiv, "video-item-container")
	}

	return videos, nil
}
func FindSearchVideoDetails(node *html.Node) (video *Video, err error) {
	var ok bool
	video = &Video{}

	video.Title, ok = htmldom.GetNodeAttributeValue(node, htmldom.AttributeTitle)
	if !ok {
		return nil, errors.ErrorFromFormattedText(Errf_AttributeIsNotFound, htmldom.AttributeTitle)
	}

	x := htmldom.GetChildNodeByTagAndClass(node, htmldom.TagDiv, "horizontal-card")
	x = htmldom.GetChildNodeByTagAndClass(x, htmldom.TagA, "video-link")

	video.PageUrl, ok = htmldom.GetNodeAttributeValue(x, htmldom.AttributeHref)
	if !ok {
		return nil, errors.ErrorFromFormattedText(Errf_AttributeIsNotFound, htmldom.AttributeHref)
	}

	err = video.ParseIdFromUrl()
	if err != nil {
		return nil, err
	}

	return video, nil
}
func FindVideoStreams(pageData []byte) (videoStreams []*VideoStream, err error) {
	var domNode *html.Node
	domNode, err = html.Parse(strings.NewReader(string(pageData)))
	if err != nil {
		return nil, err
	}
	if domNode == nil {
		return nil, errors.ErrorFromText(Err_DomNodeIsNotFound)
	}

	var node = htmldom.GetChildNodeByTag(domNode, htmldom.TagHtml) // -> <doctype>
	node = htmldom.GetSiblingNodeByTag(node, htmldom.TagHtml)      // -> <html>
	node = htmldom.GetChildNodeByTag(node, htmldom.TagHead)        // -> <head>
	node = htmldom.GetSiblingNodeByTag(node, htmldom.TagBody)      // -> <body>

	var x *html.Node
	x = htmldom.GetElementById(node, "shareBtn-title")
	vTitle := htmldom.GetTextContent(x)

	x = htmldom.GetElementById(node, "player")
	sourceNodes := htmldom.GetChildNodes(x)
	videoStreams = make([]*VideoStream, 0, len(sourceNodes))
	var stream *VideoStream
	var ok bool
	for _, sourceNode := range sourceNodes {
		if htmldom.GetNodeTagName(sourceNode) != htmldom.TagSource {
			continue
		}

		stream = new(VideoStream)

		stream.Title = vTitle

		stream.Url, ok = htmldom.GetNodeAttributeValue(sourceNode, htmldom.AttributeSrc)
		if !ok {
			return nil, errors.ErrorFromFormattedText(Errf_AttributeIsNotFound, htmldom.AttributeSrc)
		}

		stream.Resolution, ok = htmldom.GetNodeAttributeValue(sourceNode, htmldom.AttributeSize)
		if !ok {
			return nil, errors.ErrorFromFormattedText(Errf_AttributeIsNotFound, htmldom.AttributeSize)
		}

		videoStreams = append(videoStreams, stream)
	}

	return videoStreams, nil
}
