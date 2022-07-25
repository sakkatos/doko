package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

type (
	ContextKey string
)

func ContainsString(list []string, item string) bool {
	for _, curr := range list {
		if curr == item {
			return true
		}
	}
	return false
}

func PanicOnError(err error) {
	if err != nil {
		LogErrorWithStackTrace(err)
		panic(err)
	}
}

func RecoverWithStackTrace(wg *sync.WaitGroup) {
	if r := recover(); r != nil {
		LogErrorWithStackTrace(r)
	}
	wg.Done()
}

func LogErrorWithStackTrace(r interface{}) {
	var err error
	switch r := r.(type) {
	case error:
		err = r
	default:
		err = fmt.Errorf("%v", r)
	}
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	log.Errorf("error: %v, stack trace: %s", err.Error(), buf)
}

// To remove once new version of funk is released
func GetDiff(x interface{}, y interface{}, f func(interface{}) interface{}) interface{} {
	if !funk.IsCollection(x) {
		panic("First parameter must be a collection")
	}
	if !funk.IsCollection(y) {
		panic("Second parameter must be a collection")
	}

	hash := map[interface{}]interface{}{}

	xValue := reflect.ValueOf(x)
	xType := xValue.Type()

	yValue := reflect.ValueOf(y)
	yType := yValue.Type()

	if funk.NotEqual(xType, yType) {
		panic("Parameters must have the same type")
	}

	zType := reflect.SliceOf(xType.Elem())
	zSlice := reflect.MakeSlice(zType, 0, 0)

	for i := 0; i < xValue.Len(); i++ {
		v := xValue.Index(i).Interface()
		hash[f(v)] = v
	}

	for i := 0; i < yValue.Len(); i++ {
		v := yValue.Index(i).Interface()
		_, ok := hash[f(v)]
		if ok {
			delete(hash, f(v))
		}
	}

	for _, v := range hash {
		kValue := reflect.ValueOf(v)
		zSlice = reflect.Append(zSlice, kValue)
	}

	return zSlice.Interface()
}

func sendHTTPRequest(requestMethod, baseUrl, relativePath string, queryParams map[string]string, headers map[string]string, requestBody interface{}) (resp *http.Response, err error) {
	requestUrl := constructRequestUrl(baseUrl, relativePath, queryParams)
	requestByte, _ := json.Marshal(requestBody)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest(requestMethod, requestUrl, requestReader)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func Get(baseUrl, relativePath string, queryParams map[string]string, headers map[string]string) (resp *http.Response, err error) {
	return sendHTTPRequest("GET", baseUrl, relativePath, queryParams, headers, map[string]string{})
}

func Post(baseUrl, relativePath string, queryParams map[string]string, headers map[string]string, requestBody interface{}) (resp *http.Response, err error) {
	return sendHTTPRequest("POST", baseUrl, relativePath, queryParams, headers, requestBody)
}

func constructRequestUrl(server, relativePath string, queryParams map[string]string) string {
	relativeUrl, _ := url.Parse(relativePath)
	queryString := relativeUrl.Query()
	for k, v := range queryParams {
		queryString.Set(k, v)
	}
	relativeUrl.RawQuery = queryString.Encode()

	serverUrl, err := url.Parse(server)
	PanicOnError(err)

	return serverUrl.ResolveReference(relativeUrl).String()
}

func SendJsonResponse(w http.ResponseWriter, modelResponse interface{}) {
	jsonResponse, err := json.Marshal(modelResponse)
	PanicOnError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	PanicOnError(err)
}