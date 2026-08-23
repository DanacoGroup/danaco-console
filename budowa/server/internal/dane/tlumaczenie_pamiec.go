// Odpowiedzialność pliku: pamięć tłumaczeń modułu Translate w kształcie
// kontraktu `TranslationMemoryEntry` (tabela `pamiec_tlumaczen`, migracja 160)
// oraz polityka pamięci okna (`polityka_pamieci_okna`).
//
// Plik `slownik_pamiec.go` obsługuje jedną, wąską drogę tej samej tabeli: zapis
// pary zdjętej z zatwierdzonego panelu i dopasowanie przybliżone dla
// `memory.suggest`. Ten plik odpowiada za pamięć jako byt Operatora — wykaz,
// zapis wprost, usunięcie, utrzymanie i wymianę z plikiem. Dwa pliki, jedna
// tabela, dwie różne odpowiedzialności.
//
// Zapytania składane są tu wprost na `*sql.DB`, nie przez pamięć przygotowanych
// poleceń: wykaz pamięci ma cztery nieobowiązkowe zawężenia i limit, więc treść
// zapytania zależy od żądania. Pamięć przygotowanych poleceń trzymałaby
// kilkanaście wariantów jednego odczytu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// WpisPamieciTlumaczenPelny to wiersz `pamiec_tlumaczen` ze wszystkimi polami,
// których żąda kontrakt. `PanelKod` doczytywany jest złączeniem — kontrakt
// oddaje identyfikator zewnętrzny panelu, nie klucz wewnętrzny.
type WpisPamieciTlumaczenPelny struct {
	ID                int64
	Kod               string
	PanelID           *int64
	PanelKod          *string
	Jezyk             string
	SegmentZrodlowy   string
	SegmentDocelowy   string
	Projekt           *string
	Klient            *string
	Autor             *string
	KontekstPoprzedni *string
	KontekstNastepny  *string
	Zasieg            string
	Utworzono         int64
	Zaktualizowano    int64
}

// FiltrPamieciTlumaczen zawęża wykaz pamięci. Pola puste nie zawężają niczego —
// puste zawężenie jest brakiem zawężenia, nie zawężeniem do pustki.
type FiltrPamieciTlumaczen struct {
	Jezyk   string
	Fraza   string
	Projekt string
	Zasieg  string
	Limit   int
	Offset  int
}

// PolitykaPamieciOkna to wiersz `polityka_pamieci_okna` — nastawa, wedle której
// okno sięga do pamięci: zasięg par, próg dopasowania w procentach, wymóg
// zgodności kontekstu i zgoda na tłumaczenie wstępne.
type PolitykaPamieciOkna struct {
	OknoID               int64
	Zasieg               string
	Prog                 int64
	DopasowanieKontekstu bool
	WstepneTlumaczenie   bool
	Zaktualizowano       int64
}

const kolumnyWpisuPamieciTlumaczen = `w.id, w.identyfikator_zewnetrzny, w.panel_id, p.identyfikator_zewnetrzny,
	w.jezyk, w.segment_zrodlowy, w.segment_docelowy, w.projekt, w.klient, w.autor,
	w.kontekst_poprzedni, w.kontekst_nastepny, w.zasieg, w.utworzono, w.zaktualizowano`

const zrodloWpisuPamieciTlumaczen = ` FROM pamiec_tlumaczen w
	LEFT JOIN panel_tlumaczenia p ON p.id = w.panel_id`

