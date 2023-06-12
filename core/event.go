package core

import "github.com/labstack/echo/v5"

// -------------------------------------------------------------------
// Serve events data
// -------------------------------------------------------------------

type BootstrapEvent struct {
	App App
}

type TerminateEvent struct {
	App App
}

type ServeEvent struct {
	App    App
	Router *echo.Echo
}

type ApiErrorEvent struct {
	HttpContext echo.Context
	Error       error
}
