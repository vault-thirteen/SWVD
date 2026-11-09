package m

type Cookies struct {
	Cookies []Cookie `json:"cookies"`
}

func ReadCookiesFromFile(filePath string) (cs []Cookie, err error) {
	var x *Cookies
	x, err = ReadObjectFromDisk[Cookies](filePath)
	if err != nil {
		return nil, err
	}

	return x.Cookies, nil
}
