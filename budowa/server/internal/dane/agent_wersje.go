// Plik prowadzi historię wersji tożsamości eksperta: odczyt wykazu migawek oraz przywrócenie wersji wcześniejszej;
// archiwum eksperta leży w agent_archiwum.go jako ta sama implementacja repozytorium, rozdzielona wzdłuż dwóch odpowiedzialności.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// WersjaAgenta to wiersz `agent_wersja` — pełna tożsamość eksperta w jednej
// wersji wraz z wyliczoną listą pól zmienionych względem wersji poprzedniej.
type WersjaAgenta struct {
	// Identyfikator jest kluczem wiersza migawki, trzonem rodziny wersji; numer niesie tylko porządek.
	Identyfikator       int64
	Numer               int
	Nazwa               string
	Opis                string
	InstrukcjeSystemowe string
	KanalKod            *string
	Model               *string
	Transport           *string
	ParametryJSON       string
	ImieWlasne          string
	Favikon             string
	UstawieniaJSON      string
	TrybNakladki        string
	Aktywny             bool
	// Autor niesie napis, nie klucz konta: wartości to operator albo restore.
	Autor string
	// Powod jest pusty przy zwykłej zmianie; przywrócenie wpisuje tu numer
	// wersji źródłowej.
	Powod    string
	Zapisano string
	// Widocznosc i PoziomyPamieci są polami wymaganymi migawki; warstw promptu migawka nie niesie.
	Widocznosc     string
	PoziomyPamieci []string
	// ZmienionePola niesie nazwy kolumn różniące się od wersji poprzedniej; najstarsza ma listę pustą.
	ZmienionePola []string
}

// RepozytoriumWersjiAgenta jest kontraktem historii tożsamości eksperta: odczyt wykazu wersji i przywrócenie wersji wcześniejszej.
type RepozytoriumWersjiAgenta interface {
	// Wersje oddaje historię eksperta od najnowszej wersji do najstarszej.
	Wersje(ctx context.Context, kodAgenta string) ([]WersjaAgenta, error)
	// Przywroc zapisuje treść wersji `numer` jako wersję kolejną i oddaje
	// migawkę nowo powstałą.
	Przywroc(ctx context.Context, kodAgenta string, numer int) (WersjaAgenta, error)
}

const (
	kolumnyWersjiAgenta = `id, numer, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
	                       parametry_json, imie_wlasne, favikon, ustawienia_json, tryb_nakladki,
	                       aktywny, autor, powod, zapisano, widocznosc, poziomy_pamieci`

	historiaAgenta = `SELECT ` + kolumnyWersjiAgenta + ` FROM agent_wersja
	                  WHERE agent_id = (SELECT id FROM agent WHERE kod = ?)
	                  ORDER BY numer DESC`

	wersjaAgentaPoNumerze = `SELECT ` + kolumnyWersjiAgenta + ` FROM agent_wersja
	                         WHERE agent_id = (SELECT id FROM agent WHERE kod = ?) AND numer = ?`

	// Przywrócenie idzie jednym poleceniem z podniesieniem licznika wersji, tak samo jak zapis tożsamości.
	przywrocTrescWersji = `UPDATE agent
	                       SET nazwa = ?, opis = ?, instrukcje_systemowe = ?, kanal_kod = ?,
	                           model = ?, transport = ?, parametry_json = ?, imie_wlasne = ?,
	                           favikon = ?, ustawienia_json = ?, tryb_nakladki = ?,
	                           widocznosc = ?,
	                           wersja = wersja + 1,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
	                       WHERE kod = ?`

	// Wyzwalacz zakłada migawkę bez autora i powodu; ten zapis dokłada oba pola do wiersza świeżo powstałego.
	oznaczWersjePrzywrocona = `UPDATE agent_wersja
	                           SET autor = 'restore', powod = ?
	                           WHERE agent_id = (SELECT id FROM agent WHERE kod = ?)
	                             AND numer = (SELECT wersja FROM agent WHERE kod = ?)`

	numerWersjiBiezacej = `SELECT wersja FROM agent WHERE kod = ?`
)

// repozytoriumWersjiAgenta obsługuje historię i archiwum eksperta —
// dwie odpowiedzialności czytające tę samą tabelę `agent`.
type repozytoriumWersjiAgenta struct {
	zapytania *zapytania
	db        *sql.DB
}

var (
	_ RepozytoriumWersjiAgenta    = (*repozytoriumWersjiAgenta)(nil)
	_ RepozytoriumArchiwumAgentow = (*repozytoriumWersjiAgenta)(nil)
)

// noweRepozytoriumWersjiAgenta wiąże historię i archiwum eksperta z bazą.
// Wywołuje ją złożenie zestawu repozytoriów (`zestaw.go`).
func noweRepozytoriumWersjiAgenta(z *zapytania, db *sql.DB) *repozytoriumWersjiAgenta {
	return &repozytoriumWersjiAgenta{zapytania: z, db: db}
}

