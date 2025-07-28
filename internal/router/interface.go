package router

type RouterInterface interface {
	SetupRouter()
	Serve(listenAddr string) error
}
