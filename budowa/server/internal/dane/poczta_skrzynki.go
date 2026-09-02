// Odpowiedzialność pliku: trwałość skrzynek Operatora, czyli wierszy tabeli skrzynka_pocztowa, wraz z protokołem, źródłem nastaw, trybem szyfrowania i domyślnością.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SkrzynkaOperatora to wiersz tabeli skrzynka_pocztowa widziany oczami rodziny mail.*, z polami protokołu, źródła nastaw i domyślności.
type SkrzynkaOperatora struct {
	ID               int64
	Kod              string
	Adres            string
	NazwaWyswietlana *string
	Protokol         string
	Zrodlo           string
	HostOdbioru      string
	PortOdbioru      int
	HostWysylki      string
	PortWysylki      int
	Uzytkownik       string
	HasloOdwolanie   *string
	SzyfrujOdbior    bool
	SzyfrujWysylke   bool
	WeryfikujTLS     bool
	Domyslna         bool
}

// RepozytoriumSkrzynek jest kontraktem trwałości skrzynek Operatora: podpięcia, wykazu, domyślności i śladu wysyłki.
type RepozytoriumSkrzynek interface {
	// Skrzynki oddaje wszystkie podpięte skrzynki, domyślną na początku.
	Skrzynki(ctx context.Context) ([]SkrzynkaOperatora, error)
	// Skrzynka odnajduje skrzynkę po kodzie (ErrBrakWiersza, gdy jej nie ma).
	Skrzynka(ctx context.Context, kod string) (SkrzynkaOperatora, error)
	// SkrzynkaDomyslna oddaje skrzynkę domyślną, a przy jej braku jedyną podpiętą.
	SkrzynkaDomyslna(ctx context.Context) (SkrzynkaOperatora, error)
	// Zapisz zakłada skrzynkę albo nadpisuje zastaną po kodzie i oddaje wiersz po zapisie.
	Zapisz(ctx context.Context, s SkrzynkaOperatora) (SkrzynkaOperatora, error)
	// Usun odpina skrzynkę. Fałsz znaczy „nie było czego odpinać”, nie błąd.
	Usun(ctx context.Context, kod string) (bool, error)
	// ZapiszSladWysylki utrwala fakt nadania listu — także nieudanego.
	ZapiszSladWysylki(ctx context.Context, s SladWysylki) error
}

