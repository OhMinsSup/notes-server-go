package apis

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	sqlstore "github.com/OhMinsSup/notes-server-go"
	"github.com/OhMinsSup/notes-server-go/settings"
	"github.com/OhMinsSup/notes-server-go/stores"
	"github.com/OhMinsSup/notes-server-go/tools/config"
	"github.com/OhMinsSup/notes-server-go/tools/hook"
	"github.com/OhMinsSup/notes-server-go/tools/rest"
	"github.com/OhMinsSup/notes-server-go/tools/serve"
	"github.com/fatih/color"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	DefaultDataMaxOpenConns int = 120
	DefaultDataMaxIdleConns int = 20
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	// configurable parameters
	config           *config.Configuration
	serverOptions    *serve.ServeOptions
	dataMaxOpenConns int
	dataMaxIdleConns int

	// internals
	store    *stores.Store
	settings *settings.Settings

	// app event hooks
	onBeforeBootstrap *hook.Hook[*BootstrapEvent]
	onAfterBootstrap  *hook.Hook[*BootstrapEvent]
	onBeforeServe     *hook.Hook[*ServeEvent]
	onBeforeApiError  *hook.Hook[*ApiErrorEvent]
	onAfterApiError   *hook.Hook[*ApiErrorEvent]
	onTerminate       *hook.Hook[*TerminateEvent]
}

type BaseAppConfig struct {
	// configurable parameters
	Config           *config.Configuration
	ServerOptions    *serve.ServeOptions
	DataMaxOpenConns int // default to 500
	DataMaxIdleConns int // default 20
}

func NewBaseApp(config *BaseAppConfig) *BaseApp {
	app := &BaseApp{
		config:           config.Config,
		serverOptions:    config.ServerOptions,
		dataMaxOpenConns: config.DataMaxOpenConns,
		dataMaxIdleConns: config.DataMaxIdleConns,
		settings:         settings.New(config.Config),

		// app event hooks
		onBeforeBootstrap: &hook.Hook[*BootstrapEvent]{},
		onAfterBootstrap:  &hook.Hook[*BootstrapEvent]{},
		onBeforeServe:     &hook.Hook[*ServeEvent]{},
		onBeforeApiError:  &hook.Hook[*ApiErrorEvent]{},
		onAfterApiError:   &hook.Hook[*ApiErrorEvent]{},
		onTerminate:       &hook.Hook[*TerminateEvent]{},
	}

	app.registerDefaultHooks()

	return app
}

func (app *BaseApp) Bootstrap() error {
	event := &BootstrapEvent{app}

	if err := app.OnBeforeBootstrap().Trigger(event); err != nil {
		return err
	}

	// clear resources of previous core state (if any)
	if err := app.ResetBootstrapState(); err != nil {
		return err
	}

	if err := app.initDataDB(); err != nil {
		return err
	}

	if err := app.initServer(); err != nil {
		return err
	}

	if err := app.OnAfterBootstrap().Trigger(event); err != nil && app.IsDebug() {
		log.Println(err)
	}

	return nil
}

func (app *BaseApp) initDataDB() error {
	store, err := app.createStore()
	if err != nil {
		return err
	}

	app.store = store
	date := new(strings.Builder)
	log.New(date, "", log.LstdFlags).Print()
	bold := color.New(color.Bold).Add(color.FgGreen)
	bold.Printf(
		"%s Database Connection started\n",
		strings.TrimSpace(date.String()),
	)
	return nil
}

func (app *BaseApp) initServer() error {
	// create router
	router := app.createRouter()

	// start http server
	// ---
	mainAddr := app.serverOptions.HttpAddr
	if app.serverOptions.HttpsAddr != "" {
		mainAddr = app.serverOptions.HttpsAddr
	}

	mainHost, _, _ := net.SplitHostPort(mainAddr)
	log.Printf("main host: %s\n", mainHost)

	serverConfig := &http.Server{
		ReadTimeout:       10 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		// WriteTimeout: 60 * time.Second, // breaks sse!
		Handler: router,
		Addr:    mainAddr,
	}

	if app.serverOptions.BeforeServeFunc != nil {
		if err := app.serverOptions.BeforeServeFunc(serverConfig); err != nil {
			return err
		}
	}

	if app.serverOptions.ShowStartBanner {
		schema := "http"
		if app.serverOptions.HttpsAddr != "" {
			schema = "https"
		}

		date := new(strings.Builder)
		log.New(date, "", log.LstdFlags).Print()

		bold := color.New(color.Bold).Add(color.FgGreen)
		bold.Printf(
			"%s Server started at %s\n",
			strings.TrimSpace(date.String()),
			color.CyanString("%s://%s", schema, serverConfig.Addr),
		)

		regular := color.New()
		regular.Printf(" ➜ REST API: %s\n", color.CyanString("%s://%s/api/", schema, serverConfig.Addr))
		regular.Printf(" ➜ Admin UI: %s\n", color.CyanString("%s://%s/_/", schema, serverConfig.Addr))
	}

	// start HTTPS server
	if app.serverOptions.HttpsAddr != "" {
		// if httpAddr is set, start an HTTP server to redirect the traffic to the HTTPS version
		if app.serverOptions.HttpAddr != "" {
			// TODO: add a flag to disable this
		}

		return serverConfig.ListenAndServeTLS("", "")
	}

	// OR start HTTP server
	return serverConfig.ListenAndServe()
}