// warunkiPamieci składa część WHERE wraz z argumentami wedle zawężeń filtru.
func warunkiPamieci(filtr FiltrPamieciTlumaczen) (string, []any) {
	warunki := []string{}
	argumenty := []any{}
	if strings.TrimSpace(filtr.Jezyk) != "" {
		warunki = append(warunki, "w.jezyk = ?")
		argumenty = append(argumenty, filtr.Jezyk)
	}
	if strings.TrimSpace(filtr.Projekt) != "" {
		warunki = append(warunki, "w.projekt = ?")
		argumenty = append(argumenty, filtr.Projekt)
	}
	if strings.TrimSpace(filtr.Zasieg) != "" {
		warunki = append(warunki, "w.zasieg = ?")
		argumenty = append(argumenty, filtr.Zasieg)
	}
	if fraza := strings.TrimSpace(filtr.Fraza); fraza != "" {
		// Dopasowanie w dowolnym miejscu obu segmentów — tak stanowi kontrakt
		// (`query`: „dopasowanie w dowolnym miejscu").
		warunki = append(warunki, "(w.segment_zrodlowy LIKE ? OR w.segment_docelowy LIKE ?)")
		wzorzec := "%" + fraza + "%"
		argumenty = append(argumenty, wzorzec, wzorzec)
	}
	if len(warunki) == 0 {
		return "", argumenty
	}
	return " WHERE " + strings.Join(warunki, " AND "), argumenty
}

