// Odpowiedzialność pliku: historia schowka Operatora, rodzina komend clipboard.*; rdzeń nie czyta schowka maszyny, daje tylko trwałość.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// WpisSchowka to wiersz tabeli wpis_schowka, niosący treść, postać i przypięcie jednego wpisu historii.
type WpisSchowka struct {
	Kod           string
	Rodzaj        string
	Tresc         string
	Odcisk        string
	Zajawka       *string
	RozmiarBajtow int64
	Przypiety     bool
	Wrazliwy      bool
	OknoZrodlowe  *string
	Utworzono     int64
	Uzyto         *int64
	// PostacJSON niesie postać fragmentu odłożonego do schowka, nie wchodzi do rachunku odcisku treści.
	PostacJSON *string
}

// SitoSchowka zawęża odczyt całej historii schowka do wpisów spełniających wskazane kryteria tego filtru.
type SitoSchowka struct {
	Fraza        string
	Rodzaj       string
	TylkoPrzypie bool
	Granica      int
	Przesuniecie int
}

// RepozytoriumSchowka jest kontraktem historii schowka Operatora wraz z przypinaniem i usuwaniem wpisów.
type RepozytoriumSchowka interface {
	DopiszWpisSchowka(ctx context.Context, wpis WpisSchowka) (WpisSchowka, bool, error)
	WpisySchowka(ctx context.Context, sito SitoSchowka) ([]WpisSchowka, int, error)
	PrzypnijWpisSchowka(ctx context.Context, kod string, przypiety bool) (WpisSchowka, error)
	UsunWpisySchowka(ctx context.Context, kod string, zPrzypietymi bool) (int, error)
}

