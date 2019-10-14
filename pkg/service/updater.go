package service

import (
	"github.com/golang/glog"
	"github.com/tommenx/demo/pkg/store"
	"math/rand"
	"strconv"
	"time"
)

type updater struct {
	db store.RedisClient
}

type UpdaterInterface interface {
	GetData(key string) (map[string][]string, error)
}

func NewUpdater(addr string) UpdaterInterface {
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

type fakeUpdater struct {
	job  []string
	data map[string][]string
}

func NewFakeUpdater() UpdaterInterface {
	job := []string{"a", "b", "c"}
	data := make(map[string][]string)
	return &fakeUpdater{
		job:  job,
		data: data,
	}
}

func (f *fakeUpdater) GetData(key string) (map[string][]string, error) {
	rand.Seed(time.Now().UnixNano())
	for _, k := range f.job {
		num := rand.Intn(1000)
		numStr := strconv.Itoa(num)
		f.data[k] = append(f.data[k], numStr)
	}
	return f.data, nil
}