// WpisyPamieci oddaje wykaz par wraz z liczbą wszystkich pasujących. Liczba
// całkowita liczona jest bez limitu — inaczej klient nie odróżniłby „to
// wszystko" od „to pierwsza strona".
func (r *repozytoriumTlumaczen) WpisyPamieci(ctx context.Context,
	filtr FiltrPamieciTlumaczen) ([]WpisPamieciTlumaczenPelny, int, error) {

	warunek, argumenty := warunkiPamieci(filtr)

	var razem int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pamiec_tlumaczen w`+warunek, argumenty...).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć pamięci tłumaczeń: %w", err)
	}

	zapytanie := `SELECT ` + kolumnyWpisuPamieciTlumaczen + zrodloWpisuPamieciTlumaczen + warunek +
		` ORDER BY w.zaktualizowano DESC, w.id DESC`
	if filtr.Limit > 0 {
		zapytanie += " LIMIT ?"
		argumenty = append(argumenty, filtr.Limit)
		if filtr.Offset > 0 {
			zapytanie += " OFFSET ?"
			argumenty = append(argumenty, filtr.Offset)
		}
	}

	wiersze, err := r.db.QueryContext(ctx, zapytanie, argumenty...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać pamięci tłumaczeń: %w", err)
	}
	defer wiersze.Close()

	lista := []WpisPamieciTlumaczenPelny{}
	for wiersze.Next() {
		wpis, err := odczytajWpisPamieciPelny(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz pamięci tłumaczeń: %w", err)
		}
		lista = append(lista, wpis)
	}
	return lista, razem, wiersze.Err()
}

// odczytajWpisPamieciPelny składa strukturę z jednego wiersza wyniku.
func odczytajWpisPamieciPelny(wiersz skaner) (WpisPamieciTlumaczenPelny, error) {
	var wpis WpisPamieciTlumaczenPelny
	var panelID sql.NullInt64
	var panelKod, projekt, klient, autor, poprzedni, nastepny sql.NullString
	err := wiersz.Scan(&wpis.ID, &wpis.Kod, &panelID, &panelKod, &wpis.Jezyk,
		&wpis.SegmentZrodlowy, &wpis.SegmentDocelowy, &projekt, &klient, &autor,
		&poprzedni, &nastepny, &wpis.Zasieg, &wpis.Utworzono, &wpis.Zaktualizowano)
	if err != nil {
		return WpisPamieciTlumaczenPelny{}, err
	}
	wpis.PanelID = liczbaZKolumny(panelID)
	wpis.PanelKod = tekstZKolumny(panelKod)
	wpis.Projekt = tekstZKolumny(projekt)
	wpis.Klient = tekstZKolumny(klient)
	wpis.Autor = tekstZKolumny(autor)
	wpis.KontekstPoprzedni = tekstZKolumny(poprzedni)
	wpis.KontekstNastepny = tekstZKolumny(nastepny)
	return wpis, nil
}

// WpisPamieci oddaje jedną parę po kodzie zewnętrznym.
func (r *repozytoriumTlumaczen) WpisPamieci(ctx context.Context,
	kod string) (WpisPamieciTlumaczenPelny, error) {

	wiersz := r.db.QueryRowContext(ctx,
		`SELECT `+kolumnyWpisuPamieciTlumaczen+zrodloWpisuPamieciTlumaczen+` WHERE w.identyfikator_zewnetrzny = ?`, kod)
	wpis, err := odczytajWpisPamieciPelny(wiersz)
	if errors.Is(err, sql.ErrNoRows) {
		return WpisPamieciTlumaczenPelny{}, ErrBrakWiersza
	}
	if err != nil {
		return WpisPamieciTlumaczenPelny{}, fmt.Errorf("dane: nieczytelny wpis pamięci %q: %w", kod, err)
	}
	return wpis, nil
}

// ZapiszWpisPamieci zakłada parę albo nadpisuje zastaną po kodzie zewnętrznym.
// Chwila założenia nie przesuwa się przy nadpisaniu — para wniesiona rok temu
// i poprawiona dziś zostaje parą sprzed roku.
func (r *repozytoriumTlumaczen) ZapiszWpisPamieci(ctx context.Context,
	wpis WpisPamieciTlumaczenPelny) (WpisPamieciTlumaczenPelny, error) {

	if strings.TrimSpace(wpis.Kod) == "" {
		return WpisPamieciTlumaczenPelny{}, fmt.Errorf("dane: wpis pamięci bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	utworzono := wpis.Utworzono
	if utworzono == 0 {
		utworzono = teraz
	}
	zasieg := wpis.Zasieg
	if zasieg == "" {
		zasieg = "card"
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO pamiec_tlumaczen
		(identyfikator_zewnetrzny, panel_id, jezyk, segment_zrodlowy, segment_docelowy,
		 projekt, klient, autor, kontekst_poprzedni, kontekst_nastepny, zasieg,
		 utworzono, zaktualizowano)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
			panel_id = excluded.panel_id,
			jezyk = excluded.jezyk,
			segment_zrodlowy = excluded.segment_zrodlowy,
			segment_docelowy = excluded.segment_docelowy,
			projekt = excluded.projekt,
			klient = excluded.klient,
			autor = excluded.autor,
			kontekst_poprzedni = excluded.kontekst_poprzedni,
			kontekst_nastepny = excluded.kontekst_nastepny,
			zasieg = excluded.zasieg,
			zaktualizowano = excluded.zaktualizowano`,
		wpis.Kod, liczbaDoKolumny(wpis.PanelID), wpis.Jezyk, wpis.SegmentZrodlowy,
		wpis.SegmentDocelowy, tekstDoKolumny(wpis.Projekt), tekstDoKolumny(wpis.Klient),
		tekstDoKolumny(wpis.Autor), tekstDoKolumny(wpis.KontekstPoprzedni),
		tekstDoKolumny(wpis.KontekstNastepny), zasieg, utworzono, teraz)
	if err != nil {
		return WpisPamieciTlumaczenPelny{}, fmt.Errorf("dane: nie można zapisać wpisu pamięci %q: %w", wpis.Kod, err)
	}
	return r.WpisPamieci(ctx, wpis.Kod)
}

// UsunWpisPamieci kasuje parę po kodzie. Wartość logiczna mówi, czy coś realnie
// zeszło — kasowanie pary, której nie ma, nie jest usterką, ale nie jest też
// usunięciem.
func (r *repozytoriumTlumaczen) UsunWpisPamieci(ctx context.Context, kod string) (bool, error) {
	wynik, err := r.db.ExecContext(ctx,
		`DELETE FROM pamiec_tlumaczen WHERE identyfikator_zewnetrzny = ?`, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć wpisu pamięci %q: %w", kod, err)
	}
	zeszlo, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia wpisu pamięci: %w", err)
	}
	return zeszlo > 0, nil
}

