package mwchain

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MWChain struct {
	middlewares []Middleware
}

func NewChain(middlewares ...Middleware) *MWChain {
	return &MWChain{middlewares: middlewares}
}

func (m *MWChain) Add(middlewares ...Middleware) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *MWChain) Wrap(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	//first wrap with the handler-specific middleware
	h = wrapHandlerFunc(h, middlewares)

	//then wrap with the chain middleware
	return wrapHandlerFunc(h, m.middlewares)
}

func wrapHandlerFunc(h http.HandlerFunc, m []Middleware) http.HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		mw := m[i]
		if mw == nil {
			continue
		}
		h = mw(h)
	}
	return h
}