const (
	kolumnyWpisuSchowka = `identyfikator_zewnetrzny, rodzaj, tresc, odcisk, zajawka,
	                       rozmiar_bajtow, przypiety, wrazliwy, okno_zrodlowe, utworzono, uzyto,
	                       postac_json`

	// Postać odłożona ponownie nadpisuje zastaną, a postać niepodana jej nie zabiera przy tym samym zapisie.
	// Więz UNIQUE kolumny `odcisk` obejmuje całą tabelę, więc ta sama treść odłożona przez dwa konta trafia
	// w jeden wiersz; WarunekKonta przy DO UPDATE odcina zapis w wiersz cudzy.
	wstawWpisSchowka = `INSERT INTO wpis_schowka (` + kolumnyWpisuSchowka + `, konto_id)
	                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                    ON CONFLICT(odcisk) DO UPDATE SET
	                        utworzono = excluded.utworzono,
	                        okno_zrodlowe = COALESCE(excluded.okno_zrodlowe, okno_zrodlowe),
	                        wrazliwy = MAX(wrazliwy, excluded.wrazliwy),
	                        postac_json = COALESCE(excluded.postac_json, postac_json)
	                    WHERE ` + WarunekKonta

	pobierzWpisSchowkaPoOdcisku = `SELECT ` + kolumnyWpisuSchowka +
		` FROM wpis_schowka WHERE odcisk = ? AND ` + WarunekKonta

	pobierzWpisSchowka = `SELECT ` + kolumnyWpisuSchowka +
		` FROM wpis_schowka WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	przypnijWpisSchowka = `UPDATE wpis_schowka SET przypiety = ?
	                       WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunWpisSchowka = `DELETE FROM wpis_schowka
	                   WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

type repozytoriumSchowka struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumSchowka zakłada historię schowka nad wspólną bazą zestawu, gotową do zapisu i odczytu.
func noweRepozytoriumSchowka(z *zapytania, db *sql.DB) *repozytoriumSchowka {
	return &repozytoriumSchowka{zapytania: z, db: db}
}

// DopiszWpisSchowka dopisuje treść albo podnosi wpis zastany. Drugi zwracany
// wynik mówi, czy treść była już w historii — powtórzenie nie jest błędem.
func (r *repozytoriumSchowka) DopiszWpisSchowka(ctx context.Context,
	wpis WpisSchowka) (WpisSchowka, bool, error) {

	if strings.TrimSpace(wpis.Odcisk) == "" {
		return WpisSchowka{}, false, fmt.Errorf("dane: wpis schowka bez odcisku treści")
	}
	bylo := false
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzWpisSchowkaPoOdcisku)
		if err != nil {
			return err
		}
		zastany, err := odczytajWpisSchowka(
			odczyt.QueryRowContext(ctx, wpis.Odcisk, KontoOperatora(ctx)))
		switch {
		case errors.Is(err, sql.ErrNoRows):
		case err != nil:
			return fmt.Errorf("dane: nie można odczytać wpisu schowka: %w", err)
		default:
			bylo = true
			wpis.Kod = zastany.Kod
		}

		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWpisSchowka)
		if err != nil {
			return err
		}
		wynik, err := zapis.ExecContext(ctx, wpis.Kod, wpis.Rodzaj, wpis.Tresc, wpis.Odcisk,
			tekstDoKolumny(wpis.Zajawka), wpis.RozmiarBajtow, liczbaLogiczna(wpis.Przypiety),
			liczbaLogiczna(wpis.Wrazliwy), tekstDoKolumny(wpis.OknoZrodlowe), wpis.Utworzono,
			liczbaDoKolumny(wpis.Uzyto), tekstDoKolumny(wpis.PostacJSON), KontoOperatora(ctx),
			KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać wpisu schowka: %w", err)
		}
		return sprawdzTrafienieZapisu(wynik, "wpis schowka o odcisku", wpis.Odcisk)
	})
	if err != nil {
		return WpisSchowka{}, false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisSchowkaPoOdcisku)
	if err != nil {
		return WpisSchowka{}, false, err
	}
	zapisany, err := odczytajWpisSchowka(
		polecenie.QueryRowContext(ctx, wpis.Odcisk, KontoOperatora(ctx)))
	if err != nil {
		return WpisSchowka{}, false, fmt.Errorf("dane: nieczytelny wiersz wpisu schowka: %w", err)
	}
	return zapisany, bylo, nil
}

// WpisySchowka zwraca historię schowka: przypięte wpisy na czele, potem pozostałe od najnowszego wpisu.
func (r *repozytoriumSchowka) WpisySchowka(ctx context.Context,
	sito SitoSchowka) ([]WpisSchowka, int, error) {

	warunki := []string{WarunekKonta}
	argumenty := []any{KontoOperatora(ctx)}
	if fraza := strings.TrimSpace(sito.Fraza); fraza != "" {
		warunki = append(warunki, "tresc LIKE ?")
		argumenty = append(argumenty, "%"+fraza+"%")
	}
	if strings.TrimSpace(sito.Rodzaj) != "" {
		warunki = append(warunki, "rodzaj = ?")
		argumenty = append(argumenty, sito.Rodzaj)
	}
	if sito.TylkoPrzypie {
		warunki = append(warunki, "przypiety = 1")
	}
	gdzie := strings.Join(warunki, " AND ")

	wszystkich := 0
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wpis_schowka WHERE `+gdzie,
		argumenty...).Scan(&wszystkich); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wpisów schowka: %w", err)
	}

	tekst := `SELECT ` + kolumnyWpisuSchowka + ` FROM wpis_schowka WHERE ` + gdzie +
		` ORDER BY przypiety DESC, utworzono DESC, id DESC`
	if sito.Granica > 0 {
		tekst += fmt.Sprintf(" LIMIT %d", sito.Granica)
	} else if sito.Przesuniecie > 0 {
		tekst += " LIMIT -1"
	}
	if sito.Przesuniecie > 0 {
		tekst += fmt.Sprintf(" OFFSET %d", sito.Przesuniecie)
	}

	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać wpisów schowka: %w", err)
	}
	defer wiersze.Close()

	lista := []WpisSchowka{}
	for wiersze.Next() {
		wpis, err := odczytajWpisSchowka(wiersze)
		if err != nil {
			return nil, 0, err
		}
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt wpisów schowka: %w", err)
	}
	return lista, wszystkich, nil
}

