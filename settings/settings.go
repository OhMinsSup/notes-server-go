package settings

import (
	"encoding/json"
	"sync"

	"github.com/OhMinsSup/notes-server-go/tools/config"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// SecretMask is the default settings secrets replacement value
// (see Settings.RedactClone()).
const SecretMask string = "******"

// Settings defines common app configuration options.
type Settings struct {
	mux sync.RWMutex

	Logs LogsConfig `form:"logs" json:"logs"`

	RecordAuthToken TokenConfig `form:"recordAuthToken" json:"recordAuthToken"`
}

// New creates and returns a new default Settings instance.
func New(config *config.Configuration) *Settings {
	return &Settings{
		Logs: LogsConfig{
			MaxDays: 5,
		},
		RecordAuthToken: TokenConfig{
			Secret:   config.RecordAuthTokenSecret,
			Duration: 1209600, // 14 days
		},
	}
}

// Validate makes Settings validatable by implementing [validation.Validatable] interface.
func (s *Settings) Validate() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return validation.ValidateStruct(s,
		validation.Field(&s.Logs),
		validation.Field(&s.RecordAuthToken),
	)
}

// Merge merges `other` settings into the current one.
func (s *Settings) Merge(other *Settings) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	bytes, err := json.Marshal(other)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

// Clone creates a new deep copy of the current settings.
func (s *Settings) Clone() (*Settings, error) {
	clone := &Settings{}
	if err := clone.Merge(s); err != nil {
		return nil, err
	}
	return clone, nil
}

// RedactClone creates a new deep copy of the current settings,
// while replacing the secret values with `******`.
func (s *Settings) RedactClone() (*Settings, error) {
	clone, err := s.Clone()
	if err != nil {
		return nil, err
	}

	sensitiveFields := []*string{
		&clone.RecordAuthToken.Secret,
	}

	// mask all sensitive fields
	for _, v := range sensitiveFields {
		if v != nil && *v != "" {
			*v = SecretMask
		}
	}

	return clone, nil
}

// -------------------------------------------------------------------

type TokenConfig struct {
	Secret   string `form:"secret" json:"secret"`
	Duration int64  `form:"duration" json:"duration"`
}

// Validate makes TokenConfig validatable by implementing [validation.Validatable] interface.
func (c TokenConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Secret, validation.Required, validation.Length(30, 300)),
		validation.Field(&c.Duration, validation.Required, validation.Min(5), validation.Max(63072000)),
	)
}

// -------------------------------------------------------------------

type LogsConfig struct {
	MaxDays int `form:"maxDays" json:"maxDays"`
}

// Validate makes LogsConfig validatable by implementing [validation.Validatable] interface.
func (c LogsConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.MaxDays, validation.Min(0)),
	)
}