// Wersje oddaje historię eksperta od najnowszej. Ekspert bez historii nie jest
// błędem — każdy istniejący ekspert ma choć jedną migawkę, więc pusty wykaz
// znaczy „eksperta nie ma”, i to rozstrzyga sprawdzenie wywołujące.
func (r *repozytoriumWersjiAgenta) Wersje(ctx context.Context, kodAgenta string) ([]WersjaAgenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, historiaAgenta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kodAgenta)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać historii eksperta %q: %w", kodAgenta, err)
	}
	defer wiersze.Close()

	historia := []WersjaAgenta{}
	for wiersze.Next() {
		wersja, err := odczytajWersjeAgenta(wiersze)
		if err != nil {
			return nil, err
		}
		historia = append(historia, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt historii eksperta %q: %w", kodAgenta, err)
	}
	oznaczZmiany(historia)
	return historia, nil
}

// Przywroc zapisuje treść wskazanej wersji jako wersję kolejną. Wersji
// późniejszych nie usuwa — powrót jest zdarzeniem historii, nie jej korektą.
func (r *repozytoriumWersjiAgenta) Przywroc(ctx context.Context, kodAgenta string,
	numer int) (WersjaAgenta, error) {

	zrodlo, err := r.wersjaPoNumerze(ctx, kodAgenta, numer)
	if err != nil {
		return WersjaAgenta{}, err
	}
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, przywrocTrescWersji)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, zrodlo.Nazwa, zrodlo.Opis, zrodlo.InstrukcjeSystemowe,
			tekstDoKolumny(zrodlo.KanalKod), tekstDoKolumny(zrodlo.Model),
			tekstDoKolumny(zrodlo.Transport), zrodlo.ParametryJSON, zrodlo.ImieWlasne,
			zrodlo.Favikon, zrodlo.UstawieniaJSON, zrodlo.TrybNakladki, zrodlo.Widocznosc,
			kodAgenta); err != nil {
			return fmt.Errorf("dane: nie można przywrócić wersji %d eksperta %q: %w",
				numer, kodAgenta, err)
		}
		// Poziomy pamięci wracają razem z tożsamością, nie zostają z wersji bieżącej — migawka je niesie.
		if err := przywrocPoziomyPamieciWersji(ctx, transakcja, kodAgenta, zrodlo.PoziomyPamieci); err != nil {
			return err
		}
		oznacz, err := r.zapytania.wTransakcji(ctx, transakcja, oznaczWersjePrzywrocona)
		if err != nil {
			return err
		}
		powod := fmt.Sprintf("przywrocenie wersji %d", numer)
		if _, err := oznacz.ExecContext(ctx, powod, kodAgenta, kodAgenta); err != nil {
			return fmt.Errorf("dane: nie można oznaczyć wersji przywróconej eksperta %q: %w",
				kodAgenta, err)
		}
		return nil
	})
	if err != nil {
		return WersjaAgenta{}, err
	}
	biezaca, err := r.numerBiezacy(ctx, kodAgenta)
	if err != nil {
		return WersjaAgenta{}, err
	}
	return r.wersjaPoNumerze(ctx, kodAgenta, biezaca)
}

// przywrocPoziomyPamieciWersji zastępuje poziomy pamięci eksperta poziomami
// utrwalonymi w przywracanej wersji. Wycinek pusty znaczy pamięć wyłączoną
// i jest żądaniem, nie brakiem żądania.
func przywrocPoziomyPamieciWersji(ctx context.Context, transakcja *sql.Tx,
	kodAgenta string, poziomy []string) error {

	if _, err := transakcja.ExecContext(ctx,
		`DELETE FROM agent_pamiec_poziom WHERE agent_id = (SELECT id FROM agent WHERE kod = ?)`,
		kodAgenta); err != nil {
		return fmt.Errorf("dane: nie można wyczyścić poziomów pamięci eksperta %q: %w", kodAgenta, err)
	}
	for _, poziom := range poziomy {
		if _, err := transakcja.ExecContext(ctx,
			`INSERT INTO agent_pamiec_poziom (agent_id, poziom)
			 VALUES ((SELECT id FROM agent WHERE kod = ?), ?)
			 ON CONFLICT(agent_id, poziom) DO NOTHING`, kodAgenta, poziom); err != nil {
			return fmt.Errorf("dane: nie można przywrócić poziomu %q eksperta %q: %w",
				poziom, kodAgenta, err)
		}
	}
	return nil
}

// wersjaPoNumerze oddaje jedną migawkę. Brak wiersza wraca jako ErrBrakWiersza,
// żeby warstwa wyższa odróżniła „takiej wersji nie było” od „odczyt padł”.
func (r *repozytoriumWersjiAgenta) wersjaPoNumerze(ctx context.Context, kodAgenta string,
	numer int) (WersjaAgenta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wersjaAgentaPoNumerze)
	if err != nil {
		return WersjaAgenta{}, err
	}
	wersja, err := odczytajWersjeAgenta(polecenie.QueryRowContext(ctx, kodAgenta, numer))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaAgenta{}, fmt.Errorf("dane: ekspert %q nie ma wersji %d: %w",
			kodAgenta, numer, ErrBrakWiersza)
	}
	if err != nil {
		return WersjaAgenta{}, fmt.Errorf("dane: nie można odczytać wersji %d eksperta %q: %w",
			numer, kodAgenta, err)
	}
	return wersja, nil
}

