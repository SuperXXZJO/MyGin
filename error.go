package MyGin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	ErrUnsupportedMediaType        = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                   = NewHTTPError(http.StatusForbidden)
	ErrMethodNotAllowed            = NewHTTPError(http.StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrTooManyRequests             = NewHTTPError(http.StatusTooManyRequests)
	ErrBadRequest                  = NewHTTPError(http.StatusBadRequest)
	ErrBadGateway                  = NewHTTPError(http.StatusBadGateway)
	ErrInternalServerError         = NewHTTPError(http.StatusInternalServerError)
	ErrRequestTimeout              = NewHTTPError(http.StatusRequestTimeout)
	ErrServiceUnavailable          = NewHTTPError(http.StatusServiceUnavailable)
)

type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`

}

type errorString struct {
	s string
}

type errorInfo struct {
	Time     string `json:"time"`
	Alarm    string `json:"alarm"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Funcname string `json:"funcname"`
}

func (e *errorString) Error() string {
	return e.s
}

// Panic 异常
func Panic (text string) error {
	alarm("PANIC", text)
	return &errorString{text}
}


func  alarm(level string, str string) {
	// 当前时间
	currentTime := time.Now().String()

	// 定义 文件名、行号、方法名
	fileName, line, functionName := "?", 0 , "?"

	pc, fileName, line, ok := runtime.Caller(2)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}

	var msg = errorInfo{
		Time     : currentTime,
		Alarm    : level,
		Message  : str,
		Filename : fileName,
		Line     : line,
		Funcname : functionName,
	}

	jsons, errs := json.Marshal(msg)

	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}

	errorJsonInfo := string(jsons)

	fmt.Println(errorJsonInfo)


}



func Recover() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				err:= Panic(fmt.Sprintf("%s", r))
				fmt.Println(err)
			}
		}()
		c.Next()
	}
}


func  NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}
