package service

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"html/template"
	"os"
)

type RenderInterface interface {
	Render(resp *restful.Response) error
}

type render struct {
	tpl *template.Template
}

func NewRender() RenderInterface {
	tpl := template.New("htmlTplEngine")
	_, err := tpl.ParseGlob("views/*.html")
	if err != nil {
		glog.Errorf("create tpl engine error, err=%+v", err)
		os.Exit(1)
	}
	return &render{
		tpl: tpl,
	}
}

func (r *render) Render(resp *restful.Response) error {
	err := r.tpl.ExecuteTemplate(
		resp.ResponseWriter,
		"index",
		nil,
	)
	if err != nil {
		glog.Errorf("render error, err=%+v", err)
		return err
	}
	return nil
}
