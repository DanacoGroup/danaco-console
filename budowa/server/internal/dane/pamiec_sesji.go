// Odpowiedzialność pliku: konfiguracja dostępu sesji do pamięci (tabela
// `konfiguracja_pamieci_sesji`). Pamięć zostaje na poziomie sesji;
// brak wiersza oznacza ustawienie domyślne, nie odmowę dostępu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// KonfiguracjaPamieci opisuje, z których poziomów pamięci korzysta sesja
// i czy wolno jej pamięć zapisywać.
type KonfiguracjaPamieci struct {
	SesjaID         int64
	PoziomyWlaczone []PoziomPamieci
	ZapisWlaczony   bool
	Zaktualizowano  string
}

const (
	zapiszKonfiguracjePamieci = `INSERT INTO konfiguracja_pamieci_sesji
	                             (sesja_id, poziomy_wlaczone, zapis_wlaczony)
	                             VALUES (?, ?, ?)
	                             ON CONFLICT(sesja_id) DO UPDATE SET
	                                 poziomy_wlaczone = excluded.poziomy_wlaczone,
	                                 zapis_wlaczony = excluded.zapis_wlaczony,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzKonfiguracjePamieci = `SELECT sesja_id, poziomy_wlaczone, zapis_wlaczony, zaktualizowano
	                              FROM konfiguracja_pamieci_sesji WHERE sesja_id = ?`
)

// UstawKonfiguracjeSesji zapisuje konfigurację pamięci sesji.
func (r *repozytoriumPamieci) UstawKonfiguracjeSesji(ctx context.Context,
	konfiguracja KonfiguracjaPamieci) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKonfiguracjePamieci)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, konfiguracja.SesjaID,
		zlozPoziomy(konfiguracja.PoziomyWlaczone), liczbaLogiczna(konfiguracja.ZapisWlaczony))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać konfiguracji pamięci sesji %d: %w",
			konfiguracja.SesjaID, err)
	}
	return nil
}

// KonfiguracjaSesji zwraca konfigurację pamięci sesji. Drugi wynik mówi, czy
// Operator ją w ogóle ustawił.
func (r *repozytoriumPamieci) KonfiguracjaSesji(ctx context.Context,
	sesjaID int64) (KonfiguracjaPamieci, bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKonfiguracjePamieci)
	if err != nil {
		return KonfiguracjaPamieci{}, false, err
	}
	var konfiguracja KonfiguracjaPamieci
	var poziomy string
	var zapis int
	err = polecenie.QueryRowContext(ctx, sesjaID).Scan(&konfiguracja.SesjaID, &poziomy, &zapis,
		&konfiguracja.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return KonfiguracjaPamieci{}, false, nil
	}
	if err != nil {
		return KonfiguracjaPamieci{}, false,
			fmt.Errorf("dane: nieczytelna konfiguracja pamięci sesji %d: %w", sesjaID, err)
	}
	konfiguracja.PoziomyWlaczone = rozlozPoziomy(poziomy)
	konfiguracja.ZapisWlaczony = zapis == 1
	return konfiguracja, true, nil
}

// zlozPoziomy zapisuje listę poziomów jako wartość kolumny.
func zlozPoziomy(poziomy []PoziomPamieci) string {
	nazwy := make([]string, 0, len(poziomy))
	for _, poziom := range poziomy {
		if poziom != "" {
			nazwy = append(nazwy, string(poziom))
		}
	}
	return strings.Join(nazwy, ",")
}

// rozlozPoziomy odczytuje listę poziomów z wartości kolumny.
func rozlozPoziomy(wartosc string) []PoziomPamieci {
	poziomy := []PoziomPamieci{}
	for _, nazwa := range strings.Split(wartosc, ",") {
		nazwa = strings.TrimSpace(nazwa)
		if nazwa != "" {
			poziomy = append(poziomy, PoziomPamieci(nazwa))
		}
	}
	return poziomy
}
