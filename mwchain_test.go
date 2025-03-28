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
	tests := []struct {
		name                 string
		initialMiddleware    []mwchain.Middleware
		additionalMiddleware []mwchain.Middleware
		handlerMiddleware    []mwchain.Middleware
		expectedRequest      string
		expectedResponse     string
	}{
		{
			name:                 "new empty chain",
			initialMiddleware:    nil,
			additionalMiddleware: nil,
			expectedRequest:      "",
			expectedResponse:     "",
		},
		{
			name:                 "new single middleware",
			initialMiddleware:    []mwchain.Middleware{mw(1)},
			additionalMiddleware: nil,
			expectedRequest:      "1",
			expectedResponse:     "1",
		},
		{
			name:                 "new multiple middlewares",
			initialMiddleware:    []mwchain.Middleware{mw(1), mw(2)},
			additionalMiddleware: nil,
			expectedRequest:      "12",
			expectedResponse:     "21",
		},
		{
			name:                 "empty new, add single",
			initialMiddleware:    nil,
			additionalMiddleware: []mwchain.Middleware{mw(1)},
			expectedRequest:      "1",
			expectedResponse:     "1",
		},
		{
			name:                 "empty new, add multiple",
			initialMiddleware:    nil,
			additionalMiddleware: []mwchain.Middleware{mw(1), mw(2), mw(3)},
			expectedRequest:      "123",
			expectedResponse:     "321",
		},
		{
			name:                 "single new, add single",
			initialMiddleware:    []mwchain.Middleware{mw(1)},
			additionalMiddleware: []mwchain.Middleware{mw(2)},
			expectedRequest:      "12",
			expectedResponse:     "21",
		},
		{
			name:                 "single new, add multiple",
			initialMiddleware:    []mwchain.Middleware{mw(1)},
			additionalMiddleware: []mwchain.Middleware{mw(2), mw(3), mw(4)},
			expectedRequest:      "1234",
			expectedResponse:     "4321",
		},
		{
			name:                 "multiple new, add multiple",
			initialMiddleware:    []mwchain.Middleware{mw(1), mw(2), mw(3)},
			additionalMiddleware: []mwchain.Middleware{mw(4), mw(5), mw(6)},
			expectedRequest:      "123456",
			expectedResponse:     "654321",
		},
		{
			name:                 "multiple new, add multiple, multiple handler",
			initialMiddleware:    []mwchain.Middleware{mw(1), mw(2), mw(3)},
			additionalMiddleware: []mwchain.Middleware{mw(4), mw(5), mw(6)},
			handlerMiddleware:    []mwchain.Middleware{mw(7), mw(8), mw(9)},
			expectedRequest:      "123456789",
			expectedResponse:     "987654321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create a new middleware chain with the initial middlewares
			chain := mwchain.NewChain(tt.initialMiddleware...)

			// Add additional middlewares to the chain
			chain.Add(tt.additionalMiddleware...)

			h := &testHandler{}
			sut := chain.Wrap(h.ServeHTTP, tt.handlerMiddleware...)

			r, _ := http.NewRequest("GET", "http://example.com", nil)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedRequest, h.requestHeader.Get("X-Middleware-Number"))
			assert.Equal(t, tt.expectedResponse, w.Header().Get("X-Middleware-Number"))
		})
	}
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
