package server

import "net/http"

type Server interface {
	Start() error
	Stop() error
}

type HttpProxyServer interface {
	Server
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
