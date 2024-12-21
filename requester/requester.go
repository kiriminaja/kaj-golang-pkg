package requester

import (
	"io"
	"os"
	"time"

	"github.com/kiriminaja/kaj-golang-pkg/logger"

	"github.com/imroc/req/v3"
)

type reqsPkg struct {
	logField []logger.Field
	client   *req.Client
	cfg      *Config
}

func NewRequester(cfg *Config) RequesterContract {

	return &reqsPkg{
		client: buildClient(cfg),
		cfg:    cfg,
		logField: []logger.Field{
			logger.EventName("requester:log"),
		},
	}
}

func buildClient(cfg *Config) *req.Client {
	client := req.C().SetTimeout(time.Duration(cfg.Timeout) * time.Second)
	if cfg.Debug {
		client.EnableDump(&req.DumpOptions{
			Output:        os.Stdout,
			RequestHeader: true,
			ResponseBody:  true,
			RequestBody:   false,
		})
		client.DevMode()
	}
	client.SetUserAgent("Go-http-client/1.1")
	return client
}

func (r *reqsPkg) RAW() *req.Request {
	return r.client.R().
		SetHeader("Content-Type", "application/json")
}

func (r *reqsPkg) GET(url string, params map[string]string, headers map[string]string, result interface{}) (*req.Response, error) {
	client := r.client.R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetHeader("Content-Type", "application/json")

	r.logField = append(r.logField, logger.Any("url", url))
	r.logField = append(r.logField, logger.Any("headers", headers))
	r.logField = append(r.logField, logger.Any("params", params))
	r.logField = append(r.logField, logger.Any("method", "GET"))

	if r.cfg.Debug {
		client.EnableTrace()
	}
	response, err := client.SetSuccessResult(result).Get(url)
	if err != nil {
		logger.Error(logger.SetMessageFormat("Error GET request"), r.logField...)
		return nil, err
	}

	if response.IsErrorState() {
		logger.Error(logger.SetMessageFormat("Error State GET request"), r.logField...)
		return response, nil
	}
	r.logField = append(r.logField, logger.Any("duration", response.TotalTime()))
	logger.Info(logger.SetMessageFormat("Success GET request"), r.logField...)
	return response, nil
}

func (r *reqsPkg) POST(url string, body interface{}, headers map[string]string, result interface{}) (*req.Response, error) {
	client := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).SetBody(body)

	r.logField = append(r.logField, logger.Any("url", url))
	r.logField = append(r.logField, logger.Any("headers", headers))
	r.logField = append(r.logField, logger.Any("body", body))
	r.logField = append(r.logField, logger.Any("method", "POST"))

	if r.cfg.Debug {
		client.EnableTrace()
	}
	response, err := client.SetSuccessResult(result).Post(url)
	if err != nil {
		logger.Error(logger.SetMessageFormat("Error POST request"), r.logField...)
		return nil, err
	}

	if response.IsErrorState() {
		logger.Error(logger.SetMessageFormat("Error State POST request"), r.logField...)
		return response, nil
	}
	r.logField = append(r.logField, logger.Any("duration", response.TotalTime()))
	logger.Info(logger.SetMessageFormat("Success POST request"), r.logField...)
	return response, nil
}

func (r *reqsPkg) PUT(url string, body interface{}, headers map[string]string, result interface{}) (*req.Response, error) {
	client := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).SetBody(body)

	r.logField = append(r.logField, logger.Any("url", url))
	r.logField = append(r.logField, logger.Any("headers", headers))
	r.logField = append(r.logField, logger.Any("body", body))
	r.logField = append(r.logField, logger.Any("method", "PUT"))

	if r.cfg.Debug {
		client.EnableTrace()
	}
	response, err := client.SetSuccessResult(result).Put(url)
	if err != nil {
		logger.Error(logger.SetMessageFormat("Error PUT request"), r.logField...)
		return nil, err
	}

	if response.IsErrorState() {
		logger.Error(logger.SetMessageFormat("Error State PUT request"), r.logField...)
		return response, nil
	}
	r.logField = append(r.logField, logger.Any("duration", response.TotalTime()))
	logger.Info(logger.SetMessageFormat("Success PUT request"), r.logField...)
	return response, nil
}

func (r *reqsPkg) DELETE(url string, params map[string]string, headers map[string]string, result interface{}) (*req.Response, error) {
	client := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).SetQueryParams(params)

	r.logField = append(r.logField, logger.Any("url", url))
	r.logField = append(r.logField, logger.Any("headers", headers))
	r.logField = append(r.logField, logger.Any("params", params))
	r.logField = append(r.logField, logger.Any("method", "DELETE"))

	if r.cfg.Debug {
		client.EnableTrace()
	}
	response, err := client.SetSuccessResult(result).Delete(url)
	if err != nil {
		logger.Error(logger.SetMessageFormat("Error Delete request"), r.logField...)
		return nil, err
	}

	if response.IsErrorState() {
		logger.Error(logger.SetMessageFormat("Error State Delete request"), r.logField...)
		return response, nil
	}
	r.logField = append(r.logField, logger.Any("duration", response.TotalTime()))
	logger.Info(logger.SetMessageFormat("Success Delete request"), r.logField...)
	return response, nil
}

func (r *reqsPkg) Upload(url string, body, headers map[string]string,
	param, filename string, reader io.Reader, result interface{}) (*req.Response, error) {
	client := r.client.R().
		SetHeaders(headers).SetBody(body)

	r.logField = append(r.logField, logger.Any("url", url))
	r.logField = append(r.logField, logger.Any("body", body))
	r.logField = append(r.logField, logger.Any("method", "POST"))

	if r.cfg.Debug {
		client.EnableTrace()
	}
	response, err := client.SetSuccessResult(result).SetFormData(body).SetFileReader(param, filename, reader).Post(url)
	if err != nil {
		logger.Error(logger.SetMessageFormat("Error POST request"), r.logField...)
		return nil, err
	}

	if response.IsErrorState() {
		logger.Error(logger.SetMessageFormat("Error State POST request"), r.logField...)
		return response, nil
	}
	r.logField = append(r.logField, logger.Any("duration", response.TotalTime()))
	logger.Info(logger.SetMessageFormat("Success POST request"), r.logField...)
	return response, nil
}
