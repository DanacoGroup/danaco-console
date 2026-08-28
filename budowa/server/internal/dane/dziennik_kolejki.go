// Odpowiedzialność pliku: dziennik akcji silnika kolejek (tabela
// `log_akcji_kolejki`). Każda akcja, każdy obieg naprawczy i każde przejście
// stanu zostawia wpis.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WpisDziennika to wiersz tabeli `log_akcji_kolejki`, niosący jedną odnotowaną zmianę stanu danej kolejki.
type WpisDziennika struct {
	ID          int64
	KolejkaID   int64
	PozycjaID   *int64
	Akcja       string
	StanPrzed   *string
	StanPo      *string
	NumerObiegu int
	Szczegoly   *string
	Utworzono   string
}

const (
	kolumnyDziennika = `id, kolejka_id, pozycja_kolejki_id, akcja, stan_przed, stan_po,
	                    numer_obiegu, szczegoly, utworzono`

	wstawWpisDziennika = `INSERT INTO log_akcji_kolejki
	                      (kolejka_id, pozycja_kolejki_id, akcja, stan_przed, stan_po,
	                       numer_obiegu, szczegoly)
	                      VALUES (?, ?, ?, ?, ?, ?, ?)`

	polozeniePozycjiKolejki = `SELECT kolejka_id, licznik_obiegow FROM pozycja_kolejki WHERE id = ?`

	listaDziennikaKolejki = `SELECT ` + kolumnyDziennika + ` FROM log_akcji_kolejki
	                         WHERE kolejka_id = ?
	                         ORDER BY id
	                         LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`
)

// dopiszWpisDziennika zapisuje akcję w tej samej transakcji, co zmiana stanu —
// dziennik nie może rozejść się ze stanem kolejki.
func dopiszWpisDziennika(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	wpis WpisDziennika) error {

	polecenie, err := z.wTransakcji(ctx, transakcja, wstawWpisDziennika)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wpis.KolejkaID, liczbaDoKolumny(wpis.PozycjaID),
		wpis.Akcja, tekstDoKolumny(wpis.StanPrzed), tekstDoKolumny(wpis.StanPo),
		wpis.NumerObiegu, tekstDoKolumny(wpis.Szczegoly))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać akcji %q kolejki %d: %w",
			wpis.Akcja, wpis.KolejkaID, err)
	}
	return nil
}

// polozeniePozycji zwraca kolejkę i numer obiegu zlecenia — dziennik zapisuje
// jedno i drugie przy każdej zmianie.
func polozeniePozycji(ctx context.Context, z *zapytania, transakcja *sql.Tx,
	pozycjaID int64) (int64, int, error) {

	polecenie, err := z.wTransakcji(ctx, transakcja, polozeniePozycjiKolejki)
	if err != nil {
		return 0, 0, err
	}
	var kolejkaID int64
	var obieg int
	err = polecenie.QueryRowContext(ctx, pozycjaID).Scan(&kolejkaID, &obieg)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, fmt.Errorf("dane: pozycja kolejki %d nie istnieje", pozycjaID)
	}
	if err != nil {
		return 0, 0, fmt.Errorf("dane: nie można odczytać pozycji kolejki %d: %w", pozycjaID, err)
	}
	return kolejkaID, obieg, nil
}

// Dziennik zwraca wpisy kolejki w porządku chronologicznym malejącym; limit zero oznacza cały dziennik zdarzeń.
func (r *repozytoriumKolejek) Dziennik(ctx context.Context, kolejkaID int64,
	limit int) ([]WpisDziennika, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaDziennikaKolejki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kolejkaID, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dziennika kolejki %d: %w", kolejkaID, err)
	}
	defer wiersze.Close()

	lista := []WpisDziennika{}
	for wiersze.Next() {
		var wpis WpisDziennika
		var pozycja sql.NullInt64
		var przed, po, szczegoly sql.NullString
		err := wiersze.Scan(&wpis.ID, &wpis.KolejkaID, &pozycja, &wpis.Akcja, &przed, &po,
			&wpis.NumerObiegu, &szczegoly, &wpis.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wpis dziennika kolejki %d: %w", kolejkaID, err)
		}
		wpis.PozycjaID = liczbaZKolumny(pozycja)
		wpis.StanPrzed = tekstZKolumny(przed)
		wpis.StanPo = tekstZKolumny(po)
		wpis.Szczegoly = tekstZKolumny(szczegoly)
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt dziennika kolejki %d: %w", kolejkaID, err)
	}
	return lista, nil
}
