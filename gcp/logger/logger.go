package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func New() *log.Logger {
	return log.New(os.Stderr, "", 0)
}

type Severity int

const (
	Default   Severity = 0
	Debug     Severity = 100
	Info      Severity = 200
	Notice    Severity = 300
	Warning   Severity = 400
	Error     Severity = 500
	Critical  Severity = 600
	Alert     Severity = 700
	Emergency Severity = 800
)

var SeverityMap = map[Severity]string{
	Default:   "DEFAULT",
	Debug:     "DEBUG",
	Info:      "INFO",
	Notice:    "NOTICE",
	Warning:   "WARNING",
	Error:     "ERROR",
	Critical:  "CRITICAL",
	Alert:     "ALERT",
	Emergency: "EMERGENCY",
}

func (s Severity) String() string {
	return SeverityMap[s]
}

type CxEntry struct {
	Message   string    `json:"message"`
	Severity  Severity  `json:"severity,omitempty"`
	Trace     string    `json:"logging.googleapis.com/trace,omitempty"`
	Component string    `json:"component,omitempty"`
	Time      time.Time `json:"time,omitempty"`
}

func (e CxEntry) String() string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "    ")
	e.Time = time.Now()
	err := enc.Encode(e)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func CxEntryListenAndServe(addr string) *CxEntry {
	return &CxEntry{
		Severity:  Notice,
		Message:   fmt.Sprintf("ListenAndServe: ezcx listening and serving on %s", addr),
		Component: "ezcx.Server",
	}
}

func CxEntryContextDone() *CxEntry {
	return &CxEntry{
		Severity:  Notice,
		Message:   "ListenAndServe: ezcx server context is done",
		Component: "ezcx.Server",
	}
}

func CxEntryContextError(err error) *CxEntry {
	return &CxEntry{
		Severity:  Error,
		Message:   fmt.Sprintf("ListenAndServe: ezcx server context error: %s", err),
		Component: "ezcx.Server",
	}
}

func CxEntryServerError(err error) *CxEntry {
	return &CxEntry{
		Severity:  Critical,
		Message:   fmt.Sprintf("ListenAndServe: ezcx server processed a non-nil error: %s", err),
		Component: "ezcx.Server",
	}
}

func CxEntrySignalIntercepted(sig os.Signal) *CxEntry {
	return &CxEntry{
		Severity:  Notice,
		Message:   fmt.Sprintf("ListenAndServe: ezcx server intercepted the %s signal", sig),
		Component: "ezcx.Server",
	}
}

func CxEntryGracefulShutdown() *CxEntry {
	return &CxEntry{
		Severity:  Notice,
		Message:   "ListenAndServe: ezcx server was gracefully shutdown",
		Component: "ezcx.Server",
	}
}
