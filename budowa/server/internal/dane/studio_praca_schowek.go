// Odpowiedzialność pliku: warstwa danych odcinka kontroli pracy modułu Studio,
// część czwarta — dostęp Studia do HISTORII SCHOWKA platformy (tabela
// `wpis_schowka`, migracja 292) oraz do nastaw wielu wykonawców trzymanych
// zasięgami konfiguracji (tabela `ustawienie`).
//
// ── Dlaczego Studio sięga do tabeli schowka, a nie zakłada własnej ──────────
// Zlecenie stanowi wprost: rodzina `clipboard.*` istnieje i drugiej się nie
// zakłada — schowek ma być JEDEN, wspólny Operatorowi i wykonawcom, bo cały
// sens historii schowka polega na tym, że wpis odłożony w jednym miejscu daje
// się wkleić w drugim. Osobna tabela schowka Studia rozdzieliłaby jedną historię
// na dwie i „wklej wpis sprzed kilku ruchów" przestałoby działać między oknem
// pracy z dokumentem a resztą platformy.
//
// Wiersze pisze się więc DOKŁADNIE tymi samymi poleceniami, co rodzina
// `clipboard.*` (stałe z `schowek.go`), i tym samym rachunkiem odcisku — inaczej
// ta sama treść odłożona dwiema drogami zrobiłaby dwa wpisy, a warunek UNIQUE na
// odcisku istnieje właśnie po to, żeby tego nie było.
//
// ── Dlaczego nastawy wykonawców idą tabelą `ustawienie` ─────────────────────
// Nastawy pętli wykonawczej i pracy wielu agentów są w kontrakcie opisane jako
// wartości zasięgów rodziny `config.*` — z poziomem, z którego przyszły, oddanym
// w odpowiedzi. Odczyt idzie rozstrzygaczem dziewięciu poziomów (rdzeń nie zna
// ich pierwszeństwa i o nie nie pyta), a zapis tym poleceniem: `ustawienie` jest
// jedyną tabelą nastaw platformy i osobny magazyn nastaw Studia byłby drugim
// miejscem, w którym trzeba by szukać wartości obowiązującej.
package dane

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

const (
	// Najświeższy wpis historii — „wklej to, co ostatnio odłożone", gdy żądanie
	// nie wskazuje wpisu. Przypięte NIE mają tu pierwszeństwa: przypięcie chroni
	// wpis od wygaśnięcia, a nie czyni go ostatnio odłożonym.
	wpisSchowkaNajswiezszy = `SELECT ` + kolumnyWpisuSchowka +
		` FROM wpis_schowka ORDER BY utworzono DESC, id DESC LIMIT 1`
)

// DopiszWpisSchowkaStudia odkłada fragment dokumentu w historii schowka
// platformy i mówi, czy treść była już odłożona.
//
// Odcisk liczy się z rodzaju i treści razem — tak samo jak w rodzinie
// `clipboard.*`. Drugi rachunek odcisku rozjechałby warunek UNIQUE i ta sama
// treść stałaby w historii dwa razy.
func (r *repozytoriumStudia) DopiszWpisSchowkaStudia(ctx context.Context,
	wpis WpisSchowka) (WpisSchowka, bool, error) {

	if wpis.Tresc == "" {
		return WpisSchowka{}, false, fmt.Errorf("dane: wpis schowka bez treści")
	}
	if wpis.Rodzaj == "" {
		wpis.Rodzaj = shared.ClipboardEntryKindText
	}
	suma := sha256.Sum256([]byte(wpis.Rodzaj + "\x00" + wpis.Tresc))
	wpis.Odcisk = hex.EncodeToString(suma[:])
	if wpis.RozmiarBajtow == 0 {
		wpis.RozmiarBajtow = int64(len(wpis.Tresc))
	}

	bylo := false
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzWpisSchowkaPoOdcisku)
		if err != nil {
			return err
		}
		zastany, err := odczytajWpisSchowka(odczyt.QueryRowContext(ctx, wpis.Odcisk))
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
		_, err = zapis.ExecContext(ctx, wpis.Kod, wpis.Rodzaj, wpis.Tresc, wpis.Odcisk,
			tekstDoKolumny(wpis.Zajawka), wpis.RozmiarBajtow, liczbaLogiczna(wpis.Przypiety),
			liczbaLogiczna(wpis.Wrazliwy), tekstDoKolumny(wpis.OknoZrodlowe), wpis.Utworzono,
			liczbaDoKolumny(wpis.Uzyto), tekstDoKolumny(wpis.PostacJSON))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać wpisu schowka: %w", err)
		}
		return nil
	})
	if err != nil {
		return WpisSchowka{}, false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisSchowkaPoOdcisku)
	if err != nil {
		return WpisSchowka{}, false, err
	}
	zapisany, err := odczytajWpisSchowka(polecenie.QueryRowContext(ctx, wpis.Odcisk))
	if err != nil {
		return WpisSchowka{}, false, fmt.Errorf("dane: nieczytelny wiersz wpisu schowka: %w", err)
	}
	return zapisany, bylo, nil
}