const (
	kolumnySkrzynkiOperatora = `id, kod, adres, nazwa_wyswietlana, protokol, zrodlo,
	                            host_odbioru, port_odbioru, host_wysylki, port_wysylki,
	                            uzytkownik, haslo_odwolanie, szyfruj_odbior,
	                            szyfruj_wysylke, tls_weryfikacja, domyslna`

	// Domyślna idzie pierwsza, reszta po kodzie — wykaz ma ten sam porządek przy
	// każdym odczycie, a ta, którą rdzeń weźmie bez wskazania, stoi na wierzchu.
	wykazSkrzynekOperatora = `SELECT ` + kolumnySkrzynkiOperatora + `
	                          FROM skrzynka_pocztowa WHERE ` + WarunekKonta + `
	                          ORDER BY domyslna DESC, kod`

	skrzynkaOperatoraPoKodzie = `SELECT ` + kolumnySkrzynkiOperatora + `
	                             FROM skrzynka_pocztowa
	                             WHERE kod = ? AND ` + WarunekKonta

	// Bez wskazanej domyślnej pierwszeństwo ma pozycja z tego samego porządku, co
	// wykaz (`LIMIT 1`). Przy jednej podpiętej skrzynce jest to ona sama, więc nie
	// trzeba jej osobno oznaczać.
	skrzynkaOperatoraDomyslna = `SELECT ` + kolumnySkrzynkiOperatora + `
	                             FROM skrzynka_pocztowa WHERE ` + WarunekKonta + `
	                             ORDER BY domyslna DESC, kod LIMIT 1`

	// Więz UNIQUE na `kod` obejmuje całą tabelę: człon DO UPDATE bez zawężenia nadpisałby skrzynkę konta cudzego wraz z jej hasłem.
	zapiszSkrzynkeOperatora = `INSERT INTO skrzynka_pocztowa
	    (kod, adres, nazwa_wyswietlana, protokol, zrodlo, host_odbioru, port_odbioru,
	     host_wysylki, port_wysylki, uzytkownik, haslo_odwolanie, szyfruj_odbior,
	     szyfruj_wysylke, tls_weryfikacja, domyslna, konto_id)
	    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	    ON CONFLICT(kod) DO UPDATE SET
	        adres = excluded.adres,
	        nazwa_wyswietlana = excluded.nazwa_wyswietlana,
	        protokol = excluded.protokol,
	        zrodlo = excluded.zrodlo,
	        host_odbioru = excluded.host_odbioru,
	        port_odbioru = excluded.port_odbioru,
	        host_wysylki = excluded.host_wysylki,
	        port_wysylki = excluded.port_wysylki,
	        uzytkownik = excluded.uzytkownik,
	        haslo_odwolanie = excluded.haslo_odwolanie,
	        szyfruj_odbior = excluded.szyfruj_odbior,
	        szyfruj_wysylke = excluded.szyfruj_wysylke,
	        tls_weryfikacja = excluded.tls_weryfikacja,
	        domyslna = excluded.domyslna
	    WHERE ` + WarunekKonta

	// Zdjęcie domyślności z pozostałych. Indeks częściowy schematu dopuszcza jedną
	// domyślną skrzynkę, więc nadanie domyślności nowej musi zdjąć ją starej
	// w tej samej transakcji. Bez tego INSERT rozbiłby się o indeks.
	zdejmijDomyslnoscSkrzynek = `UPDATE skrzynka_pocztowa SET domyslna = 0
	                             WHERE kod <> ? AND ` + WarunekKonta

	usunSkrzynkeOperatora = `DELETE FROM skrzynka_pocztowa WHERE kod = ? AND ` + WarunekKonta
)

// repozytoriumSkrzynek stoi na wspólnej pamięci zapytań zestawu i na uchwycie
// bazy — ten drugi jest potrzebny do transakcji przy nadawaniu domyślności.
type repozytoriumSkrzynek struct {
	zapytania *zapytania
	db        *sql.DB
}

// SkrzynkiOperatora oddaje repozytorium skrzynek nad tą samą bazą, co reszta
// zestawu. Jest metodą, a nie polem struktury — wzorem `SekcjePaneli`
// i `RoleOkien` — bo repozytorium nie trzyma stanu poza wskaźnikami na wspólną
// pamięć zapytań i uchwyt bazy.
func (z *Zestaw) SkrzynkiOperatora() RepozytoriumSkrzynek {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return &repozytoriumSkrzynek{zapytania: z.zapytania, db: z.zapytania.db}
}

func (r *repozytoriumSkrzynek) Skrzynki(ctx context.Context) ([]SkrzynkaOperatora, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wykazSkrzynekOperatora)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać skrzynek Operatora: %w", err)
	}
	defer wiersze.Close()
	lista := []SkrzynkaOperatora{}
	for wiersze.Next() {
		skrzynka, err := odczytajSkrzynkeOperatora(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, skrzynka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt skrzynek Operatora: %w", err)
	}
	return lista, nil
}

func (r *repozytoriumSkrzynek) Skrzynka(ctx context.Context, kod string) (SkrzynkaOperatora, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, skrzynkaOperatoraPoKodzie)
	if err != nil {
		return SkrzynkaOperatora{}, err
	}
	skrzynka, err := odczytajSkrzynkeOperatora(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return SkrzynkaOperatora{}, ErrBrakWiersza
	}
	return skrzynka, err
}

func (r *repozytoriumSkrzynek) SkrzynkaDomyslna(ctx context.Context) (SkrzynkaOperatora, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, skrzynkaOperatoraDomyslna)
	if err != nil {
		return SkrzynkaOperatora{}, err
	}
	skrzynka, err := odczytajSkrzynkeOperatora(polecenie.QueryRowContext(ctx, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return SkrzynkaOperatora{}, ErrBrakWiersza
	}
	return skrzynka, err
}

