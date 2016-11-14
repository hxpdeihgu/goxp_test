package goxp

import (
	"log"
	"net/http"
	"os"
)

type Handler interface{}

type Goxp struct {
	*Container
	// handler数组
	handlers      []Handler
    // 当前handler
	routerHandler Handler
	logger        *log.Logger
	IP            string
	Port          string
}

type MatureDuck struct {
	*Goxp
	Router
}

func Incubate() *MatureDuck {
	d := &Goxp{Container: New(), IP: "", Port: "3030",
		logger: log.New(os.Stdout, "logger开始", 0)}
	r := NewRouter()
	d.routerHandler = r.Handle
	d.SetMap(d.logger)
	d.Use(Static("public"))
	return &MatureDuck{d, r}
}

func (d *Goxp) Run() {
	d.logger.Println("listening on", d.IP+":"+d.Port)
	d.logger.Fatalln(http.ListenAndServe(d.IP+":"+d.Port, d))
}

func (d *Goxp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.logger.Println("start 开始")
	d.createContext(w, r).Run()
	d.logger.Println("end 结束")
}

func (d *Goxp) Use(handler Handler) {
	d.handlers = append(d.handlers, handler)
}

func (d *Goxp) createContext(w http.ResponseWriter, r *http.Request) *Context {
	d.SetMap(w)
	d.SetMap(r)
	c := &Context{Container: d.Container, routerHandler: d.routerHandler, handlers: d.handlers, index: 0}
	d.SetMap(c)
	return c
}
