// Plik zawiera graf argumentów debaty: węzły, krawędzie, katalog błędów logicznych oraz
// oznaczenia postawione na węzłach. Uzasadnienie trwałości grafu i zachowania oznaczeń
// niesie rozdział roundtable_graf.go dokumentacji architektury.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WezelDebaty odwzorowuje jednostkę argumentacyjną wydobytą z wypowiedzi uczestnika
// debaty, wraz z aktem mowy i poziomem poparcia.
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

// KrawedzDebaty odwzorowuje relację między dwoma węzłami grafu argumentów wraz z jej
// rodzajem i pewnością wykrycia.
type KrawedzDebaty struct {
	Kod           string
	Okno          string
	WezelZrodlowy string
	WezelWskazany string
	Relacja       string
	Pewnosc       float64
}

// DefinicjaBleduDebaty to pozycja katalogu błędów logicznych wraz z informacją, czy
// okno debaty ją aktualnie wykrywa.
type DefinicjaBleduDebaty struct {
	Kod       string
	Nazwa     string
	Opis      string
	Wykrywany bool
}

// OznaczenieBleduDebaty odwzorowuje błąd logiczny rozpoznany na węźle grafu wraz
// z uzasadnieniem i pewnością rozpoznania.
type OznaczenieBleduDebaty struct {
	Kod          string
	Okno         string
	Wezel        string
	KodBledu     string
	Nazwa        string
	Uzasadnienie string
	Pewnosc      float64
}

// RepozytoriumDebatyGrafu jest częścią kontraktu obszaru odpowiadającą za graf
// argumentów i katalog błędów.
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

	// Krawędź własnej kolumny konta nie ma: granica dochodzi do niej przez węzeł
	// wskazany kolumną `wezel_zrodlowy`.
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

	// Oznaczenie własnej kolumny konta nie ma: granica dochodzi do niego przez
	// węzeł wskazany kolumną `wezel`.
	pobierzOznaczeniaBledowDebaty = `SELECT identyfikator_zewnetrzny, okno, wezel, kod, nazwa,
	                                        uzasadnienie, pewnosc
	                                 FROM debata_oznaczenie_bledu
	                                 WHERE okno = ?
	                                   AND EXISTS (SELECT 1 FROM debata_wezel
	                                                WHERE debata_wezel.identyfikator_zewnetrzny
	                                                      = debata_oznaczenie_bledu.wezel
	                                                  AND ` + WarunekKonta + `)
	                                 ORDER BY id ASC`

	// Brak wiersza zawężenia dla okna znaczy katalog w całości włączony, więc
	// warunek pyta najpierw, czy okno w ogóle coś zawężało.
	pobierzKatalogBledowDebaty = `SELECT k.kod, k.nazwa, k.opis,
	                                     CASE WHEN NOT EXISTS (SELECT 1 FROM debata_katalog_okna
	                                                            WHERE okno = ?)
	                                          THEN 1
	                                          WHEN EXISTS (SELECT 1 FROM debata_katalog_okna
	                                                        WHERE okno = ? AND kod = k.kod)
	                                          THEN 1 ELSE 0 END
	                                FROM debata_katalog_bledu k
	                               ORDER BY k.id ASC`

	usunKatalogOknaDebaty = `DELETE FROM debata_katalog_okna WHERE okno = ?`

	zapiszKatalogOknaDebaty = `INSERT INTO debata_katalog_okna (okno, kod) VALUES (?, ?)
	                           ON CONFLICT(okno, kod) DO NOTHING`
)

// ZastapGrafDebaty wymienia graf okna albo jednej tury w jednej transakcji. Oznaczenia
// kluczowe przechodzą przez zastąpienie po treści węzła, żeby ponowne wydobycie
// argumentów nie kasowało wyboru operatora.
func (r *repozytoriumRoundtable) ZastapGrafDebaty(ctx context.Context, okno, tura string,
	wezly []WezelDebaty, krawedzie []KrawedzDebaty) error {

	kluczowe, err := r.trescKluczowychDebaty(ctx, okno, tura)
	if err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		// Krawędzie idą całym oknem, bo relacja łączy węzły z różnych tur, nie tylko z jednej.
		// Giną przed węzłami: warunek konta dochodzi do krawędzi przez węzeł
		// źródłowy, a po usunięciu węzłów nie ma już po czym dojść.
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

// trescKluczowychDebaty zbiera treści węzłów oznaczonych jako kluczowe, żeby
// zastąpienie grafu nie zgubiło wyboru Operatora.
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

// WezelDebatyPoKodzie zwraca jeden węzeł grafu wskazany kodem zewnętrznym, albo błąd
// ErrBrakWiersza, gdy nie istnieje.
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

// WezlyDebaty zwraca węzły okna, opcjonalnie zawężone do wskazanej tury oraz wyłącznie
// do węzłów oznaczonych jako kluczowe.
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

// OznaczWezelDebaty stawia albo zdejmuje oznaczenie argumentu kluczowego na węźle
// wskazanym kodem zewnętrznym.
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

// KrawedzieDebaty zwraca wszystkie krawędzie grafu argumentów wskazanego okna,
// uporządkowane według kolejności zapisu.
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

// ZapiszOznaczenieBleduDebaty dopisuje błąd logiczny rozpoznany na węźle grafu wraz
// z jego uzasadnieniem i pewnością.
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

// OznaczeniaBledowDebaty zwraca wszystkie oznaczenia błędów logicznych postawione
// w oknie debaty, uporządkowane według zapisu.
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

// KatalogBledowDebaty zwraca katalog wraz z zakresem wykrywania w oknie.
// Okno puste oddaje katalog wnoszony migracją, w całości włączony.
func (r *repozytoriumRoundtable) KatalogBledowDebaty(ctx context.Context,
	okno string) ([]DefinicjaBleduDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKatalogBledowDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, okno)
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

// UstawKatalogBledowDebaty zapisuje zakres wykrywania dla okna. Wykaz pusty
// wraca do stanu „wykrywaj wszystko": okno bez zawężenia nie ma wiersza.
func (r *repozytoriumRoundtable) UstawKatalogBledowDebaty(ctx context.Context,
	okno string, kody []string) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunKatalogOknaDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, okno); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić katalogu błędów okna %q: %w", okno, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKatalogOknaDebaty)
		if err != nil {
			return err
		}
		for _, kod := range kody {
			if _, err := wstaw.ExecContext(ctx, okno, kod); err != nil {
				return fmt.Errorf("dane: nie można zapisać zakresu katalogu błędów okna %q: %w",
					okno, err)
			}
		}
		return nil
	})
}

// odczytajWezelDebaty składa strukturę WezelDebaty z jednego wiersza wyniku zapytania,
// niezależnie od jego źródła.
func odczytajWezelDebaty(wiersz interface{ Scan(...any) error }) (WezelDebaty, error) {
	var wezel WezelDebaty
	err := wiersz.Scan(&wezel.Kod, &wezel.Okno, &wezel.Wypowiedz, &wezel.Uczestnik, &wezel.Tura,
		&wezel.AktMowy, &wezel.Tresc, &wezel.Poparcie, &wezel.Kluczowy, &wezel.Utworzono)
	return wezel, err
}
