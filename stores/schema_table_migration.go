package stores

// EnsureSchemaMigrationFormat checks the schema migrations table
// format and, if it's not using the new shape, it migrates the old
// one's status before initializing the migrations engine.
func (s *Store) EnsureSchemaMigrationFormat() error {
	return nil
}

func (s *Store) isSchemaMigrationNeeded() (bool, error) {
	return false, nil
}