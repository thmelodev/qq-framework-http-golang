package http

type IHttpProvider interface {
	GetHttpPort() int
	GetLimitePaginacao() int
	GetLimiteRotinas() int
}