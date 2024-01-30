package goweb

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/JhonatanRSantos/gocore/pkg/gocontext"
	"github.com/JhonatanRSantos/gocore/pkg/golog"

	fiber "github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

const (
	defaultAppName      = "ms-backend-default"
	defaultSwaggerTitle = "Default Swagger UI"
	defaultSwaggerRoute = "/swagger/*"
)

type WebRoute struct {
	Method   string
	Path     string
	Handlers []func(c *fiber.Ctx) error
}

type WebServer struct {
	app           *fiber.App
	routers       []WebRoute
	swaggerConfig WebServerSwaggerConfig
}

type WebServerConfig struct {
	app           *fiber.App
	routers       []WebRoute
	swaggerConfig WebServerSwaggerConfig
}

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type RateLimiteConfig struct {
	MaxRequests         int
	MaxRequestsInterval time.Duration
}

type ProfilingConfig struct {
	EndpointPrefix string
}

type WebServerSwaggerConfig struct {
	Title string
	Route string
}

type WebServerLogger interface {
	Info(ctx context.Context, message string, opts ...golog.Options)
	Warn(ctx context.Context, message string, opts ...golog.Options)
	Debug(ctx context.Context, message string, opts ...golog.Options)
	Error(ctx context.Context, message string, opts ...golog.Options)
}

type WebServerDefaultConfig struct {
	AppName    string
	Cors       CorsConfig
	Swagger    WebServerSwaggerConfig
	RateLimite RateLimiteConfig
	Profiling  ProfilingConfig
	Logger     WebServerLogger
}

// DefaultConfig Build the web server default configurations
func DefaultConfig(config WebServerDefaultConfig) *WebServerConfig {
	app := fiber.New(fiber.Config{
		AppName:               config.AppName,
		CompressedFileSuffix:  fmt.Sprintf(".%s.gz", config.AppName),
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if config.Logger != nil {
				var message string
				if _, ok := gocontext.Get[string](c.UserContext(), "panic-error"); ok {
					message = fmt.Sprintf("Recovered from panic. Cause: %s", err)
				} else {
					message = fmt.Sprintf("Unexpected error. Cause: %s", err)
				}
				config.Logger.Error(c.Context(), message)
			}
			return c.SendStatus(http.StatusInternalServerError)
		},
		// ReadTimeout:  time.Second * 5, // max time for reading the request
		// WriteTimeout: time.Second * 5, // max time for write the response
	})

	app.Use(favicon.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.Cors.AllowOrigins, ","),
		AllowMethods:     strings.Join(config.Cors.AllowMethods, ","),
		AllowHeaders:     strings.Join(config.Cors.AllowMethods, ","),
		AllowCredentials: config.Cors.AllowCredentials,
		MaxAge:           config.Cors.MaxAge,
	}))
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			ctx := gocontext.FromContext(c.UserContext())
			if err, ok := e.(error); ok {
				gocontext.Add[string](ctx, "panic-error", err.Error())
			} else {
				gocontext.Add[string](ctx, "panic-error", fmt.Sprint(e))
			}
			c.SetUserContext(ctx)
		},
	}))

	if config.RateLimite != (RateLimiteConfig{}) {
		app.Use(limiter.New(limiter.Config{
			Max:               config.RateLimite.MaxRequests,
			Expiration:        config.RateLimite.MaxRequestsInterval,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
	}

	if config.Profiling != (ProfilingConfig{}) {
		app.Use(pprof.New(pprof.Config{
			Prefix: config.Profiling.EndpointPrefix,
		}))
	}

	return &WebServerConfig{
		app:           app,
		routers:       []WebRoute{},
		swaggerConfig: config.Swagger,
	}
}

// isNil check if the webserver config is nil
func (wsc *WebServerConfig) isNil() bool {
	return wsc.app == nil || wsc.routers == nil
}

// NewWebServer Creates a new web server based using a given configuration
func NewWebServer(config *WebServerConfig) *WebServer {
	if config == nil || config.isNil() {
		config = DefaultConfig(WebServerDefaultConfig{
			AppName: defaultAppName,
			Swagger: WebServerSwaggerConfig{
				Title: defaultSwaggerTitle,
				Route: defaultSwaggerRoute,
			},
		})
	}

	return &WebServer{
		app:           config.app,
		routers:       config.routers,
		swaggerConfig: config.swaggerConfig,
	}
}

// GetApp Get the underlying *fiber.App
func (ws *WebServer) GetApp() *fiber.App {
	return ws.app
}

// AddRoutes Add a new web route
func (ws *WebServer) AddRoutes(routes ...WebRoute) {
	if ws.routers == nil {
		ws.routers = []WebRoute{}
	}

	ws.routers = append(ws.routers, routes...)
}

// gracefulShutdown Terminates the web server
func (ws *WebServer) gracefulShutdown() {
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt)

	go func(shutdownChannel chan os.Signal, app *fiber.App) {
		<-shutdownChannel
		err := app.Shutdown()
		if err != nil {
			panic(err)
		}
	}(shutdownChannel, ws.app)
}

// swaggerUI Configurate the default swagger route
func (ws *WebServer) swaggerUI() {
	if ws.swaggerConfig.Title == "" {
		ws.swaggerConfig.Title = defaultSwaggerTitle
	}

	if ws.swaggerConfig.Route == "" {
		ws.swaggerConfig.Route = defaultSwaggerRoute
	}

	ws.app.Get(ws.swaggerConfig.Route, swagger.New(swagger.Config{
		Title:  ws.swaggerConfig.Title,
		Layout: "BaseLayout",
		Plugins: []template.JS{
			template.JS("SwaggerUIBundle.plugins.DownloadUrl"),
		},
		Presets: []template.JS{
			template.JS("SwaggerUIBundle.presets.apis"),
			template.JS("SwaggerUIStandalonePreset"),
		},
		DeepLinking:              true,
		DefaultModelsExpandDepth: 1,
		DefaultModelExpandDepth:  1,
		DefaultModelRendering:    "example",
		DocExpansion:             "list",
		SyntaxHighlight: &swagger.SyntaxHighlightConfig{
			Activate: true,
			Theme:    "agate",
		},
		ShowMutatedRequest: true,
	}))
}

// Listen Start the web server
func (ws *WebServer) Listen(address string) error {
	ws.gracefulShutdown()

	if ws.swaggerConfig != (WebServerSwaggerConfig{}) {
		ws.swaggerUI()
	}

	if ws.routers != nil && len(ws.routers) > 0 {
		for _, route := range ws.routers {
			ws.app.Add(route.Method, route.Path, route.Handlers...)
		}
	}
	return ws.app.Listen(address)
}
