package handler

type Method string

const (
	Get     Method = "GET"
	Post    Method = "POST"
	Head    Method = "HEAD"
	PreUse  Method = "PreUse"
	PostUse Method = "PostUse"
)
