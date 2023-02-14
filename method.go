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

func (p Method) String() string {
	return string(p)
}

func (p Method) Equal(m Method) bool {
	return p == m
}

func (p Method) GoString() string {
	return p.String()
}

func (p Method) MarshalJSON() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Method) UnmarshalJSON(data []byte) error {
	*p = Method(string(data))
	return nil
}
