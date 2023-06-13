package apis

import (
	"github.com/OhMinsSup/notes-server-go/services"
	"github.com/OhMinsSup/notes-server-go/stores"
	"github.com/OhMinsSup/notes-server-go/tools/hook"
)

// App은 기본 인터페이스를 정의합니다.
type App interface {
	DataURL() string
	IsDebug() bool
	Bootstrap() error
	ResetBootstrapState() error
	Store() *stores.Store
	Service() *services.Service

	// ---------------------------------------------------------------
	// App event hooks
	// ---------------------------------------------------------------

	// OnBeforeBootstrap hook is triggered before initializing the base
	// application resources (eg. before db open and initial settings load).
	OnBeforeBootstrap() *hook.Hook[*BootstrapEvent]

	// OnAfterBootstrap hook is triggered after initializing the base
	// application resources (eg. after db open and initial settings load).
	OnAfterBootstrap() *hook.Hook[*BootstrapEvent]

	// OnBeforeServe hook is triggered before serving the internal router (echo),
	// allowing you to adjust its options and attach new routes.
	OnBeforeServe() *hook.Hook[*ServeEvent]

	// OnBeforeApiError hook is triggered right before sending an error API
	// response to the client, allowing you to further modify the error data
	// or to return a completely different API response (using [hook.StopPropagation]).
	OnBeforeApiError() *hook.Hook[*ApiErrorEvent]

	// OnAfterApiError hook is triggered right after sending an error API
	// response to the client.
	// It could be used to log the final API error in external services.
	OnAfterApiError() *hook.Hook[*ApiErrorEvent]

	// OnTerminate hook is triggered when the app is in the process
	// of being terminated (eg. on SIGTERM signal).
	OnTerminate() *hook.Hook[*TerminateEvent]
}
