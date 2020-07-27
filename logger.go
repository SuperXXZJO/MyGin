package MyGin

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func DefaultLogger() HandlerFunc {
	res:= LoggerTocmd()
	return res
}

func Logger() *logrus.Logger {

	//实例化
	logger := logrus.New()
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	logger.Out = os.Stdout
	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

func LoggerTocmd() HandlerFunc {
	logger := Logger()
	return func(c *Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 请求host
		host := c.Host()
		// 状态码
		statusCode := c.Request.Header.Get("Status Code")

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			host,
			reqMethod,
			reqUri,
		)
	}
}

// 日志记录到文件
func LoggerToFile(filename string)HandlerFunc {

	src,err:=os.Create(filename)
	if err!=nil {
		fmt.Println("create file err:",err)
	}
	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Request.Header.Get("Status Code")

		// 请求host
		host := c.Host()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			host,
			reqMethod,
			reqUri,
		)
	}
}