// numerBiezacy odczytuje licznik wersji eksperta po zapisie zmiany, na podstawie kodu trwałego eksperta.
func (r *repozytoriumWersjiAgenta) numerBiezacy(ctx context.Context, kodAgenta string) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, numerWersjiBiezacej)
	if err != nil {
		return 0, err
	}
	var numer int
	err = polecenie.QueryRowContext(ctx, kodAgenta).Scan(&numer)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("dane: ekspert %q nie istnieje: %w", kodAgenta, ErrBrakWiersza)
	}
	if err != nil {
		return 0, fmt.Errorf("dane: nieznana wersja eksperta %q: %w", kodAgenta, err)
	}
	return numer, nil
}

// odczytajWersjeAgenta składa migawkę wersji eksperta wprost z jednego wiersza wyniku zapytania do bazy.
func odczytajWersjeAgenta(wiersz skaner) (WersjaAgenta, error) {
	var wersja WersjaAgenta
	var kanal, model, transport sql.NullString
	var aktywny int
	var poziomy string
	err := wiersz.Scan(&wersja.Identyfikator, &wersja.Numer, &wersja.Nazwa, &wersja.Opis, &wersja.InstrukcjeSystemowe,
		&kanal, &model, &transport, &wersja.ParametryJSON, &wersja.ImieWlasne, &wersja.Favikon,
		&wersja.UstawieniaJSON, &wersja.TrybNakladki, &aktywny, &wersja.Autor, &wersja.Powod,
		&wersja.Zapisano, &wersja.Widocznosc, &poziomy)
	if err != nil {
		return WersjaAgenta{}, err
	}
	wersja.PoziomyPamieci = poziomyMigawkiWersji(poziomy)
	wersja.KanalKod = tekstZKolumny(kanal)
	wersja.Model = tekstZKolumny(model)
	wersja.Transport = tekstZKolumny(transport)
	wersja.Aktywny = aktywny == 1
	return wersja, nil
}

// poziomyMigawkiWersji rozbiera kolumnę `poziomy_pamieci` na wycinek. Kolumna
// pusta znaczy pamięć wyłączoną w tej wersji — wycinek pusty, nie wycinek
// z jednym pustym napisem.
func poziomyMigawkiWersji(tresc string) []string {
	przyciety := strings.TrimSpace(tresc)
	if przyciety == "" {
		return nil
	}
	pozycje := strings.Split(przyciety, ",")
	wynik := make([]string, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if poziom := strings.TrimSpace(pozycja); poziom != "" {
			wynik = append(wynik, poziom)
		}
	}
	return wynik
}

// oznaczZmiany wypełnia `ZmienionePola` każdej migawki poza najstarszą.
// Wykaz przychodzi malejąco, więc poprzedniczką wiersza `i` jest wiersz `i+1`.
func oznaczZmiany(historia []WersjaAgenta) {
	for i := 0; i+1 < len(historia); i++ {
		historia[i].ZmienionePola = roznicaWersji(historia[i+1], historia[i])
	}
}

// roznicaWersji wylicza nazwy kolumn różniące dwie migawki. Nazwy są nazwami
// kolumn bazy (po polsku) — przekład na pola kontraktu należy do rdzenia,
// nie do warstwy danych.
func roznicaWersji(przed, po WersjaAgenta) []string {
	zmiany := []string{}
	porownaj := []struct {
		kolumna   string
		przed, po string
	}{
		{"nazwa", przed.Nazwa, po.Nazwa},
		{"opis", przed.Opis, po.Opis},
		{"instrukcje_systemowe", przed.InstrukcjeSystemowe, po.InstrukcjeSystemowe},
		{"kanal_kod", wartoscWskaznika(przed.KanalKod), wartoscWskaznika(po.KanalKod)},
		{"model", wartoscWskaznika(przed.Model), wartoscWskaznika(po.Model)},
		{"transport", wartoscWskaznika(przed.Transport), wartoscWskaznika(po.Transport)},
		{"parametry_json", przed.ParametryJSON, po.ParametryJSON},
		{"imie_wlasne", przed.ImieWlasne, po.ImieWlasne},
		{"favikon", przed.Favikon, po.Favikon},
		{"ustawienia_json", przed.UstawieniaJSON, po.UstawieniaJSON},
		{"tryb_nakladki", przed.TrybNakladki, po.TrybNakladki},
	}
	for _, para := range porownaj {
		if para.przed != para.po {
			zmiany = append(zmiany, para.kolumna)
		}
	}
	if przed.Aktywny != po.Aktywny {
		zmiany = append(zmiany, "aktywny")
	}
	return zmiany
}

// wartoscWskaznika sprowadza kolumnę pustą i wskaźnik pusty do jednego napisu,
// żeby porównanie nie zgłaszało zmiany tam, gdzie zmiany nie było.
func wartoscWskaznika(tekst *string) string {
	if tekst == nil {
		return ""
	}
	return *tekst
}
