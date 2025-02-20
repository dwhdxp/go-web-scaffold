package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	// 将日志格式设置为Json，并格式化时间戳
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetReportCaller(false)
}

// 记录调用API接口的日志信息
func Write(msg string, filename string) {
	setOutPutFile(logrus.InfoLevel, filename)
	logrus.Info(msg)
}

// logrus.Fields == map[string]interface{}
func Debug(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.DebugLevel, "debug")
	logrus.WithFields(fields).Debug(args)
}

func Info(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.InfoLevel, "info")
	logrus.WithFields(fields).Info(args)
}

func Warn(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.WarnLevel, "warn")
	logrus.WithFields(fields).Warn(args)
}

func Fatal(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.FatalLevel, "fatal")
	logrus.WithFields(fields).Fatal(args)
}

func Error(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.ErrorLevel, "error")
	logrus.WithFields(fields).Error(args)
}

func Panic(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.PanicLevel, "panic")
	logrus.WithFields(fields).Panic(args)
}

func Trace(fields logrus.Fields, args ...interface{}) {
	setOutPutFile(logrus.TraceLevel, "trace")
	logrus.WithFields(fields).Trace(args)
}

// 生成日志文件，格式：filename_time.log
func setOutPutFile(level logrus.Level, logName string) {
	// 1.检查输出日志文件夹是否存在
	if _, err := os.Stat("./runtime/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./runtime/log", 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s", "./runtime/log", err))
		}
	}

	// 2.生成日志文件
	timeStr := time.Now().Format("2006-01-02")
	fileName := path.Join("./runtime/log", logName+"_"+timeStr+".log")

	var err error
	os.Stderr, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("open log file err", err)
	}
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(level)
	return
}

// 以中间件形式调用
// 打印请求日志
func LoggerToFile() gin.LoggerConfig {
	// 1.检查日志文件夹是否存在
	if _, err := os.Stat("./runtime/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./runtime/log", 0777)
		if err != nil {
			panic(fmt.Errorf("create log dir '%s' error: %s", "./runtime/log", err))
		}
	}

	// 2.生成请求日志
	timeStr := time.Now().Format("2006-01-02")
	fileName := path.Join("./runtime/log", "success_"+timeStr+".log")

	os.Stderr, _ = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	var conf = gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - %s \"%s %s %s %d %s \"%s\" %s\"\n",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.ClientIP, // 请求IP
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency, // 响应时间
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		Output: io.MultiWriter(os.Stdout, os.Stderr),
	}

	return conf
}

// recover panic并打印错误日志
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 检查os.Stderr是否为nil
			if os.Stderr == nil {
				log.Printf("os.Stderr is nil, cannot write to stderr.")
			} else {
				fmt.Fprintf(os.Stderr, "panic occurred: %v\n", err)
				fmt.Fprintf(os.Stderr, "stacktrace from panic:\n%s\n", string(debug.Stack()))
			}

			// 确保日志目录存在
			logDir := "./runtime/log"
			if _, errDir := os.Stat(logDir); os.IsNotExist(errDir) {
				if errDir = os.MkdirAll(logDir, 0777); errDir != nil {
					fmt.Fprintf(os.Stderr, "create log dir '%s' error: %v\n", logDir, errDir)
					return
				}
			}

			timeStr := time.Now().Format("2006-01-02")
			fileName := path.Join("./runtime/log", "error_"+timeStr+".log")

			f, errFile := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if errFile != nil {
				fmt.Println(errFile)
				return
			}
			defer f.Close()

			timeFileStr := time.Now().Format("2006-01-02 15:04:05")
			// _, _ = f.WriteString("panic error time:" + timeFileStr + "\n")
			// _, _ = f.WriteString(fmt.Sprintf("%v", err) + "\n")
			// _, _ = f.WriteString("stacktrace from panic:" + string(debug.Stack()) + "\n")
			fmt.Fprintf(f, "panic error time:%s\n", timeFileStr)
			fmt.Fprintf(f, "%v\n", err)
			fmt.Fprintf(f, "stacktrace from panic:%s\n", string(debug.Stack()))

			// 不报错误，但返回500状态号
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  fmt.Sprintf("%v", err),
			})
			// 终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	c.Next()
}
