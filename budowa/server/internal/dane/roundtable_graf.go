// Odpowiedzialność pliku: graf argumentów debaty — węzły, krawędzie, katalog błędów
// logicznych i oznaczenia na węzłach.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type WezelDebaty struct {
	Kod       string
	Okno      string
	Wypowiedz string
	Uczestnik string
	Tura      string
	AktMowy   string
	Tresc     string
	Poparcie  int
	Kluczowy  bool
	Utworzono string
}

type KrawedzDebaty struct {
	Kod           string
	Okno          string
	WezelZrodlowy string
	WezelWskazany string
	Relacja       string
	Pewnosc       float64
}

type DefinicjaBleduDebaty struct {
	Kod       string
	Nazwa     string
	Opis      string
	Wykrywany bool
}

type OznaczenieBleduDebaty struct {
	Kod          string
	Okno         string
	Wezel        string
	KodBledu     string
	Nazwa        string
	Uzasadnienie string
	Pewnosc      float64
}

type RepozytoriumDebatyGrafu interface {
	ZastapGrafDebaty(ctx context.Context, okno, tura string,
		wezly []WezelDebaty, krawedzie []KrawedzDebaty) error
	WezelDebatyPoKodzie(ctx context.Context, kod string) (WezelDebaty, error)
	WezlyDebaty(ctx context.Context, okno, tura string, tylkoKluczowe bool) ([]WezelDebaty, error)
	OznaczWezelDebaty(ctx context.Context, kod string, kluczowy bool) error
	KrawedzieDebaty(ctx context.Context, okno string) ([]KrawedzDebaty, error)

	ZapiszOznaczenieBleduDebaty(ctx context.Context, oznaczenie OznaczenieBleduDebaty) error
	OznaczeniaBledowDebaty(ctx context.Context, okno string) ([]OznaczenieBleduDebaty, error)
	KatalogBledowDebaty(ctx context.Context, okno string) ([]DefinicjaBleduDebaty, error)
	UstawKatalogBledowDebaty(ctx context.Context, okno string, kody []string) error
}