func (app *BaseApp) ResetBootstrapState() error {
	return nil
}

func (app *BaseApp) createStore() (*stores.Store, error) {
	client := sqlstore.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}

	ctx := context.Background()

	var store *stores.Store
	params := stores.Params{
		DB:  client,
		Ctx: ctx,
	}

	store, err := stores.New(params)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (app *BaseApp) createRouter() *echo.Echo {

	// TODO: DB Migration
	router := echo.New()

	router.Debug = app.IsDebug()
	router.JSONSerializer = &rest.Serializer{
		FieldsParam: "fields",
	}

	// configure a custom router
	router.ResetRouterCreator(func(ec *echo.Echo) echo.Router {
		return echo.NewRouter(echo.RouterConfig{
			UnescapePathParamValues: true,
		})
	})

	// default middlewares
	router.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.RemoveTrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			// enable by default only for the API routes
			return !strings.HasPrefix(c.Request().URL.Path, "/api/")
		},
	}))
	router.Use(middleware.Recover())
	router.Use(middleware.Secure())
	router.Use(LoadAuthContext(app))

	// custom error handler
	router.HTTPErrorHandler = func(c echo.Context, err error) {
		if c.Response().Committed {
			if app.IsDebug() {
				log.Println("HTTPErrorHandler response was already committed:", err)
			}
			return
		}

		var apiErr *ApiError

		switch v := err.(type) {
		case *echo.HTTPError:
			if v.Internal != nil && app.IsDebug() {
				log.Println(v.Internal)
			}
			msg := fmt.Sprintf("%v", v.Message)
			apiErr = NewApiError(v.Code, msg, v)
		case *ApiError:
			if app.IsDebug() && v.RawData() != nil {
				log.Println(v.RawData())
			}
			apiErr = v
		default:
			if err != nil && app.IsDebug() {
				log.Println(err)
			}
			apiErr = NewBadRequestError("", err)
		}

		event := new(ApiErrorEvent)
		event.HttpContext = c
		event.Error = apiErr

		// send error response
		hookErr := app.OnBeforeApiError().Trigger(event, func(e *ApiErrorEvent) error {
			// @see https://github.com/labstack/echo/issues/608
			if e.HttpContext.Request().Method == http.MethodHead {
				return e.HttpContext.NoContent(apiErr.Code)
			}

			return e.HttpContext.JSON(apiErr.Code, apiErr)
		})

		// truly rare case; eg. client already disconnected
		if hookErr != nil && app.IsDebug() {
			log.Println(hookErr)
		}

		app.OnAfterApiError().Trigger(event)
	}

	// default routes
	api := router.Group("/api")
	bindHealthApi(app, api)
	bindAuthApi(app, api)

	return router

}

// -------------------------------------------------------------------
// App event hooks
// -------------------------------------------------------------------

func (app *BaseApp) OnBeforeBootstrap() *hook.Hook[*BootstrapEvent] {
	return app.onBeforeBootstrap
}

func (app *BaseApp) OnAfterBootstrap() *hook.Hook[*BootstrapEvent] {
	return app.onAfterBootstrap
}

func (app *BaseApp) OnBeforeServe() *hook.Hook[*ServeEvent] {
	return app.onBeforeServe
}

func (app *BaseApp) OnBeforeApiError() *hook.Hook[*ApiErrorEvent] {
	return app.onBeforeApiError
}

func (app *BaseApp) OnAfterApiError() *hook.Hook[*ApiErrorEvent] {
	return app.onAfterApiError
}

func (app *BaseApp) OnTerminate() *hook.Hook[*TerminateEvent] {
	return app.onTerminate
}

// -------------------------------------------------------------------

func (app *BaseApp) IsDebug() bool {
	return app.config.IsDebug
}

// Dao returns the default app Dao instance.
func (app *BaseApp) Store() *stores.Store {
	return app.store
}

func (app *BaseApp) Config() *config.Configuration {
	return app.config
}

func (app *BaseApp) registerDefaultHooks() {
	app.OnTerminate().Add(func(e *TerminateEvent) error {
		log.Println("Terminating the application...")
		app.ResetBootstrapState()
		return nil
	})
}

// Settings returns the loaded app settings.
func (app *BaseApp) Settings() *settings.Settings {
	return app.settings
}
