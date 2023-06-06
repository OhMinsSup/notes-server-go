package core

const (
	DefaultDataMaxOpenConns int = 120
	DefaultDataMaxIdleConns int = 20
	DefaultLogsMaxOpenConns int = 10
	DefaultLogsMaxIdleConns int = 2

	LocalStorageDirName string = "storage"
	LocalBackupsDirName string = "backups"
	LocalTempDirName    string = ".pb_temp_to_delete" // temp pb_data sub directory that will be deleted on each app.Bootstrap()
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	// configurable parameters
	isDebug          bool
	dataDir          string
	encryptionEnv    string
	dataMaxOpenConns int
	dataMaxIdleConns int
	logsMaxOpenConns int
	logsMaxIdleConns int
}

func NewBaseApp() *BaseApp {
	return &BaseApp{}
}

func (app *BaseApp) DB() {}

func (app *BaseApp) initLogsDB() error {}

func (app *BaseApp) initDataDB() error {}
