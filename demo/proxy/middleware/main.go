package main

import "gateway_demo/proxy/middleware"

var addr = "127.0.0.1:2002"

func main() {
	reverseProxy := func(c *middleware.HandlerFunc) {}
}
