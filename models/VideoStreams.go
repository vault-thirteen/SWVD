package m

import (
	"github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/number"
)

const (
	Err_QualityIsNotFound = "quality is not found"
)

const (
	IndexNone = -1
)

type VideoStreams struct {
	AuthorName string         `json:"authorName"`
	AuthorPage string         `json:"authorPage"`
	VideoId    string         `json:"videoId"`
	Streams    []*VideoStream `json:"streams"`
}

func ReadVideoStreamsFromFile(filePath string) (s *VideoStreams, err error) {
	return ReadObjectFromDisk[VideoStreams](filePath)
}

func (s VideoStreams) FindBestStreamByQuality(quality QualityInt) (idx int, err error) {
	// First attempt, e.g. 1080.
	idx, err = s.FindStreamByQuality(quality)
	if err == nil {
		return idx, nil
	}

	// Second attempt, e.g. 720.
	quality, err = s.LowerQuality(quality)
	if err != nil {
		return IndexNone, err
	}

	idx, err = s.FindStreamByQuality(quality)
	if err == nil {
		return idx, nil
	}

	// Third attempt, e.g. 480.
	quality, err = s.LowerQuality(quality)
	if err != nil {
		return IndexNone, err
	}

	idx, err = s.FindStreamByQuality(quality)
	if err == nil {
		return idx, nil
	}

	return IndexNone, errors.ErrorFromText(Err_QualityIsNotFound)
}
func (s VideoStreams) FindStreamByQuality(quality QualityInt) (idx int, err error) {
	var qualities = make([]QualityInt, 0, len(s.Streams))

	var q int
	for _, st := range s.Streams {
		q, err = number.ParseInt(st.Resolution)
		if err != nil {
			return IndexNone, err
		}

		qualities = append(qualities, QualityInt(q))
	}

	for i, qal := range qualities {
		if qal == quality {
			return i, nil
		}
	}

	return IndexNone, errors.ErrorFromText(Err_QualityIsNotFound)
}
func (s VideoStreams) LowerQuality(q QualityInt) (lq QualityInt, err error) {
	switch q {
	case QualityInt_1080p:
		return QualityInt_720p, nil
	case QualityInt_720p:
		return QualityInt_480p, nil
	default:
		return QualityInt_None, errors.ErrorFromText(Err_QualityIsNotFound)
	}
}