// Zapisz zakłada albo nadpisuje skrzynkę. Nadanie domyślności i zdjęcie jej
// z pozostałych idą jedną transakcją, bo opisują jedną zmianę; przerwane
// w połowie zostawiłyby dwie domyślne skrzynki albo żadnej.
func (r *repozytoriumSkrzynek) Zapisz(ctx context.Context, s SkrzynkaOperatora) (SkrzynkaOperatora, error) {
	err := wTransakcji(ctx, r.db, func(tx *sql.Tx) error {
		if s.Domyslna {
			if _, err := tx.ExecContext(ctx, zdejmijDomyslnoscSkrzynek, s.Kod, KontoOperatora(ctx)); err != nil {
				return fmt.Errorf("dane: nie można zdjąć domyślności z pozostałych skrzynek: %w", err)
			}
		}
		wynik, err := tx.ExecContext(ctx, zapiszSkrzynkeOperatora,
			s.Kod, s.Adres, tekstDoKolumny(s.NazwaWyswietlana), s.Protokol, s.Zrodlo,
			s.HostOdbioru, s.PortOdbioru, s.HostWysylki, s.PortWysylki, s.Uzytkownik,
			tekstDoKolumny(s.HasloOdwolanie), liczbaLogiczna(s.SzyfrujOdbior),
			liczbaLogiczna(s.SzyfrujWysylke), liczbaLogiczna(s.WeryfikujTLS),
			liczbaLogiczna(s.Domyslna), KontoOperatora(ctx), KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać skrzynki %q: %w", s.Kod, err)
		}
		return sprawdzTrafienieZapisu(wynik, "skrzynka", s.Kod)
	})
	if err != nil {
		return SkrzynkaOperatora{}, err
	}
	// Wraca wiersz odczytany, nie przysłany: klucz nadaje baza, kolumny mogły dopowiedzieć wartości.
	return r.Skrzynka(ctx, s.Kod)
}

func (r *repozytoriumSkrzynek) Usun(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunSkrzynkeOperatora)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można odpiąć skrzynki %q: %w", kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany wynik odpięcia skrzynki %q: %w", kod, err)
	}
	return usuniete > 0, nil
}

// odczytajSkrzynkeOperatora czyta jeden wiersz. Przyjmuje `skaner` (interfejs
// z `dane`), więc obsługuje i pojedynczy odczyt, i pętlę wykazu.
func odczytajSkrzynkeOperatora(w skaner) (SkrzynkaOperatora, error) {
	var (
		s                SkrzynkaOperatora
		nazwaWyswietlana sql.NullString
		hasloOdwolanie   sql.NullString
		szyfrujOdbior    int
		szyfrujWysylke   int
		weryfikujTLS     int
		domyslna         int
	)
	err := w.Scan(&s.ID, &s.Kod, &s.Adres, &nazwaWyswietlana, &s.Protokol, &s.Zrodlo,
		&s.HostOdbioru, &s.PortOdbioru, &s.HostWysylki, &s.PortWysylki,
		&s.Uzytkownik, &hasloOdwolanie, &szyfrujOdbior, &szyfrujWysylke,
		&weryfikujTLS, &domyslna)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SkrzynkaOperatora{}, err
		}
		return SkrzynkaOperatora{}, fmt.Errorf("dane: nieczytelny wiersz skrzynki Operatora: %w", err)
	}
	s.NazwaWyswietlana = tekstZKolumny(nazwaWyswietlana)
	s.HasloOdwolanie = tekstZKolumny(hasloOdwolanie)
	s.SzyfrujOdbior = szyfrujOdbior == 1
	s.SzyfrujWysylke = szyfrujWysylke == 1
	s.WeryfikujTLS = weryfikujTLS == 1
	s.Domyslna = domyslna == 1
	return s, nil
}
