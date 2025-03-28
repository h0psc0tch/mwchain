package mwchain_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h0psc0tch/mwchain"
	"github.com/stretchr/testify/assert"
)

func Test_MWChain(t *testing.T) {

	// mw1 := func(next http.HandlerFunc) http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		// Do something before the next handler
	// 		// Call the next handler
	// 		next(w, r)
	// 		// Do something after the next handler
	// 	}
	// }
	// mw2 := func(next http.HandlerFunc) http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		// Do something before the next handler
	// 		// Call the next handler
	// 		next(w, r)
	// 		// Do something after the next handler
	// 	}
	// }

	t.Run("new empty chain", func(t *testing.T) {
		chain := mwchain.NewChain()
		assert.NotNil(t, chain)

		h := &testHandler{}
		chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		assert.Equal(t, "", w.Header().Get("X-Middleware-Number"))
	})

	t.Run("new single middleware", func(t *testing.T) {
		chain := mwchain.NewChain(mw(1))

		h := &testHandler{}
		sut := chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, "1", h.requestHeader.Get("X-Middleware-Number"))
		assert.Equal(t, "1", w.Header().Get("X-Middleware-Number"))
	})

	t.Run("new multiple middlewares", func(t *testing.T) {
		chain := mwchain.NewChain(mw(1), mw(2))

		h := &testHandler{}
		sut := chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, "12", h.requestHeader.Get("X-Middleware-Number"))
		assert.Equal(t, "21", w.Header().Get("X-Middleware-Number"))
	})

	t.Run("empty new, add single", func(t *testing.T) {
		chain := mwchain.NewChain()

		chain.Add(mw(1))

		h := &testHandler{}
		sut := chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, "1", h.requestHeader.Get("X-Middleware-Number"))
		assert.Equal(t, "1", w.Header().Get("X-Middleware-Number"))
	})

	t.Run("empty new, add multiple", func(t *testing.T) {
		chain := mwchain.NewChain()

		chain.Add(mw(1), mw(2), mw(3))

		h := &testHandler{}
		sut := chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, "123", h.requestHeader.Get("X-Middleware-Number"))
		assert.Equal(t, "321", w.Header().Get("X-Middleware-Number"))
	})

	t.Run("single new, add single", func(t *testing.T) {
		chain := mwchain.NewChain(mw(1))

		chain.Add(mw(2))

		h := &testHandler{}
		sut := chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, "12", h.requestHeader.Get("X-Middleware-Number"))
		assert.Equal(t, "21", w.Header().Get("X-Middleware-Number"))
	})

	t.Run("single new, add multiple", func(t *testing.T) {
		chain := mwchain.NewChain(mw(1))

		chain.Add(mw(2), mw(3), mw(4))

		h := &testHandler{}
		sut := chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, "1234", h.requestHeader.Get("X-Middleware-Number"))
		assert.Equal(t, "4321", w.Header().Get("X-Middleware-Number"))
	})

	t.Run("multiple new, add multiple", func(t *testing.T) {
		chain := mwchain.NewChain(mw(1), mw(2), mw(3))

		chain.Add(mw(4), mw(5), mw(6))
		chain.Add(mw(7), mw(8), mw(9))

		h := &testHandler{}
		sut := chain.Wrap(h.ServeHTTP)

		r, _ := http.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)

		assert.Equal(t, "123456789", h.requestHeader.Get("X-Middleware-Number"))
		assert.Equal(t, "987654321", w.Header().Get("X-Middleware-Number"))
	})
}

func mw(number int) mwchain.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			//apppend current number to the REQUEST header, so order of middleware can be checked
			requestHeader := r.Header.Get("X-Middleware-Number")
			r.Header.Set("X-Middleware-Number", fmt.Sprintf("%s%d", requestHeader, number))

			next(w, r)

			//apppend current number to the RESPONSE header, so order of middleware can be checked
			responseHeader := w.Header().Get("X-Middleware-Number")
			w.Header().Set("X-Middleware-Number", fmt.Sprintf("%s%d", responseHeader, number))

		}
	}
}

type testHandler struct {
	requestHeader http.Header
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Write the response
	h.requestHeader = r.Header.Clone()
	_, _ = w.Write([]byte(""))
}
