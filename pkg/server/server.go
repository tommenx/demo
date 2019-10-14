package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/tommenx/demo/pkg/server/types"
	"github.com/tommenx/demo/pkg/service"
	"net/http"
	"sync"
)

var (
	errFailToRead   = restful.NewError(http.StatusBadRequest, "unable to read request body")
	errFailInternal = restful.NewError(http.StatusInternalServerError, "internal error")
	errFailToWrite  = restful.NewError(http.StatusInternalServerError, "unable to write response")
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
	u := service.NewFakeUpdater()
	svr := &server{
		r:    r,
		u:    u,
		port: port,
	}
	ws := new(restful.WebService)
	da := new(restful.WebService)
	da.Path("/data").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/render").To(svr.render))
	da.Route(ws.GET("/data").To(svr.update).
		Doc("update data").
		Operation("Update data").
		Writes(types.GetDataResponse{}))
	restful.Add(ws)
	restful.Add(da)
	return svr
}

func (svr *server) render(req *restful.Request, resp *restful.Response) {
	svr.lock.Lock()
	defer svr.lock.Unlock()
	_ = svr.r.Render(resp)
}

func (svr *server) update(req *restful.Request, resp *restful.Response) {
	svr.lock.Lock()
	defer svr.lock.Unlock()
	//args := &types.GetDataRequest{}
	//if err := req.ReadEntity(args); err != nil {
	//	errorResponse(resp, errFailToRead)
	//	return
	//}
	data, err := svr.u.GetData("jobs")
	if err != nil {
		errorResponse(resp, errFailInternal)
		return
	}
	info := make([]types.Data, 0)
	for name, logs := range data {
		info = append(info, types.Data{
			JobName: name,
			Values:  logs,
		})
	}
	getDataResp := &types.GetDataResponse{
		Data:   info,
		Status: "successs",
		Code:   0,
	}
	if err := resp.WriteEntity(getDataResp); err != nil {
		errorResponse(resp, errFailToWrite)
	}

}

func errorResponse(resp *restful.Response, svcErr restful.ServiceError) {
	glog.Error(svcErr.Message)
	if writeErr := resp.WriteServiceError(svcErr.Code, svcErr); writeErr != nil {
		glog.Errorf("unable to write error: %v", writeErr)
	}
}

func (svr *server) Run() {
	glog.Infof("start server, listen at :%s", svr.port)
	glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", svr.port), nil))
}