// PrzypnijWpisSchowka przypina wskazany wpis historii schowka albo zdejmuje jego wcześniejsze przypięcie.
func (r *repozytoriumSchowka) PrzypnijWpisSchowka(ctx context.Context, kod string,
	przypiety bool) (WpisSchowka, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, przypnijWpisSchowka)
	if err != nil {
		return WpisSchowka{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, liczbaLogiczna(przypiety), kod, KontoOperatora(ctx))
	if err != nil {
		return WpisSchowka{}, fmt.Errorf("dane: nie można przypiąć wpisu schowka %q: %w", kod, err)
	}
	if zmienione, err := wynik.RowsAffected(); err == nil && zmienione == 0 {
		return WpisSchowka{}, fmt.Errorf("dane: wpis schowka %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzWpisSchowka)
	if err != nil {
		return WpisSchowka{}, err
	}
	return odczytajWpisSchowka(odczyt.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
}

// UsunWpisySchowka kasuje jeden wpis albo całą historię. Bez wskazania wpisu
// i bez zgody na przypięte kasowane są wyłącznie wpisy nieprzypięte —
// przypięcie jest właśnie oświadczeniem, że wpis ma przeżyć czyszczenie.
func (r *repozytoriumSchowka) UsunWpisySchowka(ctx context.Context, kod string,
	zPrzypietymi bool) (int, error) {

	if strings.TrimSpace(kod) != "" {
		polecenie, err := r.zapytania.przygotuj(ctx, usunWpisSchowka)
		if err != nil {
			return 0, err
		}
		wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
		if err != nil {
			return 0, fmt.Errorf("dane: nie można usunąć wpisu schowka %q: %w", kod, err)
		}
		usuniete, err := wynik.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("dane: nie można policzyć usuniętych wpisów schowka: %w", err)
		}
		if usuniete == 0 {
			return 0, fmt.Errorf("dane: wpis schowka %q nie istnieje: %w", kod, ErrBrakWiersza)
		}
		return int(usuniete), nil
	}

	tekst := `DELETE FROM wpis_schowka WHERE ` + WarunekKonta
	if !zPrzypietymi {
		tekst += ` AND przypiety = 0`
	}
	wynik, err := r.db.ExecContext(ctx, tekst, KontoOperatora(ctx))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można wyczyścić historii schowka: %w", err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć usuniętych wpisów schowka: %w", err)
	}
	return int(usuniete), nil
}

// odczytajWpisSchowka przekłada wiersz tabeli na wpis historii schowka Operatora, kolumna po kolumnie.
func odczytajWpisSchowka(s skaner) (WpisSchowka, error) {
	var wpis WpisSchowka
	var zajawka, okno, postac sql.NullString
	var uzyto sql.NullInt64
	var przypiety, wrazliwy int
	err := s.Scan(&wpis.Kod, &wpis.Rodzaj, &wpis.Tresc, &wpis.Odcisk, &zajawka,
		&wpis.RozmiarBajtow, &przypiety, &wrazliwy, &okno, &wpis.Utworzono, &uzyto, &postac)
	if err != nil {
		return WpisSchowka{}, err
	}
	wpis.Zajawka = tekstZKolumny(zajawka)
	wpis.Przypiety = przypiety == 1
	wpis.Wrazliwy = wrazliwy == 1
	wpis.OknoZrodlowe = tekstZKolumny(okno)
	wpis.Uzyto = liczbaZKolumny(uzyto)
	wpis.PostacJSON = tekstZKolumny(postac)
	return wpis, nil
}
