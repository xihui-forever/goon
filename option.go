package goon

type option struct {
	// 对返回值的字段权限自动移除的 header
	PermHeader string `json:"perm_header,omitempty"`
}

func (o *option) applyOption(opts ...*option) *option {
	for _, opt := range opts {
		if opt.PermHeader != "" {
			o.PermHeader = opt.PermHeader
		}
	}

	return o
}

func WithPermHeader(header string) *App {
	return a.WithOptions(&option{PermHeader: header})
}
