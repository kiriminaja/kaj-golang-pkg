package requester

import (
	"io"

	"github.com/imroc/req/v3"
)

type Config struct {
	Timeout int
	Debug   bool
}

type RequesterContract interface {
	RAW() *req.Request
	GET(url string, params map[string]string, headers map[string]string, result interface{}) (*req.Response, error)
	POST(url string, body interface{}, headers map[string]string, result interface{}) (*req.Response, error)
	PUT(url string, body interface{}, headers map[string]string, result interface{}) (*req.Response, error)
	DELETE(url string, params map[string]string, headers map[string]string, result interface{}) (*req.Response, error)
	Upload(url string, body, headers map[string]string,
		param, filename string, reader io.Reader, result interface{}) (*req.Response, error)
}
