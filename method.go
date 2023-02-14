package goon

type Method string

const (
	Get  Method = "GET"
	Post Method = "POST"
	Head Method = "HEAD"

	Use     Method = "Use"     // 拦截器
	PreUse  Method = "PreUse"  // 前置拦截器
	PostUse Method = "PostUse" // 后置拦截器
)
