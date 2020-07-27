package MyGin

import (
	"html/template"
	"net/http"
	"strings"
)

const (
	CONNECT = http.MethodConnect
	DELETE  = http.MethodDelete
	GET     = http.MethodGet
	HEAD    = http.MethodHead
	OPTIONS = http.MethodOptions
	PATCH   = http.MethodPatch
	POST    = http.MethodPost
	// PROPFIND = "PROPFIND"
	PUT   = http.MethodPut
	TRACE = http.MethodTrace
)



type HandlerFunc func(*Context)


type Engine struct {
	*RouterGroup
	Router *Router
	Groups []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

//新建引擎
func NewEngine()*Engine {
	r := NewRouter()
	engine := &Engine{
		Router: r,
	}
	engine.RouterGroup = &RouterGroup{
		Engine: engine,
	}
	engine.Groups = []*RouterGroup{engine.RouterGroup}
	return engine

}

//defult
func Default()*Engine {
	r := NewRouter()
	engine := &Engine{
		Router: r,
	}
	engine.RouterGroup = &RouterGroup{
		Engine: engine,
	}
	engine.Groups = []*RouterGroup{engine.RouterGroup}
	engine.Use(DefaultLogger(), Recover())
	return engine

}


//解析请求的路径，查找路由映射表
func (engine *Engine)ServeHTTP(w http.ResponseWriter,r *http.Request)  {
	if !strings.Contains(r.URL.Path,"/favicon.ico") {
		var middlewares []HandlerFunc
		for _, v := range engine.Groups {
			if strings.Contains(r.URL.Path, v.Profix) {
				middlewares = append(middlewares, v.Middleware...)
			}
		}

		c := NewContext(w, r)
		c.Engine = engine
		c.Handlers = middlewares
		engine.Router.Handle(c)
	}

}


func(engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr,engine)
}




func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}