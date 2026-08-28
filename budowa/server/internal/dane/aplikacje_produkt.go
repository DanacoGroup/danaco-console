// Obszar Apps — strona budowy produktu. Produkt okna (tabela `produkt_apps`),
// etapy budowy (`etap_apps`) oraz kamienie milowe wraz ze związkiem z etapami
// (`kamien_milowy_apps`, `kamien_milowy_etap_apps`).
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// ProduktApp to wiersz tabeli `produkt_apps`: metadane produktu jednego
// okna wraz z listą platform docelowych.
type ProduktApp struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Opis           *string
	Platformy      []string
	Repozytorium   *string
	Utworzono      string
	Zaktualizowano string
}

// EtapApp to wiersz tabeli `etap_apps` — etap trackera Product Buildera,
// wraz z kolejnością i wykonawcą.
type EtapApp struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Kolejnosc      int
	Stan           string
	Wykonawca      *string
	Utworzono      string
	Zaktualizowano string
}

// KamienMilowyApp to wiersz tabeli `kamien_milowy_apps` wraz z kodami etapów,
// które się na niego składają — kontrakt oddaje kamień milowy zawsze razem
// z nimi, więc rozdzielanie tego na dwa odczyty byłoby pracą bez odbiorcy.
type KamienMilowyApp struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Termin         *int64
	Stan           string
	KodyEtapow     []string
	Utworzono      string
	Zaktualizowano string
}

