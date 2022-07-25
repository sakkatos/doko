package utils

import (
	"net/http"
	"time"
)

type ApiClient struct {
	Logger *AppLogger
}

type HttpClient interface {
	CallExternalAPI(req *http.Request) ([]byte, int, error)
	Get(baseUrl, relativePath string, queryParams map[string]string, headers map[string]string) (resp *http.Response, err error)
	Post(baseUrl, relativePath string, queryParams map[string]string, headers map[string]string, requestBody interface{}) (resp *http.Response, err error)
	GetLogger() *AppLogger
}

func (apiClient ApiClient) Post(baseUrl, relativePath string, queryParams map[string]string, headers map[string]string, requestBody interface{}) (*http.Response, error) {
	logger := apiClient.GetLogger().GetLogger()
	headers["X-Correlation-ID"] = apiClient.Logger.correlationID
	logger.Infof("POST %s%s", baseUrl, MaskNric(relativePath))

	resp, err := Post(baseUrl, relativePath, queryParams, headers, requestBody)
	if err != nil {
		logger.Infof("error POST %s%s", baseUrl, MaskNric(relativePath))
		return nil, err
	}

	logger.WithField("status", resp.StatusCode).Infof("done POST %s%s", baseUrl, MaskNric(relativePath))
	return resp, err
}

func (apiClient ApiClient) Get(baseUrl, relativePath string, queryParams map[string]string, headers map[string]string) (*http.Response, error) {
	logger := apiClient.GetLogger().GetLogger()
	headers["X-Correlation-ID"] = apiClient.Logger.correlationID
	logger.Infof("GET %s%s", baseUrl, MaskNric(relativePath))

	resp, err := Get(baseUrl, relativePath, queryParams, headers)
	if err != nil {
		logger.Infof("error GET %s%s", baseUrl, MaskNric(relativePath))
		return nil, err
	}

	logger.WithField("status", resp.StatusCode).Infof("done GET %s%s", baseUrl, MaskNric(relativePath))
	return resp, err
}

func (apiClient ApiClient) CallExternalAPI(req *http.Request) ([]byte, int, error) {
	logger := apiClient.GetLogger().GetLogger()
	logger.Debugf("calling external api %s", MaskNric(req.URL.Path))
	start := time.Now()

	response, status, err := CallExternalAPI(req)
	logger.
		WithField("status", status).
		WithField("duration", time.Since(start)).
		Infof("done calling external api %s", MaskNric(req.URL.Path))

	logger.Debug("External API Response:", string(response))
	logger.Debug("External API Error:", err)
	return response, status, err
}

func (apiClient ApiClient) GetLogger() *AppLogger {
	return apiClient.Logger
}
