package util

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	BaseFormatter logrus.Formatter
	ProjectRoot   string
}

// log message formatting
// 절대경로 말고 파일만 보이도록
// 이때, 함수명에서 패키지 경로도 지움
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry.Caller != nil {
		file := entry.Caller.File

		if f.ProjectRoot != "" && strings.Contains(file, f.ProjectRoot) {
			if idx := strings.Index(file, f.ProjectRoot); idx != -1 {
				projectRootLen := len(f.ProjectRoot)
				if idx+projectRootLen < len(file) {
					file = file[idx+projectRootLen:]
					if strings.HasPrefix(file, "/") {
						file = file[1:]
					}
				}
			}
		}

		if strings.Contains(file, "KWS_Control/") {
			if idx := strings.Index(file, "KWS_Control/"); idx != -1 {
				file = file[idx+len("KWS_Control/"):]
			}
		}

		funcName := entry.Caller.Function
		if strings.Contains(funcName, "github.com/easy-cloud-Knet/KWS_Control/") {
			funcName = strings.Replace(funcName, "github.com/easy-cloud-Knet/KWS_Control/", "", 1)
		}

		entry.Caller = &runtime.Frame{
			PC:       entry.Caller.PC,
			File:     file,
			Function: funcName,
			Line:     entry.Caller.Line,
		}
	}

	return f.BaseFormatter.Format(entry)
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)

	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := ""

	dir := filepath.Dir(currentFile)
	for {
		if strings.HasSuffix(dir, "KWS_Control") {
			projectRoot = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	formatter := &CustomFormatter{
		BaseFormatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "15:04:05",
		},
		ProjectRoot: projectRoot,
	}

	logger.SetFormatter(formatter)
	return logger
}

func GetLogger() *logrus.Logger {
	return NewLogger()
}
