package core

import (
	"path/filepath"
	"time"

	"github.com/OhMinsSup/notes-server-go/daos"
	"github.com/OhMinsSup/notes-server-go/tools/config"
	"github.com/OhMinsSup/notes-server-go/tools/hook"
	"xorm.io/xorm"
)

const (
	DefaultDataMaxOpenConns int = 120
	DefaultDataMaxIdleConns int = 20
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	// configurable parameters
	config           *config.Configuration
	dataMaxOpenConns int
	dataMaxIdleConns int

	// internals
	dao *daos.Dao

	// app event hooks
	onBeforeBootstrap *hook.Hook[*BootstrapEvent]
	onAfterBootstrap  *hook.Hook[*BootstrapEvent]
	onBeforeServe     *hook.Hook[*ServeEvent]
	onBeforeApiError  *hook.Hook[*ApiErrorEvent]
	onAfterApiError   *hook.Hook[*ApiErrorEvent]
}

type BaseAppConfig struct {
	// configurable parameters
	config           *config.Configuration
	DataMaxOpenConns int // default to 500
	DataMaxIdleConns int // default 20
}

func NewBaseApp(config *BaseAppConfig) *BaseApp {
	app := &BaseApp{
		config:           config.config,
		dataMaxOpenConns: config.DataMaxOpenConns,
		dataMaxIdleConns: config.DataMaxIdleConns,

		// app event hooks
		onBeforeBootstrap: &hook.Hook[*BootstrapEvent]{},
		onAfterBootstrap:  &hook.Hook[*BootstrapEvent]{},
		onBeforeServe:     &hook.Hook[*ServeEvent]{},
		onBeforeApiError:  &hook.Hook[*ApiErrorEvent]{},
		onAfterApiError:   &hook.Hook[*ApiErrorEvent]{},
	}
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
	return nil
}

func (app *BaseApp) initDataDB() error {
	maxOpenConns := DefaultDataMaxOpenConns
	maxIdleConns := DefaultDataMaxIdleConns
	if app.dataMaxOpenConns > 0 {
		maxOpenConns = app.dataMaxOpenConns
	}
	if app.dataMaxIdleConns > 0 {
		maxIdleConns = app.dataMaxIdleConns
	}

	db, err := connectDB(filepath.Join(app.DataURL()))
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	app.dao = app.createDao(db)

	return nil
}

func (app *BaseApp) ResetBootstrapState() error {
	return nil
}

func (app *BaseApp) createDao(db *xorm.Engine) *daos.Dao {
	dao := daos.New(db)
	return dao
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

// -------------------------------------------------------------------

func (app *BaseApp) DataURL() string {
	return app.config.DBConfigString
}

func (app *BaseApp) IsDebug() bool {
	return app.config.IsDebug
}
