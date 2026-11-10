package m

import (
	"context"
	"io"
	"net"
	"net/http"

	utls "github.com/refraction-networking/utls"
	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"
)

var customClientHelloId = utls.HelloFirefox_Auto

func FetchDataFromUrl(url string, taskSettings *TaskSettings) (data []byte, err error) {
	var httpReq *http.Request
	httpReq, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range taskSettings.Common.Cookies {
		httpReq.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}

	// uTLS.
	httpTransport := &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (nc net.Conn, err error) {
			var tcpConn net.Conn
			tcpConn, err = (&net.Dialer{}).DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			host, _, _ := net.SplitHostPort(addr)
			config := &utls.Config{ServerName: host}

			var clientHelloID = customClientHelloId
			uConn := utls.UClient(tcpConn, config, clientHelloID)
			err = uConn.Handshake()
			if err != nil {
				_ = uConn.Close()
				return nil, err
			}
			return uConn, nil
		},
	}

	var httpClient *http.Client
	if taskSettings.Common.UseCustomTransport {
		httpClient = &http.Client{Transport: httpTransport}
	} else {
		httpClient = &http.Client{}
	}

	for _, hh := range taskSettings.Common.HttpHeaders {
		if len(hh.Name) == 0 {
			continue
		}

		httpReq.Header.Set(hh.Name, hh.Value)
	}

	if taskSettings.Common.PrintUrl {
		cio.PrintMessage(url)
	}

	var httpResp *http.Response
	httpResp, err = httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := httpResp.Body.Close()
		err = errors.Combine(err, derr)
	}()

	var httpRespBody []byte
	httpRespBody, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return httpRespBody, nil
}
