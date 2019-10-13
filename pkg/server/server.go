package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/tommenx/demo/pkg/service"
	"net/http"
	"sync"
)

type server struct {
	u    service.UpdaterInterface
	r    service.RenderInterface
	lock sync.Mutex
	port string
}

type serverInterface interface {
	Run()
}

func NewServer(port string) serverInterface {
	r := service.NewRender()
	svr := &server{
		r:    r,
		port: port,
	}
	ws := new(restful.WebService)
	ws.Route(ws.GET("/").To(svr.render))
	restful.Add(ws)
	return svr
}

func (svr *server) render(req *restful.Request, resp *restful.Response) {
	svr.lock.Lock()
	defer svr.lock.Unlock()
	_ = svr.r.Render(resp)
}

func (svr *server) Run() {
	glog.Infof("start server, listen at :%s", svr.port)
	glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", svr.port), nil))
}
