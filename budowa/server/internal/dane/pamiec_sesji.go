// Odpowiedzialność pliku: konfiguracja dostępu sesji do pamięci (tabela
// `konfiguracja_pamieci_sesji`); brak wiersza oznacza ustawienie domyślne.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type KonfiguracjaPamieci struct {
	SesjaID         int64
	PoziomyWlaczone []PoziomPamieci
	ZapisWlaczony   bool
	Zaktualizowano  string
}

// Tabela `sesja` konta nie niesie: granica dochodzi przez `karta_sesji.konto_id` (migracja 407), aliasem z sesje.go.
var (
	warunekKontaSesjiPamieci = `EXISTS (SELECT 1 FROM sesja s JOIN karta_sesji k ON k.id = s.karta_sesji_id
	                                     WHERE s.id = konfiguracja_pamieci_sesji.sesja_id
	                                       AND ` + warunekKontaKartySesji + `)`

	zapiszKonfiguracjePamieci = `INSERT INTO konfiguracja_pamieci_sesji
	                             (sesja_id, poziomy_wlaczone, zapis_wlaczony)
	                             SELECT ?, ?, ?
	                              WHERE EXISTS (SELECT 1 FROM sesja s
	                                             JOIN karta_sesji k ON k.id = s.karta_sesji_id
	                                            WHERE s.id = ? AND ` + warunekKontaKartySesji + `)
	                             ON CONFLICT(sesja_id) DO UPDATE SET
	                                 poziomy_wlaczone = excluded.poziomy_wlaczone,
	                                 zapis_wlaczony = excluded.zapis_wlaczony,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE ` + warunekKontaSesjiPamieci

	pobierzKonfiguracjePamieci = `SELECT sesja_id, poziomy_wlaczone, zapis_wlaczony, zaktualizowano
	                              FROM konfiguracja_pamieci_sesji
	                              WHERE sesja_id = ? AND ` + warunekKontaSesjiPamieci
)

func (r *repozytoriumPamieci) UstawKonfiguracjeSesji(ctx context.Context,
	konfiguracja KonfiguracjaPamieci) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKonfiguracjePamieci)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, konfiguracja.SesjaID,
		zlozPoziomy(konfiguracja.PoziomyWlaczone), liczbaLogiczna(konfiguracja.ZapisWlaczony),
		konfiguracja.SesjaID, KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać konfiguracji pamięci sesji %d: %w",
			konfiguracja.SesjaID, err)
	}
	return sprawdzTrafienieZapisu(wynik, "konfiguracja pamięci sesji",
		fmt.Sprintf("%d", konfiguracja.SesjaID))
}

// KonfiguracjaSesji zwraca konfigurację pamięci sesji; drugi wynik mówi, czy Operator ją ustawił.
func (r *repozytoriumPamieci) KonfiguracjaSesji(ctx context.Context,
	sesjaID int64) (KonfiguracjaPamieci, bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKonfiguracjePamieci)
	if err != nil {
		return KonfiguracjaPamieci{}, false, err
	}
	var konfiguracja KonfiguracjaPamieci
	var poziomy string
	var zapis int
	err = polecenie.QueryRowContext(ctx, sesjaID, KontoOperatora(ctx)).Scan(&konfiguracja.SesjaID,
		&poziomy, &zapis, &konfiguracja.Zaktualizowano)
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

func zlozPoziomy(poziomy []PoziomPamieci) string {
	nazwy := make([]string, 0, len(poziomy))
	for _, poziom := range poziomy {
		if poziom != "" {
			nazwy = append(nazwy, string(poziom))
		}
	}
	return strings.Join(nazwy, ",")
}

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
