package main

import (
	"flag"
	"github.com/tommenx/demo/pkg/server"
)

var (
	port string
)

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&port, "port", "8080", "specify sever port")
}
func main() {
	svr := server.NewServer(port)
	svr.Run()
}
