package goweb

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	fiber "github.com/gofiber/fiber/v2"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

func TestWebServer(t *testing.T) {
	ws := NewWebServer(nil)
	assert.NotNil(t, ws, "failed to start webserver with nil configs")

	ws = NewWebServer(DefaultConfig(WebServerDefaultConfig{
		Swagger: WebServerSwaggerConfig{
			Title: "swagger test",
			Route: "/swagger/*",
		},
		RateLimite: RateLimiteConfig{
			MaxRequests:         10,
			MaxRequestsInterval: time.Second * 5,
		},
		Profiling: ProfilingConfig{
			EndpointPrefix: "debug",
		},
		JSONConfig: JSONConfig{
			Encoder: sonic.Marshal,
			Decoder: sonic.Unmarshal,
		},
	}))

	baseHandler := func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	}

	panicHandler := func(c *fiber.Ctx) error {
		panic(errors.New("panic webserver"))
	}

	routes := []WebRoute{
		{
			Method:   "GET",
			Path:     "/",
			Handlers: []func(c *fiber.Ctx) error{baseHandler},
		},
		{
			Method:   "GET",
			Path:     "/panic",
			Handlers: []func(c *fiber.Ctx) error{panicHandler},
		},
	}

	ws.AddRoutes(routes...)

	listener, err := net.Listen("tcp", ":0")
	assert.NoErrorf(t, err, "failed to get net listener. Cause: %s", err)

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	assert.NoErrorf(t, err, "failed to close listener. Cause: %s", err)

	go func(ws *WebServer, port int) {
		if err := ws.Listen(fmt.Sprintf(":%d", port)); err != nil {
			panic(fmt.Errorf("failed to close webserver. Cause: %s", err))
		}
	}(ws, port)

	defer func() {
		if err := ws.GetApp().Shutdown(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second * 2)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	assert.NoError(t, err, "failed to GET /. Cause: %s", err)
	defer resp.Body.Close()
	assert.Equalf(t, http.StatusOK, resp.StatusCode, "invalid status code when calling GET /. Expected 200 but got %d", resp.StatusCode)

	bs, err := io.ReadAll(resp.Body)
	assert.NoErrorf(t, err, "failed to read response body from GET / request. Cause: %s", err)
	assert.Equalf(t, "OK", string(bs), "invalid response body when calling GET /. Expected OK but got %s", bs)

	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/swagger", port))
	assert.NoErrorf(t, err, "failed to GET /. Cause: %s", err)
	defer resp.Body.Close()

	assert.Equalf(t, http.StatusOK, resp.StatusCode, "invalid status code when calling GET /. Expected 200 but got %d", resp.StatusCode)

	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/panic", port))
	assert.NoError(t, err, "failed to GET /panic. Cause: %s", err)
	defer resp.Body.Close()
	assert.Equalf(t, http.StatusInternalServerError, resp.StatusCode, "invalid status code when calling GET /panic. Expected 500 but got %d", resp.StatusCode)

	bs, err = io.ReadAll(resp.Body)
	assert.NoErrorf(t, err, "failed to read response body from GET /panic request. Cause: %s", err)
	assert.Equalf(t, "Internal Server Error", string(bs), "invalid response body when calling GET /. Expected Internal Server Error but got %s", bs)
}
