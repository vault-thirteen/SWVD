package m

type FailedUrl struct {
	Url string
	Err error
}

func NewFailedUrl(url string, err error) (fu *FailedUrl) {
	fu = &FailedUrl{
		Url: url,
		Err: err,
	}

	return fu
}
