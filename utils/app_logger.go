package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	UID_FIELD            = "uid"
	EVENT_FIELD          = "event"
	CORRELATION_ID_FIELD = "correlationID"
	SESSION_ID_FIELD     = "sessionID"
)

type AppLogger struct {
	entry         *log.Entry
	correlationID string
	sessionID     string
	uid           string
	event         string
}

func CreateAppLogger(correlationID string) *AppLogger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{})
	var loggerEntry *log.Entry

	if correlationID != "" {
		loggerEntry = logger.WithField(CORRELATION_ID_FIELD, correlationID)
	} else {
		loggerEntry = logger.WithField(CORRELATION_ID_FIELD, "-")
	}

	if GetenvAsBool("DEBUG_MODE", false) {
		logger.SetLevel(log.DebugLevel)
	}

	return &AppLogger{
		entry:         loggerEntry,
		correlationID: correlationID,
	}
}

func (logger *AppLogger) SetSessionID(sessionID string) {
	logger.sessionID = sessionID
	logger.entry = logger.entry.WithField(SESSION_ID_FIELD, sessionID)
}

func (logger *AppLogger) SetUID(uid interface{}) {
	logger.uid = fmt.Sprint(uid)
	logger.entry = logger.entry.WithField(UID_FIELD, fmt.Sprint(uid))
}

func (logger *AppLogger) SetEvent(event interface{}) {
	logger.event = fmt.Sprint(event)
	logger.entry = logger.entry.WithField(EVENT_FIELD, fmt.Sprint(event))
}

func (logger *AppLogger) GetSessionID() string {
	return logger.sessionID
}

func (logger *AppLogger) PanicOnError(err error) {
	if err != nil {
		logger.LogErrorWithStackTrace(err)
		panic(nil)
	}
}

func (logger *AppLogger) LogErrorWithStackTrace(r interface{}) {
	var err error
	switch r := r.(type) {
	case error:
		err = r
	default:
		err = fmt.Errorf("%v", r)
	}
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	logger.entry.WithField("error", err.Error()).Errorf("stack=%s", strings.ReplaceAll(string(buf), "\n", ""))
}

func (logger *AppLogger) LogErrorWithStackTraceAndNotify(r interface{}) {
	var err error
	switch r := r.(type) {
	case error:
		err = r
	default:
		err = fmt.Errorf("%v", r)
	}
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)

	logger.entry.WithField("error", err.Error()).Errorf("stack=%s", buf)

}

func (logger *AppLogger) RecoverWithStackTrace(wg *sync.WaitGroup) {
	if r := recover(); r != nil {
		logger.LogErrorWithStackTrace(r)
	}
	wg.Done()
}

func (logger *AppLogger) GetLogger() *log.Entry {
	return logger.entry
}

func (logger *AppLogger) GetCorrelationID() string {
	return logger.correlationID
}

func (logger *AppLogger) SendJsonResponse(w http.ResponseWriter, modelResponse interface{}) {
	jsonResponse, err := json.Marshal(modelResponse)
	logger.PanicOnError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	logger.PanicOnError(err)
}

func (logger *AppLogger) Debug(args ...interface{}) {
	var parsedArgs []interface{}
	for _, arg := range args {
		parsedArgs = append(parsedArgs, logger.maskSensitiveInfo(arg))
	}

	logger.GetLogger().Debug(parsedArgs...)
}

func (logger *AppLogger) Debugf(format string, args ...interface{}) {
	var parsedArgs []interface{}
	for _, arg := range args {
		parsedArgs = append(parsedArgs, logger.maskSensitiveInfo(arg))
	}

	logger.GetLogger().Debugf(format, parsedArgs...)
}

func (logger *AppLogger) Info(args ...interface{}) {
	var parsedArgs []interface{}
	for _, arg := range args {
		parsedArgs = append(parsedArgs, logger.maskSensitiveInfo(arg))
	}

	logger.GetLogger().Info(parsedArgs...)
}

func (logger *AppLogger) Infof(format string, args ...interface{}) {
	var parsedArgs []interface{}
	for _, arg := range args {
		parsedArgs = append(parsedArgs, logger.maskSensitiveInfo(arg))
	}

	logger.GetLogger().Infof(format, parsedArgs...)
}


func (logger *AppLogger) maskSensitiveInfo(info interface{}) interface{} {
	var parsedInfo interface{}
	switch info.(type) {
	case string:
		parsedInfo = MaskNric(info)
		parsedInfo = MaskEmail(parsedInfo)
		parsedInfo = MaskPhoneNumber(parsedInfo)
	case int:
		parsedInfo = MaskPhoneNumber(strconv.Itoa(info.(int)))
	default:
		parsedInfo = info
	}

	return parsedInfo
}

func MaskNric(content interface{}) interface{} {
	// https://regex101.com/r/cXoIOr/2
	var re = regexp.MustCompile(`(^|\b)[STGF]\d{4}(\d{3}[\w](\b|$))`)
	s := re.ReplaceAllString(content.(string), `$1*****$2`)

	return s
}

func RedactNric(content interface{}) interface{} {
	// https://regex101.com/r/foUK6R/1
	var re = regexp.MustCompile(`(^|\b)[STGF]\d{7}[\w](\b|$)`)
	s := re.ReplaceAllString(content.(string), `[REDACTED]`)

	return s
}

func MaskEmail(content interface{}) interface{} {
	// https://regex101.com/r/uyC2t3/3
	var re = regexp.MustCompile(`(\b|^)[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}(\b|$)`)
	s := re.ReplaceAllString(content.(string), `[REDACTED]`)

	return s
}

func MaskPhoneNumber(content interface{}) interface{} {
	// https://regex101.com/r/tx8Qyf/3
	var re = regexp.MustCompile(`(\b|^)([9|6|8](\d{7,11}|\d{3}\s?\d{4}))(\b|$)`)
	s := re.ReplaceAllString(content.(string), `[REDACTED]`)

	return s
}

func (logger *AppLogger) WithField(key string, val interface{}) *AppLogger {
	cloneLogger := Clone(logger).(*AppLogger)
	cloneLogger.entry = logger.entry.WithField(key, fmt.Sprint(val))
	return cloneLogger
}
