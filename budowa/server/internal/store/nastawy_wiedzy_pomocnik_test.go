package store_test

// Nastawy pakietu wiedza czyta sprawdzian z zewnątrz pakietu store: wiedza
// sięga po dane, a dane po store, więc sprawdzian wewnętrzny domykałby krąg
// zależności i pakiet przestałby się kompilować.

import (
	"path/filepath"
	"testing"

	"danacoconsole/server/internal/store"
)

// swiezaBazaWiedzy otwiera bazę przejechaną wszystkimi krokami migracji.
func swiezaBazaWiedzy(t *testing.T) *store.Baza {
	t.Helper()
	baza, err := store.Otworz(filepath.Join(t.TempDir(), "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })
	return baza
}
