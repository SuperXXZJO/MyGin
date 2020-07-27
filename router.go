package MyGin

import "fmt"

type Router struct {
	BigMap map[string]RouteGroup
}

func NewRouter()*Router {
	return &Router{BigMap: make(map[string]RouteGroup)}
}



func(r *Router)addRoute(method string,path string,handlers ...HandlerFunc)  {
	//判断path规范
	newpath := cleanPath(path)
	paths := SplitPath(newpath)
	lens := len(paths)
	if paths == nil {
		//todo 报错
		return
	}
	//判断有没有method
	if _,ok :=r.BigMap[method];!ok{
		r.BigMap[method] = make(RouteGroup,0)
	}
	//判断路径是否重复
	for _,v :=range r.BigMap[method] {
		ok :=v.Match(paths,lens)
		if !ok {
			panic(fmt.Errorf("路径 %s 重复",path))
		}

	}
	mod :=&Path_Mapping_Handlers{
		NowPath:  path,
		Paths:    paths,
		Handlers: handlers,
	}
	r.BigMap[method] = append(r.BigMap[method],*mod)

}

// getroute 查找路由
func(r *Router)getRoute(method string,path string) []HandlerFunc {

	for _,v :=range r.BigMap[method]{
		if v.NowPath == path{
			return v.Handlers
		}
	}
	return nil
}


// handle 实现handler
func(r *Router)Handle(c *Context)  {
	handles :=r.getRoute(c.Method,c.Path)
	if handles !=nil {
		c.Handlers = append(handles)
	}else {
		c.Handlers = append(c.Handlers, func(context *Context) {
			context.JSON(404,"404notfound...")
		})
	}
	c.Next()
}