// WpisSchowkaStudia oddaje wpis historii o wskazanym kodzie.
func (r *repozytoriumStudia) WpisSchowkaStudia(ctx context.Context,
	kod string) (WpisSchowka, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisSchowka)
	if err != nil {
		return WpisSchowka{}, err
	}
	wpis, err := odczytajWpisSchowka(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WpisSchowka{}, ErrBrakWiersza
	}
	if err != nil {
		return WpisSchowka{}, fmt.Errorf("dane: nieczytelny wiersz wpisu schowka %q: %w", kod, err)
	}
	return wpis, nil
}

// NajswiezszyWpisSchowkaStudia oddaje wpis odłożony najpóźniej. Historia pusta
// jest brakiem wiersza, nie pustą treścią: wklejenie z pustego schowka ma
// odmówić słowami, a nie wstawić nic i powiedzieć, że wstawiło.
func (r *repozytoriumStudia) NajswiezszyWpisSchowkaStudia(ctx context.Context) (WpisSchowka, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wpisSchowkaNajswiezszy)
	if err != nil {
		return WpisSchowka{}, err
	}
	wpis, err := odczytajWpisSchowka(polecenie.QueryRowContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return WpisSchowka{}, ErrBrakWiersza
	}
	if err != nil {
		return WpisSchowka{}, fmt.Errorf("dane: nieczytelny wiersz wpisu schowka: %w", err)
	}
	return wpis, nil
}

// ── Nastawy zasięgu ─────────────────────────────────────────────────────────

// ZapiszNastaweZasieguStudia zapisuje wartość nastawy na wskazanym poziomie
// zasięgu, osią platformy. Wiersz istniejący nadpisuje.
func (r *repozytoriumStudia) ZapiszNastaweZasieguStudia(ctx context.Context,
	ustawienie Ustawienie) error {

	poziom, err := poziomZasieguNaBaze(ustawienie.Poziom)
	if err != nil {
		return err
	}
	os, err := osZasieguNaBaze(shared.ConfigAxisPlatform)
	if err != nil {
		return err
	}
	rodzaj := ustawienie.RodzajWartosci
	if rodzaj == "" {
		rodzaj = domyslnyRodzajWartosci
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszUstawienie)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, poziom, ustawienie.KluczZasiegu, os, "",
		ustawienie.Klucz, tekstDoKolumny(ustawienie.Wartosc), rodzaj)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać nastawy %q na poziomie %q: %w",
			ustawienie.Klucz, poziom, err)
	}
	return nil
}

// NastawaZasieguStudia oddaje wartość nastawy zapisaną NA WSKAZANYM poziomie —
// bez rozstrzygania pierwszeństwa poziomów, bo pierwszeństwo należy do pakietu
// konfiguracji, a nie do warstwy danych.
func (r *repozytoriumStudia) NastawaZasieguStudia(ctx context.Context,
	poziomZasiegu shared.ConfigScope, kluczZasiegu, klucz string) (Ustawienie, bool, error) {

	poziom, err := poziomZasieguNaBaze(poziomZasiegu)
	if err != nil {
		return Ustawienie{}, false, err
	}
	os, err := osZasieguNaBaze(shared.ConfigAxisPlatform)
	if err != nil {
		return Ustawienie{}, false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUstawienie)
	if err != nil {
		return Ustawienie{}, false, err
	}
	ustawienie, err := odczytajUstawienie(
		polecenie.QueryRowContext(ctx, poziom, kluczZasiegu, os, "", klucz))
	if errors.Is(err, sql.ErrNoRows) {
		return Ustawienie{}, false, nil
	}
	if err != nil {
		return Ustawienie{}, false, fmt.Errorf("dane: nieczytelny wiersz nastawy %q: %w", klucz, err)
	}
	return ustawienie, true, nil
}
