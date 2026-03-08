package t39

//参数
type Request interface{}

//返回值
type Response interface{}

//过滤器接口
type Filter interface {
	Process(data Request) (Response, error)
}
