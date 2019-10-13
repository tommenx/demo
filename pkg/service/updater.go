package service

import (
	"github.com/golang/glog"
	"github.com/tommenx/demo/pkg/store"
)

type updater struct {
	db store.RedisClient
}

type UpdaterInterface interface {
	GetData(key string) (map[string][]string, error)
}

func NewExecutor(addr string) UpdaterInterface {
	db := store.NewRedis(addr)
	return &updater{
		db: db,
	}
}

func (e *updater) GetData(key string) (map[string][]string, error) {
	jobs, err := e.db.LRange(key)
	if err != nil {
		glog.Errorf("get job names error, err=%+v", err)
		return nil, err
	}
	data, err := e.db.ZRange(jobs)
	if err != nil {
		glog.Errorf("get data error, err=%+v", err)
		return nil, err
	}
	return data, nil
}
