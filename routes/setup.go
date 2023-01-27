package routes

import (
	"docker_test/src/config"
	"docker_test/src/helpers"
	"docker_test/src/logs"
	"docker_test/src/middlewares"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Counter is a struct that tracks the number of requests made within a given time interval
type Counter struct {
	sync.Mutex
	count    int
	limit    int
	interval time.Duration
	timer    *time.Timer
}

func Setup() {
	// Echo instance
	e := echo.New()

	e.Renderer = setupTemplate()

	// ApiRoot with Echo instance
	r := echoswagger.New(e, "/docs", &echoswagger.Info{
		Title:   "docker_test API",
		Version: "0.0.1",
	})

	// Middleware bind to echo instance
	r.Echo().Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${remote_ip}::${method}${path}, status=${status}, latency=${latency_human}, error=${error}\n",
		Output: io.MultiWriter(os.Stdout, logs.LogFile),
	}))
	r.Echo().Use(middleware.RecoverWithConfig(
		middleware.RecoverConfig{
			StackSize: 1 << 10, // 1 KB
		},
	))

	// Rate limit
	// r.Echo().Use(RequestLimiter(config.RATELIMITPERMINUTE, time.Minute)) // All incomming
	r.Echo().Use(middleware.RateLimiterWithConfig(configRateLimit())) // only ip

	r.Echo().Validator = &middlewares.CustomValidator{Validator: validator.New()}
	// ------------------------------ Middleware ------------------------------

	// setup route
	RoutePath(r)

	// Start server
	r.Echo().Logger.Info(e.Start(":" + os.Getenv("SERVER_PORT")))
}

func configRateLimit() middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(config.RATELIMITPERMINUTE_IP), // 10 request per second
	}
}

// SETUP TEMPLATE
type Template struct {
	templates *template.Template
}

func setupTemplate() *Template {
	t := &Template{
		templates: template.Must(template.ParseGlob(filepath.Join(helpers.RootDir(), "views", "*.html"))),
	}
	return t
}
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// END SETUP TEMPLATE

// NewCounter creates a new Counter with the given limit and interval
func NewCounter(limit int, interval time.Duration) *Counter {
	c := &Counter{
		limit:    limit,
		interval: interval,
		timer:    time.NewTimer(interval),
	}
	go func() {
		<-c.timer.C
		c.Lock()
		c.count = 0
		c.timer.Reset(interval)
		c.Unlock()
	}()
	return c
}

// Request increments the counter and returns true if the request should be allowed
func (c *Counter) Request() bool {
	c.Lock()
	defer c.Unlock()
	if c.count >= c.limit {
		return false
	}
	c.count++
	return true
}

func RequestLimiter(limit int, interval time.Duration) echo.MiddlewareFunc {
	counter := NewCounter(limit, interval)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !counter.Request() {
				return c.String(http.StatusTooManyRequests, "Too many requests")
			}
			return next(c)
		}
	}
}