const (
	kolumnyWezlaDebaty = `identyfikator_zewnetrzny, okno, wypowiedz, uczestnik, tura,
	                      akt_mowy, tresc, poparcie, kluczowy, utworzono`

	zapiszWezelDebaty = `INSERT INTO debata_wezel
	                     (identyfikator_zewnetrzny, okno, wypowiedz, uczestnik, tura,
	                      akt_mowy, tresc, poparcie, kluczowy, konto_id)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	usunWezlyDebaty = `DELETE FROM debata_wezel
	                   WHERE okno = ? AND (? = '' OR tura = ?) AND ` + WarunekKonta

	// Krawędź własnej kolumny konta nie ma: granica dochodzi przez węzeł `wezel_zrodlowy`.
	usunKrawedzieDebaty = `DELETE FROM debata_krawedz
	                       WHERE okno = ?
	                         AND EXISTS (SELECT 1 FROM debata_wezel
	                                      WHERE debata_wezel.identyfikator_zewnetrzny
	                                            = debata_krawedz.wezel_zrodlowy
	                                        AND ` + WarunekKonta + `)`

	pobierzWezelDebaty = `SELECT ` + kolumnyWezlaDebaty + `
	                      FROM debata_wezel
	                      WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzWezlyDebaty = `SELECT ` + kolumnyWezlaDebaty + `
	                      FROM debata_wezel
	                      WHERE okno = ? AND (? = '' OR tura = ?) AND (? = 0 OR kluczowy = 1)
	                        AND ` + WarunekKonta + `
	                      ORDER BY id ASC`

	oznaczWezelDebaty = `UPDATE debata_wezel SET kluczowy = ?
	                     WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	zapiszKrawedzDebaty = `INSERT INTO debata_krawedz
	                       (identyfikator_zewnetrzny, okno, wezel_zrodlowy, wezel_wskazany,
	                        relacja, pewnosc)
	                       VALUES (?, ?, ?, ?, ?, ?)`

	pobierzKrawedzieDebaty = `SELECT identyfikator_zewnetrzny, okno, wezel_zrodlowy,
	                                 wezel_wskazany, relacja, pewnosc
	                          FROM debata_krawedz
	                          WHERE okno = ?
	                            AND EXISTS (SELECT 1 FROM debata_wezel
	                                         WHERE debata_wezel.identyfikator_zewnetrzny
	                                               = debata_krawedz.wezel_zrodlowy
	                                           AND ` + WarunekKonta + `)
	                          ORDER BY id ASC`

	zapiszOznaczenieBleduDebaty = `INSERT INTO debata_oznaczenie_bledu
	                               (identyfikator_zewnetrzny, okno, wezel, kod, nazwa,
	                                uzasadnienie, pewnosc)
	                               VALUES (?, ?, ?, ?, ?, ?, ?)`

	pobierzOznaczeniaBledowDebaty = `SELECT identyfikator_zewnetrzny, okno, wezel, kod, nazwa,
	                                        uzasadnienie, pewnosc
	                                 FROM debata_oznaczenie_bledu
	                                 WHERE okno = ?
	                                   AND EXISTS (SELECT 1 FROM debata_wezel
	                                                WHERE debata_wezel.identyfikator_zewnetrzny
	                                                      = debata_oznaczenie_bledu.wezel
	                                                  AND ` + WarunekKonta + `)
	                                 ORDER BY id ASC`

	// Brak wiersza zawężenia dla okna w koncie znaczy katalog w całości włączony; `debata_katalog_bledu` jest słownikiem wspólnym instalacji.
	pobierzKatalogBledowDebaty = `SELECT k.kod, k.nazwa, k.opis,
	                                     CASE WHEN NOT EXISTS (SELECT 1 FROM debata_katalog_okna
	                                                            WHERE okno = ? AND ` + WarunekKonta + `)
	                                          THEN 1
	                                          WHEN EXISTS (SELECT 1 FROM debata_katalog_okna
	                                                        WHERE okno = ? AND kod = k.kod
	                                                          AND ` + WarunekKonta + `)
	                                          THEN 1 ELSE 0 END
	                                FROM debata_katalog_bledu k
	                               ORDER BY k.id ASC`

	usunKatalogOknaDebaty = `DELETE FROM debata_katalog_okna WHERE okno = ? AND ` + WarunekKonta

	// Para (okno, kod) jest jednoznaczna w całej tabeli: gałąź konfliktu bez warunku konta sięgałaby wiersza konta cudzego.
	zapiszKatalogOknaDebaty = `INSERT INTO debata_katalog_okna (okno, kod, konto_id)
	                           VALUES (?, ?, ` + WskazanieKonta + `)
	                           ON CONFLICT(okno, kod) DO UPDATE SET kod = excluded.kod
	                           WHERE ` + WarunekKonta
)

// Oznaczenia kluczowe przechodzą przez zastąpienie po treści węzła: ponowne wydobycie argumentów nie kasuje wyboru Operatora.
func (r *repozytoriumRoundtable) ZastapGrafDebaty(ctx context.Context, okno, tura string,
	wezly []WezelDebaty, krawedzie []KrawedzDebaty) error {

	kluczowe, err := r.trescKluczowychDebaty(ctx, okno, tura)
	if err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		// Krawędzie giną przed węzłami: warunek konta dochodzi do nich przez węzeł źródłowy.
		if tura == "" {
			wyczyscKrawedzie, err := r.zapytania.wTransakcji(ctx, transakcja, usunKrawedzieDebaty)
			if err != nil {
				return err
			}
			if _, err := wyczyscKrawedzie.ExecContext(ctx, okno, KontoOperatora(ctx)); err != nil {
				return fmt.Errorf("dane: nie można wyczyścić krawędzi debaty okna %q: %w", okno, err)
			}
		}
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunWezlyDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, okno, tura, tura, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić grafu debaty okna %q: %w", okno, err)
		}
		wstawWezel, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszWezelDebaty)
		if err != nil {
			return err
		}
		for _, wezel := range wezly {
			if kluczowe[wezel.Tresc] {
				wezel.Kluczowy = true
			}
			if wezel.Poparcie <= 0 {
				wezel.Poparcie = 1
			}
			if _, err := wstawWezel.ExecContext(ctx, wezel.Kod, okno, wezel.Wypowiedz,
				wezel.Uczestnik, wezel.Tura, wezel.AktMowy, wezel.Tresc, wezel.Poparcie,
				wezel.Kluczowy, KontoOperatora(ctx)); err != nil {
				return fmt.Errorf("dane: nie można zapisać węzła grafu debaty %q: %w", wezel.Kod, err)
			}
		}
		wstawKrawedz, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKrawedzDebaty)
		if err != nil {
			return err
		}
		for _, krawedz := range krawedzie {
			if _, err := wstawKrawedz.ExecContext(ctx, krawedz.Kod, okno, krawedz.WezelZrodlowy,
				krawedz.WezelWskazany, krawedz.Relacja, krawedz.Pewnosc); err != nil {
				return fmt.Errorf("dane: nie można zapisać krawędzi grafu debaty %q: %w",
					krawedz.Kod, err)
			}
		}
		return nil
	})
}

func (r *repozytoriumRoundtable) trescKluczowychDebaty(ctx context.Context,
	okno, tura string) (map[string]bool, error) {

	wezly, err := r.WezlyDebaty(ctx, okno, tura, true)
	if err != nil {
		return nil, err
	}
	kluczowe := make(map[string]bool, len(wezly))
	for _, wezel := range wezly {
		kluczowe[wezel.Tresc] = true
	}
	return kluczowe, nil
}

func (r *repozytoriumRoundtable) WezelDebatyPoKodzie(ctx context.Context, kod string) (WezelDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWezelDebaty)
	if err != nil {
		return WezelDebaty{}, err
	}
	wezel, err := odczytajWezelDebaty(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return WezelDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return WezelDebaty{}, fmt.Errorf("dane: nieczytelny wiersz węzła grafu debaty %q: %w", kod, err)
	}
	return wezel, nil
}

func (r *repozytoriumRoundtable) WezlyDebaty(ctx context.Context, okno, tura string,
	tylkoKluczowe bool) ([]WezelDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWezlyDebaty)
	if err != nil {
		return nil, err
	}
	zawezenie := 0
	if tylkoKluczowe {
		zawezenie = 1
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, tura, tura, zawezenie, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać węzłów grafu debaty okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	wezly := make([]WezelDebaty, 0, 16)
	for wiersze.Next() {
		wezel, err := odczytajWezelDebaty(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz węzła grafu debaty: %w", err)
		}
		wezly = append(wezly, wezel)
	}
	return wezly, wiersze.Err()
}

func (r *repozytoriumRoundtable) OznaczWezelDebaty(ctx context.Context, kod string, kluczowy bool) error {
	polecenie, err := r.zapytania.przygotuj(ctx, oznaczWezelDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kluczowy, kod, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można oznaczyć węzła grafu debaty %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

func (r *repozytoriumRoundtable) KrawedzieDebaty(ctx context.Context, okno string) ([]KrawedzDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKrawedzieDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać krawędzi grafu debaty okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	krawedzie := make([]KrawedzDebaty, 0, 16)
	for wiersze.Next() {
		var krawedz KrawedzDebaty
		if err := wiersze.Scan(&krawedz.Kod, &krawedz.Okno, &krawedz.WezelZrodlowy,
			&krawedz.WezelWskazany, &krawedz.Relacja, &krawedz.Pewnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz krawędzi grafu debaty: %w", err)
		}
		krawedzie = append(krawedzie, krawedz)
	}
	return krawedzie, wiersze.Err()
}

func (r *repozytoriumRoundtable) ZapiszOznaczenieBleduDebaty(ctx context.Context,
	oznaczenie OznaczenieBleduDebaty) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszOznaczenieBleduDebaty)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, oznaczenie.Kod, oznaczenie.Okno, oznaczenie.Wezel,
		oznaczenie.KodBledu, oznaczenie.Nazwa, oznaczenie.Uzasadnienie,
		oznaczenie.Pewnosc); err != nil {
		return fmt.Errorf("dane: nie można zapisać oznaczenia błędu %q: %w", oznaczenie.Kod, err)
	}
	return nil
}

func (r *repozytoriumRoundtable) OznaczeniaBledowDebaty(ctx context.Context,
	okno string) ([]OznaczenieBleduDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOznaczeniaBledowDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać oznaczeń błędów okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	oznaczenia := make([]OznaczenieBleduDebaty, 0, 8)
	for wiersze.Next() {
		var oznaczenie OznaczenieBleduDebaty
		if err := wiersze.Scan(&oznaczenie.Kod, &oznaczenie.Okno, &oznaczenie.Wezel,
			&oznaczenie.KodBledu, &oznaczenie.Nazwa, &oznaczenie.Uzasadnienie,
			&oznaczenie.Pewnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz oznaczenia błędu: %w", err)
		}
		oznaczenia = append(oznaczenia, oznaczenie)
	}
	return oznaczenia, wiersze.Err()
}

// Okno puste oddaje katalog wnoszony migracją, w całości włączony.
func (r *repozytoriumRoundtable) KatalogBledowDebaty(ctx context.Context,
	okno string) ([]DefinicjaBleduDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKatalogBledowDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx), okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogu błędów: %w", err)
	}
	defer wiersze.Close()

	katalog := make([]DefinicjaBleduDebaty, 0, 16)
	for wiersze.Next() {
		var pozycja DefinicjaBleduDebaty
		if err := wiersze.Scan(&pozycja.Kod, &pozycja.Nazwa, &pozycja.Opis,
			&pozycja.Wykrywany); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz katalogu błędów: %w", err)
		}
		katalog = append(katalog, pozycja)
	}
	return katalog, wiersze.Err()
}

// Wykaz pusty wraca do stanu „wykrywaj wszystko": okno bez zawężenia nie ma wiersza.
func (r *repozytoriumRoundtable) UstawKatalogBledowDebaty(ctx context.Context,
	okno string, kody []string) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunKatalogOknaDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, okno, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić katalogu błędów okna %q: %w", okno, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKatalogOknaDebaty)
		if err != nil {
			return err
		}
		for _, kod := range kody {
			wynik, err := wstaw.ExecContext(ctx, okno, kod, KontoOperatora(ctx), KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać zakresu katalogu błędów okna %q: %w",
					okno, err)
			}
			if err := sprawdzTrafienieZapisu(wynik, "zakres katalogu błędów okna "+okno, kod); err != nil {
				return err
			}
		}
		return nil
	})
}

func odczytajWezelDebaty(wiersz interface{ Scan(...any) error }) (WezelDebaty, error) {
	var wezel WezelDebaty
	err := wiersz.Scan(&wezel.Kod, &wezel.Okno, &wezel.Wypowiedz, &wezel.Uczestnik, &wezel.Tura,
		&wezel.AktMowy, &wezel.Tresc, &wezel.Poparcie, &wezel.Kluczowy, &wezel.Utworzono)
	return wezel, err
}
