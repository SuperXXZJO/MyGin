package MyGin

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request
	*Engine

	Path string
	Paths []string
	Method string
	Statuscode int
	index int
	Keys map[string]interface{}
	Handlers  []HandlerFunc
	Lock sync.RWMutex
}

func NewContext(w http.ResponseWriter,r *http.Request)*Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Paths:   SplitPath(r.URL.Path),
		Method:  r.Method,
		index:   -1,
		Keys:    nil,
	}
}

func(c *Context)Next()  {
	c.index++
	lens :=len(c.Handlers)
	for ;c.index <lens;c.index++{
		c.Handlers[c.index](c)
	}
}

func(c *Context)Abort()  {
	c.index = math.MaxInt8/2
}

func(c *Context)Set(key string,value interface{})  {
	c.Lock.Lock()
	if c.Keys== nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
	c.Lock.Unlock()
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.Lock.RLock()
	value, exists = c.Keys[key]
	c.Lock.RUnlock()
	return
}


//query

func(c *Context)Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func(c *Context)Querys()map[string][]string  {
	return c.Request.URL.Query()
}

//postform
func(c *Context)Postform(key string) string {
	return c.Request.FormValue(key)
}


//设置header
func (c *Context) SetHeader(key string,value string)  {
	c.Writer.Header().Set(key,value)
}

//状态码
func (c *Context)SetStatus(code int)  {
	c.Statuscode = code
	c.Writer.WriteHeader(code)
}


//string
func (c *Context)String(code int,str string)  {
	c.SetHeader("Content-Type","text/plain")
	c.SetStatus(code)
	b,_ :=json.Marshal(str)
	c.Writer.Write(b)
}

//json
func (c *Context)JSON(code int,values ...interface{})  {
	c.SetHeader("Content-Type","application/json")
	c.SetStatus(code)
	res,_:=json.Marshal(values)
	c.Writer.Write(res)
}


//xml
func(c *Context)XML(code int,obj interface{})  {
	c.SetStatus(code)
	header :=c.Writer.Header()
	header["Content-Type"]=[]string{"application/xml; charset=utf-8"}
	xml.NewEncoder(c.Writer).Encode(obj)
}


//html
func(c *Context)HTML(code int,html string)  {
	c.SetStatus(code)
	c.Set("Content-Type","text/html"+"; " + "charset=UTF-8")
	c.Writer.Write([]byte(html))
}



//file
func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}


//上传文件
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	f, fh, err := c.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}




//重定向
func (c *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return errors.New("invalid redirect status code")
	}
	c.Writer.Header().Set("Location", url)
	c.SetStatus(code)
	return nil
}

//cookie
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}


func(c *Context)ParseForm()(url.Values, error){
	if err:=c.Request.ParseForm();err!=nil{
		return nil, err
	}
	return c.Request.Form,nil
}

func(c *Context)Bind(obj interface{})  {
	if err:= Validate(obj);err!=nil{
		c.String(400,err.Error())
		c.Abort()
	}
	if err:= Bind(obj,c);err!=nil{
		c.String(400,err.Error())
		c.Abort()
	}
}

//host
func(c *Context)Host() string  {
	return c.Request.URL.Host
}

//copy
func(c *Context)Copy() *Context {
	cp:=&Context{
		Writer:     c.Writer,
		Request:    c.Request,
		Engine:     c.Engine,
		Path:       c.Path,
		Paths:      SplitPath(c.Path),
		Method:     c.Path,
		Statuscode: c.Statuscode,
		index:      c.index,
		Handlers:   nil,
		Keys:       map[string]interface{}{},
		Lock:       sync.RWMutex{},
	}
	return cp
}

//path
func(c *Context)FullPath()string  {
	return c.Request.URL.Path
}


func (c *Context) IsAborted() bool {
	return c.index >= math.MaxInt8/2
}

// GetBool
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

// MustGet
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}



func (c *Context) Html(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	if err := c.Engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.String(500,err.Error())
	}
}
