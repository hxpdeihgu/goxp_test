package goxp

import (
	"fmt"
	"net/http"
	"reflect"
)

//上下文容器
type Context struct {
	*Container
	handlers      []Handler
	routerHandler Handler
	index         int
}

func (c *Context) Next() {
	c.index += 1
	c.Run()
}

func (c *Context) Run() {
	// 运行所以的handlers
	for c.index < len(c.handlers) {
		if val, err := c.Invoke(c.handlers[c.index]); err != nil {
			ReturnHandler(val, c)
		}
		c.index++
	}
	// 运行所以的handler
	if c.index == len(c.handlers) {
		c.Invoke(c.routerHandler)
	}
}

type RouterContext struct {
	*Container
	handlers []Handler
	index    int
}

type Contexter interface {
	GetType(interface{}) reflect.Type
	Get(reflect.Type) reflect.Value
}

func (rc *RouterContext) Run() {
	for i := 0; i < len(rc.handlers); i++ {
		if val, err := rc.Invoke(rc.handlers[i]); err == nil {
			ReturnHandler(val, rc)
		} else {
			fmt.Println(err)
		}
	}
}

func ReturnHandler(val []reflect.Value, c Contexter) {
	rv := c.Get(c.GetType((*http.ResponseWriter)(nil)))
	res := rv.Interface().(http.ResponseWriter)
	if len(val) == 1 {
		fmt.Fprintf(res, val[0].Interface().(string))
	} else if len(val) == 2 {
		// 设置的 status code
		res.WriteHeader(int(val[0].Int()))
		fmt.Fprintf(res, val[1].Interface().(string))
	}
}
