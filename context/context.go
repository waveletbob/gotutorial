package context

import (
	"context"
	"net/http"
	"net/url"

	"github.com/waveletbob/gotutorial/utils"
)

var Key interface{} = "$.gentleman"

type Store map[interface{}]interface{}

type Context struct {
	Error    error
	Stopped  bool
	Parent   *Context
	Client   *http.Client
	Request  *http.Request
	Response *http.Response
}

func New() *Context {
	req := createRequest()
	res := createResponse(req)
	cli := &http.Client{Transport: http.DefaultTransport}
	return &Context{Request: req, Response: res, Client: cli}
}
func (c *Context) UseParent(ctx *Context) {
	c.Parent = ctx
}
func createRequest() *http.Request {
	// Create HTTP request
	req := &http.Request{
		Method:     "GET",
		URL:        &url.URL{},
		Host:       "",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Proto:      "HTTP/1.1",
		Header:     make(http.Header),
		Body:       utils.NopCloser(),
	}
	// Return shallow copy of Request with the new context
	return req.WithContext(emptyContext())
}
func (c *Context) SetCancelContext(ctx context.Context) *Context {
	golRequestContext := context.WithValue(ctx, Key, c.Value(Key))
	c.Request = c.Request.WithContext(golRequestContext)
	return c
}
func (c *Context) Value(key interface{}) interface{} {
	return c.Request.Context().Value(key)
}
func (c *Context) Clone() *Context {
	ctx := new(Context)
	*ctx = *c

	req := new(http.Request)
	*req = *c.Request
	ctx.Request = req
	c.CopyTo(ctx)

	res := new(http.Response)
	*res = *c.Response
	ctx.Response = res

	return ctx
}
func (c *Context) CopyTo(newCtx *Context) {
	store := Store{}

	for key, value := range c.getStore() {
		store[key] = value
	}

	ctx := context.WithValue(context.Background(), Key, store)
	newCtx.Request = newCtx.Request.WithContext(ctx)
}
func (c *Context) getStore() Store {
	store, ok := c.Request.Context().Value(Key).(Store)
	if !ok {
		panic("invalid request context")
	}
	return store
}

// createResponse creates a default http.Response instance.
func createResponse(req *http.Request) *http.Response {
	return &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
		Proto:      "HTTP/1.1",
		Request:    req,
		Header:     make(http.Header),
		Body:       utils.NopCloser(),
	}
}

func emptyContext() context.Context {
	return context.WithValue(context.Background(), Key, Store{})
}

func (c *Context) Set(key interface{}, value interface{}) {
	store := c.getStore()
	store[key] = value
}
