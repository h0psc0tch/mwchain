package mwchain

import "net/http"

// Middleware is a function that allows you to wrap an http.HandlerFunc with additional functionality.
// It takes an http.HandlerFunc as input and returns a new http.HandlerFunc, thus allowing the caller
// to add functionality before and/or after the original handler is called.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// MWChain is simple system for managing middleware that can then be used to wrap http.HandlerFunc.
// Middleware functions can be added to a chanin, and then the chain can be used to wrap a handler.
// The order in which the middleware is added to the chain is the order in which it will be applied to the handler.
type MWChain struct {
	middlewares []Middleware
}

// NewChain creates a new MWChain with the given middlewares.
func NewChain(middlewares ...Middleware) *MWChain {
	return &MWChain{middlewares: middlewares}
}

// Add adds the given middlewares to the chain, appending the supplid middlewares to the end of the existing chain (i.e. after any previously added middlewares).
func (m *MWChain) Add(middlewares ...Middleware) {
	m.middlewares = append(m.middlewares, middlewares...)
}

// Wrap wraps the given http.HandlerFunc with the middlewares in the chain.
// In addition, it is also possible to pass in handler-specific middlewares that will be applied first.
// The order in which the middlewares are applied is the order in which they were added to the chain, i.e. FIFO
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
