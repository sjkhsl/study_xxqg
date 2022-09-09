package utils

import (
	"net/http"

	"github.com/imroc/req/v3"
	log "github.com/sirupsen/logrus"
)

var client *req.Client

func init() {
	client = req.C()
	client.SetProxy(http.ProxyFromEnvironment)
	if log.GetLevel() == log.DebugLevel {
		client.DebugLog = true
		client = client.DevMode()
	}
	client.SetLogger(&myLog{})
	client.SetCommonHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
}

func GetClient() *req.Client {
	return client
}

type myLog struct {
}

func (m myLog) Errorf(format string, v ...interface{}) {
	log.Errorf(format, v)
}

func (m myLog) Warnf(format string, v ...interface{}) {
	log.Warnf(format, v)
}

func (m myLog) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v)
}

type LogWriter struct {
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	log.Debugln(string(p))
	return len(p), nil
}
