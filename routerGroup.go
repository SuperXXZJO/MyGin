package MyGin

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	*Engine
	Profix string
	Middleware []HandlerFunc
	father *RouterGroup
}

func(r *RouterGroup)Group(profix string)*RouterGroup {
	engine := r.Engine
	rg :=&RouterGroup{
		Engine:engine,
		Profix: r.Profix + profix,
		father: r,

	}
	engine.Groups = append(engine.Groups,rg)
	return rg
}



func(r *RouterGroup)addroute(method string,path string,handlers HandlerFunc)  {
	path = r.Profix + path
	r.Engine.Router.addRoute(method,path,handlers)
}

func(r *RouterGroup)Use(middleware ...HandlerFunc)  {
	r.Middleware = append(r.Middleware,middleware...)
}


//add
func(r *RouterGroup)ADD(method string,path string,handler HandlerFunc)  {
	r.addroute(method,path,handler)
}



//GET方法
func (r *RouterGroup)GET(path string,handler HandlerFunc){
	r.addroute(GET,path,handler)
}

//POST方法
func (r *RouterGroup)POST(path string,handler HandlerFunc){
	r.addroute(POST,path,handler)
}

//HEAD方法
func (r *RouterGroup)HEAD(path string,handler HandlerFunc){
	r.addroute(HEAD,path,handler)
}

//PUT方法
func (r *RouterGroup)PUT(path string,handler HandlerFunc){
	r.addroute(PUT,path,handler)
}

//DELETE方法
func (r *RouterGroup)DELETE(path string,handler HandlerFunc){
	r.addroute(DELETE,path,handler)
}

//CONNECT方法
func (r *RouterGroup)CONNECT(path string,handler HandlerFunc){
	r.addroute(CONNECT,path,handler)
}

//OPTIONS方法
func (r *RouterGroup)OPTIONS(path string,handler HandlerFunc){
	r.addroute(OPTIONS,path,handler)
}

//TRACE方法
func (r *RouterGroup)TRACE(path string,handler HandlerFunc){
	r.addroute(TRACE,path,handler)
}

//PATCH方法
func (r *RouterGroup)PATCH(path string,handler HandlerFunc){
	r.addroute(PATCH,path,handler)
}



func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.Profix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Paths[len(c.Paths)]
		if _, err := fs.Open(file); err != nil {

			he:= NewHTTPError(404)
			c.JSON(he.Code,he.Message)
			c.Abort()
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}


func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern:=path.Join(group.Profix,relativePath)
	group.GET(urlPattern, handler)
}