const (
	kolumnyProduktuApp = `id, identyfikator_zewnetrzny, okno, nazwa, opis, platformy,
	                      repozytorium, utworzono, zaktualizowano`

	// UPSERT po oknie — jeden produkt na okno (czoło migracji 200). Kod
	// zewnętrzny nadany przy pierwszym zapisie zostaje: `AppProduct.id` ma być
	// stały, a kolejne zapisy zmieniają metadane, nie tożsamość produktu.
	zapiszProduktApp = `INSERT INTO produkt_apps
	                    (identyfikator_zewnetrzny, okno, nazwa, opis, platformy, repozytorium)
	                    VALUES (?, ?, ?, ?, ?, ?)
	                    ON CONFLICT(okno) DO UPDATE SET
	                        nazwa = excluded.nazwa,
	                        opis = excluded.opis,
	                        platformy = excluded.platformy,
	                        repozytorium = excluded.repozytorium,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzProduktApp = `SELECT ` + kolumnyProduktuApp + ` FROM produkt_apps WHERE okno = ?`

	kolumnyEtapuApp = `id, identyfikator_zewnetrzny, okno, nazwa, kolejnosc, stan,
	                   wykonawca, utworzono, zaktualizowano`

	zapiszEtapApp = `INSERT INTO etap_apps
	                 (identyfikator_zewnetrzny, okno, nazwa, kolejnosc, stan, wykonawca)
	                 VALUES (?, ?, ?, ?, ?, ?)
	                 ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                     nazwa = excluded.nazwa,
	                     kolejnosc = excluded.kolejnosc,
	                     stan = excluded.stan,
	                     wykonawca = excluded.wykonawca,
	                     zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzEtapApp = `SELECT ` + kolumnyEtapuApp + ` FROM etap_apps
	                  WHERE identyfikator_zewnetrzny = ?`

	listaEtapowApp = `SELECT ` + kolumnyEtapuApp + ` FROM etap_apps
	                  WHERE okno = ? ORDER BY kolejnosc, id`

	kolumnyKamieniaApp = `id, identyfikator_zewnetrzny, okno, nazwa, termin, stan,
	                      utworzono, zaktualizowano`

	zapiszKamienApp = `INSERT INTO kamien_milowy_apps
	                   (identyfikator_zewnetrzny, okno, nazwa, termin, stan)
	                   VALUES (?, ?, ?, ?, ?)
	                   ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                       nazwa = excluded.nazwa,
	                       termin = excluded.termin,
	                       stan = excluded.stan,
	                       zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzKamienApp = `SELECT ` + kolumnyKamieniaApp + ` FROM kamien_milowy_apps
	                    WHERE identyfikator_zewnetrzny = ?`

	listaKamieniApp = `SELECT ` + kolumnyKamieniaApp + ` FROM kamien_milowy_apps
	                   WHERE okno = ? ORDER BY IFNULL(termin, 0), id`

	usunKamienApp = `DELETE FROM kamien_milowy_apps WHERE identyfikator_zewnetrzny = ?`

	usunEtapyKamieniaApp = `DELETE FROM kamien_milowy_etap_apps WHERE kamien_id = ?`

	wstawEtapKamieniaApp = `INSERT INTO kamien_milowy_etap_apps (kamien_id, etap_kod)
	                        VALUES (?, ?) ON CONFLICT(kamien_id, etap_kod) DO NOTHING`

	listaEtapowKamieniApp = `SELECT k.identyfikator_zewnetrzny, z.etap_kod
	                         FROM kamien_milowy_etap_apps z
	                         JOIN kamien_milowy_apps k ON k.id = z.kamien_id
	                         WHERE k.okno = ? ORDER BY z.etap_kod`
)

// ZapiszProduktApp zapisuje produkt okna, UPSERT po kolumnie okna. Kod
// zewnętrzny podany w strukturze obowiązuje wyłącznie przy pierwszym
// zapisie — przy kolejnych wiersz zostaje pod kodem nadanym wcześniej.
func (r *repozytoriumAplikacji) ZapiszProduktApp(ctx context.Context, produkt ProduktApp) (ProduktApp, error) {
	if produkt.Okno == "" {
		return ProduktApp{}, fmt.Errorf("dane: produkt aplikacji bez okna")
	}
	if produkt.Nazwa == "" {
		return ProduktApp{}, fmt.Errorf("dane: produkt aplikacji okna %q bez nazwy", produkt.Okno)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszProduktApp)
	if err != nil {
		return ProduktApp{}, err
	}
	_, err = polecenie.ExecContext(ctx, produkt.Kod, produkt.Okno, produkt.Nazwa,
		tekstDoKolumny(produkt.Opis), listaDoKolumny(produkt.Platformy),
		tekstDoKolumny(produkt.Repozytorium))
	if err != nil {
		return ProduktApp{}, fmt.Errorf("dane: nie można zapisać produktu okna %q: %w", produkt.Okno, err)
	}
	return r.ProduktApp(ctx, produkt.Okno)
}

// ProduktApp zwraca produkt zapisany dla wskazanego okna aplikacji; brak
// wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumAplikacji) ProduktApp(ctx context.Context, okno string) (ProduktApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProduktApp)
	if err != nil {
		return ProduktApp{}, err
	}
	produkt, err := odczytajProduktApp(polecenie.QueryRowContext(ctx, okno))
	if err == sql.ErrNoRows {
		return ProduktApp{}, ErrBrakWiersza
	}
	if err != nil {
		return ProduktApp{}, fmt.Errorf("dane: nieczytelny produkt okna %q: %w", okno, err)
	}
	return produkt, nil
}

// ZapiszEtapApp zapisuje jeden etap budowy produktu, UPSERT po
// identyfikatorze zewnętrznym etapu w oknie.
func (r *repozytoriumAplikacji) ZapiszEtapApp(ctx context.Context, etap EtapApp) (EtapApp, error) {
	if etap.Kod == "" || etap.Okno == "" {
		return EtapApp{}, fmt.Errorf("dane: etap aplikacji bez identyfikatora albo bez okna")
	}
	if etap.Stan == "" {
		etap.Stan = "pending"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszEtapApp)
	if err != nil {
		return EtapApp{}, err
	}
	_, err = polecenie.ExecContext(ctx, etap.Kod, etap.Okno, etap.Nazwa, etap.Kolejnosc,
		etap.Stan, tekstDoKolumny(etap.Wykonawca))
	if err != nil {
		return EtapApp{}, fmt.Errorf("dane: nie można zapisać etapu %q: %w", etap.Kod, err)
	}
	return r.EtapApp(ctx, etap.Kod)
}

// EtapApp zwraca jeden etap trackera Product Buildera po jego kodzie
// zewnętrznym, niezależnie od okna.
func (r *repozytoriumAplikacji) EtapApp(ctx context.Context, kod string) (EtapApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzEtapApp)
	if err != nil {
		return EtapApp{}, err
	}
	etap, err := odczytajEtapApp(polecenie.QueryRowContext(ctx, kod))
	if err == sql.ErrNoRows {
		return EtapApp{}, ErrBrakWiersza
	}
	if err != nil {
		return EtapApp{}, fmt.Errorf("dane: nieczytelny etap %q: %w", kod, err)
	}
	return etap, nil
}

// EtapyApp zwraca wszystkie etapy budowy produktu okna w kolejności
// trackera, tej ustawionej polem Kolejnosc.
func (r *repozytoriumAplikacji) EtapyApp(ctx context.Context, okno string) ([]EtapApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaEtapowApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać etapów okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []EtapApp{}
	for wiersze.Next() {
		etap, err := odczytajEtapApp(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz etapu okna %q: %w", okno, err)
		}
		lista = append(lista, etap)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt etapów okna %q: %w", okno, err)
	}
	return lista, nil
}

// ZapiszKamienMilowyApp zapisuje kamień milowy wraz z kompletem etapów, które
// się na niego składają — związek wymieniany jest w całości w jednej
// transakcji, bo kontrakt nadsyła `stageIds` bez trybu częściowej zmiany.
func (r *repozytoriumAplikacji) ZapiszKamienMilowyApp(ctx context.Context,
	kamien KamienMilowyApp) (KamienMilowyApp, error) {

	if kamien.Kod == "" || kamien.Okno == "" {
		return KamienMilowyApp{}, fmt.Errorf("dane: kamień milowy bez identyfikatora albo bez okna")
	}
	if kamien.Stan == "" {
		kamien.Stan = "planned"
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKamienApp)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, kamien.Kod, kamien.Okno, kamien.Nazwa,
			liczbaDoKolumny(kamien.Termin), kamien.Stan); err != nil {
			return fmt.Errorf("dane: nie można zapisać kamienia milowego %q: %w", kamien.Kod, err)
		}

		var kamienID int64
		wiersz := transakcja.QueryRowContext(ctx,
			`SELECT id FROM kamien_milowy_apps WHERE identyfikator_zewnetrzny = ?`, kamien.Kod)
		if err := wiersz.Scan(&kamienID); err != nil {
			return fmt.Errorf("dane: nie można odczytać id kamienia milowego %q: %w", kamien.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtapyKamieniaApp)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, kamienID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić etapów kamienia %q: %w", kamien.Kod, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawEtapKamieniaApp)
		if err != nil {
			return err
		}
		for _, etap := range kamien.KodyEtapow {
			if etap == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, kamienID, etap); err != nil {
				return fmt.Errorf("dane: nie można zapisać etapu %q kamienia %q: %w",
					etap, kamien.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return KamienMilowyApp{}, err
	}
	return r.KamienMilowyApp(ctx, kamien.Kod)
}

// KamienMilowyApp zwraca jeden kamień milowy budowy produktu wraz z kodami
// wszystkich jego etapów składowych.
func (r *repozytoriumAplikacji) KamienMilowyApp(ctx context.Context, kod string) (KamienMilowyApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKamienApp)
	if err != nil {
		return KamienMilowyApp{}, err
	}
	kamien, err := odczytajKamienApp(polecenie.QueryRowContext(ctx, kod))
	if err == sql.ErrNoRows {
		return KamienMilowyApp{}, ErrBrakWiersza
	}
	if err != nil {
		return KamienMilowyApp{}, fmt.Errorf("dane: nieczytelny kamień milowy %q: %w", kod, err)
	}
	etapy, err := r.etapyKamieniOkna(ctx, kamien.Okno)
	if err != nil {
		return KamienMilowyApp{}, err
	}
	kamien.KodyEtapow = etapy[kamien.Kod]
	return kamien, nil
}

// KamienieMiloweApp zwraca wszystkie kamienie milowe budowy produktu okna
// wraz z kodami etapów każdego z nich.
func (r *repozytoriumAplikacji) KamienieMiloweApp(ctx context.Context, okno string) ([]KamienMilowyApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKamieniApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kamieni milowych okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []KamienMilowyApp{}
	for wiersze.Next() {
		kamien, err := odczytajKamienApp(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kamienia milowego okna %q: %w", okno, err)
		}
		lista = append(lista, kamien)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kamieni milowych okna %q: %w", okno, err)
	}

	etapy, err := r.etapyKamieniOkna(ctx, okno)
	if err != nil {
		return nil, err
	}
	for indeks := range lista {
		lista[indeks].KodyEtapow = etapy[lista[indeks].Kod]
	}
	return lista, nil
}

// UsunKamienMilowyApp kasuje kamień milowy; związek z etapami znika kaskadą.
// Zwraca prawdę, gdy wiersz naprawdę zniknął — kontrakt oddaje `deleted`, więc
// „nie było czego kasować" nie ma prawa wrócić jako powodzenie usunięcia.
func (r *repozytoriumAplikacji) UsunKamienMilowyApp(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunKamienApp)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć kamienia milowego %q: %w", kod, err)
	}
	ile, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych kamieni %q: %w", kod, err)
	}
	return ile > 0, nil
}

// etapyKamieniOkna zwraca mapę kod kamienia → kody etapów. Jedno zapytanie na
// całe okno zamiast jednego na kamień: wykaz kamieni ciągnąłby inaczej tyle
// zapytań, ile ma pozycji.
func (r *repozytoriumAplikacji) etapyKamieniOkna(ctx context.Context, okno string) (map[string][]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaEtapowKamieniApp)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać etapów kamieni okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	mapa := map[string][]string{}
	for wiersze.Next() {
		var kamien, etap string
		if err := wiersze.Scan(&kamien, &etap); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz etapu kamienia: %w", err)
		}
		mapa[kamien] = append(mapa[kamien], etap)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt etapów kamieni okna %q: %w", okno, err)
	}
	return mapa, nil
}

// odczytajProduktApp składa pełną strukturę ProduktApp z jednego wiersza
// wyniku zapytania SQL bazy danych.
func odczytajProduktApp(wiersz skaner) (ProduktApp, error) {
	var produkt ProduktApp
	var opis, platformy, repozytorium sql.NullString
	err := wiersz.Scan(&produkt.ID, &produkt.Kod, &produkt.Okno, &produkt.Nazwa,
		&opis, &platformy, &repozytorium, &produkt.Utworzono, &produkt.Zaktualizowano)
	if err != nil {
		return ProduktApp{}, err
	}
	produkt.Opis = tekstZKolumny(opis)
	produkt.Platformy = listaZKolumny(platformy)
	produkt.Repozytorium = tekstZKolumny(repozytorium)
	return produkt, nil
}

// odczytajEtapApp składa pełną strukturę EtapApp z jednego wiersza wyniku
// wykonanego zapytania SQL bazy danych.
func odczytajEtapApp(wiersz skaner) (EtapApp, error) {
	var etap EtapApp
	var wykonawca sql.NullString
	err := wiersz.Scan(&etap.ID, &etap.Kod, &etap.Okno, &etap.Nazwa, &etap.Kolejnosc,
		&etap.Stan, &wykonawca, &etap.Utworzono, &etap.Zaktualizowano)
	if err != nil {
		return EtapApp{}, err
	}
	etap.Wykonawca = tekstZKolumny(wykonawca)
	return etap, nil
}

// odczytajKamienApp składa kamień milowy z jednego wiersza wyniku; kody etapów
// dokłada wołający z osobnego zapytania.
func odczytajKamienApp(wiersz skaner) (KamienMilowyApp, error) {
	var kamien KamienMilowyApp
	var termin sql.NullInt64
	err := wiersz.Scan(&kamien.ID, &kamien.Kod, &kamien.Okno, &kamien.Nazwa, &termin,
		&kamien.Stan, &kamien.Utworzono, &kamien.Zaktualizowano)
	if err != nil {
		return KamienMilowyApp{}, err
	}
	kamien.Termin = liczbaZKolumny(termin)
	return kamien, nil
}
