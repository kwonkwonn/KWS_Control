package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// 유저 관련 구조체랑 formatter 구조체는 분리해도 무방한데
// logger 구조체는 분리하면 해결하면 로직이 더 복잡해지는 역참조문제 발생해서
// 일단은 분리하지는 않을게요.

type CustomFormatter struct {
	BaseFormatter logrus.Formatter
	ProjectRoot   string
}

type Logger struct {
	systemLogger *logrus.Logger
	infoLogger   *logrus.Logger
	warnLogger   *logrus.Logger
	errorLogger  *logrus.Logger
}

// log message formatting
// 절대경로 말고 파일만 보이도록
// 이때, 함수명에서 패키지 경로도 지움
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b strings.Builder

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if f.BaseFormatter != nil {
		if textFormatter, ok := f.BaseFormatter.(*logrus.TextFormatter); ok && textFormatter.TimestampFormat != "" {
			timestamp = entry.Time.Format(textFormatter.TimestampFormat)
		}
	}

	level := strings.ToUpper(entry.Level.String())

	b.WriteString(timestamp)
	b.WriteString(" [")
	b.WriteString(level)
	b.WriteString("] ")
	b.WriteString(entry.Message)

	if entry.Caller != nil {
		file := entry.Caller.File

		if f.ProjectRoot != "" {
			file = strings.Replace(file, f.ProjectRoot+"/", "", 1)
			file = strings.Replace(file, f.ProjectRoot, "", 1)
		}

		// 저희 레포 이름 안바꾸겠죠..?
		// 하드코딩 해놓을게요
		file = strings.Replace(file, "KWS_Control/", "", 1)

		// funcName := entry.Caller.Function
		// funcName = strings.Replace(funcName, "github.com/easy-cloud-Knet/KWS_Control/", "", 1)

		b.WriteString(" [")
		b.WriteString(file)
		b.WriteString(":")
		b.WriteString(fmt.Sprintf("%d", entry.Caller.Line))
		b.WriteString("]")
	}

	b.WriteString("\n")
	return []byte(b.String()), nil
}

// 로그파일 뱉는 함수
func createLoggerWithFile(filename string, projectRoot string) *logrus.Logger {
	logger := logrus.New()
	// caller reporting을 비활성화하여 manual handling
	logger.SetReportCaller(false)

	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
	}

	logFile := filepath.Join(logDir, filename)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file %s: %v\n", logFile, err)
		logger.SetOutput(os.Stdout)
	} else {
		multiWriter := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(multiWriter)
	}

	formatter := &CustomFormatter{
		BaseFormatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
		ProjectRoot: projectRoot,
	}

	logger.SetFormatter(formatter)
	return logger
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	// caller reporting을 비활성화하여 manual handling
	logger.SetReportCaller(false)

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
			TimestampFormat: "2006-01-02 15:04:05",
		},
		ProjectRoot: projectRoot,
	}

	logger.SetFormatter(formatter)
	return logger
}

func NewEnhancedLogger() *Logger {
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

	return &Logger{
		systemLogger: NewLogger(),
		infoLogger:   createLoggerWithFile("info.log", projectRoot),
		warnLogger:   createLoggerWithFile("warn.log", projectRoot),
		errorLogger:  createLoggerWithFile("error.log", projectRoot),
	}
}

// 아래 info, warn error 모두 로그 출력하고
// 마지막 인자가 bool이면 save 파라미터로 들감 // 기본값 false
// printf 스타일 포매팅도 지원 (첫번째 인자가 format string이고 % 포함시)
func (l *Logger) Info(args ...interface{}) {
	message, save := parseLogArgs(args...)
	if save {
		l.infoLogger.Info(message)
	} else {
		l.systemLogger.Info(message)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	message, save := parseLogArgs(args...)
	if save {
		l.warnLogger.Warn(message)
	} else {
		l.systemLogger.Warn(message)
	}
}

func (l *Logger) Error(args ...interface{}) {
	message, save := parseLogArgs(args...)
	if save {
		l.errorLogger.Error(message)
	} else {
		l.systemLogger.Error(message)
	}
}

func (l *Logger) DebugInfo(format string, args ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		if strings.Contains(file, "KWS_Control/") {
			file = file[strings.Index(file, "KWS_Control/")+len("KWS_Control/"):]
		}

		newFormat := fmt.Sprintf("[%s:%d] %s", file, line, format)
		l.systemLogger.Infof(newFormat, args...)
	} else {
		l.systemLogger.Infof(format, args...)
	}
}

func (l *Logger) DebugWarn(format string, args ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		if strings.Contains(file, "KWS_Control/") {
			file = file[strings.Index(file, "KWS_Control/")+len("KWS_Control/"):]
		}

		newFormat := fmt.Sprintf("[%s:%d] %s", file, line, format)
		l.systemLogger.Warnf(newFormat, args...)
	} else {
		l.systemLogger.Warnf(format, args...)
	}
}

func (l *Logger) DebugError(format string, args ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		if strings.Contains(file, "KWS_Control/") {
			file = file[strings.Index(file, "KWS_Control/")+len("KWS_Control/"):]
		}

		newFormat := fmt.Sprintf("[%s:%d] %s", file, line, format)
		l.systemLogger.Errorf(newFormat, args...)
	} else {
		l.systemLogger.Errorf(format, args...)
	}
}

func (l *Logger) Println(args ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		if strings.Contains(file, "KWS_Control/") {
			file = file[strings.Index(file, "KWS_Control/")+len("KWS_Control/"):]
		}

		message := fmt.Sprint(args...)
		newMessage := fmt.Sprintf("[%s:%d] %s", file, line, message)
		l.systemLogger.Info(newMessage)
	} else {
		l.systemLogger.Info(args...)
	}
}

// 가변 인자를 파싱하여 메시지와 save 플래그를 반환
// printf 스타일 포매팅 지원: 첫번째 인자가 format string이면 fmt.Sprintf 사용
func parseLogArgs(args ...interface{}) (string, bool) {
	if len(args) == 0 {
		return "", false
	}

	save := false
	var formatArgs []interface{}
	var messageArgs []interface{}

	// bool 인자 분리
	filteredArgs := make([]interface{}, 0, len(args))
	for _, arg := range args {
		if v, ok := arg.(bool); ok {
			save = v
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	if len(filteredArgs) == 0 {
		return "", save
	}

	// printf 스타일 포매팅 감지
	if len(filteredArgs) > 1 {
		if formatStr, ok := filteredArgs[0].(string); ok && strings.Contains(formatStr, "%") {
			// printf 스타일로 처리
			formatArgs = filteredArgs[1:]
			message := fmt.Sprintf(formatStr, formatArgs...)
			return message, save
		}
	}

	messageArgs = filteredArgs
	var messageParts []string
	for _, arg := range messageArgs {
		messageParts = append(messageParts, fmt.Sprintf("%v", arg))
	}

	message := strings.Join(messageParts, " ")
	return message, save
}

// 현재 미사용중
func GetEnhancedLogger() *Logger {
	return NewEnhancedLogger()
}

func GetLogger() *Logger {
	return NewEnhancedLogger()
}
