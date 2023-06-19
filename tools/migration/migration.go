package migration

import "github.com/OhMinsSup/notes-server-go/tools/utils"

const (
	DatabaseDriverMysql    = "mysql"
	DatabaseDriverPostgres = "postgres"
	DatabaseDriverSqlite3  = "sqlite3"

	// sqlite3 is the default database driver. connect_timeout=10&binary_parameters=yes
	SqlSettingsDefaultDataSource = "file:data/database.sqlite?cache=shared&mode=rwc&_journal_mode=WAL&_txlock=immediate&connect_timeout=10&binary_parameters=yes"
)

type ReplicaLagSettings struct {
	DataSource       *string `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
	QueryAbsoluteLag *string `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
	QueryTimeLag     *string `access:"environment,write_restrictable,cloud_restrictable"` // telemetry: none
}

type SqlSettings struct {
	DriverName                        *string               `access:"environment_database,write_restrictable,cloud_restrictable"`
	DataSource                        *string               `access:"environment_database,write_restrictable,cloud_restrictable"` // telemetry: none
	DataSourceReplicas                []string              `access:"environment_database,write_restrictable,cloud_restrictable"`
	DataSourceSearchReplicas          []string              `access:"environment_database,write_restrictable,cloud_restrictable"`
	MaxIdleConns                      *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	ConnMaxLifetimeMilliseconds       *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	ConnMaxIdleTimeMilliseconds       *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	MaxOpenConns                      *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	Trace                             *bool                 `access:"environment_database,write_restrictable,cloud_restrictable"`
	AtRestEncryptKey                  *string               `access:"environment_database,write_restrictable,cloud_restrictable"` // telemetry: none
	QueryTimeout                      *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	DisableDatabaseSearch             *bool                 `access:"environment_database,write_restrictable,cloud_restrictable"`
	MigrationsStatementTimeoutSeconds *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
	ReplicaLagSettings                []*ReplicaLagSettings `access:"environment_database,write_restrictable,cloud_restrictable"` // telemetry: none
	ReplicaMonitorIntervalSeconds     *int                  `access:"environment_database,write_restrictable,cloud_restrictable"`
}

func (s *SqlSettings) SetDefaults(isUpdate bool) {
	if s.DriverName == nil {
		s.DriverName = utils.NewString(DatabaseDriverSqlite3)
	}

	if s.DataSource == nil {
		s.DataSource = utils.NewString(SqlSettingsDefaultDataSource)
	}

	if s.DataSourceReplicas == nil {
		s.DataSourceReplicas = []string{}
	}

	if s.DataSourceSearchReplicas == nil {
		s.DataSourceSearchReplicas = []string{}
	}

	if isUpdate {
		// When updating an existing configuration, ensure an encryption key has been specified.
		if s.AtRestEncryptKey == nil || *s.AtRestEncryptKey == "" {
			s.AtRestEncryptKey = utils.NewString(utils.NewRandomString(32))
		}
	} else {
		// When generating a blank configuration, leave this key empty to be generated on server start.
		s.AtRestEncryptKey = utils.NewString("")
	}

	if s.MaxIdleConns == nil {
		s.MaxIdleConns = utils.NewInt(20)
	}

	if s.MaxOpenConns == nil {
		s.MaxOpenConns = utils.NewInt(300)
	}

	if s.ConnMaxLifetimeMilliseconds == nil {
		s.ConnMaxLifetimeMilliseconds = utils.NewInt(3600000)
	}

	if s.ConnMaxIdleTimeMilliseconds == nil {
		s.ConnMaxIdleTimeMilliseconds = utils.NewInt(300000)
	}

	if s.Trace == nil {
		s.Trace = utils.NewBool(false)
	}

	if s.QueryTimeout == nil {
		s.QueryTimeout = utils.NewInt(30)
	}

	if s.DisableDatabaseSearch == nil {
		s.DisableDatabaseSearch = utils.NewBool(false)
	}

	if s.MigrationsStatementTimeoutSeconds == nil {
		s.MigrationsStatementTimeoutSeconds = utils.NewInt(100000)
	}

	if s.ReplicaLagSettings == nil {
		s.ReplicaLagSettings = []*ReplicaLagSettings{}
	}

	if s.ReplicaMonitorIntervalSeconds == nil {
		s.ReplicaMonitorIntervalSeconds = utils.NewInt(5)
	}
}