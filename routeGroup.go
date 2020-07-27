package MyGin

type RouteGroup []Path_Mapping_Handlers

type Path_Mapping_Handlers struct {
	NowPath string
	Paths []string
	Handlers []HandlerFunc
}

//判断路径是否重复
func(p *Path_Mapping_Handlers)Match(paths []string,lens int) bool {
	count :=0
	for _,v1:=range p.Paths {
		for _,v2:=range paths{
			if v1 == v2 {
				count ++
			}
		}
	}
	if count == lens {
		return false
	}
	return true

}