// UsunWpisyPamieci kasuje wskazane pary w jednej transakcji i oddaje liczbę
// wierszy, które realnie zeszły. Utrzymanie pamięci (`memory.maintain`) kasuje
// dziesiątki par naraz — pojedyncze wywołania byłyby dziesiątkami transakcji.
func (r *repozytoriumTlumaczen) UsunWpisyPamieci(ctx context.Context, kody []string) (int, error) {
	if len(kody) == 0 {
		return 0, nil
	}
	zeszlo := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := transakcja.PrepareContext(ctx,
			`DELETE FROM pamiec_tlumaczen WHERE identyfikator_zewnetrzny = ?`)
		if err != nil {
			return fmt.Errorf("dane: nie można przygotować usunięcia par pamięci: %w", err)
		}
		defer polecenie.Close()
		for _, kod := range kody {
			wynik, err := polecenie.ExecContext(ctx, kod)
			if err != nil {
				return fmt.Errorf("dane: nie można usunąć pary pamięci %q: %w", kod, err)
			}
			if liczba, err := wynik.RowsAffected(); err == nil {
				zeszlo += int(liczba)
			}
		}
		return nil
	})
	return zeszlo, err
}

// PolitykaPamieci oddaje politykę okna. Brak wiersza wraca jako ErrBrakWiersza:
// „polityki nie ustawiono" to nie to samo, co polityka wyzerowana.
func (r *repozytoriumTlumaczen) PolitykaPamieci(ctx context.Context,
	oknoID int64) (PolitykaPamieciOkna, error) {

	var polityka PolitykaPamieciOkna
	var kontekst, wstepne int64
	err := r.db.QueryRowContext(ctx,
		`SELECT okno_id, zasieg, prog, dopasowanie_kontekstu, wstepne_tlumaczenie, zaktualizowano
		   FROM polityka_pamieci_okna WHERE okno_id = ?`, oknoID).
		Scan(&polityka.OknoID, &polityka.Zasieg, &polityka.Prog, &kontekst, &wstepne,
			&polityka.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return PolitykaPamieciOkna{}, ErrBrakWiersza
	}
	if err != nil {
		return PolitykaPamieciOkna{}, fmt.Errorf("dane: nieczytelna polityka pamięci okna %d: %w", oknoID, err)
	}
	polityka.DopasowanieKontekstu = kontekst == 1
	polityka.WstepneTlumaczenie = wstepne == 1
	return polityka, nil
}

// ZapiszPolitykePamieci zakłada politykę okna albo nadpisuje zastaną.
func (r *repozytoriumTlumaczen) ZapiszPolitykePamieci(ctx context.Context,
	polityka PolitykaPamieciOkna) (PolitykaPamieciOkna, error) {

	teraz := time.Now().UnixMilli()
	_, err := r.db.ExecContext(ctx, `INSERT INTO polityka_pamieci_okna
		(okno_id, zasieg, prog, dopasowanie_kontekstu, wstepne_tlumaczenie, zaktualizowano)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(okno_id) DO UPDATE SET
			zasieg = excluded.zasieg,
			prog = excluded.prog,
			dopasowanie_kontekstu = excluded.dopasowanie_kontekstu,
			wstepne_tlumaczenie = excluded.wstepne_tlumaczenie,
			zaktualizowano = excluded.zaktualizowano`,
		polityka.OknoID, polityka.Zasieg, polityka.Prog,
		wartoscLogicznaDoKolumny(polityka.DopasowanieKontekstu),
		wartoscLogicznaDoKolumny(polityka.WstepneTlumaczenie), teraz)
	if err != nil {
		return PolitykaPamieciOkna{}, fmt.Errorf("dane: nie można zapisać polityki pamięci okna %d: %w",
			polityka.OknoID, err)
	}
	return r.PolitykaPamieci(ctx, polityka.OknoID)
}

// wartoscLogicznaDoKolumny przekłada `bool` na kolumnę INTEGER z warunkiem
// CHECK na 0 albo 1 — schemat modułu Translate nie zna typu logicznego.
func wartoscLogicznaDoKolumny(wartosc bool) int64 {
	if wartosc {
		return 1
	}
	return 0
}
