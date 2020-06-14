package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"
)

//目标定位是 tcp、http通用的中间件

const abortIndex int8 = math.MaxInt8 / 2 //最多63个中间件

type HandlerFunc func(*SliceRouterContext)

//router 结构体
type SliceRouter struct {
	groups []*SliceGroup
}

// group 结构体
type SliceGroup struct {
	*SliceRouter
	path     string
	handlers []HandlerFunc
}

//router 上下文
type SliceRouterContext struct {
	Rw  http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	*SliceGroup
	index int8
}

func newSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouterContext {
	newSliceGroup := &SliceGroup{}

	//最长url前缀匹配
	matchUrlLen := 0
	for _, group := range r.groups {
		if strings.HasPrefix(req.RequestURI, group.path) {
			pathLen := len(group.path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				*newSliceGroup = *group //浅拷贝数组指针
			}
		}
	}
	c := &SliceRouterContext{
		Rw:         rw,
		Req:        req,
		Ctx:        req.Context(),
		SliceGroup: newSliceGroup,
	}
	c.Reset()
	return c
}

func (c *SliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *SliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

type SliceRouterHandler struct {
	coreFunc func(*SliceRouterContext) http.Handler
	router   *SliceRouter
}

func (w *SliceRouterHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := newSliceRouterContext(rw, req, w.router)
	if w.coreFunc != nil {
		c.handlers = append(c.handlers, func(c *SliceRouterContext) {
			w.coreFunc(c).ServeHTTP(rw, req)
		})
		c.Reset()
		c.Next()
	}
}

func NewSliceRouterHandle(coreFunc func(*SliceRouterContext) http.Handler, router *SliceRouter) *SliceRouterHandler {
	return &SliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}

// 构造router
func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

//创建group
func (g *SliceRouter) Group(path string) *SliceGroup {
	return &SliceGroup{
		SliceRouter: g,
		path:        path,
	}
}

//构建回调方法
func (g *SliceGroup) Use(middlewares ...HandlerFunc) *SliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	exitsFlag := false
	for _, oldGroup := range g.SliceRouter.groups {
		if oldGroup == g {
			exitsFlag = true
		}
	}
	if !exitsFlag {
		g.SliceRouter.groups = append(g.SliceRouter.groups, g)
	}
	return g
}

//从最先加入中间件开始回调
func (c *SliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

//跳出中间件方法
func (c *SliceRouterContext) Abort() {
	c.index = abortIndex
}

//是否跳出了回调
func (c *SliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

// 重置回调
func (c *SliceRouterContext) Reset() {
	c.index = -1
}